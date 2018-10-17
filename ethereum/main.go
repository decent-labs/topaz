package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/decentorganization/topaz/ethereum/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	multihash "github.com/multiformats/go-multihash"
)

// ConnectionHandler defines what's necessary for an ethereum tx to happen
type ConnectionHandler struct {
	Auth       *bind.TransactOpts
	Blockchain *ethclient.Client
}

// Connect service to blockchain
func (api *ConnectionHandler) Connect(connection string, privateKey string) error {
	blockchain, err := ethclient.Dial(connection)
	if err != nil {
		return err
	}

	pkecdsa, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return err
	}

	auth := bind.NewKeyedTransactor(pkecdsa)

	api.Auth = auth
	api.Blockchain = blockchain

	return nil
}

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
	log.Println("starting ethereum store service handler")

	decoder := json.NewDecoder(r.Body)

	var data StoreRequest
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("error reading ethereum service request body: %s", err.Error()),
			http.StatusBadRequest,
		)
		return
	}

	address := data.Address
	hash := data.Hash

	log.Printf("dir hash to store: %s", hash)
	log.Printf("address to store it at: %s", address)

	m, err := multihash.FromB58String(hash)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("error converting base58 encoded hash to multihash format: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	dm, err := multihash.Decode(m)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("error decoding multihash to expanded digest: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	var digest [32]byte
	copy(digest[:], dm.Digest)
	var code = uint8(dm.Code)
	var length = uint8(dm.Length)

	contract, err := contracts.NewClientCapture(common.HexToAddress(address), api.Blockchain)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("error instantiating contract from address: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	transaction, err := contract.Store(&bind.TransactOpts{
		From:   api.Auth.From,
		Signer: api.Auth.Signer,
	}, digest, code, length)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("error creating store transaction: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	log.Printf("ethereum transaction: %s", transaction.Hash().Hex())

	sr := StoreResponse{transaction.Hash().Hex()}
	w.Header().Set("Content-Type", "application/vnd.api+json")
	err = json.NewEncoder(w).Encode(sr)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("error encoding ethereum service tx response: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	log.Println("finished with ethereum store service handler")
}

// DeployResponse defines what gets returned on deploy route
type DeployResponse struct {
	Tx   string
	Addr string
}

// Deploy handles the api request
func (api *ConnectionHandler) Deploy(w http.ResponseWriter, r *http.Request) {
	log.Println("starting ethereum deploy service handler")

	address, transaction, _, err := contracts.DeployClientCapture(api.Auth, api.Blockchain)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("error deploying new contract: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	dr := DeployResponse{transaction.Hash().Hex(), address.Hex()}
	w.Header().Set("Content-Type", "application/vnd.api+json")
	err = json.NewEncoder(w).Encode(dr)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("error encoding new deployment results response: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	log.Println("finished with ethereum deploy service handler")
}

func main() {
	c := new(ConnectionHandler)
	err := c.Connect(os.Getenv("CONN"), os.Getenv("PRIVKEY"))
	if err != nil {
		log.Fatalf("could not connect to ethereum blockchain: %s", err.Error())
	}

	http.HandleFunc("/deploy", c.Deploy)
	http.HandleFunc("/store", c.Store)

	log.Println("wake up, ethereum...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
