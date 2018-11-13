package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/decentorganization/topaz/api/services"
	"github.com/decentorganization/topaz/shared/models"
)

// NewApp allows a user to create a new entry in the 'apps' table
func NewApp(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	requestApp := new(models.App)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestApp)

	requestApp.UserID, _ = strconv.Atoi(r.Header.Get("userId"))
	responseStatus, app := services.NewApp(requestApp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseStatus)
	w.Write(app)
}
