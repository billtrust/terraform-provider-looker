package looker

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupCreate,
		Read:   resourceGroupRead,
		Update: resourceGroupUpdate,
		Delete: resourceGroupDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGroupImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceGroupCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerSDK)
	groupName := d.Get("name").(string)

	writeGroup := apiclient.WriteGroup{
		Name: &groupName,
	}

	group, err := client.CreateGroup(writeGroup, "", nil)
	if err != nil {
		return err
	}

	groupID := *group.Id
	d.SetId(groupID)

	return resourceGroupRead(d, m)
}

func resourceGroupRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerSDK)

	groupID := d.Id()

	group, err := client.Group(groupID, "", nil)
	if err != nil {
		return err
	}

	if err = d.Set("name", group.Name); err != nil {
		return err
	}

	return nil
}

func resourceGroupUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerSDK)

	groupID := d.Id()

	groupName := d.Get("name").(string)
	writeGroup := apiclient.WriteGroup{
		Name: &groupName,
	}
	_, err := client.UpdateGroup(groupID, writeGroup, "", nil)
	if err != nil {
		return err
	}

	return resourceGroupRead(d, m)
}

func resourceGroupDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerSDK)

	groupID := d.Id()

	_, err := client.DeleteGroup(groupID, nil)
	if err != nil {
		return err
	}

	return nil
}

func resourceGroupImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceGroupRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
