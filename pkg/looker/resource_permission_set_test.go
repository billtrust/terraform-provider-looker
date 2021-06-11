package looker

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_PermissionSet(t *testing.T) {
	name1 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	name2 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: permissionSetConfig(name1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("looker_permission_set.test", "name", name1),
					resource.TestCheckResourceAttr("looker_permission_set.test", "permissions.#", "1"),
				),
			},
			{
				Config: permissionSetConfig(name2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("looker_permission_set.test", "name", name2),
					resource.TestCheckResourceAttr("looker_permission_set.test", "permissions.#", "1"),
				),
			},
			{
				ResourceName:      "looker_permission_set.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func permissionSetConfig(name string) string {
	return fmt.Sprintf(`
	resource "looker_permission_set" "test" {
		name = "%s"
		permissions = ["test"]
	}
	`, name)
}
