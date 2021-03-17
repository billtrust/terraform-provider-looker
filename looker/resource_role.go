package looker

import (
	"strings"

	"github.com/Foxtel-DnA/looker-go-sdk/client/role"

	apiclient "github.com/Foxtel-DnA/looker-go-sdk/client"
	"github.com/Foxtel-DnA/looker-go-sdk/models"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceRoleCreate,
		Read:   resourceRoleRead,
		Update: resourceRoleUpdate,
		Delete: resourceRoleDelete,
		Exists: resourceRoleExists,
		Importer: &schema.ResourceImporter{
			State: resourceRoleImport,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"permission_set_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"model_set_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceRoleCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	permissionSetID, err := getIDFromString(d.Get("permission_set_id").(string))
	if err != nil {
		return err
	}

	modelSetID, err := getIDFromString(d.Get("model_set_id").(string))
	if err != nil {
		return err
	}

	params := role.NewCreateRoleParams()
	params.Body = &models.Role{}
	params.Body.Name = d.Get("name").(string)
	params.Body.PermissionSetID = permissionSetID
	params.Body.ModelSetID = modelSetID

	result, err := client.Role.CreateRole(params)
	if err != nil {
		return err
	}

	d.SetId(getStringFromID(result.Payload.ID))

	return resourceRoleRead(d, m)
}

func resourceRoleRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	ID, err := getIDFromString(d.Id())
	if err != nil {
		return err
	}

	params := role.NewRoleParams()
	params.RoleID = ID

	result, err := client.Role.Role(params)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", result.Payload.Name)
	d.Set("permission_set_id", result.Payload.PermissionSetID)
	d.Set("model_set_id", result.Payload.ModelSetID)

	return nil
}

func resourceRoleUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	ID, err := getIDFromString(d.Id())
	if err != nil {
		return err
	}

	permissionSetID, err := getIDFromString(d.Get("permission_set_id").(string))
	if err != nil {
		return err
	}

	modelSetID, err := getIDFromString(d.Get("model_set_id").(string))
	if err != nil {
		return err
	}

	params := role.NewUpdateRoleParams()
	params.RoleID = ID
	params.Body = &models.Role{}
	params.Body.Name = d.Get("name").(string)
	params.Body.PermissionSetID = permissionSetID
	params.Body.ModelSetID = modelSetID

	_, err = client.Role.UpdateRole(params)
	if err != nil {
		return err
	}

	return resourceRoleRead(d, m)
}

func resourceRoleDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	ID, err := getIDFromString(d.Id())
	if err != nil {
		return err
	}

	params := role.NewDeleteRoleParams()
	params.RoleID = ID

	_, err = client.Role.DeleteRole(params)
	if err != nil {
		return err
	}

	return nil
}

func resourceRoleExists(d *schema.ResourceData, m interface{}) (b bool, e error) {
	// Exists - This is called to verify a resource still exists. It is called prior to Read,
	// and lowers the burden of Read to be able to assume the resource exists.
	client := m.(*apiclient.Looker)

	ID, err := getIDFromString(d.Id())
	if err != nil {
		return false, err
	}

	params := role.NewRoleParams()
	params.RoleID = ID

	_, err = client.Role.Role(params)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func resourceRoleImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceRoleRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
