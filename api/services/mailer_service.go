package services

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/decentorganization/topaz/shared/models"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendWelcomeEmail ...
func SendWelcomeEmail(to string) {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	if apiKey == "" {
		return
	}

	f := mail.NewEmail(os.Getenv("SENDGRID_FROM_NAME"), os.Getenv("SENDGRID_FROM_EMAIL"))
	t := mail.NewEmail("", to)

	email := mail.NewV3MailInit(f, "", t)
	email.SetTemplateID(os.Getenv("SENDGRID_WELCOME_EMAIL_ID"))

	p := email.Personalizations[0]
	p.SetDynamicTemplateData("Sender_Name", os.Getenv("SENDGRID_SENDER_NAME"))
	p.SetDynamicTemplateData("Sender_Address", os.Getenv("SENDGRID_SENDER_ADDRESS"))
	p.SetDynamicTemplateData("Sender_City", os.Getenv("SENDGRID_SENDER_CITY"))
	p.SetDynamicTemplateData("Sender_State", os.Getenv("SENDGRID_SENDER_STATE"))
	p.SetDynamicTemplateData("Sender_Zip", os.Getenv("SENDGRID_SENDER_ZIP"))

	client := sendgrid.NewSendClient(apiKey)
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
func SendPasswordResetEmail(to, tokenURL string) {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	if apiKey == "" {
		return
	}

	f := mail.NewEmail(os.Getenv("SENDGRID_FROM_NAME"), os.Getenv("SENDGRID_FROM_EMAIL"))
	t := mail.NewEmail("", to)

	email := mail.NewV3MailInit(f, "", t)
	email.SetTemplateID(os.Getenv("SENDGRID_PASSWORD_RESET_ID"))

	p := email.Personalizations[0]
	p.SetDynamicTemplateData("Sender_Name", os.Getenv("SENDGRID_SENDER_NAME"))
	p.SetDynamicTemplateData("Sender_Address", os.Getenv("SENDGRID_SENDER_ADDRESS"))
	p.SetDynamicTemplateData("Sender_City", os.Getenv("SENDGRID_SENDER_CITY"))
	p.SetDynamicTemplateData("Sender_State", os.Getenv("SENDGRID_SENDER_STATE"))
	p.SetDynamicTemplateData("Sender_Zip", os.Getenv("SENDGRID_SENDER_ZIP"))
	p.SetDynamicTemplateData("password_reset_token", tokenURL)
	p.SetDynamicTemplateData("user_email", to)

	client := sendgrid.NewSendClient(apiKey)
	response, err := client.Send(email)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}

// CreateNewEmailOnList sleeps, so be sure to call it with a goroutine
func CreateNewEmailOnList(to *models.SendgridEmails, list string) {
	sgr, ok := addEmailToContacts(to)
	if ok {
		time.Sleep(5 * time.Second)
		addContactToList(sgr.PersistedRecipients[0], list)
	}
}

func addEmailToContacts(to *models.SendgridEmails) (*models.SendgridNewContactRespone, bool) {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	if apiKey == "" {
		return nil, false
	}

	t, _ := json.Marshal(&to)

	create := sendgrid.GetRequest(apiKey, "/v3/contactdb/recipients", os.Getenv("SENDGRID_API_ROOT"))
	create.Method = "POST"
	create.Body = t

	createRes, err := sendgrid.API(create)

	if err != nil {
		fmt.Println(err)
		return nil, false
	}

	sgr := new(models.SendgridNewContactRespone)
	if err = json.Unmarshal([]byte(createRes.Body), &sgr); err != nil {
		fmt.Println(err)
		return nil, false
	}

	if len(sgr.PersistedRecipients) == 0 {
		fmt.Println("didn't create the new contact")
		return nil, false
	}

	return sgr, true
}

func addContactToList(contactID, listID string) {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	if apiKey == "" {
		return
	}

	listReq := fmt.Sprintf("/v3/contactdb/lists/%s/recipients/%s", listID, contactID)
	list := sendgrid.GetRequest(apiKey, listReq, os.Getenv("SENDGRID_API_ROOT"))
	list.Method = "POST"

	_, err := sendgrid.API(list)

	if err != nil {
		fmt.Println(err)
	}
}
