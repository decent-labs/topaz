package routers

import (
	"github.com/decentorganization/topaz/api/authentication"
	"github.com/decentorganization/topaz/api/controllers"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// SetBatchesRoutes provisions routes for 'hashes' activity
func SetBatchesRoutes(r *mux.Router) *mux.Router {
	s := r.PathPrefix("/apps/{appId}/batches").Subrouter()

	// Get all batches
	s.Handle("", negroni.New(
		negroni.HandlerFunc(authentication.Auth),
		negroni.HandlerFunc(controllers.GetBatches),
	)).Methods("GET")

	// Get a batch
	s.Handle("/{id}", negroni.New(
		negroni.HandlerFunc(authentication.Auth),
		negroni.HandlerFunc(controllers.GetBatch),
	)).Methods("GET")

	return r
}
