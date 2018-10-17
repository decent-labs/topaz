package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func createUserHandler(w http.ResponseWriter, r *http.Request) {
}

// StoreResponse is what gets returned
type StoreResponse struct {
	Hash string
}

func storeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("starting /store handler")

	// we need to buffer the body if we want to read it here and send it
	// in the request.
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// you can reassign the body if you need to parse it as multipart
	r.Body = ioutil.NopCloser(bytes.NewReader(body))

	// create a new url from the raw RequestURI sent by the client
	url := fmt.Sprintf("%s://%s:%s", os.Getenv("STORE_PROXY"), os.Getenv("STORE_HOST"), os.Getenv("STORE_PORT"))

	proxyReq, err := http.NewRequest(r.Method, url, bytes.NewReader(body))

	// We may want to filter some headers, otherwise we could just use a shallow copy
	// proxyReq.Header = req.Header
	proxyReq.Header = make(http.Header)
	for h, val := range r.Header {
		proxyReq.Header[h] = val
	}

	httpClient := http.Client{}

	resp, err := httpClient.Do(proxyReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	sr := new(StoreResponse)
	json.NewDecoder(resp.Body).Decode(sr)

	log.Printf("IPFS hash as told to api: %s", sr.Hash)

	w.Header().Set("Content-Type", "application/vnd.api+json")
	json.NewEncoder(w).Encode(sr)

	log.Println("finished with /store handler")
}

func verifyHandler(w http.ResponseWriter, r *http.Request) {
}

func retrieveHandler(w http.ResponseWriter, r *http.Request) {
}

func main() {
	http.HandleFunc("/create_user", createUserHandler)
	http.HandleFunc("/store", storeHandler)
	http.HandleFunc("/verify", verifyHandler)
	http.HandleFunc("/retrieve", retrieveHandler)

	log.Println("Wake up, api...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
