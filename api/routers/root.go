package routers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// SetRootRoute ...
func SetRootRoute(r *mux.Router) *mux.Router {
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello Topaz :)")
	})

	return r
}
