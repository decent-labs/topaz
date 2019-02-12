package routers

import (
	"github.com/decentorganization/topaz/api/authentication"
	"github.com/decentorganization/topaz/api/controllers"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// SetAppsRoutes provisions routes for 'apps'.
// All methods scoped to user
func SetAppsRoutes(r *mux.Router) *mux.Router {
	s := r.PathPrefix("/apps").Subrouter()

	// Create new app
	s.Handle("", negroni.New(
		negroni.HandlerFunc(authentication.Auth),
		negroni.HandlerFunc(controllers.CreateApp),
	)).Methods("POST")

	// Get single app
	s.Handle("/{id}", negroni.New(
		negroni.HandlerFunc(authentication.Auth),
		negroni.HandlerFunc(controllers.GetApp),
	)).Methods("GET")

	// Get all apps
	s.Handle("", negroni.New(
		negroni.HandlerFunc(authentication.Auth),
		negroni.HandlerFunc(controllers.GetApps),
	)).Methods("GET")

	return r
}
