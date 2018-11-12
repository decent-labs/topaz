package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/urfave/negroni"

	"github.com/decentorganization/topaz/api/routers"
	"github.com/decentorganization/topaz/api/settings"
)

func main() {
	i, err := strconv.Atoi(os.Getenv("STARTUP_SLEEP"))
	if err != nil {
		log.Fatalf("missing valid STARTUP_SLEEP environment variable: %s", err.Error())
	}
	time.Sleep(time.Duration(i) * time.Second)

	settings.Init()
	router := routers.InitRoutes()
	n := negroni.Classic()
	n.UseHandler(router)

	log.Println("wake up, api...")
	log.Fatal(http.ListenAndServe(":8080", n))
}
