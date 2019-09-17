package routerv1

import (
	"github.com/decentorganization/topaz/api/authentication"
	"github.com/decentorganization/topaz/api/controllers"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// SetAppsRoutes ...
func SetAppsRoutes(r *mux.Router, c controllers.AppController) *mux.Router {
	s := r.PathPrefix("/apps").Subrouter()

	// Create new app
	s.Handle("", negroni.New(
		negroni.HandlerFunc(authentication.APIAuth),
		negroni.HandlerFunc(c.CreateApp),
	)).Methods("POST")

	// Get all apps
	s.Handle("", negroni.New(
		negroni.HandlerFunc(authentication.APIAuth),
		negroni.HandlerFunc(c.GetApps),
	)).Methods("GET")

	// Get single app
	s.Handle("/{id}", negroni.New(
		negroni.HandlerFunc(authentication.APIAuth),
		negroni.HandlerFunc(c.GetApp),
	)).Methods("GET")

	return r
}
