package routers

import (
	"github.com/decentorganization/topaz/api/controllers"
	"github.com/decentorganization/topaz/api/core/authentication"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func SetObjectsRoutes(router *mux.Router) *mux.Router {
	router.Handle("/trust",
		negroni.New(
			negroni.HandlerFunc(authentication.App),
			negroni.HandlerFunc(controllers.Trust),
		)).Methods("POST")
	return router
}
