package contracts

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func deployContract(t *testing.T) (*backends.SimulatedBackend, *bind.TransactOpts, common.Address, *types.Transaction, *ClientCapture) {
	key, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(key)
	alloc := make(core.GenesisAlloc)
	alloc[auth.From] = core.GenesisAccount{Balance: big.NewInt(1000000000)}
	blockchain := backends.NewSimulatedBackend(alloc, 4500000)
	address, transaction, contract, err := DeployClientCapture(auth, blockchain)
	blockchain.Commit()
	if err != nil {
		t.Fatalf("Failed to deploy the ClientCapture contract: %v", err)
	}
	return blockchain, auth, address, transaction, contract
}

// Test client capture contract gets deployed correctly
func TestDeployClientCapture(t *testing.T) {
	_, _, address, _, _ := deployContract(t)
	if len(address.Bytes()) == 0 {
		t.Error("Expected a valid deployment address. Received empty address byte array instead")
	}
}

// Test HashCaptured event is emitted correctly
func TestHashCapturedEvent(t *testing.T) {
	blockchain, auth, _, _, contract := deployContract(t)

	ch := make(chan *ClientCaptureHashCaptured)
	contract.WatchHashCaptured(&bind.WatchOpts{}, ch)

	var digest = "hello world"
	var hashFunction, size uint8 = 1, 2

	var digestArr [32]byte
	copy(digestArr[:], digest)

	contract.Store(&bind.TransactOpts{
		From:   auth.From,
		Signer: auth.Signer,
	}, digestArr, hashFunction, size)
	blockchain.Commit()

	var hashCaptureEvent *ClientCaptureHashCaptured = <-ch

	if hashCaptureEvent.Digest != digestArr {
		t.Errorf("Expected digest to be: %s. Got: %s", digest, string(hashCaptureEvent.Digest[:32]))
	}

	if hashCaptureEvent.HashFunction != 1 {
		t.Errorf("Expected hash function to be: %d. Got: %d", hashFunction, hashCaptureEvent.HashFunction)
	}

	if hashCaptureEvent.Size != 2 {
		t.Errorf("Expected hash function to be: %d. Got: %d", size, hashCaptureEvent.Size)
	}
}
