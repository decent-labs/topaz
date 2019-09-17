package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/decentorganization/topaz/api/services"
	"github.com/decentorganization/topaz/shared/models"
	"github.com/gorilla/mux"
)

type appController struct {
	service services.AppService
}

// AppController ...
type AppController interface {
	CreateApp(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
	GetApps(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
	GetApp(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

// CreateApp ...
func (c appController) CreateApp(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ra := new(models.App)
	d := json.NewDecoder(r.Body)
	d.Decode(&ra)

	rs, ar := c.service.CreateApp(
		r.Context().Value(models.AuthUser).(*models.User),
		ra,
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(rs)
	w.Write(ar)
}

// GetApps ...
func (c appController) GetApps(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	rs, as := c.service.GetApps(
		r.Context().Value(models.AuthUser).(*models.User),
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(rs)
	w.Write(as)
}

// GetApp ...
func (c appController) GetApp(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	rs, ra := c.service.GetApp(
		r.Context().Value(models.AuthUser).(*models.User),
		mux.Vars(r)["id"],
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(rs)
	w.Write(ra)
}

// NewAppController ...
func NewAppController(s services.AppService) AppController {
	return appController{s}
}
