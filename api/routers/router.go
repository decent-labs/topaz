package routers

import (
	"github.com/gorilla/mux"
)

// InitRoutes provisions our router with routes for various models
func InitRoutes() *mux.Router {
	router := mux.NewRouter()
	router = SetAuthenticationRoutes(router)
	router = SetTestRoutes(router)
	router = SetUsersRoutes(router)
	router = SetAppsRoutes(router)
	router = SetObjectsRoutes(router)
	return router
}
