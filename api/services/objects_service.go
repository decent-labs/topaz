package services

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/decentorganization/topaz/shared/database"
	"github.com/decentorganization/topaz/shared/models"
)

func Trust(appID uint, hash *models.Hash) (int, []byte) {
	hb, err := hex.DecodeString(hash.HashHex)
	if err != nil {
		return http.StatusBadRequest, []byte("cannot decode hex hash")
	}

	if len(hb) != 32 {
		return http.StatusBadRequest, []byte("invalid hash length")
	}

	t := time.Now().Unix()

	o := models.Object{
		AppID:         appID,
		UnixTimestamp: t,
	}

	if err := o.CreateObject(database.Manager); err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}

	h := models.Hash{
		Hash:          hb,
		UnixTimestamp: t,
		ObjectID:      &o.ID,
	}

	if err := h.CreateHash(database.Manager); err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}

	response, _ := json.Marshal(h)
	return http.StatusOK, response
}

func Verify(appID uint, hash []byte) (int, []byte) {
	so := models.Object{
		AppID: appID,
	}

	sh := models.Hash{
		Hash: hash,
	}

	hs := new(models.Hashes)
	if err := hs.GetVerifiedHashes(database.Manager, &so, &sh); err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}

	response, _ := json.Marshal(hs)
	return http.StatusOK, response
}

func Report(appID uint, body []byte) (int, []byte) {
	var f interface{}
	if err := json.Unmarshal(body, &f); err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}

	m := f.(map[string]interface{})
	start := int(m["start"].(float64))
	end := int(m["end"].(float64))

	os := new(models.Objects)
	if err := os.GetObjectsByTimestamps(database.Manager, appID, start, end); err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}

	response, _ := json.Marshal(os)
	return http.StatusOK, response
}
