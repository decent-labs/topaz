package main

import (
	"encoding/json"
	"log"
	"net/http"

	"topaz.io/topaz-ethereum/contracts"
)

// DeployResponse defines what gets returned on deploy route
type DeployResponse struct {
	Tx   string
	Addr string
}

// Deploy handles the api request
func (api *ConnectionHandler) Deploy(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		transaction, address := deploy(api)
		dr := DeployResponse{transaction, address}
		w.Header().Set("Content-Type", "application/vnd.api+json")
		json.NewEncoder(w).Encode(dr)
	default:
		http.Error(w, "only POST allowed", http.StatusInternalServerError)
	}
}

// Deploy contract
func deploy(conn *ConnectionHandler) (string, string) {
	address, transaction, _, err := contracts.DeployClientCapture(conn.Auth, conn.Blockchain)
	if err != nil {
		log.Fatal(err)
	}

	return transaction.Hash().Hex(), address.Hex()
}
