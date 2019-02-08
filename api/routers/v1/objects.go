package routers

import (
	"github.com/decentorganization/topaz/api/auth"
	"github.com/decentorganization/topaz/api/controllers"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// SetObjectsRoutes provisions routes for 'object' activity
func SetObjectsRoutes(r *mux.Router) *mux.Router {
	s := r.PathPrefix("/objects").Subrouter()

	// Create new object
	s.Handle("", negroni.New(
		negroni.HandlerFunc(auth.Auth),
		negroni.HandlerFunc(controllers.CreateObject),
	)).Methods("POST")

	// Get all objects
	s.Handle("", negroni.New(
		negroni.HandlerFunc(auth.Auth),
		negroni.HandlerFunc(controllers.GetObjects),
	)).Methods("GET")

	return r
}
