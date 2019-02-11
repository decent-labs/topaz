package controllers

import (
	"net/http"

	"github.com/decentorganization/topaz/api/services"
	"github.com/decentorganization/topaz/shared/models"
	"github.com/gorilla/mux"
)

func buildObjectContext(r *http.Request, oid string) *models.Object {
	o := models.Object{
		App: &models.App{
			ID: mux.Vars(r)["appId"],
			User: &models.User{
				ID: r.Context().Value(models.UserID).(string),
			},
		},
	}

	if oid != "" {
		o.ID = oid
	}

	return &o
}

// CreateObject ...
func CreateObject(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	o := buildObjectContext(r, "")
	h, ro := services.CreateObject(o)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(h)
	w.Write(ro)
}

// GetObjects ...
func GetObjects(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	o := buildObjectContext(r, "")
	h, ros := services.GetObjects(o)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(h)
	w.Write(ros)
}

// GetObject ...
func GetObject(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	o := buildObjectContext(r, mux.Vars(r)["id"])
	h, ro := services.GetObject(o)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(h)
	w.Write(ro)
}
