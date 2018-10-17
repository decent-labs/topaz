package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

// Find all users who are due for a flush and call the 'flush' service for them.
func main() {
	i, err := strconv.Atoi(os.Getenv("STARTUP_SLEEP"))
	if err != nil {
		log.Fatalf("missing valid STARTUP_SLEEP environment variable: %s", err.Error())
	}
	time.Sleep(time.Duration(i) * time.Second)

	conn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
		os.Getenv("PQ_HOST"),
		os.Getenv("PQ_PORT"),
		os.Getenv("PQ_USER"),
		os.Getenv("PQ_NAME"),
	)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Fatalf("couldn't even pretend to open database connection: %s", err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("couldn't ping database: %s", err.Error())
	}

	i, err = strconv.Atoi(os.Getenv("DISPATCH_RATE"))
	if err != nil {
		log.Fatalf("missing valid DISPATCH_RATE environment variable: %s", err.Error())
	}

	log.Println("Wake up, dispatch...")

	for range time.Tick(time.Duration(i) * time.Second) {
		log.Println("tick")

		stmt := `
			select distinct u.id
			from users u
				inner join objects o on o.user_id = u.id
			where ((u.flushed_at is null) or (now() - u.flushed_at >= u.interval))
				and (o.flush_id is null)
		`

		rows, err := db.Query(stmt)
		if err != nil {
			log.Printf("error executing user-finding query: %s", err.Error())
			return
		}
		defer rows.Close()

		for rows.Next() {
			var id string

			err = rows.Scan(&id)
			if err != nil {
				log.Printf("error scanning row into user id var: %s", err.Error())
				continue
			}

			url := fmt.Sprintf("http://%s:%s", os.Getenv("FLUSH_HOST"), os.Getenv("FLUSH_PORT"))
			sr := strings.NewReader(id)
			_, err = http.Post(url, "application/octet-stream", sr)
			if err != nil {
				log.Printf("error dispatching user id '%s' to flush service: %s", id, err.Error())
				continue
			}
		}

		log.Println("tock")
	}
}
