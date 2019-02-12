package services

import (
	"encoding/json"
	"net/http"

	"github.com/decentorganization/topaz/api/authorization"
	"github.com/decentorganization/topaz/shared/database"
	"github.com/decentorganization/topaz/shared/ethereum"
	"github.com/decentorganization/topaz/shared/models"
)

// CreateApp ...
func CreateApp(u *models.User, ra *models.App) (int, []byte) {
	a, ok := authorization.AuthorizeApps(u)
	if !ok {
		return http.StatusUnauthorized, []byte("")
	}

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

	r, err := json.Marshal(&a)
	if err != nil {
		return http.StatusInternalServerError, []byte("")
	}

	return http.StatusOK, r
}

// GetApp ...
func GetApp(u *models.User, aid string) (int, []byte) {
	a, ok := authorization.AuthorizeApps(u)
	if !ok {
		return http.StatusUnauthorized, []byte("")
	}

	a.ID = aid
	if err := a.GetApp(database.Manager); err != nil {
		return http.StatusUnauthorized, []byte("")
	}

	r, err := json.Marshal(&a)
	if err != nil {
		return http.StatusInternalServerError, []byte("")
	}

	return http.StatusOK, r
}

// GetApps ...
func GetApps(u *models.User) (int, []byte) {
	a, ok := authorization.AuthorizeApps(u)
	if !ok {
		return http.StatusUnauthorized, []byte("")
	}

	as := new(models.Apps)
	if err := as.GetApps(a, database.Manager); err != nil {
		return http.StatusUnauthorized, []byte("")
	}

	r, err := json.Marshal(&as)
	if err != nil {
		return http.StatusInternalServerError, []byte("")
	}

	return http.StatusOK, r
}
