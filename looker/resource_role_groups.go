package looker

import (
	"strings"

	"github.com/Foxtel-DnA/looker-go-sdk/client/role"

	apiclient "github.com/Foxtel-DnA/looker-go-sdk/client"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceRoleGroups() *schema.Resource {
	return &schema.Resource{
		Create: resourceRoleGroupsCreate,
		Read:   resourceRoleGroupsRead,
		Update: resourceRoleGroupsUpdate,
		Delete: resourceRoleGroupsDelete,
		Exists: resourceRoleGroupsExists,
		Importer: &schema.ResourceImporter{
			State: resourceRoleGroupsImport,
		},

		Schema: map[string]*schema.Schema{
			"role_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"group_ids": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceRoleGroupsCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	ID, err := getIDFromString(d.Get("role_id").(string))
	if err != nil {
		return err
	}

	var groupIDs []int64
	for _, sGroupID := range d.Get("group_ids").(*schema.Set).List() {
		iGroupID, err := getIDFromString(sGroupID.(string))
		if err != nil {
			return err
		}

		groupIDs = append(groupIDs, iGroupID)
	}

	params := role.NewSetRoleGroupsParams()
	params.RoleID = ID
	params.Body = groupIDs

	_, err = client.Role.SetRoleGroups(params)
	if err != nil {
		return err
	}

	d.SetId(d.Get("role_id").(string))

	return resourceRoleGroupsRead(d, m)
}

func resourceRoleGroupsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	ID, err := getIDFromString(d.Id())
	if err != nil {
		return err
	}

	params := role.NewRoleGroupsParams()
	params.RoleID = ID

	result, err := client.Role.RoleGroups(params)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("role_id", ID)
	d.Set("group_ids", result.Payload)

	return nil
}

func resourceRoleGroupsUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	ID, err := getIDFromString(d.Id())
	if err != nil {
		return err
	}

	var groupIDs []int64
	for _, sGroupID := range d.Get("group_ids").(*schema.Set).List() {
		iGroupID, err := getIDFromString(sGroupID.(string))
		if err != nil {
			return err
		}

		groupIDs = append(groupIDs, iGroupID)
	}

	params := role.NewSetRoleGroupsParams()
	params.RoleID = ID
	params.Body = groupIDs

	_, err = client.Role.SetRoleGroups(params)
	if err != nil {
		return err
	}

	return resourceRoleGroupsRead(d, m)
}

func resourceRoleGroupsDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	ID, err := getIDFromString(d.Id())
	if err != nil {
		return err
	}

	params := role.NewSetRoleGroupsParams()
	params.RoleID = ID
	params.Body = []int64{}

	_, err = client.Role.SetRoleGroups(params)
	if err != nil {
		return err
	}

	return nil
}

func resourceRoleGroupsExists(d *schema.ResourceData, m interface{}) (b bool, e error) {
	// Exists - This is called to verify a resource still exists. It is called prior to Read,
	// and lowers the burden of Read to be able to assume the resource exists.
	client := m.(*apiclient.Looker)

	ID, err := getIDFromString(d.Id())
	if err != nil {
		return false, err
	}

	params := role.NewRoleGroupsParams()
	params.RoleID = ID

	_, err = client.Role.RoleGroups(params)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func resourceRoleGroupsImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceRoleGroupsRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
