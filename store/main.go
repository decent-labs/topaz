package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	shell "github.com/ipfs/go-ipfs-api"
	_ "github.com/lib/pq"
)

var sh *shell.Shell
var db *sql.DB

// StoreResponse defines what gets returned on store route
type StoreResponse struct {
	Hash string
}

// Take the request body and store it in IPFS, then store the resulting hash in the 'objects' table.
func requestHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("starting store service handler")

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("error reading store service request body: %s", err.Error()),
			http.StatusBadRequest,
		)
		return
	}

	br := bytes.NewReader(b)

	hash, err := sh.Add(br)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("error adding file to ipfs through shell: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	stmt := fmt.Sprintf("insert into objects (hash, user_id) values ('%s', '%s')",
		hash,
		os.Getenv("TOPAZ_USER"),
	)

	_, err = db.Exec(stmt)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("error adding ipfs hash to objects table: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	log.Printf("processed incoming data with hash: %s", hash)

	sr := StoreResponse{hash}
	w.Header().Set("Content-Type", "application/vnd.api+json")
	err = json.NewEncoder(w).Encode(sr)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("error encoding store response from store service: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	log.Println("finished with store service handler")
}

func main() {
	shConn := fmt.Sprintf("%s:%s", os.Getenv("IPFS_HOST"), os.Getenv("IPFS_PORT"))
	sh = shell.NewShell(shConn)

	dbConn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
		os.Getenv("PQ_HOST"),
		os.Getenv("PQ_PORT"),
		os.Getenv("PQ_USER"),
		os.Getenv("PQ_NAME"),
	)

	_db, err := sql.Open("postgres", dbConn)
	if err != nil {
		log.Fatalf("couldn't even pretend to open database connection: %s", err.Error())
	}
	defer _db.Close()

	err = _db.Ping()
	if err != nil {
		log.Fatalf("couldn't ping database: %s", err.Error())
	}

	db = _db

	http.HandleFunc("/", requestHandler)

	log.Println("wake up, store...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
