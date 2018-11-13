package routers

import (
	"github.com/decentorganization/topaz/api/auth"
	"github.com/decentorganization/topaz/api/controllers"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// SetObjectsRoutes provisions routes for 'object' activity
func SetObjectsRoutes(router *mux.Router) *mux.Router {
	router.Handle("/trust",
		negroni.New(
			negroni.HandlerFunc(auth.App),
			negroni.HandlerFunc(controllers.Trust),
		)).Methods("POST")
	return router
}
