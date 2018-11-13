package main

import (
	"log"

	"github.com/decentorganization/topaz/shared/database"
	"github.com/decentorganization/topaz/shared/ipfs"
	"github.com/decentorganization/topaz/shared/models"
)

func getAppsToBatch() models.Apps {
	apps := new(models.Apps)

	if err := apps.GetAppsToBatch(database.Manager); err != nil {
		log.Printf("couldn't get apps to batch: %s", err.Error())
	}

	return *apps
}

func getObjectsToBatch(app models.App) models.Objects {
	objects := new(models.Objects)

	if err := objects.GetObjectsByAppID(database.Manager, app.ID); err != nil {
		log.Printf("couldn't get objects to batch: %s", err.Error())
	}

	return *objects
}

func main() {
	var apps models.Apps
	var objects models.Objects

	apps = getAppsToBatch()

	for _, app := range apps {
		root, err := ipfs.NewObject("unixfs-dir")
		if err != nil {
			log.Printf("couldn't create new ipfs object: %s", err.Error())
		}

		objects = getObjectsToBatch(app)

		for _, object := range objects {
			root, err = ipfs.PatchLink(root, object.Hash, object.Hash, true)
			if err != nil {
				log.Printf("couldn't patchlink ipfs object: %s", err.Error())
			}
		}

		log.Printf("Batch complete. App: %d, Root: %s", app.ID, root)
	}
}
