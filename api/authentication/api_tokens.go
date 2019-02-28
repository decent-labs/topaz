package authentication

import (
	"context"
	"net/http"
	"strings"

	"github.com/decentorganization/topaz/shared/database"
	"github.com/decentorganization/topaz/shared/models"
)

// APIAuth ...
func APIAuth(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	// get first "Authorization" header
	tok := req.Header["Authorization"][0]

	// strip "BEARER " off of it
	if len(tok) > 6 && strings.ToUpper(tok[0:7]) == "BEARER " {
		tok = tok[7:]
	}

	a := models.APIToken{Token: tok}
	if err := a.GetAPITokenFromToken(database.Manager); err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	u := models.User{ID: a.UserID}
	if err := u.GetUser(database.Manager); err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	ctx := context.WithValue(req.Context(), models.AuthUser, &u)
	next(rw, req.WithContext(ctx))
}
