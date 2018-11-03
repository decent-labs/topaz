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

func Login(w http.ResponseWriter, r *http.Request) {
	requestUser := new(models.User)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestUser)

	dbConn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable", os.Getenv("PQ_HOST"), os.Getenv("PQ_PORT"), os.Getenv("PQ_USER"), os.Getenv("PQ_NAME"))
	db, err := gorm.Open("postgres", dbConn)
	if err != nil {
		log.Fatalf("couldn't even pretend to open database connection: %s", err.Error())
	}
	defer db.Close()

	responseStatus, token := services.Login(requestUser, db)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseStatus)
	w.Write(token)
}

func RefreshToken(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	requestUser := new(models.User)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestUser)

	responseStatus, token := services.RefreshToken(requestUser)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseStatus)
	w.Write(token)
}

func Logout(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	err := services.Logout(r)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
