package routers

import (
	"github.com/decentorganization/topaz/api/controllers"
	"github.com/decentorganization/topaz/api/core/authentication"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func SetUsersRoutes(router *mux.Router) *mux.Router {
	router.HandleFunc("/users", controllers.NewUser).Methods("POST")
	router.Handle("/users/{id}",
		negroni.New(
			negroni.HandlerFunc(authentication.RequireTokenAuthentication),
			negroni.HandlerFunc(controllers.EditUser),
		)).Methods("PUT", "PATCH")
	return router
}
