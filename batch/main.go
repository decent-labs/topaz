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
		fmt.Println("Had trouble getting apps eligible for new proof:", err.Error())
		return
	}

	for _, a := range *apps {
		hs := new(models.Hashes)
		if err := hs.GetHashesByApp(database.Manager, &a); err != nil {
			fmt.Println("Had trouble getting available hashes to proof:", err.Error())
			continue
		}

		if len(*hs) == 0 {
			continue
		}

		ms := hs.MakeMerkleLeafs()
		root, err := ms.GetMerkleRoot()
		if err != nil {
			fmt.Println("Had trouble creating merkle root:", err.Error())
			continue
		}

		addr := os.Getenv("ETH_CONTRACT_ADDRESS")
		if addr == "" {
			fmt.Println("Ethereum contract address not set")
			continue
		}

		tx, err := ethereum.Store(addr, root)
		if err != nil {
			fmt.Println("Had trouble storing hash in Ethereum transation:", err.Error())
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
			fmt.Println("Had trouble creating proof:", err.Error())
			continue
		}

		if err := hs.UpdateWithProof(database.Manager, &p.ID); err != nil {
			fmt.Println("Had trouble updating hashes with proof:", err.Error())
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
