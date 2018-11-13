package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	dotenv "github.com/joho/godotenv"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
)

var db *sql.DB

func main() {
	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}

	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		log.Println("couldn't execute migrations")
		log.Fatal(err)
	}

	fmt.Printf("Applied %d migrations!\n", n)
}

func init() {
	err := dotenv.Load(".env")
	if err != nil {
		log.Fatal("couldn't load .env file")
	}

	dbConn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
	)

	_db, err := sql.Open("postgres", dbConn)
	if err != nil {
		log.Println("error opening lazy connection to DB... weird")
		log.Fatal(err)
	}

	_db.Ping()
	if err != nil {
		log.Println("error pinging DB for first connection")
		log.Fatal(err)
	}

	db = _db
}
