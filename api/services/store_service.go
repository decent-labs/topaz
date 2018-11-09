package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/decentorganization/topaz/models"
)

func store(data []byte) (string, error) {
	url := fmt.Sprintf("http://%s:8080", os.Getenv("STORE_HOST"))
	storeReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return "", err
	}

	httpClient := http.Client{}
	resp, err := httpClient.Do(storeReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	sr := new(models.StoreResponse)
	if err := json.NewDecoder(resp.Body).Decode(sr); err != nil {
		return "", err
	}

	return sr.Hash, nil
}
