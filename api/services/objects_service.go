package services

import (
	"encoding/json"
	"net/http"

	"github.com/decentorganization/topaz/api/core/database"
	"github.com/decentorganization/topaz/models"
)

func Trust(newObject *models.Object) (int, []byte) {
	if len(newObject.DataBlob) == 0 {
		return http.StatusBadRequest, []byte("no data")
	}

	hash, err := store(newObject.DataBlob)
	if err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}

	o := models.Object{
		DataBlob: newObject.DataBlob,
		Hash:     hash,
		AppID:    newObject.AppID,
	}

	if err := o.CreateObject(database.Manager); err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}

	response, _ := json.Marshal(o)
	return http.StatusOK, response
}
