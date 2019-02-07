package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/decentorganization/topaz/api/services"
	"github.com/decentorganization/topaz/shared/models"
)

// AdminLogin returns the result of an attempted login by an admin user
func AdminLogin(w http.ResponseWriter, r *http.Request) {
	requestUser := new(models.User)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestUser)

	responseStatus, token := services.AdminLogin(requestUser)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseStatus)
	w.Write(token)
}

// AdminRefreshToken returns the result of an attempted token refresh by an admin user
func AdminRefreshToken(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	requestUser := new(models.User)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestUser)

	responseStatus, token := services.AdminRefreshToken(requestUser)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseStatus)
	w.Write(token)
}

// AdminLogout returns the result of an attempted logout by an admin user
func AdminLogout(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Set("Content-Type", "application/json")

	err := services.AdminLogout(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

// AppLogin returns the result of an attempted login by an 'app'
func AppLogin(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	uid := r.Context().Value(models.UserID).(uint)

	requestApp := new(models.App)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestApp)
	requestApp.UserID = uid

	responseStatus, token := services.AppLogin(requestApp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseStatus)
	w.Write(token)
}

// AppRefreshToken returns the result of an attempted token refresh by an 'app'
func AppRefreshToken(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	aid := r.Context().Value(models.AppID).(uint)

	requestApp := new(models.App)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestApp)
	requestApp.ID = aid

	responseStatus, token := services.AppRefreshToken(requestApp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseStatus)
	w.Write(token)
}
