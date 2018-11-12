package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

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
	err := dotenv.Load("../.env")
	if err != nil {
		log.Fatal("couldn't load .env file")
	}

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
