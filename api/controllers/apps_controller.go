package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/decentorganization/topaz/api/services"
	"github.com/decentorganization/topaz/shared/models"
	"github.com/gorilla/mux"
)

func buildAppContext(r *http.Request, aid string) *models.App {
	a := models.App{
		User: &models.User{
			ID: r.Context().Value(models.UserID).(string),
		},
	}

	if aid != "" {
		a.ID = aid
	}

	return &a
}

// CreateApp ...
func CreateApp(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	a := buildAppContext(r, "")

	ra := new(models.App)
	d := json.NewDecoder(r.Body)
	d.Decode(&ra)

	rs, ar := services.CreateApp(a, ra)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(rs)
	w.Write(ar)
}

// GetApps ...
func GetApps(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	a := buildAppContext(r, "")

	rs, as := services.GetApps(a)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(rs)
	w.Write(as)
}

// GetApp ...
func GetApp(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	a := buildAppContext(r, mux.Vars(r)["id"])

	rs, ra := services.GetApp(a)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(rs)
	w.Write(ra)
}
