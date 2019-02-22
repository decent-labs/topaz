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
		return http.StatusInternalServerError, []byte("contact topaz tech support")
	}

	token := authentication.NewResetToken(u.ID, time.Duration(timeout)*time.Hour, pwhash, secret)

	SendPasswordResetEmail(u.Email, token)

	return http.StatusOK, []byte("")
}

// ResetPasswordValidatePassword ...
func ResetPasswordValidatePassword(rp *models.ResetPassword) (int, []byte) {
	secret := []byte(os.Getenv("PASSWORD_RESET_SECRET"))

	userID, err := authentication.VerifyResetToken(rp.Token, getPasswordHash, secret)
	if err != nil {
		return http.StatusBadRequest, []byte("")
	}

	if len(rp.NewPassword) == 0 {
		return http.StatusBadRequest, []byte("set the new password")
	}

	hp, err := authentication.HashPassword(rp.NewPassword)
	if err != nil {
		return http.StatusInternalServerError, []byte("")
	}

	u := models.User{ID: userID}
	if err := u.GetUser(database.Manager); err != nil {
		return http.StatusInternalServerError, []byte("")
	}

	u.Password = hp

	if err := u.UpdateUser(database.Manager); err != nil {
		return http.StatusInternalServerError, []byte("")
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
