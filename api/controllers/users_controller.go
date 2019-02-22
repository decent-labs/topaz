package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/decentorganization/topaz/api/services"
	"github.com/decentorganization/topaz/shared/models"
)

// CreateUser ...
func CreateUser(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ru := new(models.User)
	d := json.NewDecoder(r.Body)
	d.Decode(&ru)

	rs, u := services.CreateUser(ru)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(rs)
	w.Write(u)
}

// GetUser ...
func GetUser(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	h, u := services.GetUser(
		r.Context().Value(models.AuthUser).(*models.User),
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(h)
	w.Write(u)
}

// EditUser ...
func EditUser(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ru := new(models.User)
	d := json.NewDecoder(r.Body)
	d.Decode(&ru)

	rs, u := services.EditUser(
		r.Context().Value(models.AuthUser).(*models.User),
		ru,
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(rs)
	w.Write(u)
}
