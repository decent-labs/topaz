package services

import (
	"encoding/json"
	"net/http"

	"github.com/decentorganization/topaz/api/authorization"
	"github.com/decentorganization/topaz/shared/database"
	"github.com/decentorganization/topaz/shared/models"
)

// CreateAPIToken ...
func CreateAPIToken(u *models.User) (int, []byte) {
	a, ok := authorization.AuthorizeAPITokens(u)
	if !ok {
		return http.StatusUnauthorized, []byte("")
	}

	if err := a.CreateAPIToken(database.Manager); err != nil {
		return http.StatusInternalServerError, []byte("")
	}

	r, err := json.Marshal(&a)
	if err != nil {
		return http.StatusInternalServerError, []byte("")
	}

	return http.StatusCreated, r
}

// GetAPITokens ...
func GetAPITokens(u *models.User) (int, []byte) {
	t, ok := authorization.AuthorizeAPITokens(u)
	if !ok {
		return http.StatusUnauthorized, []byte("")
	}

	ts := new(models.APITokens)
	if err := ts.GetAPITokens(t, database.Manager); err != nil {
		return http.StatusUnauthorized, []byte("")
	}

	r, err := json.Marshal(&ts)
	if err != nil {
		return http.StatusInternalServerError, []byte("")
	}

	return http.StatusOK, r
}

// GetAPIToken ...
func GetAPIToken(u *models.User, tid string) (int, []byte) {
	t, ok := authorization.AuthorizeAPITokens(u)
	if !ok {
		return http.StatusUnauthorized, []byte("")
	}

	t.ID = tid
	if err := t.GetAPIToken(database.Manager); err != nil {
		return http.StatusUnauthorized, []byte("")
	}

	r, err := json.Marshal(&t)
	if err != nil {
		return http.StatusInternalServerError, []byte("")
	}

	return http.StatusOK, r
}
