package authentication

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/decentorganization/topaz/shared/database"
	"github.com/decentorganization/topaz/shared/models"
	jwt "github.com/dgrijalva/jwt-go"
)

// UserAuth ...
func UserAuth(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
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

	res := claims[string(models.AuthUser)]
	if res == nil {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	u := models.User{ID: res.(string)}
	if err := u.GetUser(database.Manager); err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	ctx := context.WithValue(req.Context(), models.AuthUser, &u)
	next(rw, req.WithContext(ctx))
}

// DoubleAuth ...
func DoubleAuth(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	ru := make(map[string]interface{})
	if err := json.Unmarshal(bodyBytes, &ru); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	u := req.Context().Value(models.AuthUser).(*models.User)
	if ok := CheckPasswordHash(ru["password"].(string), u.Password); !ok {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	next(rw, req)
}
