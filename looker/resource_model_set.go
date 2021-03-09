package looker

import (
	"strings"

	"github.com/billtrust/looker-go-sdk/client/role"
	"github.com/billtrust/looker-go-sdk/models"

	apiclient "github.com/billtrust/looker-go-sdk/client"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceModelSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceModelSetCreate,
		Read:   resourceModelSetRead,
		Update: resourceModelSetUpdate,
		Delete: resourceModelSetDelete,
		Exists: resourceModelSetExists,
		Importer: &schema.ResourceImporter{
			State: resourceModelSetImport,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"models": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceModelSetCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	var modelNames []string
	for _, modelName := range d.Get("models").(*schema.Set).List() {
		modelNames = append(modelNames, modelName.(string))
	}

	params := role.NewCreateModelSetParams()
	params.Body = &models.ModelSet{}
	params.Body.Name = d.Get("name").(string)
	params.Body.Models = modelNames

	result, err := client.Role.CreateModelSet(params)
	if err != nil {
		return err
	}

	d.SetId(getStringFromID(result.Payload.ID))

	return resourceModelSetRead(d, m)
}

func resourceModelSetRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	ID, err := getIDFromString(d.Id())
	if err != nil {
		return err
	}

	params := role.NewModelSetParams()
	params.ModelSetID = ID

	result, err := client.Role.ModelSet(params)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", result.Payload.Name)
	d.Set("models", result.Payload.Models)

	return nil
}

func resourceModelSetUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	ID, err := getIDFromString(d.Id())
	if err != nil {
		return err
	}

	var modelNames []string
	for _, modelName := range d.Get("models").(*schema.Set).List() {
		modelNames = append(modelNames, modelName.(string))
	}

	params := role.NewUpdateModelSetParams()
	params.ModelSetID = ID
	params.Body = &models.ModelSet{}
	params.Body.Name = d.Get("name").(string)
	params.Body.Models = modelNames

	_, err = client.Role.UpdateModelSet(params)
	if err != nil {
		return err
	}

	return resourceModelSetRead(d, m)
}

func resourceModelSetDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	ID, err := getIDFromString(d.Id())
	if err != nil {
		return err
	}

	params := role.NewDeleteModelSetParams()
	params.ModelSetID = ID

	_, err = client.Role.DeleteModelSet(params)
	if err != nil {
		return err
	}

	return nil
}

func resourceModelSetExists(d *schema.ResourceData, m interface{}) (b bool, e error) {
	// Exists - This is called to verify a resource still exists. It is called prior to Read,
	// and lowers the burden of Read to be able to assume the resource exists.
	client := m.(*apiclient.Looker)

	ID, err := getIDFromString(d.Id())
	if err != nil {
		return false, err
	}

	params := role.NewModelSetParams()
	params.ModelSetID = ID

	_, err = client.Role.ModelSet(params)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func resourceModelSetImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceModelSetRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
