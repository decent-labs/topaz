package services

import (
	"net/http"

	"github.com/decentorganization/topaz/models"
)

func NewUser(requestUser *models.User) (int, []byte) {
	// validate input maybe
	// save record in database
	// deploy new smart contract
	// return metadata to user
	return http.StatusOK, []byte("")
}

func EditUser(requestUser *models.User) (int, []byte) {
	return http.StatusOK, []byte("")
}
