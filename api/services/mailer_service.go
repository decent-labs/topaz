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
func SendPasswordResetEmail(to, tokenURL string) {
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

// CreateNewEmailOnList sleeps, so be sure to call it with a goroutine
func CreateNewEmailOnList(to *models.SendgridEmails, list string) {
	sgr, ok := addEmailToContacts(to)
	if ok {
		time.Sleep(5 * time.Second)
		addContactToList(sgr.PersistedRecipients[0], list)
	}
}

func addEmailToContacts(to *models.SendgridEmails) (*models.SendgridNewContactRespone, bool) {
	t, _ := json.Marshal(&to)

	create := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), "/v3/contactdb/recipients", os.Getenv("SENDGRID_API_ROOT"))
	create.Method = "POST"
	create.Body = t

	createRes, err := sendgrid.API(create)

	fmt.Println(createRes.StatusCode)
	fmt.Println(createRes.Body)
	fmt.Println(createRes.Headers)

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

func addContactToList(contactID, listID string) bool {
	listReq := fmt.Sprintf("/v3/contactdb/lists/%s/recipients/%s", listID, contactID)
	list := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), listReq, os.Getenv("SENDGRID_API_ROOT"))
	list.Method = "POST"

	listRes, err := sendgrid.API(list)

	fmt.Println(listRes.StatusCode)
	fmt.Println(listRes.Body)
	fmt.Println(listRes.Headers)

	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}
