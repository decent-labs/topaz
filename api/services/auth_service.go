package services

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/decentorganization/topaz/api/auth"
	"github.com/decentorganization/topaz/shared/database"
	"github.com/decentorganization/topaz/shared/models"
)

// AdminLogin attempts to authenticate an admin user
func AdminLogin(u *models.User) (int, []byte) {
	suppliedPassword := u.Password

	if err := u.GetUser(database.Manager); err != nil {
		return http.StatusUnauthorized, []byte("")
	}

	if auth := auth.InitJWTAuthenticationBackend().Authenticate(suppliedPassword, u.Password); auth == false {
		return http.StatusUnauthorized, []byte("")
	}

	return AdminRefreshToken(u)
}

// AdminRefreshToken attempts to refresh a JWT for an admin user
func AdminRefreshToken(u *models.User) (int, []byte) {
	return okToken(
		auth.InitJWTAuthenticationBackend().GenerateAdminToken(
			strconv.FormatUint(uint64(u.ID), 10)))
}

// AdminLogout attempts to logout an admin user
func AdminLogout(r *http.Request) error {
	token, err := auth.InitJWTAuthenticationBackend().GetToken(r)
	if err != nil {
		return err
	}

	tokenString := r.Header.Get("Authorization")
	return auth.InitJWTAuthenticationBackend().Logout(tokenString, token)
}

// AppLogin attempts to authenticate an 'app' user
func AppLogin(a *models.App) (int, []byte) {
	if err := a.GetApp(database.Manager); err != nil {
		return http.StatusUnauthorized, []byte("")
	}

	return AppRefreshToken(a)
}

// AppRefreshToken attempts to refresh a JWT for an 'app' user
func AppRefreshToken(a *models.App) (int, []byte) {
	return okToken(
		auth.InitJWTAuthenticationBackend().GenerateAppToken(
			strconv.FormatUint(uint64(a.ID), 10)))
}

func okToken(token string, err error) (int, []byte) {
	if err != nil {
		return http.StatusInternalServerError, []byte("")
	}

	response, err := json.Marshal(models.TokenAuthentication{Token: token})
	if err != nil {
		return http.StatusInternalServerError, []byte("")
	}

	return http.StatusOK, response
}
