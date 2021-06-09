package looker

import (
	"strconv"

	apiclient "github.com/billtrust/looker-go-sdk/client"
	"github.com/billtrust/looker-go-sdk/client/role"
	"github.com/billtrust/looker-go-sdk/client/session"
	"github.com/billtrust/looker-go-sdk/models"
)

func updateSession(client *apiclient.LookerAPI30Reference, mode string) error {
	params := session.NewUpdateSessionParams()
	params.Body = &models.APISession{}
	params.Body.WorkspaceID = mode

	_, err := client.Session.UpdateSession(params)
	if err != nil {
		return err
	}

	return nil
}

func getRoleIds(roleNames []string, client *apiclient.LookerAPI30Reference) ([]int64, error) {
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
