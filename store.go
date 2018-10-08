package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
	_ "github.com/lib/pq"
)

// Take the request body and store it in IPFS, then store the resulting hash in the `objects` table.
func requestHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	shConn := fmt.Sprintf("%s:%s", os.Getenv("IPFS_HOST"), os.Getenv("IPFS_PORT"))
	sh := shell.NewShell(shConn)
	br := bytes.NewReader(b)

	hash, err := sh.Add(br)
	if err != nil {
		log.Fatal(err)
	}

	dbConn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
		os.Getenv("PQ_HOST"),
		os.Getenv("PQ_PORT"),
		os.Getenv("PQ_USER"),
		os.Getenv("PQ_NAME"),
	)

	db, err := sql.Open("postgres", dbConn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}

	id := strings.TrimSpace(fmt.Sprintf("%s", out))

	stmt := fmt.Sprintf("insert into objects (id, hash, user_id) values ('%s', '%s', '%s')",
		id,
		hash,
		os.Getenv("TOPAZ_USER"),
	)

	_, err = db.Exec(stmt)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("'%s' is now in the queue.", hash)
}

func main() {
	http.HandleFunc("/", requestHandler)

	log.Println("Wake up, store...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
