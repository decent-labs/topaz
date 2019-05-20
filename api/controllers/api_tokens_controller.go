package controllers

import (
	"net/http"

	"github.com/decentorganization/topaz/api/services"
	"github.com/decentorganization/topaz/shared/models"
	"github.com/gorilla/mux"
)

// CreateAPIToken ...
func CreateAPIToken(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	rs, ar := services.CreateAPIToken(
		r.Context().Value(models.AuthUser).(*models.User),
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(rs)
	w.Write(ar)
}

// GetAPITokens ...
func GetAPITokens(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	rs, as := services.GetAPITokens(
		r.Context().Value(models.AuthUser).(*models.User),
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(rs)
	w.Write(as)
}

// GetAPIToken ...
func GetAPIToken(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	rs, ra := services.GetAPIToken(
		r.Context().Value(models.AuthUser).(*models.User),
		mux.Vars(r)["id"],
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(rs)
	w.Write(ra)
}
