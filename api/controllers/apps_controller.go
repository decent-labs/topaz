package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/decentorganization/topaz/api/services"
	"github.com/decentorganization/topaz/shared/models"
	"github.com/gorilla/mux"
)

// CreateApp ...
func CreateApp(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ra := new(models.App)
	d := json.NewDecoder(r.Body)
	d.Decode(&ra)

	rs, ar := services.CreateApp(
		r.Context().Value(models.AuthUser).(*models.User),
		ra,
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(rs)
	w.Write(ar)
}

// GetApps ...
func GetApps(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	rs, as := services.GetApps(
		r.Context().Value(models.AuthUser).(*models.User),
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(rs)
	w.Write(as)
}

// GetApp ...
func GetApp(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	rs, ra := services.GetApp(
		r.Context().Value(models.AuthUser).(*models.User),
		mux.Vars(r)["id"],
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(rs)
	w.Write(ra)
}
