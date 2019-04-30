package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/decentorganization/topaz/shared/database"
	"github.com/decentorganization/topaz/shared/ethereum"
	"github.com/decentorganization/topaz/shared/models"
	multihash "github.com/multiformats/go-multihash"
)

// AppHashesBundle ...
type AppHashesBundle struct {
	App    models.App
	Hashes models.Hashes
}

// FullCollection ...
type FullCollection map[string]*AppHashesBundle

func mainLoop() {
	hwa := new(models.HashesWithApp)
	if err := hwa.GetHashesForProofing(database.Manager); err != nil {
		fmt.Println("Had trouble getting hashes for new proof:", err.Error())
		return
	}

	if len(*hwa) == 0 {
		fmt.Println("No hashes to proof")
		return
	}

	fullCollection := make(map[string]*AppHashesBundle)

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
			fullCollection[app.ID] = &AppHashesBundle{App: app}
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
			App:            &bundle.App,
			MerkleRoot:     rootString,
			EthTransaction: tx,
			UnixTimestamp:  ut,
		}

		if err := p.CreateProof(database.Manager); err != nil {
			fmt.Println("Had trouble creating proof:", err.Error())
			continue
		}

		if err := bundle.Hashes.UpdateWithProof(database.Manager, &p.ID); err != nil {
			fmt.Println("Had trouble updating hashes with proof:", err.Error())
			continue
		}
	}
}

func main() {
	mainLoop()

	i, _ := strconv.Atoi(os.Getenv("BATCH_TICKER"))
	tick := time.Tick(time.Duration(i) * time.Second)
	for {
		select {
		case <-tick:
			mainLoop()
		}
	}
}
