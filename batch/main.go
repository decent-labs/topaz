package main

import (
	"fmt"
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

	for _, obj := range objs {
		obj.ProofID = &p.ID
		if err := obj.UpdateObject(database.Manager); err != nil {
			return p, err
		}
	}

	return p, nil
}

func main() {
	apps, err := getAppsToBatch()
	if err != nil {
		panic(fmt.Sprintf("didn't get apps to batch: %s" + err.Error()))
	}

	for _, a := range apps {
		b, err := makeBatch(a)
		if err != nil {
			panic("didn't create batch record")
		}

		objs, err := getObjectsToBatch(a)
		if err != nil {
			panic("couldn't get objects to bacth")
		}

		if len(objs) == 0 {
			// no objects to batch
			continue
		}

		root, err := createTree(objs)
		if err != nil {
			panic("couldn't create hash tree")
		}

		tx, err := ethereum.Store(a.EthAddress, root)
		if err != nil {
			panic("couldn't store in contract")
		}

		_, err = makeProof(objs, b, root, tx)
		if err != nil {
			panic("couldn't create proof")
		}
	}
}
