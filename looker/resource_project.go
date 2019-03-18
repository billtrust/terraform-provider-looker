package looker

import (
	"fmt"
	"strings"

	"github.com/bmccarthy/looker-go-sdk/client/project"

	apiclient "github.com/bmccarthy/looker-go-sdk/client"
	"github.com/bmccarthy/looker-go-sdk/models"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,
		Exists: resourceProjectExists,
		Importer: &schema.ResourceImporter{
			State: resourceProjectImport,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if strings.Contains(v, " ") {
						errs = append(errs, fmt.Errorf("%q must not contain any spaces, got: %q", key, v))
					}
					return
				},
			},
		},
	}
}

func resourceProjectCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerAPI30Reference)

	err := updateSession(client, "dev")
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	params := project.NewCreateProjectParams()
	params.Body = &models.Project{}
	params.Body.Name = name

	_, err = client.Project.CreateProject(params)
	if err != nil {
		return err
	}

	d.SetId(name)

	return resourceProjectRead(d, m)
}

func resourceProjectRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerAPI30Reference)

	params := project.NewProjectParams()
	params.ProjectID = d.Id()

	result, err := client.Project.Project(params)
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

func resourceProjectUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerAPI30Reference)

	name := d.Get("name").(string)

	params := project.NewUpdateProjectParams()
	params.ProjectID = d.Id()
	params.Body = &models.Project{}
	params.Body.Name = name

	_, err := client.Project.UpdateProject(params)
	if err != nil {
		return err
	}

	return resourceProjectRead(d, m)
}

func resourceProjectDelete(d *schema.ResourceData, m interface{}) error {
	// TODO: Looker doesn't appear to support deleting projects from the API
	return nil
}

func resourceProjectExists(d *schema.ResourceData, m interface{}) (b bool, e error) {
	// Exists - This is called to verify a resource still exists. It is called prior to Read,
	// and lowers the burden of Read to be able to assume the resource exists.
	client := m.(*apiclient.LookerAPI30Reference)

	// TODO Not sure if we should always set session to "dev" instead of "production" when checking if it exists? will dev always show all dev+prod projects?
	err := updateSession(client, "dev")
	if err != nil {
		return false, err
	}

	params := project.NewProjectParams()
	params.ProjectID = d.Id()

	_, err = client.Project.Project(params)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func resourceProjectImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceProjectRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
