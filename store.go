package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	shell "github.com/ipfs/go-ipfs-api"
	_ "github.com/lib/pq"
)

// Take the request body and store it in IPFS, then store the resulting hash in the `queue` table.
func requestHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	sh := shell.NewShell("localhost:5001")
	br := bytes.NewReader(b)

	hash, err := sh.Add(br)
	if err != nil {
		log.Fatal(err)
	}

	dbConn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s",
		"localhost",
		"5432",
		"postgres",
		"postgres",
		"disable",
	)

	db, err := sql.Open("postgres", dbConn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	insert := fmt.Sprintf("INSERT INTO queue (hash) VALUES ('%s') RETURNING hash", hash)

	_, err = db.Exec(insert)
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
