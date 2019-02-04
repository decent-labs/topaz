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

	o := models.Object{
		AppID: appID,
	}

	if err := o.MakeUUID(); err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}

	h := models.Hash{
		Hash:          hb,
		UnixTimestamp: time.Now().Unix(),
		Object:        o,
	}

	if err := h.CreateHash(database.Manager); err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}

	o.Hashes = append(o.Hashes, h)

	response, _ := json.Marshal(&o)
	return http.StatusOK, response
}

func TrustUpdate(appID uint, uuid string, hash *models.Hash) (int, []byte) {
	hb, err := hex.DecodeString(hash.HashHex)
	if err != nil {
		return http.StatusBadRequest, []byte("cannot decode hex hash")
	}

	if len(hb) != 32 {
		return http.StatusBadRequest, []byte("invalid hash length")
	}

	o := models.Object{
		AppID: appID,
		UUID:  uuid,
	}

	if err := o.FindObject(database.Manager); err != nil {
		return http.StatusBadRequest, []byte("invalid object")
	}

	h := models.Hash{
		Hash:          hb,
		UnixTimestamp: time.Now().Unix(),
		Object:        o,
	}

	if err := h.CreateHash(database.Manager); err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}

	if err := o.FindFullObject(database.Manager); err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}

	response, _ := json.Marshal(&o)
	return http.StatusOK, response
}

func Verify(appID uint, uuid string) (int, []byte) {
	o := models.Object{
		AppID: appID,
		UUID:  uuid,
	}

	if err := o.FindObject(database.Manager); err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}

	if err := o.FindFullObject(database.Manager); err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}

	for _, h := range o.Hashes {
		if h.Proof == nil {
			continue
		}

		if err := h.Proof.CheckValidity(); err != nil {
			return http.StatusInternalServerError, []byte(err.Error())
		}
	}

	response, _ := json.Marshal(&o)
	return http.StatusOK, response
}

func Report(appID uint, start int, end int) (int, []byte) {
	sa := new(models.App)
	sa.ID = appID

	hs := new(models.Hashes)
	if err := hs.GetHashesByTimestamps(database.Manager, sa, start, end); err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}

	response, _ := json.Marshal(hs)
	return http.StatusOK, response
}
