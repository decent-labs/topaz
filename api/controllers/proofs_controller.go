package controllers

import (
	"net/http"

	"github.com/decentorganization/topaz/api/services"
	"github.com/decentorganization/topaz/shared/models"
	"github.com/gorilla/mux"
)

// GetProofs ...
func GetProofs(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	h, p := services.GetProofs(
		r.Context().Value(models.AuthUser).(*models.User),
		mux.Vars(r)["appId"],
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(h)
	w.Write(p)
}

// GetProof ...
func GetProof(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	h, p := services.GetProof(
		r.Context().Value(models.AuthUser).(*models.User),
		mux.Vars(r)["appId"],
		mux.Vars(r)["id"],
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(h)
	w.Write(p)
}
