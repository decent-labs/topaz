package services

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/decentorganization/topaz/shared/database"
	"github.com/decentorganization/topaz/shared/ethereum"
	"github.com/decentorganization/topaz/shared/models"
)

func appAuthContext(u models.User, aid string) (*models.App, bool) {
	a := models.App{User: &u}

	if aid != "" {
		a.ID = aid

		if err := a.GetApp(database.Manager); err != nil {
			return nil, false
		}
	}

	return &a, true
}

// CreateApp ...
func CreateApp(ctx context.Context, ra *models.App) (int, []byte) {
	a, ok := appAuthContext(ctx.Value(models.AuthUser).(models.User), "")
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
func GetApp(ctx context.Context, aid string) (int, []byte) {
	a, ok := appAuthContext(ctx.Value(models.AuthUser).(models.User), aid)
	if !ok {
		return http.StatusUnauthorized, []byte("")
	}

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
func GetApps(ctx context.Context) (int, []byte) {
	a, ok := appAuthContext(ctx.Value(models.AuthUser).(models.User), "")
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
