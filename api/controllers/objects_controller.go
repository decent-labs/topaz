package controllers

import (
	"net/http"

	"github.com/decentorganization/topaz/api/services"
	"github.com/gorilla/mux"
)

// CreateObject ...
func CreateObject(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	h, o := services.CreateObject(r.Context(), mux.Vars(r)["appId"])

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(h)
	w.Write(o)
}

// GetObjects ...
func GetObjects(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	h, os := services.GetObjects(r.Context(), mux.Vars(r)["appId"])

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(h)
	w.Write(os)
}

// GetObject ...
func GetObject(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	h, o := services.GetObject(r.Context(), mux.Vars(r)["appId"], mux.Vars(r)["id"])

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(h)
	w.Write(o)
}
