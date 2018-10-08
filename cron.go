package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	conn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
		os.Getenv("PQ_HOST"),
		os.Getenv("PQ_PORT"),
		os.Getenv("PQ_USER"),
		os.Getenv("PQ_NAME"),
	)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	stmt := "SELECT id FROM users WHERE now() - flushed_at > interval"

	rows, err := db.Query(stmt)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var id string

		err = rows.Scan(&id)
		if err != nil {
			log.Fatal(err)
		}

		log.Println(id)
	}
}
