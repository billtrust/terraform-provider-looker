package looker

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_dsRoleUsers(t *testing.T) {
	name1 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: dsRoleUserConfig(name1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("looker_role.role_users_test", "name", name1),
				),
			},
		},
	})
}

func dsRoleUserConfig(name string) string {
	return fmt.Sprintf(`
	resource "looker_model_set" "role_users_test" {
		name = "%s"
		models = ["test"]
	}
	resource "looker_permission_set" "role_users_test" {
		name = "%s"
		permissions = ["test"]
	}
	resource "looker_role" "role_users_test" {
		name = "%s"
		permission_set_id = looker_permission_set.role_users_test.id
		model_set_id = looker_model_set.role_users_test.id
	}
	data "looker_role_users" "role_users_test" {
		role_id = looker_role.role_users_test.id
	}
	`, name, name, name)
}
