package routers

import (
	"github.com/gorilla/mux"
)

// InitRoutes ...
func InitRoutes() *mux.Router {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1/").Subrouter()
	s = SetAuthRoutes(s)
	s = SetUsersRoutes(s)
	s = SetAppsRoutes(s)
	s = SetObjectsRoutes(s)
	s = SetHashesRoutes(s)
	s = SetBatchesRoutes(s)
	s = SetProofsRoutes(s)
	return r
}
