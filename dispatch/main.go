package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

// Find all users who are due for a flush and call the 'flush' service for them.
func main() {
	time.Sleep(time.Millisecond * 5000)

	conn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
		os.Getenv("PQ_HOST"),
		os.Getenv("PQ_PORT"),
		os.Getenv("PQ_USER"),
		os.Getenv("PQ_NAME"),
	)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for range time.Tick(time.Millisecond * 500) {
		go func() {
			stmt := "select id from users where flushed_at is null or now() - flushed_at >= interval"

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

				url := fmt.Sprintf("http://%s:%s", os.Getenv("FLUSH_HOST"), os.Getenv("FLUSH_PORT"))
				sr := strings.NewReader(id)
				_, err = http.Post(url, "application/octet-stream", sr)
				if err != nil {
					log.Fatal(err)
				}
			}
		}()
	}
}
