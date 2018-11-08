package routers

import (
	"github.com/decentorganization/topaz/api/controllers"
	"github.com/decentorganization/topaz/api/core/authentication"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func SetAuthenticationRoutes(router *mux.Router) *mux.Router {
	sadmin := router.PathPrefix("/auth/admin").Subrouter()

	sadmin.HandleFunc("/token", controllers.AdminLogin).Methods("POST")
	sadmin.Handle("/refresh-token",
		negroni.New(
			negroni.HandlerFunc(authentication.Admin),
			negroni.HandlerFunc(controllers.AdminRefreshToken),
		)).Methods("GET")
	sadmin.Handle("/logout",
		negroni.New(
			negroni.HandlerFunc(authentication.Admin),
			negroni.HandlerFunc(controllers.AdminLogout),
		)).Methods("GET")

	sapp := router.PathPrefix("/auth/app").Subrouter()

	sapp.Handle("/token",
		negroni.New(
			negroni.HandlerFunc(authentication.Admin),
			negroni.HandlerFunc(controllers.AppLogin),
		)).Methods("POST")

	return router
}
