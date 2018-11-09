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

	m "github.com/decentorganization/topaz/models"
	shell "github.com/ipfs/go-ipfs-api"
)

var sh *shell.Shell

// Take the request body and store it in IPFS, then store the resulting hash in the 'objects' table.
func requestHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("starting store service handler")

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("error reading store service request body: %s", err.Error()),
			http.StatusBadRequest,
		)
		return
	}

	br := bytes.NewReader(b)

	hash, err := sh.Add(br)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("error adding file to ipfs through shell: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	log.Printf("processed incoming data with hash: %s", hash)

	sr := m.StoreResponse{Hash: hash}
	w.Header().Set("Content-Type", "application/vnd.api+json")
	err = json.NewEncoder(w).Encode(sr)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("error encoding store response from store service: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	log.Println("finished with store service handler")
}

func main() {
	i, err := strconv.Atoi(os.Getenv("STARTUP_SLEEP"))
	if err != nil {
		log.Fatalf("missing valid STARTUP_SLEEP environment variable: %s", err.Error())
	}
	time.Sleep(time.Duration(i) * time.Second)

	shConn := fmt.Sprintf("%s:%s", os.Getenv("IPFS_HOST"), os.Getenv("IPFS_PORT"))
	sh = shell.NewShell(shConn)

	http.HandleFunc("/", requestHandler)

	log.Println("wake up, store...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
