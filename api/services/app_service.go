package services

import (
	"encoding/json"
	"net/http"

	"github.com/decentorganization/topaz/shared/database"
	"github.com/decentorganization/topaz/shared/ethereum"
	"github.com/decentorganization/topaz/shared/models"
)

// CreateApp ...
func CreateApp(a *models.App, ra *models.App) (int, []byte) {
	if len(ra.Name) == 0 || ra.Interval < 30 {
		return http.StatusBadRequest, []byte("bad name or interval")
	}

	addr, err := ethereum.Deploy()
	if err != nil {
		return http.StatusInternalServerError, []byte("")
	}

	a.Name = ra.Name
	a.Interval = ra.Interval
	a.EthAddress = addr

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
