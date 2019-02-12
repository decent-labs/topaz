package authorization

import (
	"github.com/decentorganization/topaz/shared/database"
	"github.com/decentorganization/topaz/shared/models"
)

// AuthorizeBatches ...
func AuthorizeBatches(u *models.User, aid string) (*models.Batch, bool) {
	a := models.App{
		ID:   aid,
		User: u,
	}

	if err := a.GetApp(database.Manager); err != nil {
		return nil, false
	}

	b := models.Batch{App: &a}

	return &b, true
}
