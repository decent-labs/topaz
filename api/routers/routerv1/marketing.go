package routerv1

import (
	"github.com/decentorganization/topaz/api/controllers"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// SetMarketingRoutes ...
func SetMarketingRoutes(r *mux.Router) *mux.Router {
	s := r.PathPrefix("/marketing").Subrouter()

	s.Handle("/emails", negroni.New(
		negroni.HandlerFunc(controllers.CreateMarketingEmail),
	)).Methods("POST")

	return r
}
