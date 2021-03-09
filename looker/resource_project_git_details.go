package looker

import (
	"strings"

	"github.com/billtrust/looker-go-sdk/client/project"
	"github.com/billtrust/looker-go-sdk/models"
	"github.com/go-openapi/strfmt"

	apiclient "github.com/billtrust/looker-go-sdk/client"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceProjectGitDetails() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectGitDetailsCreate,
		Read:   resourceProjectGitDetailsRead,
		Delete: resourceProjectGitDetailsDelete,
		Update: resourceProjectGitDetailsDelete,
		Exists: resourceProjectGitDetailsExists,
		Importer: &schema.ResourceImporter{
			State: resourceProjectGitDetailsImport,
		},

		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"git_remote_url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func setProjectGitDetails(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	err := updateSession(client, "dev")
	if err != nil {
		return err
	}

	gitRemoteURL := strfmt.URI(d.Get("git_remote_url").(string))

	projectID := d.Get("project_id").(string)
	params := project.NewUpdateProjectParams()
	params.ProjectID = projectID
	params.Body = &models.Project{}
	params.Body.GitRemoteURL = gitRemoteURL.String()

	_, err = client.Project.UpdateProject(params)
	if err != nil {
		return err
	}

	return nil
}

func resourceProjectGitDetailsCreate(d *schema.ResourceData, m interface{}) error {
	err := setProjectGitDetails(d, m)
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

func resourceProjectGitDetailsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

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
	d.Set("git_remote_url", result.Payload.GitRemoteURL)

	return nil
}

func resourceProjectGitDetailsUpdate(d *schema.ResourceData, m interface{}) error {
	err := setProjectGitDetails(d, m)
	if err != nil {
		return err
	}

	return resourceProjectGitDetailsRead(d, m)
}

func resourceProjectGitDetailsDelete(d *schema.ResourceData, m interface{}) error {
	// TODO: Deleting this resource should set the git fields back to blank values. not implementing this yet since leaving the values does not have any negative effect
	return nil
}

func resourceProjectGitDetailsExists(d *schema.ResourceData, m interface{}) (b bool, e error) {
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

func resourceProjectGitDetailsImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceProjectGitDetailsRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
