package services

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/decentorganization/topaz/shared/database"
	"github.com/decentorganization/topaz/shared/ipfs"
	"github.com/decentorganization/topaz/shared/models"
)

// Trust adds data to ipfs and creates a new 'object' in the database
func Trust(appId uint, dataBlob []byte) (int, []byte) {
	if len(dataBlob) == 0 {
		return http.StatusBadRequest, []byte("no data")
	}

	hash, err := ipfs.Add(dataBlob)
	if err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}

	o := models.Object{
		DataBlob:      dataBlob,
		Hash:          hash,
		AppID:         appId,
		UnixTimestamp: time.Now().Unix(),
	}

	if err := o.CreateObject(database.Manager); err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}

	response, _ := json.Marshal(o)
	return http.StatusOK, response
}

func Verify(appId uint, dataBlob []byte) (int, []byte) {
	hash, err := ipfs.Hash(dataBlob)
	if err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}

	o := models.Object{
		AppID: appId,
		Hash:  hash,
	}

	os := new(models.Objects)
	if err := os.GetObjects(database.Manager, &o); err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}

	response, _ := json.Marshal(os)
	return http.StatusOK, response
}

func Report(appId uint, body []byte) (int, []byte) {
	var f interface{}
	if err := json.Unmarshal(body, &f); err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}

	m := f.(map[string]interface{})
	start := int(m["start"].(float64))
	end := int(m["end"].(float64))

	os := new(models.Objects)
	if err := os.GetObjectsByTimestamps(database.Manager, appId, start, end); err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}

	response, _ := json.Marshal(os)
	return http.StatusOK, response
}
