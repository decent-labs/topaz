package routers

import (
	"github.com/decentorganization/topaz/api/authentication"
	"github.com/decentorganization/topaz/api/controllers"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// SetUsersRoutes provisions routes for 'user' activity
func SetUsersRoutes(router *mux.Router) *mux.Router {
	router.HandleFunc("/users", controllers.NewUser).Methods("POST")
	router.Handle("/users/{id}",
		negroni.New(
			negroni.HandlerFunc(authentication.Auth),
			negroni.HandlerFunc(controllers.EditUser),
		)).Methods("PUT", "PATCH")
	return router
}
