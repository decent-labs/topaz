package routerv1

import (
	"github.com/decentorganization/topaz/api/authentication"
	"github.com/decentorganization/topaz/api/controllers"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// SetAppsRoutes ...
func SetAppsRoutes(r *mux.Router) *mux.Router {
	s := r.PathPrefix("/apps").Subrouter()

	// Create new app
	s.Handle("", negroni.New(
		negroni.HandlerFunc(authentication.APIAuth),
		negroni.HandlerFunc(controllers.CreateApp),
	)).Methods("POST")

	// Get all apps
	s.Handle("", negroni.New(
		negroni.HandlerFunc(authentication.APIAuth),
		negroni.HandlerFunc(controllers.GetApps),
	)).Methods("GET")

	// Get single app
	s.Handle("/{id}", negroni.New(
		negroni.HandlerFunc(authentication.APIAuth),
		negroni.HandlerFunc(controllers.GetApp),
	)).Methods("GET")

	return r
}
