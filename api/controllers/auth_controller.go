package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/decentorganization/topaz/api/models"
	"github.com/decentorganization/topaz/api/services"
)

func AdminLogin(w http.ResponseWriter, r *http.Request) {
	requestUser := new(models.User)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestUser)

	responseStatus, token := services.AdminLogin(requestUser)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseStatus)
	w.Write(token)
}

func AdminRefreshToken(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	requestUser := new(models.User)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestUser)

	responseStatus, token := services.AdminRefreshToken(requestUser)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseStatus)
	w.Write(token)
}

func AdminLogout(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	err := services.AdminLogout(r)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func AppLogin(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	requestApp := new(models.App)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestApp)

	requestApp.UserID, _ = strconv.Atoi(r.Header.Get("userId"))
	responseStatus, token := services.AppLogin(requestApp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseStatus)
	w.Write(token)
}

func AppRefreshToken(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	requestApp := new(models.App)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestApp)

	s, _ := strconv.ParseUint(r.Header.Get("appId"), 10, 32)
	requestApp.ID = uint(s)
	responseStatus, token := services.AppRefreshToken(requestApp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseStatus)
	w.Write(token)
}
