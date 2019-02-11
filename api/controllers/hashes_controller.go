package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/decentorganization/topaz/api/services"
	"github.com/decentorganization/topaz/shared/models"
	"github.com/gorilla/mux"
)

// CreateHash ...
func CreateHash(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	rh := new(models.Hash)
	d := json.NewDecoder(r.Body)
	d.Decode(&rh)

	he, h := services.CreateHash(
		r.Context().Value(models.AuthUser).(*models.User),
		mux.Vars(r)["appId"],
		mux.Vars(r)["objectId"],
		rh,
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(he)
	w.Write(h)
}

// GetHashes ...
func GetHashes(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	he, h := services.GetHashes(
		r.Context().Value(models.AuthUser).(*models.User),
		mux.Vars(r)["appId"],
		mux.Vars(r)["objectId"],
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(he)
	w.Write(h)
}

// GetHash ...
func GetHash(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	he, h := services.GetHash(
		r.Context().Value(models.AuthUser).(*models.User),
		mux.Vars(r)["appId"],
		mux.Vars(r)["objectId"],
		mux.Vars(r)["id"],
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(he)
	w.Write(h)
}
