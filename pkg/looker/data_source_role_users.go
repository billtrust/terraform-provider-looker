package looker

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

var dsRoleUsersSchema = map[string]*schema.Schema{
	"role_id": {
		Type:     schema.TypeString,
		Required: true,
	},
	"users": {
		Type:     schema.TypeList,
		Elem:     &schema.Schema{Type: schema.TypeString},
		Computed: true,
	},
}

func dsRoleUsers() *schema.Resource {
	return &schema.Resource{
		Read:   dsReadRoleUsers,
		Schema: dsRoleUsersSchema,
	}
}

func dsReadRoleUsers(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerSDK)

	roleID := d.Get("role_id").(string)

	request := apiclient.RequestRoleUsers{RoleId: roleID}

	users, err := client.RoleUsers(request, nil)
	if err != nil {
		return err
	}

	var userEmails []string
	for _, user := range users {
		userEmails = append(userEmails, *user.Email)
	}

	d.SetId(roleID)
	return d.Set("users", userEmails)
}
