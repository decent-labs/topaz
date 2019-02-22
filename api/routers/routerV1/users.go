package routerV1

import (
	"github.com/decentorganization/topaz/api/authentication"
	"github.com/decentorganization/topaz/api/controllers"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// SetUsersRoutes ...
func SetUsersRoutes(r *mux.Router) *mux.Router {
	s := r.PathPrefix("/users").Subrouter()

	s.Handle("", negroni.New(
		negroni.HandlerFunc(controllers.CreateUser),
	)).Methods("POST")

	s.Handle("/me", negroni.New(
		negroni.HandlerFunc(authentication.Auth),
		negroni.HandlerFunc(authentication.DoubleAuth),
		negroni.HandlerFunc(controllers.EditUser),
	)).Methods("PUT", "PATCH")

	return r
}
