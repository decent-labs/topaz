package services

import (
	"encoding/json"
	"net/http"

	"github.com/decentorganization/topaz/api/authorization"
	"github.com/decentorganization/topaz/shared/database"
	"github.com/decentorganization/topaz/shared/models"
)

// GetBatches ...
func GetBatches(u *models.User, aid string) (int, []byte) {
	b, ok := authorization.AuthorizeBatches(u, aid)
	if !ok {
		return http.StatusUnauthorized, []byte("")
	}

	bs := new(models.Batches)
	if err := bs.GetBatches(b, database.Manager); err != nil {
		return http.StatusUnauthorized, []byte("")
	}

	r, err := json.Marshal(&bs)
	if err != nil {
		return http.StatusInternalServerError, []byte("")
	}

	return http.StatusOK, r
}

// GetBatch ...
func GetBatch(u *models.User, aid string, bid string) (int, []byte) {
	b, ok := authorization.AuthorizeBatches(u, aid)
	if !ok {
		return http.StatusUnauthorized, []byte("")
	}

	b.ID = bid
	if err := b.GetBatch(database.Manager); err != nil {
		return http.StatusUnauthorized, []byte("")
	}

	r, err := json.Marshal(&b)
	if err != nil {
		return http.StatusInternalServerError, []byte("")
	}

	return http.StatusOK, r
}
