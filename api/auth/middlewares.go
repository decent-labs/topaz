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

func auth(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc, id key) {
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

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		var resource interface{}
		if resource = claims[string(id)]; resource == nil {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		if id == userID {
			uid, err := strconv.ParseUint(resource.(string), 10, 64)
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
		} else if id == appID {
			aid, err := strconv.ParseUint(resource.(string), 10, 64)
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

		req.Header.Del(string(id))
		req.Header.Add(string(id), resource.(string))
		next(rw, req)
	} else {
		rw.WriteHeader(http.StatusUnauthorized)
	}
}

// Admin ...
func Admin(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	auth(rw, req, next, userID)
}

// App ...
func App(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	auth(rw, req, next, appID)
}
