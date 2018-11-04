package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/decentorganization/topaz/api/services"
	"github.com/decentorganization/topaz/models"
	"github.com/jinzhu/gorm"
)

func NewUser(w http.ResponseWriter, r *http.Request) {
	requestUser := new(models.User)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestUser)

	dbConn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable", os.Getenv("PQ_HOST"), os.Getenv("PQ_PORT"), os.Getenv("PQ_USER"), os.Getenv("PQ_NAME"))
	db, err := gorm.Open("postgres", dbConn)
	if err != nil {
		log.Fatalf("couldn't even pretend to open database connection: %s", err.Error())
	}
	defer db.Close()

	responseStatus, user := services.NewUser(requestUser, db)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseStatus)
	w.Write(user)
}

func EditUser(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	requestUser := new(models.User)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestUser)

	responseStatus, user := services.EditUser(requestUser)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseStatus)
	w.Write(user)
}
