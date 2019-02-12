package authorization

import (
	"github.com/decentorganization/topaz/shared/models"
)

// AuthorizeApps ...
func AuthorizeApps(u *models.User) (*models.App, bool) {
	a := models.App{User: u}

	return &a, true
}
