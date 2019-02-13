package routers

import (
	"github.com/gorilla/mux"
)

// InitRoutes ...
func InitRoutes() *mux.Router {
	r := mux.NewRouter()

	s := r.PathPrefix("/v1/").Subrouter()

	s = SetAuthRoutes(s)
	s = SetUsersRoutes(s)
	s = SetAppsRoutes(s)
	s = SetObjectsRoutes(s)
	s = SetHashesRoutes(s)
	s = SetProofsRoutes(s)

	return r
}
