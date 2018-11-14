package main

import (
	"log"
	"time"

	"github.com/decentorganization/topaz/shared/database"
	"github.com/decentorganization/topaz/shared/ethereum"
	"github.com/decentorganization/topaz/shared/ipfs"
	"github.com/decentorganization/topaz/shared/models"
)

func getAppsToBatch() models.Apps {
	apps := new(models.Apps)
	if err := apps.GetAppsToBatch(database.Manager); err != nil {
		log.Fatalf("couldn't get apps to batch: %s", err.Error())
	}

	return *apps
}

func getObjectsToBatch(app models.App) models.Objects {
	objects := new(models.Objects)
	if err := objects.GetObjectsByAppID(database.Manager, app.ID); err != nil {
		log.Fatalf("couldn't get objects to batch: %s", err.Error())
	}

	return *objects
}

func createTree(objs models.Objects) string {
	root, err := ipfs.NewObject("unixfs-dir")
	if err != nil {
		log.Fatalf("couldn't create new ipfs object: %s", err.Error())
	}

	for _, obj := range objs {
		root, err = ipfs.PatchLink(root, obj.Hash, obj.Hash, true)
		if err != nil {
			log.Fatalf("couldn't patchlink ipfs object: %s", err.Error())
		}
	}

	return root
}

func storeInCaptureContract(ethAddress string, root string) string {
	tx, err := ethereum.Store(ethAddress, root)
	if err != nil {
		log.Fatalf("couldn't store hash in Capture contract: %s", err.Error())
	}

	return tx
}

func batch(a models.App, objs models.Objects, root string, tx string) {
	b := models.Batch{
		DirectoryHash:  root,
		EthTransaction: tx,
		App:            a,
	}

	if err := b.CreateBatch(database.Manager); err != nil {
		log.Fatalf("couln't create batch in database: %s", err.Error())
	}

	for _, obj := range objs {
		obj.BatchID = &b.ID

		if err := obj.UpdateObject(database.Manager); err != nil {
			log.Printf("couldn't update object in database: %s", err.Error())
		}
	}

	ut := time.Now().Unix()
	a.LastBatched = &ut
	if err := a.UpdateApp(database.Manager); err != nil {
		log.Printf("couldn't update app in database: %s", err.Error())
	}
}

func main() {
	for _, a := range getAppsToBatch() {
		objs := getObjectsToBatch(a)
		root := createTree(objs)
		tx := storeInCaptureContract(a.EthAddress, root)

		batch(a, objs, root, tx)
	}
}
