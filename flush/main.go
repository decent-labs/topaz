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

// Message is what gest posted to ethereum service
type Message struct {
	Address string
	Hash    string
}

// TXResp is what gets returned
type TXResp struct {
	TX string
}

var sh *shell.Shell
var db *sql.DB

func requestHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("starting flush service handler")

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("error reading flush service request body: %s", err.Error()),
			http.StatusBadRequest,
		)
		return
	}

	userID := string(b)

	flushStmt := fmt.Sprintf(
		"insert into flushes (user_id) values ('%s') returning id, created_at;",
		userID,
	)

	log.Printf("Beginning flush for user '%s'.", userID)

	flushRows, err := db.Query(flushStmt)
	if err != nil {
		log.Fatal(err)
	}
	defer flushRows.Close()

	var flushID string
	var flushCreatedAt string

	for flushRows.Next() {
		err = flushRows.Scan(&flushID, &flushCreatedAt)
		if err != nil {
			log.Fatal(err)
		}
	}

	objStmt := fmt.Sprintf("select id, hash from objects where user_id = '%s' AND flush_id is null;",
		userID,
	)

	objRows, err := db.Query(objStmt)
	if err != nil {
		log.Fatal(err)
	}
	defer objRows.Close()

	dir, err := sh.NewObject("unixfs-dir")
	if err != nil {
		log.Fatal(err)
	}

	for objRows.Next() {
		var id string
		var hash string

		err = objRows.Scan(&id, &hash)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Flushing object '%s' with hash '%s'.", id, hash)

		dir, err = sh.PatchLink(dir, hash, hash, true)
		if err != nil {
			log.Fatal(err)
		}

		objUpdateStmt := fmt.Sprintf("update objects set flush_id = '%s' where id = '%s'", flushID, id)

		_, err := db.Exec(objUpdateStmt)
		if err != nil {
			log.Fatal(err)
		}
	}

	flushUpdateStmt := fmt.Sprintf("update flushes set hash = '%s' where id = '%s'", dir, flushID)

	_, err = db.Exec(flushUpdateStmt)
	if err != nil {
		log.Fatal(err)
	}

	userUpdateStmt := fmt.Sprintf("update users set flushed_at = '%s' where id = '%s'",
		flushCreatedAt,
		userID,
	)

	_, err = db.Exec(userUpdateStmt)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Finished flush '%s' for user '%s'.", flushID, userID)

	log.Printf("Dir: %s", dir)

	url := fmt.Sprintf("http://%s:%s/store", os.Getenv("ETH_HOST"), os.Getenv("ETH_PORT"))

	m := Message{os.Getenv("ETH_ADDRESS"), dir}

	b, err := json.Marshal(m)
	if err != nil {
		log.Println("could not create JSON out of Message")
		log.Fatal(err)
	}

	r := bytes.NewReader(b)
	resp, err := http.Post(url, "application/octet-stream", r)
	if err != nil {
		log.Println("failed posting data to ethereum service")
		log.Fatal(err)
	}
	defer resp.Body.Close()

	txresp := new(TXResp)
	json.NewDecoder(resp.Body).Decode(txresp)

	log.Println("finished with flush service handler")
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

	log.Println("Wake up, flush...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
