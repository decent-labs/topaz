package services

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/decentorganization/topaz/shared/database"
	"github.com/decentorganization/topaz/shared/models"
)

func authObject(a *models.App, uid string) error {
	if err := a.FindApp(database.Manager); err != nil {
		return errors.New("")
	}

	if a.UserID != uid {
		return errors.New("")
	}

	return nil
}

func CreateObject(o *models.Object, uid string) (int, []byte) {
	a := models.App{ID: o.AppID}
	if err := authObject(&a, uid); err != nil {
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

func GetObjects(o *models.Object, uid string) (int, []byte) {
	a := models.App{ID: o.AppID}
	if err := authObject(&a, uid); err != nil {
		return http.StatusUnauthorized, []byte("")
	}

	os := new(models.Objects)
	if err := os.GetObjects(&a, database.Manager); err != nil {
		return http.StatusInternalServerError, []byte("")
	}

	r, err := json.Marshal(&os)
	if err != nil {
		return http.StatusInternalServerError, []byte("")
	}

	return http.StatusOK, r
}
