package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/decentorganization/topaz/models"
)

func deploy() (string, error) {
	url := fmt.Sprintf("http://%s:8080/deploy", os.Getenv("ETHEREUM_HOST"))
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	deployResponse := new(models.DeployResponse)
	if err := json.NewDecoder(resp.Body).Decode(deployResponse); err != nil {
		return "", err
	}

	return deployResponse.Addr, nil
}
