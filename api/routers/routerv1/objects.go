package routerv1

import (
	"github.com/decentorganization/topaz/api/authentication"
	"github.com/decentorganization/topaz/api/controllers"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// SetObjectsRoutes ...
func SetObjectsRoutes(r *mux.Router) *mux.Router {
	s := r.PathPrefix("/apps/{appId}/objects").Subrouter()

	// Create new object
	s.Handle("", negroni.New(
		negroni.HandlerFunc(authentication.APIAuth),
		negroni.HandlerFunc(controllers.CreateObject),
	)).Methods("POST")

	// Get all objects
	s.Handle("", negroni.New(
		negroni.HandlerFunc(authentication.APIAuth),
		negroni.HandlerFunc(controllers.GetObjects),
	)).Methods("GET")

	// Get an object
	s.Handle("/{id}", negroni.New(
		negroni.HandlerFunc(authentication.APIAuth),
		negroni.HandlerFunc(controllers.GetObject),
	)).Methods("GET")

	return r
}
