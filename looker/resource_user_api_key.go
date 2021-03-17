package looker

import (
	"fmt"
	"strings"

	apiclient "github.com/Foxtel-DnA/looker-go-sdk/client"
	"github.com/Foxtel-DnA/looker-go-sdk/client/user"
	"github.com/Foxtel-DnA/looker-go-sdk/models"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceUserAPIKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserAPIKeyCreate,
		Read:   resourceUserAPIKeyRead,
		Delete: resourceUserAPIKeyDelete,
		Exists: resourceUserAPIKeyExists,
		Importer: &schema.ResourceImporter{
			State: resourceUserAPIKeyImport,
		},

		Schema: map[string]*schema.Schema{
			"user_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"client_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceUserAPIKeyCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	sUserID := d.Get("user_id").(string)

	iUserID, err := getIDFromString(sUserID)
	if err != nil {
		return err
	}

	params := user.NewCreateUserCredentialsApi3Params()
	params.UserID = iUserID
	params.Body = &models.CredentialsApi3{}

	resp, err := client.User.CreateUserCredentialsApi3(params)
	if err != nil {
		return err
	}

	id := sUserID + ":" + getStringFromID(resp.Payload.ID)
	d.SetId(id)

	return resourceUserAPIKeyRead(d, m)
}

func resourceUserAPIKeyRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	id := strings.Split(d.Id(), ":")
	if len(id) != 2 {
		return fmt.Errorf("ID Should be two strings separated by a colon (:)")
	}

	sUserID := id[0]
	sAPIID := id[1]

	iUserID, err := getIDFromString(sUserID)
	if err != nil {
		return err
	}

	iAPIID, err := getIDFromString(sAPIID)
	if err != nil {
		return err
	}

	params := user.NewUserCredentialsApi3Params()
	params.UserID = iUserID
	params.CredentialsApi3ID = iAPIID

	resp, err := client.User.UserCredentialsApi3(params)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("user_id", sUserID)
	d.Set("client_id", resp.Payload.ClientID)

	return nil
}

func resourceUserAPIKeyDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	id := strings.Split(d.Id(), ":")
	if len(id) != 2 {
		return fmt.Errorf("ID Should be two strings separated by a colon (:)")
	}

	sUserID := id[0]
	sAPIID := id[1]

	iUserID, err := getIDFromString(sUserID)
	if err != nil {
		return err
	}

	iAPIID, err := getIDFromString(sAPIID)
	if err != nil {
		return err
	}

	params := user.NewDeleteUserCredentialsApi3Params()
	params.UserID = iUserID
	params.CredentialsApi3ID = iAPIID

	_, err = client.User.DeleteUserCredentialsApi3(params)
	if err != nil {
		return err
	}

	return nil
}

func resourceUserAPIKeyExists(d *schema.ResourceData, m interface{}) (b bool, e error) {
	// Exists - This is called to verify a resource still exists. It is called prior to Read,
	// and lowers the burden of Read to be able to assume the resource exists.
	client := m.(*apiclient.Looker)

	id := strings.Split(d.Id(), ":")
	if len(id) != 2 {
		return false, fmt.Errorf("ID Should be two strings separated by a colon (:)")
	}

	sUserID := id[0]
	sAPIID := id[1]

	iUserID, err := getIDFromString(sUserID)
	if err != nil {
		return false, err
	}

	iAPIID, err := getIDFromString(sAPIID)
	if err != nil {
		return false, err
	}

	params := user.NewUserCredentialsApi3Params()
	params.UserID = iUserID
	params.CredentialsApi3ID = iAPIID

	_, err = client.User.UserCredentialsApi3(params)

	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func resourceUserAPIKeyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceUserAPIKeyRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
