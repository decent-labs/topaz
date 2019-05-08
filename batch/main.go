package main

import (
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"syscall"
	"time"

	"github.com/decentorganization/topaz/shared/database"
	"github.com/decentorganization/topaz/shared/ethereum"
	"github.com/decentorganization/topaz/shared/models"
	"github.com/decentorganization/topaz/shared/redis"
	"github.com/joho/godotenv"
)

type appHashesBundle struct {
	App      models.App
	Hashes   models.Hashes
	TimeLeft int
}

type fullCollection []*appHashesBundle

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
	appMap := make(map[string]*appHashesBundle)

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

		if appMap[app.ID] == nil {
			appMap[app.ID] = &appHashesBundle{App: app}
			appMap[app.ID].TimeLeft = ha.TimeLeft
		}

		appMap[app.ID].Hashes = append(appMap[app.ID].Hashes, hash)
	}

	var fullCollection fullCollection
	for _, bundle := range appMap {
		fullCollection = append(fullCollection, bundle)
	}

	sort.Slice(fullCollection, func(i, j int) bool {
		return fullCollection[i].TimeLeft < fullCollection[j].TimeLeft
	})

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

func submitBlockchainTransactions(root []byte, nonce uint64, gasPrice, networkID *big.Int) (string, error) {
	tx, err := ethereum.Store(root, nonce, gasPrice, networkID)
	if err != nil {
		fmt.Println("Had trouble storing hash in Ethereum transation:", err.Error())
	}
	return tx, err
}

func makeProofModel(root []byte, app models.App) models.Proof {
	ut := time.Now().Unix()
	app.LastProofed = &ut

	p := models.Proof{
		App:                 &app,
		MerkleRootMultiHash: root,
		UnixTimestamp:       ut,
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

func makeProof(bundle *appHashesBundle, nonce uint64, gasPrice, networkID *big.Int) {
	root, err := makeMerkleRoot(bundle.Hashes)
	if err != nil {
		return
	}

	tx, err := submitBlockchainTransactions(root, nonce, gasPrice, networkID)
	if err != nil {
		return
	}

	p := makeProofModel(root, bundle.App)

	bt, err := createBlockchainTransaction(&p, tx)
	if err != nil {
		return
	}

	saveProofData(&p, *bt, bundle.Hashes)
}

func makeProofs(fullCollection fullCollection) {
	nonce, err := ethereum.GetCurrentNonce()
	if err != nil {
		return
	}

	gasPrice, err := ethereum.GetSuggestedGasPrice()
	if err != nil {
		return
	}

	networkID, err := ethereum.GetNetworkID()
	if err != nil {
		return
	}

	for _, bundle := range fullCollection {
		makeProof(bundle, nonce, gasPrice, networkID)
		nonce++
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
