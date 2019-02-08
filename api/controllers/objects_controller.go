package controllers

import (
	"net/http"

	"github.com/decentorganization/topaz/api/services"
	"github.com/decentorganization/topaz/shared/models"
	"github.com/gorilla/mux"
)

func buildContext(r *http.Request) (string, models.Object) {
	uid := r.Context().Value(models.UserID).(string)
	o := models.Object{AppID: mux.Vars(r)["appId"]}
	return uid, o
}

func CreateObject(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	uid, o := buildContext(r)
	h, ro := services.CreateObject(&o, uid)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(h)
	w.Write(ro)
}

func GetObjects(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	uid, o := buildContext(r)
	h, ros := services.GetObjects(&o, uid)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(h)
	w.Write(ros)
}

func GetObject(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	uid, o := buildContext(r)
	o.ID = mux.Vars(r)["id"]
	h, ro := services.GetObject(&o, uid)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(h)
	w.Write(ro)
}
