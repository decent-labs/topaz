package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	conn := "host=localhost port=5432 user=postgres dbname=postgres sslmode=disable"

	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Println(err)
	} else {
		defer db.Close()

		_, err := db.Exec(`CREATE TABLE users (
			"uuid" uuid primary key
		);`)
		if err != nil {
			log.Println(err)
		}

		_, err = db.Exec(`CREATE TABLE queue (
			"uuid" uuid primary key,
			"hash" varchar(255),
			"user" uuid references users(uuid)
		);`)
		if err != nil {
			log.Println(err)
		}
	}
}
