package services

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/decentorganization/topaz/api/authentication"
	"github.com/decentorganization/topaz/shared/database"
	"github.com/decentorganization/topaz/shared/models"
)

// UpdatePassword ...
func UpdatePassword(u *models.User, rp *models.UpdatePassword) (int, []byte) {
	if len(rp.NewPassword) == 0 {
		errS, _ := json.Marshal(models.Exception{Message: "set a new password"})
		return http.StatusBadRequest, errS
	}

	hp, err := authentication.HashPassword(rp.NewPassword)
	if err != nil {
		errS, _ := json.Marshal(models.Exception{Message: "contact Topaz support"})
		return http.StatusInternalServerError, errS
	}

	u.Password = hp

	if err := u.UpdateUser(database.Manager); err != nil {
		errS, _ := json.Marshal(models.Exception{Message: "contact Topaz support"})
		return http.StatusInternalServerError, errS
	}

	response, _ := json.Marshal(&u)
	return http.StatusOK, response
}

// ResetPasswordGenerateToken ...
func ResetPasswordGenerateToken(u *models.User) (int, []byte) {
	if err := u.GetUserWithEmail(database.Manager); err != nil {
		return http.StatusOK, []byte("")
	}

	pwhash, err := getPasswordHash(u.ID)
	if err != nil {
		return http.StatusOK, []byte("")
	}

	secret := []byte(os.Getenv("PASSWORD_RESET_SECRET"))

	timeout, err := strconv.Atoi(os.Getenv("PASSWORD_RESET_TIMEOUT"))
	if err != nil {
		errS, _ := json.Marshal(models.Exception{Message: "contact Topaz support"})
		return http.StatusInternalServerError, errS
	}

	token := authentication.NewResetToken(u.ID, time.Duration(timeout)*time.Hour, pwhash, secret)

	go SendPasswordResetEmail(u.Email, token)

	return http.StatusOK, []byte("")
}

// ResetPasswordValidatePassword ...
func ResetPasswordValidatePassword(rp *models.ResetPassword) (int, []byte) {
	secret := []byte(os.Getenv("PASSWORD_RESET_SECRET"))

	userID, err := authentication.VerifyResetToken(rp.Token, getPasswordHash, secret)
	if err != nil {
		errS, _ := json.Marshal(models.Exception{Message: "invalid password reset token"})
		return http.StatusBadRequest, errS
	}

	if len(rp.NewPassword) == 0 {
		errS, _ := json.Marshal(models.Exception{Message: "set a new password"})
		return http.StatusBadRequest, errS
	}

	hp, err := authentication.HashPassword(rp.NewPassword)
	if err != nil {
		errS, _ := json.Marshal(models.Exception{Message: "contact Topaz support"})
		return http.StatusInternalServerError, errS
	}

	u := models.User{ID: userID}
	if err := u.GetUser(database.Manager); err != nil {
		errS, _ := json.Marshal(models.Exception{Message: "contact Topaz support"})
		return http.StatusInternalServerError, errS
	}

	u.Password = hp

	if err := u.UpdateUser(database.Manager); err != nil {
		errS, _ := json.Marshal(models.Exception{Message: "contact Topaz support"})
		return http.StatusInternalServerError, errS
	}

	response, _ := json.Marshal(&u)
	return http.StatusOK, response
}

func getPasswordHash(userID string) ([]byte, error) {
	u := models.User{ID: userID}
	if err := u.GetUser(database.Manager); err != nil {
		return nil, err
	}
	return []byte(u.Password), nil
}
