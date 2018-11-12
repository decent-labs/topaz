package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/decentorganization/topaz/api/core/database"
	"github.com/decentorganization/topaz/api/core/ipfs"
	"github.com/decentorganization/topaz/api/models"
)

func getAppsToFlush() models.Apps {
	apps := new(models.Apps)

	if err := apps.GetAppsToFlush(database.Manager); err != nil {
		log.Printf("couldn't get apps to flush: %s", err.Error())
	}

	return *apps
}

func getObjectsToFlush(app models.App) models.Objects {
	objects := new(models.Objects)

	if err := objects.GetObjectsByAppID(database.Manager, app.ID); err != nil {
		log.Printf("couldn't get objects to flush: %s", err.Error())
	}

	return *objects
}

func main() {
	var apps models.Apps
	var objects models.Objects

	apps = getAppsToFlush()

	for _, app := range apps {
		root, err := ipfs.NewObject("unixfs-dir")
		if err != nil {
			log.Printf("couldn't create new ipfs object: %s", err.Error())
		}

		objects = getObjectsToFlush(app)

		for _, object := range objects {
			root, err = ipfs.PatchLink(root, object.Hash, object.Hash, true)
			if err != nil {
				log.Printf("couldn't patchlink ipfs object: %s", err.Error())
			}
		}

		log.Printf("Flush complete. App: %d, Root: %s", app.ID, root)
	}
}

func init() {
	i, err := strconv.Atoi(os.Getenv("STARTUP_SLEEP"))
	if err != nil {
		log.Fatalf("missing valid STARTUP_SLEEP environment variable: %s", err.Error())
	}
	time.Sleep(time.Duration(i) * time.Second)
}
