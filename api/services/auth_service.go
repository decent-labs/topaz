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
	authBackend := authentication.InitJWTAuthenticationBackend()

	var u *models.User
	if err := db.Where("email = ?", requestUser.Email).First(&u).Error; err != nil {
		return http.StatusNotFound, []byte("")
	}

	if authBackend.Authenticate(requestUser.Password, u.Password) {
		token, err := authBackend.GenerateToken(string(u.ID))
		if err != nil {
			return http.StatusInternalServerError, []byte("")
		}
		response, _ := json.Marshal(parameters.TokenAuthentication{Token: token})
		return http.StatusOK, response
	}

	return http.StatusUnauthorized, []byte("")
}

func RefreshToken(requestUser *models.User) []byte {
	authBackend := authentication.InitJWTAuthenticationBackend()
	token, err := authBackend.GenerateToken(requestUser.UUID)
	if err != nil {
		panic(err)
	}
	response, err := json.Marshal(parameters.TokenAuthentication{Token: token})
	if err != nil {
		panic(err)
	}
	return response
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
