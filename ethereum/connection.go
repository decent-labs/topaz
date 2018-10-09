package main

import (
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

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
