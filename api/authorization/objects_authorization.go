package authorization

import (
	"github.com/decentorganization/topaz/shared/database"
	"github.com/decentorganization/topaz/shared/models"
)

// AuthorizeObjects ...
func AuthorizeObjects(u *models.User, aid string) (*models.Object, bool) {
	a := models.App{
		ID:   aid,
		User: u,
	}

	if err := a.GetApp(database.Manager); err != nil {
		return nil, false
	}

	o := models.Object{App: &a}

	return &o, true
}
