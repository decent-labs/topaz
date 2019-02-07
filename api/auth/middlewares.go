package auth

import (
	"context"
	"errors"
	"net/http"

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

	r := res.(string)
	if err := verifyAuth(rID, r); err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	ctx := context.WithValue(req.Context(), rID, r)
	next(rw, req.WithContext(ctx))
}

// Admin ...
func Admin(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	auth(rw, req, next, models.UserID)
}

// App ...
func App(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	auth(rw, req, next, models.AppID)
}

func verifyAuth(rID models.AuthKey, r string) error {
	switch rID {
	case models.UserID:
		return verifyUser(r)
	case models.AppID:
		return verifyApp(r)
	}
	return errors.New("unknown resource ID")
}

func verifyUser(r string) error {
	u := new(models.User)
	u.ID = r
	return u.FindUser(database.Manager)
}

func verifyApp(r string) error {
	a := new(models.App)
	a.ID = r
	return a.FindApp(database.Manager)
}
