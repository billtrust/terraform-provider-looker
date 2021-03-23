package looker

import (
	"log"

	apiclient "github.com/Foxtel-DnA/looker-go-sdk/client"
	"github.com/Foxtel-DnA/looker-go-sdk/client/api_auth"

	"github.com/go-openapi/strfmt"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/hashicorp/terraform/terraform"

	"github.com/hashicorp/terraform/helper/schema"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("LOOKER_API_CLIENT_ID", nil),
				Description: "Client ID to authenticate with Looker",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("LOOKER_API_CLIENT_SECRET", nil),
				Description: "Client Secret to authenticate with Looker",
			},
			"base_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("LOOKER_API_BASE_URL", nil),
				Description: "Looker API Base URL",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"looker_user":                    resourceUser(),
			"looker_user_email":              resourceUserEmail(),
			"looker_user_roles":              resourceUserRoles(),
			"looker_user_api_key":            resourceUserAPIKey(),
			"looker_permission_set":          resourcePermissionSet(),
			"looker_model_set":               resourceModelSet(),
			"looker_group":                   resourceGroup(),
			"looker_role":                    resourceRole(),
			"looker_role_groups":             resourceRoleGroups(),
			"looker_main_space":              resourceMainSpace(),
			"looker_child_space":             resourceChildSpace(),
			"looker_content_metadata_access": resourceContentMetadataAccess(),
			"looker_connection":              resourceConnection(),
			"looker_project":                 resourceProject(),
			"looker_git_deploy_key":          resourceGitDeployKey(),
			"looker_project_git_details":     resourceProjectGitDetails(),
			"looker_user_attribute":          resourceUserAttribute(),
			"looker_theme":                   resourceTheme(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	transport := httptransport.New(d.Get("base_url").(string), "/api/3.1/", nil)
	client := apiclient.New(transport, strfmt.Default)

	clientID := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)

	pd := api_auth.NewLoginParams()
	pd.ClientID = &clientID
	pd.ClientSecret = &clientSecret

	resp, err := client.APIAuth.Login(pd)

	if err != nil {
		return nil, err
	}

	token := resp.Payload.AccessToken
	log.Println("[INFO] token " + token)

	authInfoWriter := httptransport.APIKeyAuth("Authorization", "header", "token "+token)
	transport.DefaultAuthentication = authInfoWriter

	authClient := apiclient.New(transport, strfmt.Default)

	return authClient, nil
}
