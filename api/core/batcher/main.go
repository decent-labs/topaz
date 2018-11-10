package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	log.Println("BATCH: BEGINNING MAIN.")
}

func init() {
	log.Println("BATCH: HELLO, WORLD!")

	i, err := strconv.Atoi(os.Getenv("STARTUP_SLEEP"))
	if err != nil {
		log.Fatalf("missing valid STARTUP_SLEEP environment variable: %s", err.Error())
	}
	time.Sleep(time.Duration(i) * time.Second)

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

	log.Println("BATCH: 'db' INSTANCE CREATED!")

	stmt := `
		select distinct a.id
		from apps a
			inner join objects o on o.app_id = a.id
		where ((a.last_flushed is null) or (now() - a.lashed_flush >= a.interval))
			and (o.flush_id is null)
	`

	rows, err := db.Query(stmt)
	if err != nil {
		log.Printf("error executing app-finding query: %s", err.Error())
		return
	}
	defer rows.Close()

	log.Println("BATCH: 'apps' FOUND.")

	for rows.Next() {
		var id string

		err = rows.Scan(&id)
		if err != nil {
			log.Printf("error: %s", err.Error())
			continue
		}

		log.Printf("app found: %s", id)
	}

	// TODO: Send app ids (id int) to "Flush" for processing.

	log.Println("BATCH: BEGINNING FLUSH ROUTINE.")
}
