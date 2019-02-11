package services

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/decentorganization/topaz/shared/database"
	"github.com/decentorganization/topaz/shared/models"
)

func objectAuthContext(u models.User, aid string, oid string) (*models.Object, bool) {
	a := models.App{
		ID:   aid,
		User: &u,
	}

	if err := a.GetApp(database.Manager); err != nil {
		return nil, false
	}

	o := models.Object{App: &a}

	if oid != "" {
		o.ID = oid

		if err := o.GetObject(database.Manager); err != nil {
			return nil, false
		}
	}

	return &o, true
}

// CreateObject ...
func CreateObject(ctx context.Context, aid string) (int, []byte) {
	o, ok := objectAuthContext(ctx.Value(models.AuthUser).(models.User), aid, "")
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

// GetObject ...
func GetObject(ctx context.Context, aid string, oid string) (int, []byte) {
	o, ok := objectAuthContext(ctx.Value(models.AuthUser).(models.User), aid, oid)
	if !ok {
		return http.StatusUnauthorized, []byte("")
	}

	r, err := json.Marshal(&o)
	if err != nil {
		return http.StatusInternalServerError, []byte("")
	}

	return http.StatusOK, r
}

// GetObjects ...
func GetObjects(ctx context.Context, aid string) (int, []byte) {
	o, ok := objectAuthContext(ctx.Value(models.AuthUser).(models.User), aid, "")
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
