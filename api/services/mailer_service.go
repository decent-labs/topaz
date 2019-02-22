package services

import (
	"fmt"
	"os"
	"strconv"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendWelcomeEmail ...
func SendWelcomeEmail(to string) {
	f := mail.NewEmail(os.Getenv("SENDGRID_FROM_NAME"), os.Getenv("SENDGRID_FROM_EMAIL"))
	t := mail.NewEmail("", to)

	email := mail.NewV3MailInit(f, "", t)
	email.SetTemplateID(os.Getenv("SENDGRID_WELCOME_EMAIL_ID"))

	asm := mail.NewASM()
	asmID, _ := strconv.Atoi(os.Getenv("SENDGRID_ACCOUNT_ASM"))
	asm.SetGroupID(asmID)
	asm.AddGroupsToDisplay(asmID)
	email.SetASM(asm)

	p := email.Personalizations[0]
	p.SetDynamicTemplateData("Sender_Name", os.Getenv("SENDGRID_SENDER_NAME"))
	p.SetDynamicTemplateData("Sender_Address", os.Getenv("SENDGRID_SENDER_ADDRESS"))
	p.SetDynamicTemplateData("Sender_City", os.Getenv("SENDGRID_SENDER_CITY"))
	p.SetDynamicTemplateData("Sender_State", os.Getenv("SENDGRID_SENDER_STATE"))
	p.SetDynamicTemplateData("Sender_Zip", os.Getenv("SENDGRID_SENDER_ZIP"))

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(email)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}

// SendPasswordResetEmail ...
func SendPasswordResetEmail(to, token string) {
	f := mail.NewEmail(os.Getenv("SENDGRID_FROM_NAME"), os.Getenv("SENDGRID_FROM_EMAIL"))
	t := mail.NewEmail("", to)

	email := mail.NewV3MailInit(f, "", t)
	email.SetTemplateID(os.Getenv("SENDGRID_PASSWORD_RESET_ID"))

	asm := mail.NewASM()
	asmID, _ := strconv.Atoi(os.Getenv("SENDGRID_ACCOUNT_ASM"))
	asm.SetGroupID(asmID)
	asm.AddGroupsToDisplay(asmID)
	email.SetASM(asm)

	p := email.Personalizations[0]
	p.SetDynamicTemplateData("Sender_Name", os.Getenv("SENDGRID_SENDER_NAME"))
	p.SetDynamicTemplateData("Sender_Address", os.Getenv("SENDGRID_SENDER_ADDRESS"))
	p.SetDynamicTemplateData("Sender_City", os.Getenv("SENDGRID_SENDER_CITY"))
	p.SetDynamicTemplateData("Sender_State", os.Getenv("SENDGRID_SENDER_STATE"))
	p.SetDynamicTemplateData("Sender_Zip", os.Getenv("SENDGRID_SENDER_ZIP"))
	p.SetDynamicTemplateData("password_reset_token", token)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(email)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}
