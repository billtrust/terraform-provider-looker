package looker

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func resourceModelSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceModelSetCreate,
		Read:   resourceModelSetRead,
		Update: resourceModelSetUpdate,
		Delete: resourceModelSetDelete,
		Importer: &schema.ResourceImporter{
			State: resourceModelSetImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"models": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceModelSetCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerSDK)

	modelSetName := d.Get("name").(string)

	var modelNames []string
	for _, modelName := range d.Get("models").(*schema.Set).List() {
		modelNames = append(modelNames, modelName.(string))
	}

	writeModelSet := apiclient.WriteModelSet{
		Name:   &modelSetName,
		Models: &modelNames,
	}

	modelSet, err := client.CreateModelSet(writeModelSet, nil)
	if err != nil {
		return err
	}

	modelSetID := *modelSet.Id
	d.SetId(modelSetID)

	return resourceModelSetRead(d, m)
}

func resourceModelSetRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerSDK)

	modelSetID := d.Id()

	modelSet, err := client.ModelSet(modelSetID, "", nil)
	if err != nil {
		return err
	}

	if err = d.Set("name", modelSet.Name); err != nil {
		return err
	}
	if err = d.Set("models", modelSet.Models); err != nil {
		return err
	}

	return nil
}

func resourceModelSetUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerSDK)

	modelSetID := d.Id()
	modelSetName := d.Get("name").(string)
	var modelNames []string
	for _, modelName := range d.Get("models").(*schema.Set).List() {
		modelNames = append(modelNames, modelName.(string))
	}
	writeModelSet := apiclient.WriteModelSet{
		Name:   &modelSetName,
		Models: &modelNames,
	}
	_, err := client.UpdateModelSet(modelSetID, writeModelSet, nil)
	if err != nil {
		return err
	}

	return resourceModelSetRead(d, m)
}

func resourceModelSetDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerSDK)

	modelSetID := d.Id()

	_, err := client.DeleteModelSet(modelSetID, nil)
	if err != nil {
		return err
	}

	return nil
}
func resourceModelSetImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceModelSetRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
