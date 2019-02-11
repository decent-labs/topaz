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
	uid := r.Context().Value(models.UserID).(string)

	ra := new(models.App)
	d := json.NewDecoder(r.Body)
	d.Decode(&ra)
	ra.UserID = uid

	rs, a := services.CreateApp(ra)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(rs)
	w.Write(a)
}

// GetApps ...
func GetApps(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	uid := r.Context().Value(models.UserID).(string)

	rs, as := services.GetApps(uid)

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
