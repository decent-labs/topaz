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

func getAppsToBatch() (models.Apps, error) {
	apps := new(models.Apps)
	err := apps.GetAppsToBatch(database.Manager)
	return *apps, err
}

func makeBatch(a models.App) (models.Batch, error) {
	ut := time.Now().Unix()

	b := models.Batch{
		AppID:         a.ID,
		UnixTimestamp: ut,
	}

	if err := b.CreateBatch(database.Manager); err != nil {
		return b, err
	}

	a.LastProofed = &ut
	if err := a.UpdateApp(database.Manager); err != nil {
		return b, err
	}

	return b, nil
}

func getHashesToBatch(app models.App) (models.Hashes, error) {
	hashes := new(models.Hashes)
	err := hashes.GetHashesByApp(database.Manager, &app)
	return *hashes, err
}

func makeProof(hashes models.Hashes, batch models.Batch, root string, tx string) (models.Proof, error) {
	p := models.Proof{
		BatchID:        batch.ID,
		MerkleRoot:     root,
		EthTransaction: tx,
	}

	if err := p.CreateProof(database.Manager); err != nil {
		return p, err
	}

	if err := hashes.UpdateProof(database.Manager, &p.ID); err != nil {
		return p, err
	}

	return p, nil
}

func newHashesFlow(a models.App, hs models.Hashes) {
	root, err := hs.GetMerkleRoot()
	if err != nil {
		fmt.Errorf("couldn't create hash tree: " + err.Error())
		return
	}

	tx, err := ethereum.Store(a.EthAddress, root)
	if err != nil {
		fmt.Errorf("couldn't store in contract: " + err.Error())
		return
	}

	b, err := makeBatch(a)
	if err != nil {
		fmt.Errorf("didn't create batch record: " + err.Error())
		return
	}

	_, err = makeProof(hs, b, root, tx)
	if err != nil {
		fmt.Errorf("couldn't create proof: " + err.Error())
		return
	}
}

func noHashesFlow(a models.App) {
	_, err := makeBatch(a)
	if err != nil {
		fmt.Errorf("didn't create batch record: " + err.Error())
		return
	}
}

func mainLoop() {
	apps, err := getAppsToBatch()
	if err != nil {
		fmt.Errorf("didn't get apps to batch: " + err.Error())
		return
	}

	for _, a := range apps {
		hashes, err := getHashesToBatch(a)
		if err != nil {
			fmt.Errorf("couldn't get objects to bacth: " + err.Error())
			return
		}

		if len(hashes) > 0 {
			newHashesFlow(a, hashes)
		} else {
			noHashesFlow(a)
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
