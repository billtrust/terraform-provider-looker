package looker

import (
	"strings"

	apiclient "github.com/Foxtel-DnA/looker-go-sdk/client"
	"github.com/Foxtel-DnA/looker-go-sdk/client/user"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceUserRoles() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserRolesCreate,
		Read:   resourceUserRolesRead,
		Update: resourceUserRolesUpdate,
		Delete: resourceUserRolesDelete,
		Exists: resourceUserRolesExists,
		Importer: &schema.ResourceImporter{
			State: resourceUserRolesImport,
		},

		Schema: map[string]*schema.Schema{
			"user_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"role_names": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceUserRolesCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	sUserID := d.Get("user_id").(string)

	iUserID, err := getIDFromString(sUserID)
	if err != nil {
		return err
	}

	var roleNames []string
	for _, roleName := range d.Get("role_names").(*schema.Set).List() {
		roleNames = append(roleNames, roleName.(string))
	}

	// TODO: if role name does not exist, what should it do? through an error? try to create the role?
	roleIds, err := getRoleIds(roleNames, client)

	if err != nil {
		return err
	}

	params := user.NewSetUserRolesParams()
	params.UserID = iUserID
	params.Body = roleIds

	_, err = client.User.SetUserRoles(params)
	if err != nil {
		return err
	}

	d.SetId(sUserID)

	return resourceUserRolesRead(d, m)
}

func resourceUserRolesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	userID, err := getIDFromString(d.Id())
	if err != nil {
		return err
	}

	params := user.NewUserRolesParams()
	params.UserID = userID

	rolesResult, err := client.User.UserRoles(params)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			d.SetId("")
			return nil
		}
		return err
	}

	roleNames := []string{}
	for _, role := range rolesResult.Payload {
		roleNames = append(roleNames, role.Name)
	}

	d.Set("user_id", userID)
	d.Set("role_names", roleNames)

	return nil
}

func resourceUserRolesUpdate(d *schema.ResourceData, m interface{}) error {
	// TODO: There is no functional difference between "Creating" and "updating" which roles a user has.  Only settings the roles currently assigned to the user. Is this the correct implemenation in thise case?
	return resourceUserRolesCreate(d, m)
}

func resourceUserRolesDelete(d *schema.ResourceData, m interface{}) error {
	// TODO: Delete really just removes all the roles from the user.  Is this the correct way to implement delete in this case?
	client := m.(*apiclient.Looker)

	userID, err := getIDFromString(d.Id())
	if err != nil {
		return err
	}

	params := user.NewSetUserRolesParams()
	params.UserID = userID
	params.Body = []int64{}

	_, err = client.User.SetUserRoles(params)
	if err != nil {
		return err
	}

	return nil
}

func resourceUserRolesExists(d *schema.ResourceData, m interface{}) (b bool, e error) {
	// Exists - This is called to verify a resource still exists. It is called prior to Read,
	// and lowers the burden of Read to be able to assume the resource exists.
	client := m.(*apiclient.Looker)

	userID, err := getIDFromString(d.Id())
	if err != nil {
		return false, err
	}

	// TODO: as long as the user exists, we will say the "roles" exist. Not sure if this is correct though?
	params := user.NewUserParams()
	params.UserID = userID

	_, err = client.User.User(params)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func resourceUserRolesImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceUserRolesRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
