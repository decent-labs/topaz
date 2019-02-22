package services

import (
	"encoding/json"
	"net/http"

	"github.com/decentorganization/topaz/api/authentication"
	"github.com/decentorganization/topaz/shared/database"
	"github.com/decentorganization/topaz/shared/models"
)

// UpdatePassword ...
func UpdatePassword(u *models.User, rp *models.UpdatePassword) (int, []byte) {
	if len(rp.NewPassword) == 0 {
		return http.StatusBadRequest, []byte("set the new password")
	}

	hp, err := authentication.HashPassword(rp.NewPassword)
	if err != nil {
		return http.StatusInternalServerError, []byte("")
	}

	u.Password = hp

	if err := u.UpdateUser(database.Manager); err != nil {
		return http.StatusInternalServerError, []byte("")
	}

	response, _ := json.Marshal(&u)
	return http.StatusOK, response
}
