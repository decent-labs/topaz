package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/decentorganization/topaz/api/core/database"
	"github.com/decentorganization/topaz/models"
)

func NewApp(newApp *models.App) (int, []byte) {
	if len(newApp.Name) == 0 {
		return http.StatusBadRequest, []byte("bad name")
	}

	a := models.App{Name: newApp.Name, Interval: newApp.Interval, UserID: newApp.UserID}

	url := fmt.Sprintf("http://%s:8080/deploy", os.Getenv("ETHEREUM_HOST"))
	resp, err := http.Get(url)
	if err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}
	defer resp.Body.Close()

	deployResponse := new(models.DeployResponse)
	if err := json.NewDecoder(resp.Body).Decode(deployResponse); err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}

	a.EthAddress = deployResponse.Addr

	if err := database.Manager.Create(&a).Error; err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}

	response, _ := json.Marshal(a)
	return http.StatusOK, response
}
