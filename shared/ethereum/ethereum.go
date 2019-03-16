package ethereum

import (
	"fmt"
	"log"
	"os"

	"github.com/decentorganization/topaz/shared/ethereum/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	multihash "github.com/multiformats/go-multihash"
)

var auth *bind.TransactOpts
var blockchain *ethclient.Client

// Store takes a Capture Contract address and a hash to store on it
func Store(address, hash string) (string, error) {
	m, err := multihash.FromB58String(hash)
	if err != nil {
		return "", err
	}

	dm, err := multihash.Decode(m)
	if err != nil {
		return "", err
	}

	var digest [32]byte
	copy(digest[:], dm.Digest)
	var code = uint8(dm.Code)
	var length = uint8(dm.Length)

	contract, err := contracts.NewClientCapture(common.HexToAddress(address), blockchain)
	if err != nil {
		return "", err
	}

	transaction, err := contract.Store(&bind.TransactOpts{
		From:   auth.From,
		Signer: auth.Signer,
	}, digest, code, length)
	if err != nil {
		return "", err
	}

	return transaction.Hash().Hex(), nil
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

	pkecdsa, err := crypto.HexToECDSA(os.Getenv("GETH_PRIVATE_KEY"))
	if err != nil {
		log.Fatalf(err.Error())
	}

	a := bind.NewKeyedTransactor(pkecdsa)

	blockchain = bc
	auth = a
}
