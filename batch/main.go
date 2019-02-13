package main

import (
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
		hs := new(models.Hashes)
		if err := hs.GetHashesByApp(database.Manager, &a); err != nil {
			continue
		}

		if len(*hs) == 0 {
			continue
		}

		ms := hs.MakeMerkleLeafs()
		root, err := ms.GetMerkleRoot()
		if err != nil {
			continue
		}

		addr := os.Getenv("ETH_CONTRACT_ADDRESS")
		if addr == "" {
			fmt.Println("Ethereum contract address not set")
			continue
		}

		tx, err := ethereum.Store(addr, root)
		if err != nil {
			continue
		}

		ut := time.Now().Unix()
		a.LastProofed = &ut

		p := models.Proof{
			App:            &a,
			MerkleRoot:     root,
			EthTransaction: tx,
			UnixTimestamp:  ut,
		}

		if err := p.CreateProof(database.Manager); err != nil {
			continue
		}

		if err := hs.UpdateWithProof(database.Manager, &p.ID); err != nil {
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
