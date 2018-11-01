package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB
var httpClient = http.Client{}

type CreateUserRequest struct {
	Name string
}

// TODO: Generate user-scope key for creating apps and managing settings
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("starting /create-user handler")

	var ur CreateUserRequest
	jd := json.NewDecoder(r.Body)
	err := jd.Decode(&ur)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("error executing create new user request: %s", err.Error()),
			http.StatusBadRequest,
		)
		return
	}

	u := User{Name: ur.Name}
	db.Create(&u)

	je := json.NewEncoder(w)
	err = je.Encode(u)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("error encoding new user: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	log.Println("finished with /create-user handler")
}

type CreateAppRequest struct {
	Interval int
	Name     string
	UserID   int
}

// TODO: Associate with user that makes request and generate api key for app
func createAppHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("starting /create-app handler")

	var ar CreateAppRequest
	jd := json.NewDecoder(r.Body)
	err := jd.Decode(&ar)
	if err != nil {
		log.Fatal(err)
	}

	a := App{Interval: ar.Interval, Name: ar.Name}
	db.Create(&a)

	je := json.NewEncoder(w)
	err = je.Encode(a)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("finished with /create-app handler")
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
	url := fmt.Sprintf("http://%s:8080", os.Getenv("STORE_HOST"))

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
	i, err := strconv.Atoi(os.Getenv("STARTUP_SLEEP"))
	if err != nil {
		log.Fatalf("missing valid STARTUP_SLEEP environment variable: %s", err.Error())
	}
	time.Sleep(time.Duration(i) * time.Second)

	dbConn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
		os.Getenv("PQ_HOST"),
		os.Getenv("PQ_PORT"),
		os.Getenv("PQ_USER"),
		os.Getenv("PQ_NAME"),
	)

	db, err = gorm.Open("postgres", dbConn)
	if err != nil {
		log.Fatalf("couldn't even pretend to open database connection: %s", err.Error())
	}
	defer db.Close()

	http.HandleFunc("/create-user", createUserHandler)
	http.HandleFunc("/create-app", createAppHandler)

	http.HandleFunc("/store", storeHandler)
	http.HandleFunc("/verify", verifyHandler)
	http.HandleFunc("/report", reportHandler)

	log.Println("wake up, api...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
