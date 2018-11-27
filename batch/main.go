package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/decentorganization/topaz/shared/database"
	"github.com/decentorganization/topaz/shared/ethereum"
	"github.com/decentorganization/topaz/shared/ipfs"
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

func createTree(objs models.Objects) (string, error) {
	root, err := ipfs.NewObject("unixfs-dir")
	if err != nil {
		return root, err
	}

	for _, obj := range objs {
		root, err = ipfs.PatchLink(root, obj.Hash, obj.Hash, true)
		if err != nil {
			return root, err
		}
	}

	return root, nil
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
	root, err := createTree(os)
	if err != nil {
		panic("couldn't create hash tree")
	}

	tx, err := ethereum.Store(a.EthAddress, root)
	if err != nil {
		panic("couldn't store in contract")
	}

	b, err := makeBatch(a)
	if err != nil {
		panic("didn't create batch record")
	}

	_, err = makeProof(os, b, root, tx)
	if err != nil {
		panic("couldn't create proof")
	}
}

func noObjectsFlow(a models.App) {
	_, err := makeBatch(a)
	if err != nil {
		panic("didn't create batch record")
	}
}

func mainLoop() {
	apps, err := getAppsToBatch()
	if err != nil {
		panic(fmt.Sprintf("didn't get apps to batch: %s" + err.Error()))
	}

	for _, a := range apps {
		objs, err := getObjectsToBatch(a)
		if err != nil {
			panic("couldn't get objects to bacth")
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
