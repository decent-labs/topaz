package main

import (
	"log"
	"net/http"
	"os"

	"github.com/decentorganization/topaz/api/migrations"
	"github.com/decentorganization/topaz/api/routers"
	"github.com/decentorganization/topaz/api/settings"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
)

func main() {
	r := routers.InitRoutes()
	n := negroni.Classic()
	c := cors.AllowAll()

	n.Use(c)
	n.UseHandler(r)

	p := os.Getenv("API_PORT")

	log.Println("topaz listening on", p)
	log.Fatal(http.ListenAndServe(":"+p, n))
}

func init() {
	godotenv.Load(".env")

	settings.GenerateRootContent()
	migrations.Attempt()
}
