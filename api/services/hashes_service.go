package services

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/decentorganization/topaz/api/authorization"
	"github.com/decentorganization/topaz/shared/database"
	"github.com/decentorganization/topaz/shared/models"
)

// CreateHash ...
func CreateHash(u *models.User, aid string, oid string, rh *models.Hash) (int, []byte) {
	h, ok := authorization.AuthorizeHashes(u, aid, oid)
	if !ok {
		return http.StatusUnauthorized, []byte("")
	}

	hb, err := hex.DecodeString(rh.HashHex)
	if err != nil {
		return http.StatusBadRequest, []byte("cannot decode hex hash")
	}

	if len(hb) != 32 {
		return http.StatusBadRequest, []byte("invalid hash length")
	}

	h.Hash = hb
	h.UnixTimestamp = time.Now().Unix()

	if err := h.CreateHash(database.Manager); err != nil {
		return http.StatusInternalServerError, []byte("")
	}

	r, err := json.Marshal(&h)
	if err != nil {
		return http.StatusInternalServerError, []byte("")
	}

	return http.StatusOK, r
}

// GetHashes ...
func GetHashes(u *models.User, aid string, oid string) (int, []byte) {
	h, ok := authorization.AuthorizeHashes(u, aid, oid)
	if !ok {
		return http.StatusUnauthorized, []byte("")
	}

	hs := new(models.Hashes)
	if err := hs.GetHashes(h, database.Manager); err != nil {
		return http.StatusUnauthorized, []byte("")
	}

	r, err := json.Marshal(&hs)
	if err != nil {
		return http.StatusInternalServerError, []byte("")
	}

	return http.StatusOK, r
}

// GetHash ...
func GetHash(u *models.User, aid string, oid string, hid string) (int, []byte) {
	h, ok := authorization.AuthorizeHashes(u, aid, oid)
	if !ok {
		return http.StatusUnauthorized, []byte("")
	}

	h.ID = hid
	if err := h.GetHash(database.Manager); err != nil {
		return http.StatusUnauthorized, []byte("")
	}

	r, err := json.Marshal(&h)
	if err != nil {
		return http.StatusInternalServerError, []byte("")
	}

	return http.StatusOK, r
}
