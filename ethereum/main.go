package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	c := new(ConnectionHandler)
	c.Connect(os.Getenv("CONN"), os.Getenv("PRIVKEY"))

	http.HandleFunc("/", c.Index)
	http.HandleFunc("/deploy", c.Deploy)
	http.HandleFunc("/store", c.Store)

	log.Println("topaz-ethereum running")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
