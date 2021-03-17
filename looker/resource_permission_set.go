package looker

import (
	"strings"

	"github.com/Foxtel-DnA/looker-go-sdk/client/role"

	apiclient "github.com/Foxtel-DnA/looker-go-sdk/client"
	"github.com/Foxtel-DnA/looker-go-sdk/models"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePermissionSet() *schema.Resource {
	return &schema.Resource{
		Create: resourcePermissionSetCreate,
		Read:   resourcePermissionSetRead,
		Update: resourcePermissionSetUpdate,
		Delete: resourcePermissionSetDelete,
		Exists: resourcePermissionSetExists,
		Importer: &schema.ResourceImporter{
			State: resourcePermissionSetImport,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"permissions": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourcePermissionSetCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	var permissions []string
	for _, permission := range d.Get("permissions").(*schema.Set).List() {
		permissions = append(permissions, permission.(string))
	}

	params := role.NewCreatePermissionSetParams()
	params.Body = &models.PermissionSet{}
	params.Body.Name = d.Get("name").(string)
	params.Body.Permissions = permissions

	result, err := client.Role.CreatePermissionSet(params)
	if err != nil {
		return err
	}

	d.SetId(getStringFromID(result.Payload.ID))

	return resourcePermissionSetRead(d, m)
}

func resourcePermissionSetRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	ID, err := getIDFromString(d.Id())
	if err != nil {
		return err
	}

	params := role.NewPermissionSetParams()
	params.PermissionSetID = ID

	result, err := client.Role.PermissionSet(params)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", result.Payload.Name)
	d.Set("permissions", result.Payload.Permissions)

	return nil
}

func resourcePermissionSetUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	ID, err := getIDFromString(d.Id())
	if err != nil {
		return err
	}

	var permissions []string
	for _, permission := range d.Get("permissions").(*schema.Set).List() {
		permissions = append(permissions, permission.(string))
	}

	params := role.NewUpdatePermissionSetParams()
	params.PermissionSetID = ID
	params.Body = &models.PermissionSet{}
	params.Body.Name = d.Get("name").(string)
	params.Body.Permissions = permissions

	_, err = client.Role.UpdatePermissionSet(params)
	if err != nil {
		return err
	}

	return resourcePermissionSetRead(d, m)
}

func resourcePermissionSetDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	ID, err := getIDFromString(d.Id())
	if err != nil {
		return err
	}

	params := role.NewDeletePermissionSetParams()
	params.PermissionSetID = ID

	_, err = client.Role.DeletePermissionSet(params)
	if err != nil {
		return err
	}

	return nil
}

func resourcePermissionSetExists(d *schema.ResourceData, m interface{}) (b bool, e error) {
	// Exists - This is called to verify a resource still exists. It is called prior to Read,
	// and lowers the burden of Read to be able to assume the resource exists.
	client := m.(*apiclient.Looker)

	ID, err := getIDFromString(d.Id())
	if err != nil {
		return false, err
	}

	params := role.NewPermissionSetParams()
	params.PermissionSetID = ID

	_, err = client.Role.PermissionSet(params)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func resourcePermissionSetImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourcePermissionSetRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
