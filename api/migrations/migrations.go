package migrations

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gobuffalo/packr"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
)

var db *sql.DB

// Attempt ...
func Attempt() {
	migrations := &migrate.PackrMigrationSource{
		Box: packr.NewBox("sql"),
	}

	m, err := merkleRootMigration(migrations)
	if err != nil {
		log.Println("couldn't execute merkle root migration setup scripts")
		log.Fatal(err)
	}

	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		log.Println("couldn't execute migrations")
		log.Fatal(err)
	}

	n = n + m

	fmt.Printf("Applied %d migrations!\n", n)
}

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("couldn't load dotenv:", err.Error())
	}

	dbConn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
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
