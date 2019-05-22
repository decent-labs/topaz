package services

import (
	"net/http"
	"os"

	"github.com/decentorganization/topaz/shared/models"
)

// CreateMarketingEmail ...
func CreateMarketingEmail(ewl *models.EmailWithList) int {
	if len(ewl.Email) == 0 {
		return http.StatusBadRequest
	}

	var mes models.SendgridEmails
	mes = append(mes, models.SendgridEmail{Email: ewl.Email})

	go CreateNewEmailOnList(&mes, os.Getenv("SENDGRID_MARKETING_UPDATES_LIST"))

	return http.StatusAccepted
}
