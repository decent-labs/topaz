package services

import (
	"encoding/json"
	"net/http"

	auth "github.com/decentorganization/topaz/api/core/authentication"
	"github.com/decentorganization/topaz/api/core/database"
	"github.com/decentorganization/topaz/models"
)

func NewUser(newUser *models.User) (int, []byte) {
	if len(newUser.Email) == 0 || len(newUser.Password) == 0 || len(newUser.Name) == 0 {
		return http.StatusBadRequest, []byte("bad email, password, or name")
	}

	hp, err := auth.HashPassword(newUser.Password)
	if err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}

	u := models.User{Name: newUser.Name, Email: newUser.Email, Password: hp}
	if err := database.Manager.Create(&u).Error; err != nil {
		return http.StatusInternalServerError, []byte(err.Error())
	}

	response, _ := json.Marshal(u)
	return http.StatusOK, response
}

func EditUser(requestUser *models.User) (int, []byte) {
	return http.StatusOK, []byte("")
}
