package routers

import (
	"github.com/decentorganization/topaz/api/auth"
	"github.com/decentorganization/topaz/api/controllers"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// SetAppsRoutes provisions routes for 'apps'
func SetAppsRoutes(router *mux.Router) *mux.Router {
	router.Handle("/apps",
		negroni.New(
			negroni.HandlerFunc(auth.Admin),
			negroni.HandlerFunc(controllers.NewApp),
		)).Methods("POST")
	return router
}
