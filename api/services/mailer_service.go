package services

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

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

// CreateNewMarketingEmail ...
func CreateNewMarketingEmail(to *models.SendgridEmails) bool {
	sgr, ok := addEmailToContacts(to)
	if !ok {
		return false
	}

	if ok = addContactToList(sgr.PersistedRecipients[0], os.Getenv("SENDGRID_MARKETING_UPDATES_LIST")); !ok {
		return false
	}

	return true
}

// CreateNewAPIUserEmail ...
func CreateNewAPIUserEmail(to *models.SendgridEmails) bool {
	sgr, ok := addEmailToContacts(to)
	if !ok {
		return false
	}

	if ok = addContactToList(sgr.PersistedRecipients[0], os.Getenv("SENDGRID_API_USERS_LIST")); !ok {
		return false
	}

	return true
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
