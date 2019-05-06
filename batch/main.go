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
	isBatching, err := redis.GetBool(currentlyBatching)

	if err != nil {
		fmt.Println("error getting batch status from redis:", err)
		return false
	}

	if isBatching == true {
		fmt.Println("batch process is already executing")
		return false
	}

	return true
}

func mainLoop() {
	if !safeBatch() {
		return
	}

	if err := redis.SetValue(currentlyBatching, true); err != nil {
		fmt.Println("err telling redis that a batch is starting:", err)
		return
	}

	hwa := new(models.HashesWithApp)
	if err := hwa.GetHashesForProofing(database.Manager); err != nil {
		fmt.Println("Had trouble getting hashes for new proof:", err.Error())

		if err := redis.SetValue(currentlyBatching, false); err != nil {
			fmt.Println("error telling redis that we're done batching")
		}

		return
	}

	if len(*hwa) == 0 {
		fmt.Println("No hashes to proof")

		if err := redis.SetValue(currentlyBatching, false); err != nil {
			fmt.Println("error telling redis that we're done batching")
		}

		return
	}

	fullCollection := make(map[string]*appHashesBundle)

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

	for _, bundle := range fullCollection {
		ms := bundle.Hashes.MakeMerkleLeafs()
		root, err := ms.GetMerkleRoot()
		if err != nil {
			fmt.Println("Had trouble creating merkle root:", err.Error())
			continue
		}

		tx, err := ethereum.Store(root)
		if err != nil {
			fmt.Println("Had trouble storing hash in Ethereum transation:", err.Error())
			continue
		}

		ut := time.Now().Unix()
		bundle.App.LastProofed = &ut

		var rootMultihash multihash.Multihash = root
		rootString := rootMultihash.B58String()

		p := models.Proof{
			App:           &bundle.App,
			MerkleRoot:    rootString,
			UnixTimestamp: ut,
		}

		bcNetwork := new(models.BlockchainNetwork)
		if err := bcNetwork.GetBlockchainNetworkFromName(database.Manager, "ethereum goerli"); err != nil {
			fmt.Println("Had trouble getting blockchain network:", err.Error())
			continue
		}

		bt := models.BlockchainTransaction{
			Proof:               &p,
			BlockchainNetworkID: bcNetwork.ID,
			TransactionHash:     tx,
		}

		dbtx := database.Manager.Begin()

		if err := bt.CreateBlockchainTransaction(dbtx); err != nil {
			fmt.Println("Had trouble creating blockchain transaction record:", err.Error())
			dbtx.Rollback()
			continue
		}

		if err := bundle.Hashes.UpdateWithProof(dbtx, &p.ID); err != nil {
			fmt.Println("Had trouble updating hashes with proof:", err.Error())
			dbtx.Rollback()
			continue
		}

		dbtx.Commit()
	}

	if err := redis.SetValue(currentlyBatching, false); err != nil {
		fmt.Println("error telling redis that we're done batching")
	}

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
