package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/decentorganization/topaz/api/services"
	"github.com/decentorganization/topaz/shared/models"
)

// NewApp allows a user to create a new entry in the 'apps' table
func NewApp(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	uid := r.Context().Value(models.UserID).(string)

	requestApp := new(models.App)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestApp)
	requestApp.UserID = uid
	rs, a := services.CreateApp(ra)


	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(rs)
	w.Write(a)
}
