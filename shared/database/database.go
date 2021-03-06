package database

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"

	// gorm requires a "dialect" is imported to communicate with postgres
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Manager is used to access our database across the application
var Manager *gorm.DB

func init() {
	godotenv.Load()

	dbConn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
	)

	manager, err := gorm.Open("postgres", dbConn)
	if err != nil {
		log.Fatalf("couldn't even pretend to open database connection: %s", err.Error())
	}
	Manager = manager

	Manager.LogMode(os.Getenv("GO_ENV") != "production")

	if err := Manager.DB().Ping(); err != nil {
		log.Fatal(err)
	}

	maxOpenConns, _ := strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONNS"))
	maxIdleConns, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNS"))
	Manager.DB().SetMaxOpenConns(maxOpenConns)
	Manager.DB().SetMaxIdleConns(maxIdleConns)
}
