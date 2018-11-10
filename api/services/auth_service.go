package services

import (
	"encoding/json"
	"net/http"
	"strconv"

	auth "github.com/decentorganization/topaz/api/core/authentication"
	"github.com/decentorganization/topaz/api/core/database"
	"github.com/decentorganization/topaz/api/models"
)

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

func AdminRefreshToken(u *models.User) (int, []byte) {
	return okToken(
		auth.InitJWTAuthenticationBackend().GenerateAdminToken(
			strconv.FormatUint(uint64(u.ID), 10)))
}

func AdminLogout(r *http.Request) error {
	token, err := auth.InitJWTAuthenticationBackend().GetToken(r)
	if err != nil {
		return err
	}

	tokenString := r.Header.Get("Authorization")
	return auth.InitJWTAuthenticationBackend().Logout(tokenString, token)
}

func AppLogin(a *models.App) (int, []byte) {
	if err := a.GetApp(database.Manager); err != nil {
		return http.StatusUnauthorized, []byte("")
	}

	return AppRefreshToken(a)
}

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
