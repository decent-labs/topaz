package routers

import (
	"github.com/gorilla/mux"
)

// InitRoutes provisions our router with routes for various models
func InitRoutes() *mux.Router {
	r := mux.NewRouter()
	r = r.PathPrefix("/api/v1/").Subrouter()
	r = SetAuthenticationRoutes(r)
	r = SetTestRoutes(r)
	r = SetUsersRoutes(r)
	r = SetAppsRoutes(r)
	r = SetObjectsRoutes(r)
	return r
}
