package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/decentorganization/topaz/api/services"
	"github.com/decentorganization/topaz/shared/models"
)

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
