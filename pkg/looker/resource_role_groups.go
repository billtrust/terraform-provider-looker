package looker

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func resourceRoleGroups() *schema.Resource {
	return &schema.Resource{
		Create: resourceRoleGroupsCreate,
		Read:   resourceRoleGroupsRead,
		Update: resourceRoleGroupsUpdate,
		Delete: resourceRoleGroupsDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRoleGroupsImport,
		},

		Schema: map[string]*schema.Schema{
			"role_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_ids": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func resourceRoleGroupsCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerSDK)

	roleIDString := d.Get("role_id").(string)

	roleID, err := strconv.ParseInt(roleIDString, 10, 64)
	if err != nil {
		return err
	}

	var groupIDs []int64
	for _, groupID := range d.Get("group_ids").(*schema.Set).List() {
		groupIDs = append(groupIDs, groupID.(int64))
	}

	_, err = client.SetRoleGroups(roleID, groupIDs, nil)
	if err != nil {
		return err
	}

	d.SetId(roleIDString)

	return resourceRoleGroupsRead(d, m)
}

func resourceRoleGroupsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerSDK)

	roleID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return err
	}

	groups, err := client.RoleGroups(roleID, "", nil)
	if err != nil {
		return err
	}

	var groupIDs []int64
	for _, group := range groups {
		groupIDs = append(groupIDs, *group.Id)
	}

	if err = d.Set("role_id", roleID); err != nil {
		return err
	}

	if err = d.Set("group_ids", groupIDs); err != nil {
		return err
	}

	return nil
}

func resourceRoleGroupsUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerSDK)

	roleID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return err
	}

	var groupIDs []int64
	for _, groupID := range d.Get("group_ids").(*schema.Set).List() {
		groupIDs = append(groupIDs, groupID.(int64))
	}

	_, err = client.SetRoleGroups(roleID, groupIDs, nil)
	if err != nil {
		return err
	}

	return resourceRoleGroupsRead(d, m)
}

func resourceRoleGroupsDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerSDK)

	roleID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return err
	}

	groupIDs := []int64{}
	_, err = client.SetRoleGroups(roleID, groupIDs, nil)
	if err != nil {
		return err
	}

	return nil
}

func resourceRoleGroupsImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceRoleGroupsRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
