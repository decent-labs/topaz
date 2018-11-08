package routers

import (
	"github.com/decentorganization/topaz/api/controllers"
	"github.com/decentorganization/topaz/api/core/authentication"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func SetAuthenticationRoutes(router *mux.Router) *mux.Router {
	router.HandleFunc("/token-auth", controllers.AdminLogin).Methods("POST")
	router.Handle("/refresh-token-auth",
		negroni.New(
			negroni.HandlerFunc(authentication.Admin),
			negroni.HandlerFunc(controllers.AdminRefreshToken),
		)).Methods("GET")
	router.Handle("/logout",
		negroni.New(
			negroni.HandlerFunc(authentication.Admin),
			negroni.HandlerFunc(controllers.AdminLogout),
		)).Methods("GET")
	return router
}
