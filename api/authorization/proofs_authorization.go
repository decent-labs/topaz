package authorization

import (
	"github.com/decentorganization/topaz/shared/database"
	"github.com/decentorganization/topaz/shared/models"
)

// AuthorizeProofs ...
func AuthorizeProofs(u *models.User, aid string, bid string) (*models.Proof, bool) {
	a := models.App{
		ID:   aid,
		User: u,
	}

	if err := a.GetApp(database.Manager); err != nil {
		return nil, false
	}

	b := models.Batch{
		ID:  bid,
		App: &a,
	}

	if err := b.GetBatch(database.Manager); err != nil {
		return nil, false
	}

	p := models.Proof{Batch: &b}

	return &p, true
}
