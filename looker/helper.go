package looker

import (
	"encoding/json"
	"strconv"

	apiclient "github.com/Foxtel-DnA/looker-go-sdk/client"
	"github.com/Foxtel-DnA/looker-go-sdk/client/role"
	"github.com/Foxtel-DnA/looker-go-sdk/client/session"
	"github.com/Foxtel-DnA/looker-go-sdk/models"
	"github.com/hashicorp/terraform/helper/schema"
)

func updateSession(client *apiclient.Looker, mode string) error {
	params := session.NewUpdateSessionParams()
	params.Body = &models.APISession{}
	params.Body.WorkspaceID = mode

	_, err := client.Session.UpdateSession(params)
	if err != nil {
		return err
	}

	return nil
}

func getStringArray(d *schema.ResourceData, key string) []string {
	scope := []string{}
	for _, s := range d.Get(key).([]interface{}) {
		scope = append(scope, s.(string))
	}
	return scope
}

func getRoleIds(roleNames []string, client *apiclient.Looker) ([]int64, error) {
	rolesOK, err := client.Role.AllRoles(role.NewAllRolesParams())
	if err != nil {
		return nil, err
	}

	roleIds := []int64{}
	for _, roleName := range roleNames {
		for _, role := range rolesOK.Payload {
			if role.Name == roleName {
				roleIds = append(roleIds, role.ID)
			}
		}
	}

	return roleIds, nil
}

func getRoleNames(roleIDs []int64, client *apiclient.Looker) ([]string, error) {
	rolesOK, err := client.Role.AllRoles(role.NewAllRolesParams())
	if err != nil {
		return nil, err
	}

	roleNames := []string{}
	for _, roleID := range roleIDs {
		for _, role := range rolesOK.Payload {
			if role.ID == roleID {
				roleNames = append(roleNames, role.Name)
			}
		}
	}

	return roleNames, nil
}

func getIDFromString(id string) (int64, error) {
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0, err
	}

	return i, nil
}

func getStringFromID(i int64) string {
	return strconv.FormatInt(i, 10)
}

func getJSONString(s interface{}) (string, error) {
	bytes, err := json.Marshal(s)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
