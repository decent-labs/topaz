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

func getAppsToProof() (*models.Apps, error) {
	apps := new(models.Apps)
	err := apps.GetAppsToProof(database.Manager)
	return apps, err
}

func getHashesToProof(app *models.App) (*models.Hashes, error) {
	hashes := new(models.Hashes)
	err := hashes.GetHashesByApp(database.Manager, app)
	return hashes, err
}

func makeProof(hashes *models.Hashes, app *models.App, root string, tx string) (*models.Proof, error) {
	ut := time.Now().Unix()

	app.LastProofed = &ut

	if err := app.UpdateApp(database.Manager); err != nil {
		return nil, err
	}

	p := models.Proof{
		AppID:          app.ID,
		MerkleRoot:     root,
		EthTransaction: tx,
		UnixTimestamp:  ut,
	}

	if err := p.CreateProof(database.Manager); err != nil {
		return nil, err
	}

	if err := hashes.UpdateProof(database.Manager, &p.ID); err != nil {
		return nil, err
	}

	return &p, nil
}

func newHashesFlow(a *models.App, hs *models.Hashes) {
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

	_, err = makeProof(hs, a, root, tx)
	if err != nil {
		fmt.Errorf("couldn't create proof: " + err.Error())
		return
	}
}

func mainLoop() {
	apps, err := getAppsToProof()
	if err != nil {
		fmt.Errorf("didn't get apps to proof: " + err.Error())
		return
	}

	for _, a := range *apps {
		hashes, err := getHashesToProof(&a)
		if err != nil {
			fmt.Errorf("couldn't get objects to proof: " + err.Error())
			return
		}

		if len(*hashes) > 0 {
			newHashesFlow(&a, hashes)
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
