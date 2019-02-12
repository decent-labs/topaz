package routers

import (
	"github.com/gorilla/mux"
)

// InitRoutes provisions our router with routes for various models
func InitRoutes() *mux.Router {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1/").Subrouter()
	s = SetAuthRoutes(s)
	s = SetUsersRoutes(s)
	s = SetAppsRoutes(s)
	s = SetObjectsRoutes(s)
	s = SetHashesRoutes(s)
	s = SetBatchesRoutes(s)
	return r
}
