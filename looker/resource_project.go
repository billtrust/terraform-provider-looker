package looker

import (
	"fmt"
	"strings"

	"github.com/billtrust/looker-go-sdk/client/project"

	apiclient "github.com/billtrust/looker-go-sdk/client"
	"github.com/billtrust/looker-go-sdk/models"
	"github.com/hashicorp/terraform/helper/schema"
)

func getProject(projectID string, client *apiclient.Looker) (*project.ProjectOK, error) {
	params := project.NewProjectParams()
	params.ProjectID = projectID

	result, err := client.Project.Project(params)
	return result, err
}

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
	client := m.(*apiclient.Looker)

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
	client := m.(*apiclient.Looker)

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
	client := m.(*apiclient.Looker)

	err := updateSession(client, "dev")
	if err != nil {
		return err
	}

	name := d.Get("name").(string)

	params := project.NewUpdateProjectParams()
	params.ProjectID = d.Id()
	params.Body = &models.Project{}
	params.Body.Name = name

	_, err = client.Project.UpdateProject(params)
	if err != nil {
		// looker gives "An error has occured" 500 error even though the name correctly is updated
		if !strings.Contains(err.Error(), "An error has occurred.") {
			return err
		}

		// the project might be updated correctly, check by getting the project with the new name
		_, readError := getProject(name, client)
		if readError != nil {
			return err // return original error for the update
		}

		// looks like update of the name was succesful because we can query for it now.  Update the ID to be the new name and read it back
		d.SetId(name)

		return resourceProjectRead(d, m)
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
	client := m.(*apiclient.Looker)

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
