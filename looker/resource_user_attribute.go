package looker

import (
	"strings"

	"github.com/Foxtel-DnA/looker-go-sdk/client/user_attribute"

	apiclient "github.com/Foxtel-DnA/looker-go-sdk/client"
	"github.com/Foxtel-DnA/looker-go-sdk/models"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceUserAttribute() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserAttributeCreate,
		Read:   resourceUserAttributeRead,
		Update: resourceUserAttributeUpdate,
		Delete: resourceUserAttributeDelete,
		Exists: resourceUserAttributeExists,
		Importer: &schema.ResourceImporter{
			State: resourceUserAttributeImport,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				// TODO: hard code to only allow "advanced_filter_string" for now
			},
			"label": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceUserAttributeCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	params := user_attribute.NewCreateUserAttributeParams()
	params.Body = &models.UserAttribute{}
	params.Body.Name = d.Get("name").(string)
	params.Body.Type = d.Get("type").(string)
	params.Body.Label = d.Get("label").(string)

	result, err := client.UserAttribute.CreateUserAttribute(params)
	if err != nil {
		return err
	}

	d.SetId(getStringFromID(result.Payload.ID))

	return resourceUserAttributeRead(d, m)
}

func resourceUserAttributeRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	ID, err := getIDFromString(d.Id())
	if err != nil {
		return err
	}

	params := user_attribute.NewUserAttributeParams()
	params.UserAttributeID = ID

	result, err := client.UserAttribute.UserAttribute(params)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", result.Payload.Name)
	d.Set("type", result.Payload.Type)
	d.Set("label", result.Payload.Label)

	return nil
}

func resourceUserAttributeUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	ID, err := getIDFromString(d.Id())
	if err != nil {
		return err
	}

	params := user_attribute.NewUpdateUserAttributeParams()
	params.UserAttributeID = ID
	params.Body = &models.UserAttribute{}
	params.Body.Name = d.Get("name").(string)
	params.Body.Type = d.Get("type").(string)
	params.Body.Label = d.Get("label").(string)

	_, err = client.UserAttribute.UpdateUserAttribute(params)
	if err != nil {
		return err
	}

	return resourceUserAttributeRead(d, m)
}

func resourceUserAttributeDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	ID, err := getIDFromString(d.Id())
	if err != nil {
		return err
	}

	params := user_attribute.NewDeleteUserAttributeParams()
	params.UserAttributeID = ID

	_, err = client.UserAttribute.DeleteUserAttribute(params)
	if err != nil {
		return err
	}

	return nil
}

func resourceUserAttributeExists(d *schema.ResourceData, m interface{}) (b bool, e error) {
	// Exists - This is called to verify a resource still exists. It is called prior to Read,
	// and lowers the burden of Read to be able to assume the resource exists.
	client := m.(*apiclient.Looker)

	ID, err := getIDFromString(d.Id())
	if err != nil {
		return false, err
	}

	params := user_attribute.NewUserAttributeParams()
	params.UserAttributeID = ID

	_, err = client.UserAttribute.UserAttribute(params)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func resourceUserAttributeImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceUserAttributeRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
