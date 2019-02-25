package routers

import (
	v1 "github.com/decentorganization/topaz/api/routers/routerV1"
	"github.com/gorilla/mux"
)

// InitRoutes ...
func InitRoutes() *mux.Router {
	r := mux.NewRouter()
	r = SetRootRoute(r)

	sv1 := r.PathPrefix("/v1/").Subrouter()
	sv1 = v1.SetAuthRoutes(sv1)
	sv1 = v1.SetUsersRoutes(sv1)
	sv1 = v1.SetAppsRoutes(sv1)
	sv1 = v1.SetObjectsRoutes(sv1)
	sv1 = v1.SetHashesRoutes(sv1)
	sv1 = v1.SetProofsRoutes(sv1)
	sv1 = v1.SetMarketingRoutes(sv1)

	return r
}
