package routers

import (
	"github.com/decentorganization/topaz/api/auth"
	"github.com/decentorganization/topaz/api/controllers"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// SetObjectsRoutes provisions routes for 'object' activity
func SetObjectsRoutes(r *mux.Router) *mux.Router {
	s := r.PathPrefix("/apps/{appId}/objects").Subrouter()

	// Create new object
	s.Handle("", negroni.New(
		negroni.HandlerFunc(auth.Auth),
		negroni.HandlerFunc(controllers.CreateObject),
	)).Methods("POST")

	// Get an object
	s.Handle("/{id}", negroni.New(
		negroni.HandlerFunc(auth.Auth),
		negroni.HandlerFunc(controllers.GetObject),
	)).Methods("GET")

	// Get all objects
	s.Handle("", negroni.New(
		negroni.HandlerFunc(auth.Auth),
		negroni.HandlerFunc(controllers.GetObjects),
	)).Methods("GET")

	return r
}
