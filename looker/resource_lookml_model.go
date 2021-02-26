package looker

import (
	"github.com/billtrust/looker-go-sdk/client/lookml_model"
	"log"
	"strings"

	apiclient "github.com/billtrust/looker-go-sdk/client"
	"github.com/billtrust/looker-go-sdk/models"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceModel() *schema.Resource {
	return &schema.Resource{
		Create: resourceModelCreate,
		Read:   resourceModelRead,
		Update: resourceModelUpdate,
		Delete: resourceModelDelete,
		Exists: resourceModelExists,
		Importer: &schema.ResourceImporter{
			State: resourceModelImport,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"allowed_db_connection_names": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"project_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceModelCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerAPI30Reference)

	params := lookml_model.NewCreateLookmlModelParams()
	params.Body = &models.LookmlModel{}
	params.Body.Name = d.Get("name").(string)
	params.Body.ProjectName = d.Get("project_name").(string)
	var connectionNames []string
	for _, modelName := range d.Get("allowed_db_connection_names").(*schema.Set).List() {
		connectionNames = append(connectionNames, modelName.(string))
	}
	params.Body.AllowedDbConnectionNames = connectionNames

	model, err := client.LookmlModel.CreateLookmlModel(params)
	if err != nil {
		log.Printf("[WARN] Error creating a model., %s", err.Error())
		return err
	}

	d.SetId(model.Payload.Name)

	return resourceModelRead(d, m)
}

func resourceModelRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerAPI30Reference)

	modelID := d.Id()

	params := lookml_model.NewLookmlModelParams()
	params.LookmlModelName = modelID

	model, err := client.LookmlModel.LookmlModel(params)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(modelID)
	d.Set("name", model.Payload.Name)
	d.Set("project_name", model.Payload.ProjectName)
	d.Set("allowed_db_connection_names", model.Payload.AllowedDbConnectionNames)

	return nil
}

func resourceModelUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerAPI30Reference)

	modelID := d.Id()

	params := lookml_model.NewUpdateLookmlModelParams()
	params.LookmlModelName = modelID
	params.Body = &models.LookmlModel{}
	params.Body.Name = d.Get("name").(string)
	params.Body.ProjectName = d.Get("project_name").(string)
	var connectionNames []string
	for _, modelName := range d.Get("allowed_db_connection_names").(*schema.Set).List() {
		connectionNames = append(connectionNames, modelName.(string))
	}
	params.Body.AllowedDbConnectionNames = connectionNames

	_, err := client.LookmlModel.UpdateLookmlModel(params)
	if err != nil {
		return err
	}

	return resourceModelRead(d, m)
}

func resourceModelDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerAPI30Reference)

	modelID := d.Id()

	params := lookml_model.NewDeleteLookmlModelParams()
	params.LookmlModelName = modelID

	_, err := client.LookmlModel.DeleteLookmlModel(params)
	if err != nil {
		return err
	}

	return nil
}

func resourceModelExists(d *schema.ResourceData, m interface{}) (b bool, e error) {
	// Exists - This is called to verify a resource still exists. It is called prior to Read,
	// and lowers the burden of Read to be able to assume the resource exists.
	client := m.(*apiclient.LookerAPI30Reference)

	modelID := d.Id()

	params := lookml_model.NewLookmlModelParams()
	params.LookmlModelName = modelID

	_, err := client.LookmlModel.LookmlModel(params)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func resourceModelImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceModelRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
