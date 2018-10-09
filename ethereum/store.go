package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/multiformats/go-multihash"
	"topaz.io/topaz-ethereum/contracts"
)

// StoreRequest defines what a valid request body looks like
type StoreRequest struct {
	Address string
	Hash    string
}

// StoreResponse defines what gets returned on store route
type StoreResponse struct {
	Tx string
}

// Store handles the api request
func (api ConnectionHandler) Store(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		decoder := json.NewDecoder(r.Body)

		var data StoreRequest
		err := decoder.Decode(&data)
		if err != nil {
			panic(err)
		}

		address := data.Address
		hash := data.Hash

		m, err := multihash.FromB58String(hash)
		if err != nil {
			panic(err)
		}

		dm, err := multihash.Decode(m)
		if err != nil {
			panic(err)
		}

		var digest [32]byte
		copy(digest[:], dm.Digest)
		var code = uint8(dm.Code)
		var length = uint8(dm.Length)

		transaction := store(api, address, digest, code, length)
		sr := StoreResponse{transaction}
		w.Header().Set("Content-Type", "application/vnd.api+json")
		json.NewEncoder(w).Encode(sr)
	default:
		http.Error(w, "only POST allowed", http.StatusInternalServerError)
	}
}

// Store data using contract
func store(conn ConnectionHandler, address string, digest [32]byte, hashFunction uint8, size uint8) string {
	contract, err := contracts.NewClientCapture(common.HexToAddress(address), conn.Blockchain)
	if err != nil {
		log.Fatalf("Unable to bind to deployed instance of contract:%v\n", err)
	}

	transaction, err := contract.Store(&bind.TransactOpts{
		From:   conn.Auth.From,
		Signer: conn.Auth.Signer,
	}, digest, hashFunction, size)
	if err != nil {
		log.Fatal(err)
	}

	return transaction.Hash().Hex()
}
