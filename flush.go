package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	shell "github.com/ipfs/go-ipfs-api"
	_ "github.com/lib/pq"
)

func flush(p []string) {
	sh := shell.NewShell("localhost:5001")

	dir, err := sh.NewObject("unixfs-dir")
	if err != nil {
		fmt.Println(err)
	} else {
		for i, h := range p {
			dirHash, err := sh.PatchLink(dir, h, h, true)
			if err != nil {
				fmt.Println(err)
			} else {
				log.Printf("%d: %s\n", i, dirHash)
			}
		}
	}
}

func getQueue() {
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

	rows, err := db.Query("SELECT hash FROM queue")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	sh := shell.NewShell("localhost:5001")

	dir, err := sh.NewObject("unixfs-dir")
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var hash string

		err = rows.Scan(&hash)
		if err != nil {
			log.Fatal(err)
		}

		dirHash, err := sh.PatchLink(dir, hash, hash, true)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("%s is being pulled from the queue.", dirHash)
	}
}

func main() {
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for t := range ticker.C {
			fmt.Println("Flush at", t)
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

			rows, err := db.Query("SELECT hash FROM queue")
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()

			sh := shell.NewShell("localhost:5001")

			dir, err := sh.NewObject("unixfs-dir")
			if err != nil {
				log.Fatal(err)
			}

			for rows.Next() {
				var hash string

				err = rows.Scan(&hash)
				if err != nil {
					log.Fatal(err)
				}

				dirHash, err := sh.PatchLink(dir, hash, hash, true)
				if err != nil {
					log.Fatal(err)
				}

				log.Println(dirHash)
			}
		}
	}()

	log.Println("Wake up, flush...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
