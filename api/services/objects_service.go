package services

import (
	"encoding/json"
	"net/http"

	"github.com/decentorganization/topaz/api/authorization"
	"github.com/decentorganization/topaz/shared/database"
	"github.com/decentorganization/topaz/shared/models"
)

// CreateObject ...
func CreateObject(u *models.User, aid string) (int, []byte) {
	o, ok := authorization.AuthorizeObjects(u, aid)
	if !ok {
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

// GetObjects ...
func GetObjects(u *models.User, aid string) (int, []byte) {
	o, ok := authorization.AuthorizeObjects(u, aid)
	if !ok {
		return http.StatusUnauthorized, []byte("")
	}

	os := new(models.Objects)
	if err := os.GetObjects(o, database.Manager); err != nil {
		return http.StatusUnauthorized, []byte("")
	}

	r, err := json.Marshal(&os)
	if err != nil {
		return http.StatusInternalServerError, []byte("")
	}

	return http.StatusOK, r
}

// GetObject ...
func GetObject(u *models.User, aid string, oid string) (int, []byte) {
	o, ok := authorization.AuthorizeObjects(u, aid)
	if !ok {
		return http.StatusUnauthorized, []byte("")
	}

	o.ID = oid
	if err := o.GetObject(database.Manager); err != nil {
		return http.StatusUnauthorized, []byte("")
	}

	r, err := json.Marshal(&o)
	if err != nil {
		return http.StatusInternalServerError, []byte("")
	}

	return http.StatusOK, r
}
