package services

import (
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendWelcomeEmail ...
func SendWelcomeEmail(to string) {
	f := mail.NewEmail(os.Getenv("SENDGRID_FROM_NAME"), os.Getenv("SENDGRID_FROM_EMAIL"))
	t := mail.NewEmail("", to)

	email := mail.NewV3MailInit(f, "", t)
	email.SetTemplateID(os.Getenv("SENDGRID_WELCOME_EMAIL_ID"))

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	_, err := client.Send(email)

	if err != nil {
		fmt.Println(err)
	}
}
