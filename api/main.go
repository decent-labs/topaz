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

	"github.com/urfave/negroni"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	m "github.com/decentorganization/topaz/models"

	"github.com/decentorganization/topaz/api/auth"
	"github.com/decentorganization/topaz/api/routers"
	"github.com/decentorganization/topaz/api/settings"
)

var db *gorm.DB
var httpClient = http.Client{}

// StoreHandler takes data and does the thing
func StoreHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("starting /store handler")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		json.NewEncoder(w).Encode(m.Exception{Message: "error reading /store request body: %s" + err.Error()})
		return
	}

	r.Body = ioutil.NopCloser(bytes.NewReader(body))

	url := fmt.Sprintf("http://%s:8080", os.Getenv("STORE_HOST"))
	proxyReq, err := http.NewRequest(r.Method, url, bytes.NewReader(body))
	if err != nil {
		json.NewEncoder(w).Encode(m.Exception{Message: "error creating proxy /store request to store service: %s" + err.Error()})
		return
	}

	resp, err := httpClient.Do(proxyReq)
	if err != nil {
		json.NewEncoder(w).Encode(m.Exception{Message: "error executing store service request: %s" + err.Error()})
		return
	}
	defer resp.Body.Close()

	sr := new(m.StoreResponse)
	if err := json.NewDecoder(resp.Body).Decode(sr); err != nil {
		json.NewEncoder(w).Encode(m.Exception{Message: "error decoding store service json response: %s" + err.Error()})
		return
	}

	if err := json.NewEncoder(w).Encode(sr); err != nil {
		json.NewEncoder(w).Encode(m.Exception{Message: "error encoding store response from api: %s" + err.Error()})
		return
	}

	log.Println("finished with /store handler")
}

// VerifyHandler takes data and verifies it's authenticity
func VerifyHandler(w http.ResponseWriter, r *http.Request) {
}

// ReportHandler takes a date range and returns metadata from within it
func ReportHandler(w http.ResponseWriter, r *http.Request) {
}

func main() {
	i, err := strconv.Atoi(os.Getenv("STARTUP_SLEEP"))
	if err != nil {
		log.Fatalf("missing valid STARTUP_SLEEP environment variable: %s", err.Error())
	}
	time.Sleep(time.Duration(i) * time.Second)

	dbConn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable", os.Getenv("PQ_HOST"), os.Getenv("PQ_PORT"), os.Getenv("PQ_USER"), os.Getenv("PQ_NAME"))

	db, err := gorm.Open("postgres", dbConn)
	if err != nil {
		log.Fatalf("couldn't even pretend to open database connection: %s", err.Error())
	}
	defer db.Close()

	settings.Init()
	router := routers.InitRoutes()
	n := negroni.Classic()
	n.UseHandler(router)

	http.HandleFunc("/store", auth.AuthApp(StoreHandler))
	http.HandleFunc("/verify", auth.AuthApp(VerifyHandler))
	http.HandleFunc("/report", auth.AuthApp(ReportHandler))

	log.Println("wake up, api...")
	log.Fatal(http.ListenAndServe(":8080", n))
}
