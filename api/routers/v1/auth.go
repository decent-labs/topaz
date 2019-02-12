package routers

import (
	"github.com/decentorganization/topaz/api/authentication"
	"github.com/decentorganization/topaz/api/controllers"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// SetAuthRoutes ...
func SetAuthRoutes(r *mux.Router) *mux.Router {
	s := r.PathPrefix("/auth").Subrouter()

	s.HandleFunc("/login", controllers.Login).Methods("POST")

	s.Handle("/refresh-token", negroni.New(
		negroni.HandlerFunc(authentication.Auth),
		negroni.HandlerFunc(controllers.RefreshToken),
	)).Methods("GET")

	s.Handle("/logout", negroni.New(
		negroni.HandlerFunc(authentication.Auth),
		negroni.HandlerFunc(controllers.Logout),
	)).Methods("GET")

	return r
}
