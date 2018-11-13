package auth

import (
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

func auth(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc, id string) {
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
		if resource = claims[id]; resource == nil {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		req.Header.Del(id)
		req.Header.Add(id, resource.(string))
		next(rw, req)
	} else {
		rw.WriteHeader(http.StatusUnauthorized)
	}
}

// Admin ...
func Admin(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	auth(rw, req, next, "userId")
}

// App ...
func App(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	auth(rw, req, next, "appId")
}
