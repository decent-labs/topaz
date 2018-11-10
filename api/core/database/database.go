package database

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var Manager *gorm.DB

func init() {
	dbConn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable", os.Getenv("PQ_HOST"), os.Getenv("PQ_PORT"), os.Getenv("PQ_USER"), os.Getenv("PQ_NAME"))

	var err error
	Manager, err = gorm.Open("postgres", dbConn)
	if err != nil {
		log.Fatalf("couldn't even pretend to open database connection: %s", err.Error())
	}

	Manager.LogMode(true)

	if err := Manager.DB().Ping(); err != nil {
		log.Fatal(err)
	}
}
