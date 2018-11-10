package routers

import (
	"github.com/gorilla/mux"
)

func InitRoutes() *mux.Router {
	router := mux.NewRouter()
	router = SetAuthenticationRoutes(router)
	router = SetHelloRoutes(router)
	router = SetUsersRoutes(router)
	router = SetAppsRoutes(router)
	return router
}
