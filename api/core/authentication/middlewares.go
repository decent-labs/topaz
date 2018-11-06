package authentication

import (
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	request "github.com/dgrijalva/jwt-go/request"
)

func RequireTokenAuthentication(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	authBackend := InitJWTAuthenticationBackend()

	token, err := request.ParseFromRequest(req, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return authBackend.PublicKey, nil
	})

	if err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !token.Valid {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	if authBackend.IsInBlacklist(req.Header.Get("Authorization")) {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		req.Header.Del("userId")
		req.Header.Add("userId", claims["user_id"].(string))
		next(rw, req)
	} else {
		rw.WriteHeader(http.StatusUnauthorized)
	}
}
