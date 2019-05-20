package routerv1

import (
	"github.com/decentorganization/topaz/api/authentication"
	"github.com/decentorganization/topaz/api/controllers"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// SetProofsRoutes ...
func SetProofsRoutes(r *mux.Router) *mux.Router {
	s := r.PathPrefix("/apps/{appId}/proofs").Subrouter()

	// Get all proofs
	s.Handle("", negroni.New(
		negroni.HandlerFunc(authentication.APIAuth),
		negroni.HandlerFunc(controllers.GetProofs),
	)).Methods("GET")

	// Get a proof
	s.Handle("/{id}", negroni.New(
		negroni.HandlerFunc(authentication.APIAuth),
		negroni.HandlerFunc(controllers.GetProof),
	)).Methods("GET")

	return r
}
