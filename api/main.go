package main

import (
	"log"
	"net/http"
	"os"

	"github.com/decentorganization/topaz/api/routers"
	"github.com/decentorganization/topaz/api/settings"
	"github.com/joho/godotenv"
	"github.com/urfave/negroni"
)

func main() {
	r := routers.InitRoutes()
	n := negroni.Classic()
	n.UseHandler(r)

	p := os.Getenv("API_PORT")

	log.Println("topaz listening on", p)
	log.Fatal(http.ListenAndServe(":"+p, n))
}

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("couldn't load dotenv: %s", err.Error())
	}

	settings.Init()
}
