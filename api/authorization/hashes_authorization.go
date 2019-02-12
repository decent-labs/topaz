package authorization

import (
	"github.com/decentorganization/topaz/shared/database"
	"github.com/decentorganization/topaz/shared/models"
)

// AuthorizeHashes ...
func AuthorizeHashes(u *models.User, aid string, oid string) (*models.Hash, bool) {
	a := models.App{
		ID:   aid,
		User: u,
	}

	if err := a.GetApp(database.Manager); err != nil {
		return nil, false
	}

	o := models.Object{
		ID:  oid,
		App: &a,
	}

	if err := o.GetObject(database.Manager); err != nil {
		return nil, false
	}

	h := models.Hash{Object: &o}

	return &h, true
}
