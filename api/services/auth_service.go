package services

import (
	"encoding/json"
	"net/http"

	"github.com/decentorganization/topaz/api/auth"
	"github.com/decentorganization/topaz/shared/database"
	"github.com/decentorganization/topaz/shared/models"
)

// Login attempts to authenticate an admin user
func Login(u *models.User) (int, []byte) {
	suppliedPassword := u.Password

	if err := u.GetUser(database.Manager); err != nil {
		return http.StatusUnauthorized, []byte("")
	}

	if auth := auth.InitJWTAuthenticationBackend().Authenticate(suppliedPassword, u.Password); auth == false {
		return http.StatusUnauthorized, []byte("")
	}

	return RefreshToken(u)
}

// RefreshToken attempts to refresh a JWT for an admin user
func RefreshToken(u *models.User) (int, []byte) {
	token, err := auth.InitJWTAuthenticationBackend().GenerateToken(u.ID)
	if err != nil {
		return http.StatusInternalServerError, []byte("")
	}

	response, err := json.Marshal(models.TokenAuth{Token: token})
	if err != nil {
		return http.StatusInternalServerError, []byte("")
	}

	return http.StatusOK, response
}

// Logout attempts to logout an admin user
func Logout(r *http.Request) error {
	token, err := auth.InitJWTAuthenticationBackend().GetToken(r)
	if err != nil {
		return err
	}

	tokenString := r.Header.Get("Authorization")
	return auth.InitJWTAuthenticationBackend().Logout(tokenString, token)
}
