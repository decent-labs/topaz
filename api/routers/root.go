package routers

import (
	"encoding/json"
	"net/http"

	"github.com/decentorganization/topaz/api/settings"
	"github.com/gorilla/mux"
)

// SetRootRoute ...
func SetRootRoute(r *mux.Router) *mux.Router {
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		src, _ := json.Marshal(&settings.Rc)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(src)
	}).Methods("GET")

	return r
}
