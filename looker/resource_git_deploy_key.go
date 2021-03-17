package looker

import (
	"strings"

	"github.com/Foxtel-DnA/looker-go-sdk/client/project"

	apiclient "github.com/Foxtel-DnA/looker-go-sdk/client"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGitDeployKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitDeployKeyCreate,
		Read:   resourceGitDeployKeyRead,
		Delete: resourceGitDeployKeyDelete,
		Exists: resourceGitDeployKeyExists,
		Importer: &schema.ResourceImporter{
			State: resourceGitDeployKeyImport,
		},

		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ssh_deploy_key": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceGitDeployKeyCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	err := updateSession(client, "dev")
	if err != nil {
		return err
	}

	projectID := d.Get("project_id").(string)
	params := project.NewCreateGitDeployKeyParams()
	params.ProjectID = projectID

	_, err = client.Project.CreateGitDeployKey(params)
	if err != nil {
		return err
	}

	d.SetId(projectID)

	return resourceGitDeployKeyRead(d, m)
}

func resourceGitDeployKeyRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	err := updateSession(client, "dev")
	if err != nil {
		return err
	}

	projectID := d.Id()

	params := project.NewGitDeployKeyParams()
	params.ProjectID = projectID

	result, err := client.Project.GitDeployKey(params)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("project_id", projectID)

	// the payload is a string with 3 values separated by spaces.  The first index contains "ssh-rsa", the second index includes the key, the third index contains the project id
	// the project id doesn't appear to be part of the key though since when adding it to github, it is ignored it seems
	sshKey := strings.Fields(result.Payload)
	d.Set("ssh_deploy_key", sshKey[0]+" "+sshKey[1])

	return nil
}

func resourceGitDeployKeyDelete(d *schema.ResourceData, m interface{}) error {
	// TODO There is no way to delete a git deploy key, possibly put this into the project resource (but there is no way to delete project either)
	return nil
}

func resourceGitDeployKeyExists(d *schema.ResourceData, m interface{}) (b bool, e error) {
	// Exists - This is called to verify a resource still exists. It is called prior to Read,
	// and lowers the burden of Read to be able to assume the resource exists.
	client := m.(*apiclient.Looker)

	// TODO Not sure if we should always set session to "dev" instead of "production" when checking if it exists? will dev always show all dev+prod projects?
	err := updateSession(client, "dev")
	if err != nil {
		return false, err
	}

	params := project.NewGitDeployKeyParams()
	params.ProjectID = d.Id()

	_, err = client.Project.GitDeployKey(params)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func resourceGitDeployKeyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceGitDeployKeyRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
