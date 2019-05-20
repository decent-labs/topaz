package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/decentorganization/topaz/api/services"
	"github.com/decentorganization/topaz/shared/models"
)

// CreateMarketingEmail ...
func CreateMarketingEmail(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	m := new(models.SendgridEmail)
	d := json.NewDecoder(r.Body)
	d.Decode(&m)

	rs := services.CreateMarketingEmail(m)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(rs)
}
