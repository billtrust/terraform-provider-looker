package looker

import (
	"strconv"

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
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func resourceUserRolesCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerSDK)

	userIDString := d.Get("user_id").(string)

	userID, err := strconv.ParseInt(userIDString, 10, 64)
	if err != nil {
		return err
	}

	var roleIDs []int64
	for _, roleID := range d.Get("role_ids").(*schema.Set).List() {
		roleIDs = append(roleIDs, roleID.(int64))
	}

	_, err = client.SetUserRoles(userID, roleIDs, "", nil)
	if err != nil {
		return err
	}

	d.SetId(userIDString)

	return resourceUserRolesRead(d, m)
}

func resourceUserRolesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerSDK)

	userID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return err
	}

	request := apiclient.RequestUserRoles{UserId: userID}

	userRoles, err := client.UserRoles(request, nil)
	if err != nil {
		return err
	}

	var roleIDs []int64
	for _, role := range userRoles {
		roleIDs = append(roleIDs, *role.Id)
	}

	if err = d.Set("user_id", userID); err != nil {
		return err
	}

	if err = d.Set("role_ids", roleIDs); err != nil {
		return err
	}

	return nil
}

func resourceUserRolesUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerSDK)

	userID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return err
	}

	var roleIDs []int64
	for _, roleID := range d.Get("role_ids").(*schema.Set).List() {
		roleIDs = append(roleIDs, roleID.(int64))
	}

	_, err = client.SetUserRoles(userID, roleIDs, "", nil)
	if err != nil {
		return err
	}

	return resourceUserRolesRead(d, m)
}

func resourceUserRolesDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerSDK)

	userID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return err
	}

	roleIDs := []int64{}
	_, err = client.SetUserRoles(userID, roleIDs, "", nil)
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
