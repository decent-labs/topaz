package routers

import (
	"github.com/decentorganization/topaz/api/auth"
	"github.com/decentorganization/topaz/api/controllers"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// SetAuthenticationRoutes ...
func SetAuthenticationRoutes(router *mux.Router) *mux.Router {
	sadmin := router.PathPrefix("/auth/admin").Subrouter()

	sadmin.HandleFunc("/token", controllers.AdminLogin).Methods("POST")
	sadmin.Handle("/refresh-token",
		negroni.New(
			negroni.HandlerFunc(auth.Admin),
			negroni.HandlerFunc(controllers.AdminRefreshToken),
		)).Methods("GET")
	sadmin.Handle("/logout",
		negroni.New(
			negroni.HandlerFunc(auth.Admin),
			negroni.HandlerFunc(controllers.AdminLogout),
		)).Methods("GET")



	return router
}
