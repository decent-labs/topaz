package services

import (
	"encoding/json"
	"net/http"

	"github.com/decentorganization/topaz/shared/database"
	"github.com/decentorganization/topaz/shared/models"
)

func authObject(o *models.Object) bool {
	uid := o.App.User.ID

	if err := o.App.GetApp(database.Manager); err != nil {
		return false
	}

	if o.App.User.ID != uid {
		return false
	}

	return true
}

// CreateObject ...
func CreateObject(o *models.Object) (int, []byte) {
	if ok := authObject(o); !ok {
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

// GetObject ...
func GetObject(o *models.Object) (int, []byte) {
	if ok := authObject(o); !ok {
		return http.StatusUnauthorized, []byte("")
	}

	if err := o.GetObject(database.Manager); err != nil {
		return http.StatusUnauthorized, []byte("")
	}

	r, err := json.Marshal(&o)
	if err != nil {
		return http.StatusInternalServerError, []byte("")
	}

	return http.StatusOK, r
}

// GetObjects ...
func GetObjects(o *models.Object) (int, []byte) {
	if ok := authObject(o); !ok {
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
