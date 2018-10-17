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
	time.Sleep(5 * time.Second)

	conn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
		os.Getenv("PQ_HOST"),
		os.Getenv("PQ_PORT"),
		os.Getenv("PQ_USER"),
		os.Getenv("PQ_NAME"),
	)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Println("error getting into DB")
		log.Fatal(err)
	}
	defer db.Close()

	log.Println("Wake up, dispatch...")

	for range time.Tick(10 * time.Second) {
		go func() {
			stmt := `
				select distinct u.id
				from users u
					inner join objects o on o.user_id = u.id
				where ((u.flushed_at is null) or (now() - u.flushed_at >= u.interval))
					and (o.flush_id is null)
			`

			rows, err := db.Query(stmt)
			if err != nil {
				log.Println("error running query to find users")
				log.Fatal(err)
			}

			for rows.Next() {
				var id string

				err = rows.Scan(&id)
				if err != nil {
					log.Println("error reading results from query")
					log.Fatal(err)
				}

				log.Printf("handling current id: '%s'", id)

				url := fmt.Sprintf("http://%s:%s", os.Getenv("FLUSH_HOST"), os.Getenv("FLUSH_PORT"))
				sr := strings.NewReader(id)
				_, err = http.Post(url, "application/octet-stream", sr)
				if err != nil {
					log.Println("error posting user id to flush service")
					log.Fatal(err)
				}
			}
		}()
	}
}
