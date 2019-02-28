package services

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/decentorganization/topaz/api/authentication"
	"github.com/decentorganization/topaz/shared/database"
	"github.com/decentorganization/topaz/shared/models"
)

// CreateUser ...
func CreateUser(ru *models.User) (int, []byte) {
	if len(ru.Email) == 0 || len(ru.Password) == 0 {
		errS, _ := json.Marshal(models.Exception{Message: "email and password must be set"})
		return http.StatusBadRequest, errS
	}

	hp, err := authentication.HashPassword(ru.Password)
	if err != nil {
		errS, _ := json.Marshal(models.Exception{Message: "contact Topaz support"})
		return http.StatusInternalServerError, errS
	}

	u := models.User{
		Name:     ru.Name,
		Email:    ru.Email,
		Password: hp,
	}

	if err := u.CreateUser(database.Manager); err != nil {
		if err = u.GetUserWithEmail(database.Manager); err == nil {
			errS, _ := json.Marshal(models.Exception{Message: "the email already exists"})
			return http.StatusConflict, errS
		}
		return http.StatusInternalServerError, []byte("")
	}

	sg := models.SendgridEmail{Email: u.Email}
	var mes models.SendgridEmails
	mes = append(mes, sg)

	go CreateNewEmailOnList(&mes, os.Getenv("SENDGRID_API_USERS_LIST"))
	go SendWelcomeEmail(u.Email)

	response, _ := json.Marshal(&u)
	return http.StatusCreated, response
}

// GetUser ...
func GetUser(u *models.User) (int, []byte) {
	res, err := json.Marshal(&u)
	if err != nil {
		return http.StatusInternalServerError, []byte("")
	}

	return http.StatusOK, res
}

// EditUser ...
func EditUser(r *models.User, ru *models.User) (int, []byte) {
	if ru.Name != nil {
		r.Name = ru.Name
	}

	if ru.Email != "" {
		r.Email = ru.Email
	}

	if err := r.UpdateUser(database.Manager); err != nil {
		return http.StatusInternalServerError, []byte("")
	}

	res, err := json.Marshal(&r)
	if err != nil {
		return http.StatusInternalServerError, []byte("")
	}

	return http.StatusOK, res
}
