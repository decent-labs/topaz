package services

import (
	"encoding/json"
	"net/http"

	"github.com/decentorganization/topaz/shared/database"
	"github.com/decentorganization/topaz/shared/ethereum"
	"github.com/decentorganization/topaz/shared/models"
)

// NewApp creates a new 'app' and attempts to store it in the database
func NewApp(newApp *models.App) (int, []byte) {
	if len(newApp.Name) == 0 || newApp.Interval < 30 {
		return http.StatusBadRequest, []byte("bad name or interval")
	}

	addr, err := ethereum.Deploy()
	if err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}

	a := models.App{
		UserID:     newApp.UserID,
		Name:       newApp.Name,
		Interval:   newApp.Interval,
		EthAddress: addr,
	}

	if err := a.CreateApp(database.Manager); err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}

	response, _ := json.Marshal(a)
	return http.StatusOK, response
}
