package services

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/decentorganization/topaz/api/api/parameters"
	"github.com/decentorganization/topaz/api/core/authentication"
	"github.com/decentorganization/topaz/api/core/database"
	"github.com/decentorganization/topaz/models"
	jwt "github.com/dgrijalva/jwt-go"
	request "github.com/dgrijalva/jwt-go/request"
)

func AdminLogin(requestUser *models.User) (int, []byte) {
	u := new(models.User)
	if err := database.Manager.Where("email = ?", requestUser.Email).First(&u).Error; err != nil {
		return http.StatusUnauthorized, []byte("")
	}

	authBackend := authentication.InitJWTAuthenticationBackend()

	if authBackend.Authenticate(requestUser.Password, u.Password) {
		return makeAdminToken(authBackend, strconv.FormatUint(uint64(u.ID), 10))
	}

	return http.StatusUnauthorized, []byte("")
}

func RefreshToken(requestUser *models.User) (int, []byte) {
	authBackend := authentication.InitJWTAuthenticationBackend()
	return makeAdminToken(authBackend, string(requestUser.ID))
}

func Logout(req *http.Request) error {
	authBackend := authentication.InitJWTAuthenticationBackend()
	tokenRequest, err := request.ParseFromRequest(req, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
		return authBackend.PublicKey, nil
	})
	if err != nil {
		return err
	}
	tokenString := req.Header.Get("Authorization")
	return authBackend.Logout(tokenString, tokenRequest)
}

func makeAdminToken(authBackend *authentication.JWTAuthenticationBackend, id string) (int, []byte) {
	return makeToken(authBackend.GenerateAdminToken(id))
}

func makeToken(token string, err error) (int, []byte) {
	if err != nil {
		return http.StatusInternalServerError, []byte("")
	}

	response, err := json.Marshal(parameters.TokenAuthentication{Token: token})
	if err != nil {
		return http.StatusInternalServerError, []byte("")
	}

	return http.StatusOK, response
}
