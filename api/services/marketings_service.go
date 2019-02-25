package services

import (
	"net/http"

	"github.com/decentorganization/topaz/shared/models"
)

// CreateMarketingEmail ...
func CreateMarketingEmail(me *models.MarketingSiteEmail) int {
	if len(me.Email) == 0 {
		return http.StatusBadRequest
	}

	var mes models.MarketingSiteEmails
	mes = append(mes, *me)

	if ok := CreateNewMarketingEmail(mes); !ok {
		return http.StatusInternalServerError
	}

	return http.StatusOK
}
