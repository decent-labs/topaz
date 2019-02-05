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

	router.Handle("/trust/{uuid}",
		negroni.New(
			negroni.HandlerFunc(auth.App),
			negroni.HandlerFunc(controllers.TrustUpdate),
		)).Methods("POST")

	router.Handle("/verify/{hash}",
		negroni.New(
			negroni.HandlerFunc(auth.App),
			negroni.HandlerFunc(controllers.Verify),
		)).Methods("GET")

	router.Handle("/report/{start}/{end}",
		negroni.New(
			negroni.HandlerFunc(auth.App),
			negroni.HandlerFunc(controllers.Report),
		)).Methods("GET")
	return router
}