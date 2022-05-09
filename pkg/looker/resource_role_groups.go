package looker

import (
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
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceRoleGroupsCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerSDK)

	roleID := d.Get("role_id").(string)

	var groupIDs []string
	for _, groupID := range d.Get("group_ids").(*schema.Set).List() {
		groupIDs = append(groupIDs, groupID.(string))
	}

	_, err := client.SetRoleGroups(roleID, groupIDs, nil)
	if err != nil {
		return err
	}

	d.SetId(roleID)

	return resourceRoleGroupsRead(d, m)
}

func resourceRoleGroupsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerSDK)

	roleID := d.Id()

	groups, err := client.RoleGroups(roleID, "", nil)
	if err != nil {
		return err
	}

	var groupIDs []string
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

	roleID := d.Id()

	var groupIDs []string
	for _, groupID := range d.Get("group_ids").(*schema.Set).List() {
		groupIDs = append(groupIDs, groupID.(string))
	}

	_, err := client.SetRoleGroups(roleID, groupIDs, nil)
	if err != nil {
		return err
	}

	return resourceRoleGroupsRead(d, m)
}

func resourceRoleGroupsDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerSDK)

	roleID := d.Id()

	groupIDs := []string{}
	_, err := client.SetRoleGroups(roleID, groupIDs, nil)
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
