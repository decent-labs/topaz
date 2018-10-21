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

var httpClient = http.Client{}

type CreateUserRequest struct {
	Name string
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("starting /create-user handler")

	var u CreateUserRequest
	jd := json.NewDecoder(r.Body)
	err := jd.Decode(&u)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(u.Name)
}

func createAppHandler(w http.ResponseWriter, r *http.Request) {
}

type StoreResponse struct {
	Hash string
}

func storeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("starting /store handler")

	// we need to buffer the body if we want to read it here and send it in the request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("error reading /store request body: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	// you can reassign the body if you need to parse it as multipart
	r.Body = ioutil.NopCloser(bytes.NewReader(body))

	// create a new url from the raw RequestURI sent by the client
	url := fmt.Sprintf("%s://%s:%s", os.Getenv("STORE_PROXY"), os.Getenv("STORE_HOST"), os.Getenv("STORE_PORT"))

	proxyReq, err := http.NewRequest(r.Method, url, bytes.NewReader(body))
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("error creating proxy /store request to store service: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	// We may want to filter some headers, otherwise we could just use a shallow copy
	// proxyReq.Header = req.Header
	proxyReq.Header = make(http.Header)
	for h, val := range r.Header {
		proxyReq.Header[h] = val
	}

	resp, err := httpClient.Do(proxyReq)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("error executing store service request: %s", err.Error()),
			http.StatusBadGateway,
		)
		return
	}
	defer resp.Body.Close()

	sr := new(StoreResponse)
	err = json.NewDecoder(resp.Body).Decode(sr)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("error decoding store service json response: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	err = json.NewEncoder(w).Encode(sr)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("error encoding store response from api: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	log.Println("finished with /store handler")
}

func verifyHandler(w http.ResponseWriter, r *http.Request) {
}

func reportHandler(w http.ResponseWriter, r *http.Request) {
}

func main() {
	http.HandleFunc("/create-user", createUserHandler)
	http.HandleFunc("/create-app", createAppHandler)
	http.HandleFunc("/store", storeHandler)
	http.HandleFunc("/verify", verifyHandler)
	http.HandleFunc("/report", reportHandler)

	log.Println("wake up, api...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
