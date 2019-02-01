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
		Object:        &o,
	}

	if err := h.CreateHash(database.Manager); err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}

	response, _ := h.MarshalJSON()
	return http.StatusOK, response
}

func Verify(appID uint, hash string) (int, []byte) {
	sa := new(models.App)
	sa.ID = appID

	hs := new(models.Hashes)
	if err := hs.GetVerifiedHashes(database.Manager, sa, hash); err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}

	response, _ := json.Marshal(hs)
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
