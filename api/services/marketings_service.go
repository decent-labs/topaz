package services

import (
	"net/http"
	"os"

	"github.com/decentorganization/topaz/shared/models"
)

// CreateMarketingEmail ...
func CreateMarketingEmail(me *models.SendgridEmail) int {
	if len(me.Email) == 0 {
		return http.StatusBadRequest
	}

	var mes models.SendgridEmails
	mes = append(mes, *me)

	go CreateNewEmailOnList(&mes, os.Getenv("SENDGRID_MARKETING_UPDATES_LIST"))

	return http.StatusAccepted
}
