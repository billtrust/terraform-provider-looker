package looker

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func resourceRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceRoleCreate,
		Read:   resourceRoleRead,
		Update: resourceRoleUpdate,
		Delete: resourceRoleDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRoleImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"permission_set_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"model_set_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceRoleCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerSDK)

	roleName := d.Get("name").(string)
	permissionSetID := d.Get("permission_set_id").(int64)
	modelSetID := d.Get("model_set_id").(int64)

	writeRole := apiclient.WriteRole{
		Name:            &roleName,
		PermissionSetId: &permissionSetID,
		ModelSetId:      &modelSetID,
	}

	role, err := client.CreateRole(writeRole, nil)
	if err != nil {
		return err
	}

	roleID := *role.Id
	d.SetId(strconv.Itoa(int(roleID)))

	return resourceRoleRead(d, m)
}

func resourceRoleRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerSDK)

	roleID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return err
	}

	role, err := client.Role(roleID, nil)
	if err != nil {
		return err
	}

	if err = d.Set("name", role.Name); err != nil {
		return err
	}
	if err = d.Set("permission_set_id", role.PermissionSetId); err != nil {
		return err
	}
	if err = d.Set("model_set_id", role.ModelSetId); err != nil {
		return err
	}

	return nil
}

func resourceRoleUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerSDK)

	roleID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return err
	}

	roleName := d.Get("name").(string)
	permissionSetID := d.Get("permission_set_id").(int64)
	modelSetID := d.Get("model_set_id").(int64)
	writeRole := apiclient.WriteRole{
		Name:            &roleName,
		PermissionSetId: &permissionSetID,
		ModelSetId:      &modelSetID,
	}
	_, err = client.UpdateRole(roleID, writeRole, nil)
	if err != nil {
		return err
	}

	return resourceRoleRead(d, m)
}

func resourceRoleDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerSDK)

	roleID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return err
	}

	_, err = client.DeleteRole(roleID, nil)
	if err != nil {
		return err
	}

	return nil
}

func resourceRoleImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceRoleRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
