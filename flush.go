package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
	_ "github.com/lib/pq"
)

// Given a specific user in our system, link any queued objects to an IPFS directory.
func flush(userId string) {
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

	flushId := strings.TrimSpace(fmt.Sprintf("%s", out))

	flushStmt := fmt.Sprintf(
		"insert into flushes (id, user_id) values ('%s', '%s') returning created_at;",
		flushId,
		userId,
	)

	flushRows, err := db.Query(flushStmt)
	if err != nil {
		log.Fatal(err)
	}
	defer flushRows.Close()

	var flushCreatedAt string

	for flushRows.Next() {
		err = flushRows.Scan(&flushCreatedAt)
		if err != nil {
			log.Fatal(err)
		}
	}

	objStmt := fmt.Sprintf("selt id, hash from objects where user_id = '%s' AND flush_id is null;",
		userId,
	)

	objRows, err := db.Query(objStmt)
	if err != nil {
		log.Fatal(err)
	}
	defer objRows.Close()

	shConn := fmt.Sprintf("%s:%s", os.Getenv("IPFS_HOST"), os.Getenv("IPFS_PORT"))
	sh := shell.NewShell(shConn)

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

		dir, err = sh.PatchLink(dir, hash, hash, true)
		if err != nil {
			log.Fatal(err)
		}

		objUpdateStmt := fmt.Sprintf("update objects set flush_id = '%s' where id = '%s'", flushId, id)

		_, err := db.Exec(objUpdateStmt)
		if err != nil {
			log.Fatal(err)
		}
	}

	flushUpdateStmt := fmt.Sprintf("update flushes set hash = '%s' where id = '%s'", dir, flushId)

	_, err = db.Exec(flushUpdateStmt)
	if err != nil {
		log.Fatal(err)
	}

	userUpdateStmt := fmt.Sprintf("update users set flushed_at = '%s' where id = '%s'",
		flushCreatedAt,
		userId,
	)

	_, err = db.Exec(userUpdateStmt)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

}
