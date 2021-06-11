package looker

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func resourceUserAttribute() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserAttributeCreate,
		Read:   resourceUserAttributeRead,
		Update: resourceUserAttributeUpdate,
		Delete: resourceUserAttributeDelete,
		Importer: &schema.ResourceImporter{
			State: resourceUserAttributeImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"label": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceUserAttributeCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerSDK)
	userAttributeName := d.Get("name").(string)
	userAttributeLabel := d.Get("label").(string)
	userAttributeType := d.Get("type").(string)

	writeUserAttribute := apiclient.WriteUserAttribute{
		Name:  &userAttributeName,
		Label: &userAttributeLabel,
		Type:  &userAttributeType,
	}

	userAttribute, err := client.CreateUserAttribute(writeUserAttribute, "", nil)
	if err != nil {
		return err
	}

	userAttributeID := *userAttribute.Id
	d.SetId(strconv.Itoa(int(userAttributeID)))

	return resourceUserAttributeRead(d, m)
}

func resourceUserAttributeRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerSDK)

	userAttributeID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return err
	}

	userAttribute, err := client.UserAttribute(userAttributeID, "", nil)
	if err != nil {
		return err
	}

	if err = d.Set("name", userAttribute.Name); err != nil {
		return err
	}
	if err = d.Set("type", userAttribute.Type); err != nil {
		return err
	}
	if err = d.Set("label", userAttribute.Label); err != nil {
		return err
	}

	return nil
}

func resourceUserAttributeUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerSDK)

	userAttributeID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return err
	}

	userAttributeName := d.Get("name").(string)
	userAttributeType := d.Get("type").(string)
	userAttributeLabel := d.Get("type").(string)

	writeUserAttribute := apiclient.WriteUserAttribute{
		Name:  &userAttributeName,
		Label: &userAttributeLabel,
		Type:  &userAttributeType,
	}

	_, err = client.UpdateUserAttribute(userAttributeID, writeUserAttribute, "", nil)
	if err != nil {
		return err
	}

	return resourceUserAttributeRead(d, m)
}

func resourceUserAttributeDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerSDK)

	userAttributeID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return err
	}

	_, err = client.DeleteUserAttribute(userAttributeID, nil)
	if err != nil {
		return err
	}

	return nil
}

func resourceUserAttributeImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceUserAttributeRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
