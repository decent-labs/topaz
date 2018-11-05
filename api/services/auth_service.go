package services

import (
	"encoding/json"
	"net/http"

	"github.com/decentorganization/topaz/api/api/parameters"
	"github.com/decentorganization/topaz/api/core/authentication"
	"github.com/decentorganization/topaz/models"
	jwt "github.com/dgrijalva/jwt-go"
	request "github.com/dgrijalva/jwt-go/request"
	"github.com/jinzhu/gorm"
)

func Login(requestUser *models.User, db *gorm.DB) (int, []byte) {
	u := new(models.User)
	if err := db.Where("email = ?", requestUser.Email).First(&u).Error; err != nil {
		return http.StatusNotFound, []byte("")
	}

	authBackend := authentication.InitJWTAuthenticationBackend()

	if authBackend.Authenticate(requestUser.Password, u.Password) {
		return makeToken(authBackend, string(u.ID))
	}

	return http.StatusUnauthorized, []byte("")
}

func RefreshToken(requestUser *models.User) (int, []byte) {
	authBackend := authentication.InitJWTAuthenticationBackend()
	return makeToken(authBackend, string(requestUser.ID))
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

func makeToken(authBackend *authentication.JWTAuthenticationBackend, id string) (int, []byte) {
	token, err := authBackend.GenerateToken(id)
	if err != nil {
		return http.StatusInternalServerError, []byte("")
	}

	response, err := json.Marshal(parameters.TokenAuthentication{Token: token})
	if err != nil {
		return http.StatusInternalServerError, []byte("")
	}

	return http.StatusOK, response
}
