package auth

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/decentorganization/topaz/shared/database"
	"github.com/decentorganization/topaz/shared/models"
	jwt "github.com/dgrijalva/jwt-go"
)

func auth(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc, rID models.AuthKey) {
	token, err := InitJWTAuthenticationBackend().GetToken(req)

	if err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !token.Valid {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	if InitJWTAuthenticationBackend().IsInBlacklist(req.Header.Get("Authorization")) {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	res := claims[string(rID)]
	if res == nil {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	r, err := strconv.ParseUint(res.(string), 10, 64)
	if err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}
	ru := uint(r)

	if err := verifyAuth(rID, ru); err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	req.Header.Del(string(rID))
	req.Header.Add(string(rID), res.(string))
	next(rw, req)
}

// Admin ...
func Admin(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	auth(rw, req, next, models.UserID)
}

// App ...
func App(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	auth(rw, req, next, models.AppID)
}

func verifyAuth(rID models.AuthKey, r uint) error {
	switch rID {
	case models.UserID:
		return verifyUser(r)
	case models.AppID:
		return verifyApp(r)
	}
	return errors.New("unknown resource ID")
}

func verifyUser(r uint) error {
	u := new(models.User)
	u.ID = r
	return u.FindUser(database.Manager)
}

func verifyApp(r uint) error {
	a := new(models.App)
	a.ID = r
	return a.FindApp(database.Manager)
}
