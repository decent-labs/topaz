package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/decentorganization/topaz/api/services"
	"github.com/decentorganization/topaz/shared/models"
)

func CreateObject(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	uid := r.Context().Value(models.UserID).(string)

	ro := new(models.Object)
	d := json.NewDecoder(r.Body)
	d.Decode(&ro)

	h, o := services.CreateObject(ro, uid)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(h)
	w.Write(o)
}


