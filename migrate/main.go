package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
)

func main() {
	time.Sleep(5 * time.Second)

	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}

	conn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
		os.Getenv("PQ_HOST"),
		os.Getenv("PQ_PORT"),
		os.Getenv("PQ_USER"),
		os.Getenv("PQ_NAME"),
	)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Println("error opening lazy connection to DB... weird")
		log.Fatal(err)
	}

	db.Ping()
	if err != nil {
		log.Println("error pinging DB for first connection")
		log.Fatal(err)
	}

	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		log.Println("error executing migrations")
		log.Fatal(err)
	}

	fmt.Printf("Applied %d migrations!\n", n)
}
