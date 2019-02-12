package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/decentorganization/topaz/shared/database"
	"github.com/decentorganization/topaz/shared/ethereum"
	"github.com/decentorganization/topaz/shared/models"
)

func mainLoop() {
	apps := new(models.Apps)
	if err := apps.GetAppsToProof(database.Manager); err != nil {
		return
	}

	for _, a := range *apps {
		hashes := new(models.Hashes)
		if err := hashes.GetHashesByApp(database.Manager, app); err != nil {
			continue
		}

		if len(*hashes) = 0 {
			continue
		}
	
		root, err := hs.GetMerkleRoot()
		if err != nil {
			continue
		}

		tx, err := ethereum.Store(a.EthAddress, root)
		if err != nil {
			continue
		}

		ut := time.Now().Unix()
		a.LastProofed = &ut
		
		p := models.Proof{
			App:            a,
			MerkleRoot:     root,
			EthTransaction: tx,
			UnixTimestamp:  ut,
		}

		if err := p.CreateProof(database.Manager); err != nil {
			continue
		}

		if err := hashes.UpdateProof(database.Manager, &p.ID); err != nil {
			continue
		}
	}
}

func main() {
	i, _ := strconv.Atoi(os.Getenv("BATCH_TICKER"))
	tick := time.Tick(time.Duration(i) * time.Second)
	for {
		select {
		case <-tick:
			mainLoop()
		}
	}
}
