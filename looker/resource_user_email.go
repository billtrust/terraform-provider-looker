package looker

import (
	"log"
	"strings"

	apiclient "github.com/Foxtel-DnA/looker-go-sdk/client"
	"github.com/Foxtel-DnA/looker-go-sdk/client/user"
	"github.com/Foxtel-DnA/looker-go-sdk/models"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceUserEmail() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserEmailCreate,
		Read:   resourceUserEmailRead,
		Update: resourceUserEmailUpdate,
		Delete: resourceUserEmailDelete,
		Exists: resourceUserEmailExists,
		Importer: &schema.ResourceImporter{
			State: resourceUserEmailImport,
		},

		Schema: map[string]*schema.Schema{
			"user_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceUserEmailCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	params := user.NewCreateUserCredentialsEmailParams()

	sUserID := d.Get("user_id").(string)

	iUserID, err := getIDFromString(sUserID)
	if err != nil {
		return err
	}

	params.UserID = iUserID
	params.Body = &models.CredentialsEmail{}
	params.Body.Email = d.Get("email").(string)

	json, err := getJSONString(params)
	log.Printf("This is my first log message %s", json)

	_, err = client.User.CreateUserCredentialsEmail(params)
	if err != nil {
		return err
	}

	d.SetId(sUserID)

	return resourceUserEmailRead(d, m)
}

func resourceUserEmailRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	userID, err := getIDFromString(d.Id())
	if err != nil {
		return err
	}

	params := user.NewUserCredentialsEmailParams()
	params.UserID = userID

	user, err := client.User.UserCredentialsEmail(params)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("email", user.Payload.Email)

	return nil
}

func resourceUserEmailUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	userID, err := getIDFromString(d.Id())
	if err != nil {
		return err
	}

	params := user.NewUpdateUserCredentialsEmailParams()
	params.UserID = userID
	params.Body = &models.CredentialsEmail{}
	params.Body.Email = d.Get("email").(string)

	_, err = client.User.UpdateUserCredentialsEmail(params)
	if err != nil {
		return err
	}

	return resourceUserEmailRead(d, m)
}

func resourceUserEmailDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	userID, err := getIDFromString(d.Id())
	if err != nil {
		return err
	}

	params := user.NewDeleteUserCredentialsEmailParams()
	params.UserID = userID

	_, err = client.User.DeleteUserCredentialsEmail(params)
	if err != nil {
		return err
	}

	return nil
}

func resourceUserEmailExists(d *schema.ResourceData, m interface{}) (b bool, e error) {
	// Exists - This is called to verify a resource still exists. It is called prior to Read,
	// and lowers the burden of Read to be able to assume the resource exists.
	client := m.(*apiclient.Looker)

	userID, err := getIDFromString(d.Id())
	if err != nil {
		return false, err
	}

	params := user.NewUserCredentialsEmailParams()
	params.UserID = userID

	_, err = client.User.UserCredentialsEmail(params)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func resourceUserEmailImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceUserEmailRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
