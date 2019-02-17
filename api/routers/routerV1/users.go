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

	s.HandleFunc("", controllers.CreateUser).Methods("POST")

	s.Handle("/{id}", negroni.New(
		negroni.HandlerFunc(authentication.Auth),
		negroni.HandlerFunc(controllers.EditUser),
	)).Methods("PUT", "PATCH")

	return r
}
