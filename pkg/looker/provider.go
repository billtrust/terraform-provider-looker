package looker

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/looker-open-source/sdk-codegen/go/rtl"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func Provider() *schema.Provider {
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
			"looker_user":           resourceUser(),
			"looker_user_roles":     resourceUserRoles(),
			"looker_permission_set": resourcePermissionSet(),
			"looker_model_set":      resourceModelSet(),
			"looker_group":          resourceGroup(),
			"looker_role":           resourceRole(),
			"looker_role_groups":    resourceRoleGroups(),
			"looker_user_attribute": resourceUserAttribute(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	baseUrl := d.Get("base_url").(string)
	clientID := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)

	apiSettings := rtl.ApiSettings{
		BaseUrl:      baseUrl,
		ClientId:     clientID,
		ClientSecret: clientSecret,
		ApiVersion:   "4.0",
	}
	authSession := rtl.NewAuthSession(apiSettings)
	client := apiclient.NewLookerSDK(authSession)

	return client, nil
}
