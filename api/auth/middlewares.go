package auth

import (
	"net/http"
	"strconv"

	"github.com/decentorganization/topaz/shared/database"
	"github.com/decentorganization/topaz/shared/models"
	jwt "github.com/dgrijalva/jwt-go"
)

type key string

const (
	userID key = "userId"
	appID  key = "appId"
)

func auth(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc, rID key) {
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

	if rID == userID {
		uid, err := strconv.ParseUint(r, 10, 64)
		if err != nil {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}
		u := new(models.User)
		u.ID = uint(uid)
		if err := u.FindUser(database.Manager); err != nil {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}
	} else if rID == appID {
		aid, err := strconv.ParseUint(r, 10, 64)
		if err != nil {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}
		a := new(models.App)
		a.ID = uint(aid)
		if err := a.FindApp(database.Manager); err != nil {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	req.Header.Del(string(rID))
	req.Header.Add(string(rID), r)
	next(rw, req)
}

// Admin ...
func Admin(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	auth(rw, req, next, userID)
}

// App ...
func App(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	auth(rw, req, next, appID)
}
