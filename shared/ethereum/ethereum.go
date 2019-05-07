package ethereum

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/joho/godotenv"
)

var from common.Address
var client *ethclient.Client
var privateKey *ecdsa.PrivateKey

// GetCurrentNonce ...
func GetCurrentNonce() (uint64, error) {
	nonce, err := client.PendingNonceAt(context.Background(), from)
	if err != nil {
		fmt.Println("couldn't get nonce:", err)
	}
	return nonce, err
}

// GetSuggestedGasPrice ...
func GetSuggestedGasPrice() (*big.Int, error) {
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		fmt.Println("couldn't get the suggested gas price", err)
	}
	return gasPrice, err
}

// Store takes a hash and puts it in a transaction
func Store(hash []byte, nonce uint64, gasPrice *big.Int) (string, error) {
	to := from
	value := big.NewInt(0)

	baseFee, err := strconv.Atoi(os.Getenv("GETH_BASE_GAS"))
	if err != nil {
		fmt.Println("couldn't get the geth base gas fee:", err)
		return "", err
	}

	byteFee, err := strconv.Atoi(os.Getenv("GETH_BYTE_COST"))
	if err != nil {
		fmt.Println("couldn't get the geth byte cost fee:", err)
		return "", err
	}

	gasLimit := uint64(baseFee + (byteFee * len(hash)))

	newTx := types.NewTransaction(nonce, to, value, gasLimit, gasPrice, hash)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		fmt.Println("couldn't get the chainID", err)
		return "", err
	}

	signedTx, err := types.SignTx(newTx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		fmt.Println("couldn't sign the transaction:", err)
		return "", err
	}

	ts := types.Transactions{signedTx}
	rawTx := ts.GetRlp(0)

	var tx *types.Transaction

	rlp.DecodeBytes(rawTx, &tx)

	err = client.SendTransaction(context.Background(), tx)
	if err != nil {
		fmt.Println("couldn't send the transaction:", err)
		return "", err
	}

	return tx.Hash().Hex(), nil
}

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("couldn't load dotenv:", err.Error())
	}

	conn := fmt.Sprintf("%s:%s", os.Getenv("GETH_HOST"), os.Getenv("GETH_PORT"))
	bc, err := ethclient.Dial(conn)
	if err != nil {
		log.Fatalf(err.Error())
	}

	privateKey, err = crypto.HexToECDSA(os.Getenv("GETH_PRIVATE_KEY"))
	if err != nil {
		log.Fatalf(err.Error())
	}

	a := bind.NewKeyedTransactor(privateKey)

	from = a.From
	client = bc
}
