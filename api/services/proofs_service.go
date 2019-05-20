package services

import (
	"encoding/json"
	"net/http"

	"github.com/decentorganization/topaz/api/authorization"
	"github.com/decentorganization/topaz/shared/database"
	"github.com/decentorganization/topaz/shared/models"
)

// GetProofs ...
func GetProofs(u *models.User, aid string) (int, []byte) {
	p, ok := authorization.AuthorizeProofs(u, aid)
	if !ok {
		return http.StatusUnauthorized, []byte("")
	}

	ps := new(models.Proofs)
	if err := ps.GetProofs(p, database.Manager); err != nil {
		return http.StatusUnauthorized, []byte("")
	}

	r, err := json.Marshal(&ps)
	if err != nil {
		return http.StatusInternalServerError, []byte("")
	}

	return http.StatusOK, r
}

// GetProof ...
func GetProof(u *models.User, aid string, pid string) (int, []byte) {
	p, ok := authorization.AuthorizeProofs(u, aid)
	if !ok {
		return http.StatusUnauthorized, []byte("")
	}

	p.ID = pid
	if err := p.GetFullProof(database.Manager); err != nil {
		return http.StatusUnauthorized, []byte("")
	}

	r, err := json.Marshal(&p)
	if err != nil {
		return http.StatusInternalServerError, []byte("")
	}

	return http.StatusOK, r
}
