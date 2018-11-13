package services

import (
	"encoding/json"
	"net/http"

	"github.com/decentorganization/topaz/shared/database"
	"github.com/decentorganization/topaz/shared/ipfs"
	"github.com/decentorganization/topaz/shared/models"
)

// Trust adds data to ipfs and creates a new 'object' in the database
func Trust(newObject *models.Object) (int, []byte) {
	if len(newObject.DataBlob) == 0 {
		return http.StatusBadRequest, []byte("no data")
	}

	hash, err := ipfs.Add(newObject.DataBlob)
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

func Verify(o *models.Object) (int, []byte) {
	hash, err := ipfs.Hash(o.DataBlob)
	if err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}
	o.Hash = hash

	os := new(models.Objects)
	if err := os.GetObjectsByHash(database.Manager, o); err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}

	response, _ := json.Marshal(os)
	return http.StatusOK, response
}
