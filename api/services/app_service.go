package services

import (
	"encoding/json"
	"net/http"

	"github.com/decentorganization/topaz/shared/database"
	"github.com/decentorganization/topaz/shared/ethereum"
	"github.com/decentorganization/topaz/shared/models"
)

func CreateApp(newApp *models.App) (int, []byte) {
	if len(newApp.Name) == 0 || newApp.Interval < 30 {
		return http.StatusBadRequest, []byte("bad name or interval")
	}

	addr, err := ethereum.Deploy()
	if err != nil {
		return http.StatusInternalServerError, []byte("")
	}

	a := models.App{
		UserID:     newApp.UserID,
		Name:       newApp.Name,
		Interval:   newApp.Interval,
		EthAddress: addr,
	}

	if err := a.CreateApp(database.Manager); err != nil {
		return http.StatusInternalServerError, []byte("")
	}

	r, _ := json.Marshal(a)
	return http.StatusOK, r
}

// GetApps ...
func GetApps(a *models.App) (int, []byte) {
	as := new(models.Apps)
	if err := as.GetApps(a, database.Manager); err != nil {
		return http.StatusUnauthorized, []byte("")
	}

	r, _ := json.Marshal(as)
	return http.StatusOK, r
}

// GetApp ...
func GetApp(a *models.App) (int, []byte) {
	if err := a.GetApp(database.Manager); err != nil {
		return http.StatusUnauthorized, []byte("")
	}

	r, _ := json.Marshal(a)
	return http.StatusOK, r
}
