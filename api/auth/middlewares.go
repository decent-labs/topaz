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

	if !verifyAuth(rID, r) {
		rw.WriteHeader(http.StatusUnauthorized)
		return
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

func verifyAuth(rID key, r string) bool {
	if rID == userID {
		return verifyUser(r)
	} else if rID == appID {
		return verifyApp(r)
	}
	return false
}

func verifyUser(r string) bool {
	id, ok := parseID(r)
	if !ok {
		return false
	}

	u := new(models.User)
	u.ID = id
	if err := u.FindUser(database.Manager); err != nil {
		return false
	}
	return true
}

func verifyApp(r string) bool {
	id, ok := parseID(r)
	if !ok {
		return false
	}

	a := new(models.App)
	a.ID = id
	if err := a.FindApp(database.Manager); err != nil {
		return false
	}
	return true
}

func parseID(r string) (uint, bool) {
	id, err := strconv.ParseUint(r, 10, 64)
	if err != nil {
		return 0, false
	}
	return uint(id), true
}
