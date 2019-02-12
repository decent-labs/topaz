package controllers

import (
	"net/http"

	"github.com/decentorganization/topaz/api/services"
	"github.com/decentorganization/topaz/shared/models"
	"github.com/gorilla/mux"
)

// GetBatches ...
func GetBatches(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	h, bs := services.GetBatches(
		r.Context().Value(models.AuthUser).(*models.User),
		mux.Vars(r)["appId"],
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(h)
	w.Write(bs)
}

// GetBatch ...
func GetBatch(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	h, b := services.GetBatch(
		r.Context().Value(models.AuthUser).(*models.User),
		mux.Vars(r)["appId"],
		mux.Vars(r)["id"],
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(h)
	w.Write(b)
}
