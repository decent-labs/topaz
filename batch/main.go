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

	a.LastBatched = &ut
	if err := a.UpdateApp(database.Manager); err != nil {
		return b, err
	}

	return b, nil
}

func getObjectsToBatch(app models.App) (models.Objects, error) {
	objects := new(models.Objects)
	err := objects.GetObjectsByAppID(database.Manager, app.ID)
	return *objects, err
}

func makeProof(objs models.Objects, batch models.Batch, root string, tx string) (models.Proof, error) {
	p := models.Proof{
		BatchID:        batch.ID,
		DirectoryHash:  root,
		EthTransaction: tx,
	}

	if err := p.CreateProof(database.Manager); err != nil {
		return p, err
	}

	if err := objs.UpdateProof(database.Manager, &p.ID); err != nil {
		return p, err
	}

	return p, nil
}

func newObjectsFlow(a models.App, os models.Objects) {
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

	_, err = makeProof(os, b, root, tx)
	if err != nil {
		fmt.Errorf("couldn't create proof: " + err.Error())
		return
	}
}

func noObjectsFlow(a models.App) {
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
		objs, err := getObjectsToBatch(a)
		if err != nil {
			fmt.Errorf("couldn't get objects to bacth: " + err.Error())
			return
		}

		if len(objs) > 0 {
			newObjectsFlow(a, objs)
		} else {
			noObjectsFlow(a)
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
