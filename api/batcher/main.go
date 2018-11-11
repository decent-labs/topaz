package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/decentorganization/topaz/api/batcher/database"
	"github.com/decentorganization/topaz/api/batcher/models"
	shell "github.com/ipfs/go-ipfs-api"
)

var sh *shell.Shell

func flush(app models.App) {
	objects := new(models.Objects)

	if err := objects.GetObjectsByAppID(database.Manager, app.ID); err != nil {
		log.Printf("couldn't get objects to flush: %s", err.Error())
	}

	for _, object := range *objects {
		log.Println(object)
	}

	// f := new(models.Flush)

	// if err := f.CreateFlush(database.Manager); err != nil {
	// 	log.Printf("couldn't create new flush: %s", err.Error())
	// 	return
	// }

	// defer db.Close()

	// log.Println("BATCH: BEGINNING 'flush()'.")

	// stmt := fmt.Sprintf("insert into flushes (app_id, created_at) values ('%d', now()) returning id;", id)

	// rows, err := db.Query(stmt)
	// if err != nil {
	// 	log.Printf("couldn't execute new flush statement: %s", err.Error())
	// 	return
	// }
	// defer rows.Close()

	// for rows.Next() {
	// 	var id int

	// 	err = rows.Scan(&id)
	// 	if err != nil {
	// 		log.Printf("couldn't scan flush id: %s", err.Error())
	// 		continue
	// 	}

	// 	log.Printf("created flush id: %d", id)
	// }
}

func getAppsToFlush() models.Apps {
	apps := new(models.Apps)

	if err := apps.GetAppsToFlush(database.Manager); err != nil {
		log.Printf("couldn't get apps to flush: %s", err.Error())
	}

	return *apps
}

func main() {
	// apps := new(models.Apps)

	// if err := apps.GetAppsToFlush(database.Manager); err != nil {
	// 	log.Printf("couldn't get apps to flush: %s", err.Error())
	// }

	var apps models.Apps

	apps = getAppsToFlush()

	for _, app := range apps {
		log.Println(app)
		// flush(app)
	}
}

func init() {
	i, err := strconv.Atoi(os.Getenv("STARTUP_SLEEP"))
	if err != nil {
		log.Fatalf("missing valid STARTUP_SLEEP environment variable: %s", err.Error())
	}
	time.Sleep(time.Duration(i) * time.Second)

	conn := fmt.Sprintf("%s:%s", os.Getenv("IPFS_HOST"), os.Getenv("IPFS_PORT"))
	sh = shell.NewShell(conn)
}
