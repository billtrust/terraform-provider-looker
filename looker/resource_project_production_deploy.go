package looker

import (
	"strings"

	"github.com/billtrust/looker-go-sdk/client/project"

	apiclient "github.com/billtrust/looker-go-sdk/client"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceProjectProductionDeploy() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectProductionDeployCreate,
		Read:   resourceProjectProductionDeployRead,
		Delete: resourceProjectProductionDeployDelete,
		Update: resourceProjectProductionDeployUpdate,
		Exists: resourceProjectProductionDeployExists,
		Importer: &schema.ResourceImporter{
			State: resourceProjectProductionDeployImport,
		},

		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func setProjectProductionDeploy(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerAPI30Reference)

	err := updateSession(client, "dev")
	if err != nil {
		return err
	}

	projectID := d.Get("project_id").(string)
	params := project.NewDeployToProductionParams()
	params.ProjectID = projectID

	_, _, err = client.Project.DeployToProduction(params)
	if err != nil {
		return err
	}

	return nil
}

func resourceProjectProductionDeployCreate(d *schema.ResourceData, m interface{}) error {
	err := setProjectProductionDeploy(d, m)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			d.SetId("")
			return nil
		}

		return err
	}

	d.SetId(d.Get("project_id").(string))

	return resourceProjectGitDetailsRead(d, m)
}

func resourceProjectProductionDeployRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerAPI30Reference)

	err := updateSession(client, "dev")
	if err != nil {
		return err
	}

	projectID := d.Id()

	params := project.NewProjectParams()
	params.ProjectID = projectID

	result, err := client.Project.Project(params)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("project_id", result.Payload.ID)

	return nil
}

func resourceProjectProductionDeployUpdate(d *schema.ResourceData, m interface{}) error {
	err := setProjectProductionDeploy(d, m)
	if err != nil {
		return err
	}

	return resourceProjectGitDetailsRead(d, m)
}

func resourceProjectProductionDeployDelete(d *schema.ResourceData, m interface{}) error {
	// TODO: Deleting this resource should set the git fields back to blank values. not implementing this yet since leaving the values does not have any negative effect
	return nil
}

func resourceProjectProductionDeployExists(d *schema.ResourceData, m interface{}) (b bool, e error) {
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

func resourceProjectProductionDeployImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceProjectProductionDeployRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
