package looker

import (
	"strings"

	apiclient "github.com/billtrust/looker-go-sdk/client"
	"github.com/billtrust/looker-go-sdk/client/user"
	"github.com/billtrust/looker-go-sdk/models"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,
		Exists: resourceUserExists,
		Importer: &schema.ResourceImporter{
			State: resourceUserImport,
		},

		Schema: map[string]*schema.Schema{
			"first_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"last_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceUserCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	params := user.NewCreateUserParams()
	params.Body = &models.User{}
	params.Body.FirstName = d.Get("first_name").(string)
	params.Body.LastName = d.Get("last_name").(string)

	user, err := client.User.CreateUser(params)
	if err != nil {
		return err
	}

	d.SetId(getStringFromID(user.Payload.ID))

	return resourceUserRead(d, m)
}

func resourceUserRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	userID, err := getIDFromString(d.Id())
	if err != nil {
		return err
	}

	params := user.NewUserParams()
	params.UserID = userID

	user, err := client.User.User(params)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("first_name", user.Payload.FirstName)
	d.Set("last_name", user.Payload.LastName)

	return nil
}

func resourceUserUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	userID, err := getIDFromString(d.Id())
	if err != nil {
		return err
	}

	params := user.NewUpdateUserParams()
	params.UserID = userID
	params.Body = &models.User{}
	params.Body.FirstName = d.Get("first_name").(string)
	params.Body.LastName = d.Get("last_name").(string)

	_, err = client.User.UpdateUser(params)
	if err != nil {
		return err
	}

	return resourceUserRead(d, m)
}

func resourceUserDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	userID, err := getIDFromString(d.Id())
	if err != nil {
		return err
	}

	params := user.NewDeleteUserParams()
	params.UserID = userID

	_, err = client.User.DeleteUser(params)
	if err != nil {
		return err
	}

	return nil
}

func resourceUserExists(d *schema.ResourceData, m interface{}) (b bool, e error) {
	// Exists - This is called to verify a resource still exists. It is called prior to Read,
	// and lowers the burden of Read to be able to assume the resource exists.
	client := m.(*apiclient.Looker)

	userID, err := getIDFromString(d.Id())
	if err != nil {
		return false, err
	}

	params := user.NewUserParams()
	params.UserID = userID

	_, err = client.User.User(params)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func resourceUserImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceUserRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
