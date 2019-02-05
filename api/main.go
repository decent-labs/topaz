package main

import (
	"log"
	"net/http"

	"github.com/decentorganization/topaz/api/routers/v1"
	"github.com/decentorganization/topaz/api/settings"
	"github.com/joho/godotenv"
	"github.com/urfave/negroni"
)

func main() {
	router := routers.InitRoutes()
	n := negroni.Classic()
	n.UseHandler(router)

	log.Println("Wake up, Topaz... :)")
	log.Fatal(http.ListenAndServe(":8080", n))
}

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("couldn't load dotenv: %s", err.Error())
	}

	settings.Init()
}
