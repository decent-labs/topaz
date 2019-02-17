package routerV1

import (
	"github.com/decentorganization/topaz/api/authentication"
	"github.com/decentorganization/topaz/api/controllers"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// SetHashesRoutes ...
func SetHashesRoutes(r *mux.Router) *mux.Router {
	s := r.PathPrefix("/apps/{appId}/objects/{objectId}/hashes").Subrouter()

	// Create new hash
	s.Handle("", negroni.New(
		negroni.HandlerFunc(authentication.Auth),
		negroni.HandlerFunc(controllers.CreateHash),
	)).Methods("POST")

	// Get all hashes
	s.Handle("", negroni.New(
		negroni.HandlerFunc(authentication.Auth),
		negroni.HandlerFunc(controllers.GetHashes),
	)).Methods("GET")

	// Get an hash
	s.Handle("/{id}", negroni.New(
		negroni.HandlerFunc(authentication.Auth),
		negroni.HandlerFunc(controllers.GetHash),
	)).Methods("GET")

	return r
}
