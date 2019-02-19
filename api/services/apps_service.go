package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/decentorganization/topaz/api/authorization"
	"github.com/decentorganization/topaz/shared/database"
	"github.com/decentorganization/topaz/shared/models"
)

// CreateApp ...
func CreateApp(u *models.User, ra *models.App) (int, []byte) {
	a, ok := authorization.AuthorizeApps(u)
	if !ok {
		return http.StatusUnauthorized, []byte("")
	}

	minI, err := strconv.Atoi(os.Getenv("MIN_APP_INTERVAL"))
	if err != nil {
		return http.StatusInternalServerError, []byte("contact Topaz tech support")
	}

	if len(ra.Name) == 0 || ra.Interval < minI {
		e := fmt.Sprintf("name must not be 0 characters, and interval must be >= %d seconds", minI)
		return http.StatusBadRequest, []byte(e)
	}

	a.Name = ra.Name
	a.Interval = ra.Interval

	if err := a.CreateApp(database.Manager); err != nil {
		return http.StatusInternalServerError, []byte("")
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
