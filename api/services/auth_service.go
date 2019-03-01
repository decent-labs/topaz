package services

import (
	"encoding/json"
	"net/http"

	"github.com/decentorganization/topaz/api/authentication"
	"github.com/decentorganization/topaz/shared/database"
	"github.com/decentorganization/topaz/shared/models"
)

// Login attempts to authenticate an admin user
func Login(u *models.User) (int, []byte) {
	suppliedPassword := u.Password

	if err := u.GetUserWithEmail(database.Manager); err != nil {
		errS, _ := json.Marshal(models.Exception{Message: "bad email or password"})
		return http.StatusUnauthorized, errS
	}

	if ok := authentication.InitJWTAuthenticationBackend().Authenticate(suppliedPassword, u.Password); !ok {
		errS, _ := json.Marshal(models.Exception{Message: "bad email or password"})
		return http.StatusUnauthorized, errS
	}

	return RefreshToken(u)
}

// RefreshToken attempts to refresh a JWT for an admin user
func RefreshToken(u *models.User) (int, []byte) {
	token, err := authentication.InitJWTAuthenticationBackend().GenerateToken(u.ID)
	if err != nil {
		return http.StatusInternalServerError, []byte("")
	}

	response, err := json.Marshal(&models.TokenAuth{Token: token})
	if err != nil {
		return http.StatusInternalServerError, []byte("")
	}

	return http.StatusOK, response
}

// Logout attempts to logout an admin user
func Logout(r *http.Request) error {
	token, err := authentication.InitJWTAuthenticationBackend().GetToken(r)
	if err != nil {
		return err
	}

	tokenString := r.Header.Get("Authorization")
	return authentication.InitJWTAuthenticationBackend().Logout(tokenString, token)
}
