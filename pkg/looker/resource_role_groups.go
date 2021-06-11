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
				Elem:     &schema.Schema{Type: schema.TypeString},
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
		gID, err := strconv.ParseInt(groupID.(string), 10, 64)
		if err != nil {
			return err
		}
		groupIDs = append(groupIDs, gID)
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

	var groupIDs []string
	for _, group := range groups {
		gID := strconv.Itoa(int(*group.Id))
		groupIDs = append(groupIDs, gID)
	}

	if err = d.Set("role_id", strconv.Itoa(int(roleID))); err != nil {
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
		gID, err := strconv.ParseInt(groupID.(string), 10, 64)
		if err != nil {
			return err
		}
		groupIDs = append(groupIDs, gID)
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
