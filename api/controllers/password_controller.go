package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/decentorganization/topaz/api/services"
	"github.com/decentorganization/topaz/shared/models"
)

// UpdatePassword ...
func UpdatePassword(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	rp := new(models.UpdatePassword)
	d := json.NewDecoder(r.Body)
	d.Decode(&rp)

	rs, u := services.UpdatePassword(
		r.Context().Value(models.AuthUser).(*models.User),
		rp,
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(rs)
	w.Write(u)
}
