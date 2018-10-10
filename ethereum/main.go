package main

import (
	"encoding/json"
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

		log.Printf("hash to store: '%s'", hash)
		log.Printf("address to store it: '%s'\n", address)

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

// ConnectionHandler defines what's necessary for an ethereum tx to happen
type ConnectionHandler struct {
	Auth       *bind.TransactOpts
	Blockchain *ethclient.Client
}

// Connect service to blockchain
func (conn *ConnectionHandler) Connect(connection string, privateKey string) {
	blockchain, err := ethclient.Dial(connection)
	if err != nil {
		log.Fatalf("Unable to connect to network:%v\n", err)
	}

	pkecdsa, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		log.Fatal(err)
	}

	auth := bind.NewKeyedTransactor(pkecdsa)

	conn.Auth = auth
	conn.Blockchain = blockchain
}

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

func main() {
	c := new(ConnectionHandler)
	c.Connect(os.Getenv("CONN"), os.Getenv("PRIVKEY"))

	http.HandleFunc("/", c.Index)
	http.HandleFunc("/deploy", c.Deploy)
	http.HandleFunc("/store", c.Store)

	log.Println("wake up, ethereum")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
