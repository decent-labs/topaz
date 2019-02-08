package auth

import (
	"context"
	"net/http"

	"github.com/decentorganization/topaz/shared/database"
	"github.com/decentorganization/topaz/shared/models"
	jwt "github.com/dgrijalva/jwt-go"
)

func Auth(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
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

	res := claims[string(models.UserID)]
	if res == nil {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	u := models.User{ID: res.(string)}
	if err := u.FindUser(database.Manager); err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	ctx := context.WithValue(req.Context(), models.UserID, u.ID)
	next(rw, req.WithContext(ctx))
}
