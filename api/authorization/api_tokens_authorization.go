package authorization

import (
	"github.com/decentorganization/topaz/shared/models"
)

// AuthorizeAPITokens ...
func AuthorizeAPITokens(u *models.User) (*models.APIToken, bool) {
	a := models.APIToken{User: u}

	return &a, true
}
