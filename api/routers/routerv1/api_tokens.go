package routerv1

import (
	"github.com/decentorganization/topaz/api/authentication"
	"github.com/decentorganization/topaz/api/controllers"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// SetAPITokensRoutes ...
func SetAPITokensRoutes(r *mux.Router) *mux.Router {
	s := r.PathPrefix("/tokens").Subrouter()

	// Create new api token
	s.Handle("", negroni.New(
		negroni.HandlerFunc(authentication.UserAuth),
		negroni.HandlerFunc(controllers.CreateAPIToken),
	)).Methods("POST")

	// Get all api tokens
	s.Handle("", negroni.New(
		negroni.HandlerFunc(authentication.UserAuth),
		negroni.HandlerFunc(controllers.GetAPITokens),
	)).Methods("GET")

	// Get single api token
	s.Handle("/{id}", negroni.New(
		negroni.HandlerFunc(authentication.UserAuth),
		negroni.HandlerFunc(controllers.GetAPIToken),
	)).Methods("GET")

	return r
}
