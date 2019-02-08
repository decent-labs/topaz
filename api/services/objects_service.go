package services

import (
	"encoding/json"
	"net/http"

	"github.com/decentorganization/topaz/shared/database"
	"github.com/decentorganization/topaz/shared/models"
)

func CreateObject(o *models.Object, uid string) (int, []byte) {
	a := models.App{ID: o.AppID}
	if err := a.FindApp(database.Manager); err != nil {
		return http.StatusUnauthorized, []byte("")
	}

	if a.UserID != uid {
		return http.StatusUnauthorized, []byte("")
	}

	if err := o.CreateObject(database.Manager); err != nil {
		return http.StatusInternalServerError, []byte("")
	}

	r, err := json.Marshal(&o)
	if err != nil {
		return http.StatusInternalServerError, []byte("")
	}

	return http.StatusOK, r
}
