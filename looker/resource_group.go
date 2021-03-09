package looker

import (
	"strings"

	"github.com/billtrust/looker-go-sdk/client/group"

	"github.com/billtrust/looker-go-sdk/models"

	apiclient "github.com/billtrust/looker-go-sdk/client"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupCreate,
		Read:   resourceGroupRead,
		Update: resourceGroupUpdate,
		Delete: resourceGroupDelete,
		Exists: resourceGroupExists,
		Importer: &schema.ResourceImporter{
			State: resourceGroupImport,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceGroupCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	params := group.NewCreateGroupParams()
	params.Body = &models.Group{}
	params.Body.Name = d.Get("name").(string)

	result, err := client.Group.CreateGroup(params)
	if err != nil {
		return err
	}

	d.SetId(getStringFromID(result.Payload.ID))

	return resourceGroupRead(d, m)
}

func resourceGroupRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	ID, err := getIDFromString(d.Id())
	if err != nil {
		return err
	}

	params := group.NewGroupParams()
	params.GroupID = ID

	result, err := client.Group.Group(params)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", result.Payload.Name)

	return nil
}

func resourceGroupUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	ID, err := getIDFromString(d.Id())
	if err != nil {
		return err
	}

	params := group.NewUpdateGroupParams()
	params.GroupID = ID
	params.Body = &models.Group{}
	params.Body.Name = d.Get("name").(string)

	_, err = client.Group.UpdateGroup(params)
	if err != nil {
		return err
	}

	return resourceGroupRead(d, m)
}

func resourceGroupDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	ID, err := getIDFromString(d.Id())
	if err != nil {
		return err
	}

	params := group.NewDeleteGroupParams()
	params.GroupID = ID

	_, err = client.Group.DeleteGroup(params)
	if err != nil {
		return err
	}

	return nil
}

func resourceGroupExists(d *schema.ResourceData, m interface{}) (b bool, e error) {
	// Exists - This is called to verify a resource still exists. It is called prior to Read,
	// and lowers the burden of Read to be able to assume the resource exists.
	client := m.(*apiclient.Looker)

	ID, err := getIDFromString(d.Id())
	if err != nil {
		return false, err
	}

	params := group.NewGroupParams()
	params.GroupID = ID

	_, err = client.Group.Group(params)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func resourceGroupImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceGroupRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
