package main

import (
	"encoding/json"
	"net/http"
)

// IndexResponse defines what get returned on index route
type IndexResponse struct {
	Info string
}

// Index handles the api request
func (api ConnectionHandler) Index(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ir := IndexResponse{"topaz ethereum service"}
		w.Header().Set("Content-Type", "application/vnd.api+json")
		json.NewEncoder(w).Encode(ir)
	default:
		http.Error(w, "only GET allowed", http.StatusInternalServerError)
	}
}
