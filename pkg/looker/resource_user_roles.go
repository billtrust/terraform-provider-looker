package looker

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func resourceUserRoles() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserRolesCreate,
		Read:   resourceUserRolesRead,
		Update: resourceUserRolesUpdate,
		Delete: resourceUserRolesDelete,
		Importer: &schema.ResourceImporter{
			State: resourceUserRolesImport,
		},

		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role_ids": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceUserRolesCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerSDK)

	userID := d.Get("user_id").(string)

	var roleIDs []string
	for _, roleID := range d.Get("role_ids").(*schema.Set).List() {
		roleIDs = append(roleIDs, roleID.(string))
	}

	_, err := client.SetUserRoles(userID, roleIDs, "", nil)
	if err != nil {
		return err
	}

	d.SetId(userID)

	return resourceUserRolesRead(d, m)
}

func resourceUserRolesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerSDK)

	userID := d.Id()

	request := apiclient.RequestUserRoles{UserId: userID}

	userRoles, err := client.UserRoles(request, nil)
	if err != nil {
		return err
	}

	var roleIDs []string
	for _, role := range userRoles {
		roleIDs = append(roleIDs, *role.Id)
	}

	if err = d.Set("user_id", d.Id()); err != nil {
		return err
	}

	if err = d.Set("role_ids", roleIDs); err != nil {
		return err
	}

	return nil
}

func resourceUserRolesUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerSDK)

	userID := d.Id()

	var roleIDs []string
	for _, roleID := range d.Get("role_ids").(*schema.Set).List() {
		roleIDs = append(roleIDs, roleID.(string))
	}

	_, err := client.SetUserRoles(userID, roleIDs, "", nil)
	if err != nil {
		return err
	}

	return resourceUserRolesRead(d, m)
}

func resourceUserRolesDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerSDK)

	userID := d.Id()

	roleIDs := []string{}
	_, err := client.SetUserRoles(userID, roleIDs, "", nil)
	if err != nil {
		return err
	}

	return nil
}

func resourceUserRolesImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceUserRolesRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
