package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/decentorganization/topaz/shared/database"
	"github.com/decentorganization/topaz/shared/ethereum"
	"github.com/decentorganization/topaz/shared/models"
	"github.com/decentorganization/topaz/shared/redis"
	"github.com/joho/godotenv"
	"github.com/multiformats/go-multihash"
)

type appHashesBundle struct {
	App    models.App
	Hashes models.Hashes
}

type fullCollection map[string]*appHashesBundle

var currentlyBatching = "currently_batching"
var afterBatchSleep = 1000

func safeBatch() bool {
	isBatching, _ := redis.GetBool(currentlyBatching)

	if isBatching == true {
		fmt.Println("batch process is already executing")
	}

	return !isBatching
}

func updateBatchingState(newState bool) error {
	err := redis.SetValue(currentlyBatching, newState)
	if err != nil {
		fmt.Println("error changing redis batching state to", newState, ":", err)
	}
	return err
}

func getAllHashes() (*models.HashesWithApp, error) {
	hwa := new(models.HashesWithApp)
	err := hwa.GetHashesForProofing(database.Manager)
	if err != nil {
		fmt.Println("Had trouble getting hashes for new proof:", err.Error())
	}
	return hwa, err
}

func makeCollection(hwa *models.HashesWithApp) fullCollection {
	fullCollection := make(fullCollection)

	for _, ha := range *hwa {
		hash := models.Hash{
			ID:            ha.HashID,
			CreatedAt:     ha.HashCreatedAt,
			UpdatedAt:     ha.HashUpdatedAt,
			DeletedAt:     ha.HashDeletedAt,
			MultiHash:     ha.HashMultiHash,
			UnixTimestamp: ha.HashUnixTimestamp,
			ObjectID:      ha.HashObjectID,
			ProofID:       ha.HashProofID,
		}

		app := models.App{
			ID:          ha.AppID,
			CreatedAt:   ha.AppCreatedAt,
			UpdatedAt:   ha.AppUpdatedAt,
			DeletedAt:   ha.AppDeletedAt,
			Interval:    ha.AppInterval,
			Name:        ha.AppName,
			LastProofed: ha.AppLastProofed,
			UserID:      ha.AppUserID,
		}

		if fullCollection[app.ID] == nil {
			fullCollection[app.ID] = &appHashesBundle{App: app}
		}

		fullCollection[app.ID].Hashes = append(fullCollection[app.ID].Hashes, hash)
	}

	return fullCollection
}

func makeMerkleRoot(hashes models.Hashes) ([]byte, error) {
	ms := hashes.MakeMerkleLeafs()
	root, err := ms.GetMerkleRoot()
	if err != nil {
		fmt.Println("Had trouble creating merkle root:", err.Error())
	}
	return root, err
}

func submitBlockchainTransactions(root []byte) (string, error) {
	tx, err := ethereum.Store(root)
	if err != nil {
		fmt.Println("Had trouble storing hash in Ethereum transation:", err.Error())
	}
	return tx, err
}

func makeProofModel(root []byte, app models.App) models.Proof {
	ut := time.Now().Unix()
	app.LastProofed = &ut

	var rootMultihash multihash.Multihash = root
	rootString := rootMultihash.B58String()

	p := models.Proof{
		App:           &app,
		MerkleRoot:    rootString,
		UnixTimestamp: ut,
	}

	return p
}

func createBlockchainTransaction(p *models.Proof, tx string) (*models.BlockchainTransaction, error) {
	bcNetwork := new(models.BlockchainNetwork)
	if err := bcNetwork.GetBlockchainNetworkFromName(database.Manager, "ethereum goerli"); err != nil {
		fmt.Println("Had trouble getting blockchain network:", err.Error())
		return nil, err
	}

	bt := models.BlockchainTransaction{
		Proof:               p,
		BlockchainNetworkID: bcNetwork.ID,
		TransactionHash:     tx,
	}

	return &bt, nil
}

func saveProofData(p *models.Proof, bt models.BlockchainTransaction, hashes models.Hashes) error {
	dbtx := database.Manager.Begin()

	if err := bt.CreateBlockchainTransaction(dbtx); err != nil {
		fmt.Println("Had trouble creating blockchain transaction record:", err.Error())
		dbtx.Rollback()
		return err
	}

	if err := hashes.UpdateWithProof(dbtx, &p.ID); err != nil {
		fmt.Println("Had trouble updating hashes with proof:", err.Error())
		dbtx.Rollback()
		return err
	}

	dbtx.Commit()

	return nil
}

func makeProofs(fullCollection fullCollection) {
	for _, bundle := range fullCollection {
		root, err := makeMerkleRoot(bundle.Hashes)
		if err != nil {
			continue
		}

		tx, err := submitBlockchainTransactions(root)
		if err != nil {
			continue
		}

		// make sure the app LastProofed is actually getting updated
		p := makeProofModel(root, bundle.App)

		bt, err := createBlockchainTransaction(&p, tx)
		if err != nil {
			continue
		}

		if err := saveProofData(&p, *bt, bundle.Hashes); err != nil {
			continue
		}
	}
}

func mainLoop() {
	if !safeBatch() {
		return
	}

	if err := updateBatchingState(true); err != nil {
		return
	}

	hwa, err := getAllHashes()
	if err != nil {
		return
	}

	if len(*hwa) > 0 {
		fullCollection := makeCollection(hwa)
		makeProofs(fullCollection)
	}

	updateBatchingState(false)
	time.Sleep(time.Duration(afterBatchSleep) * time.Millisecond)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("couldn't load dotenv:", err.Error())
	}

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		sig := <-gracefulStop
		fmt.Printf("caught sig: %+v\n", sig)

		tick := time.Tick(time.Duration((afterBatchSleep / 2)) * time.Millisecond)
		for {
			select {
			case <-tick:
				if safeBatch() {
					os.Exit(0)
				}
			}
		}
	}()

	i, _ := strconv.Atoi(os.Getenv("BATCH_TICKER"))
	tick := time.Tick(time.Duration(i) * time.Second)
	for {
		select {
		case <-tick:
			mainLoop()
		}
	}
}
