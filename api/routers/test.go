package routers

import (
	"github.com/decentorganization/topaz/api/auth"
	"github.com/decentorganization/topaz/api/controllers"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// SetTestRoutes provisions routes for our tests
func SetTestRoutes(router *mux.Router) *mux.Router {
	router.Handle("/test/hello",
		negroni.New(
			negroni.HandlerFunc(auth.Admin),
			negroni.HandlerFunc(controllers.TestController),
		)).Methods("GET")

	return router
}
