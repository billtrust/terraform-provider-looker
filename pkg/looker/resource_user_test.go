package looker

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_User(t *testing.T) {
	name1 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	name2 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	email := "test2@example.com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: userConfig(name1, name1, email),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("looker_user.test", "first_name", name1),
					resource.TestCheckResourceAttr("looker_user.test", "last_name", name1),
					resource.TestCheckResourceAttr("looker_user.test", "email", email),
				),
			},
			{
				Config: userConfig(name2, name2, email),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("looker_user.test", "first_name", name2),
					resource.TestCheckResourceAttr("looker_user.test", "last_name", name2),
					resource.TestCheckResourceAttr("looker_user.test", "email", email),
				),
			},
			{
				ResourceName:      "looker_user.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func userConfig(firstName, lastName, email string) string {
	return fmt.Sprintf(`
	resource "looker_user" "test" {
		first_name = "%s"
		last_name = "%s"
		email = "%s"
	}
	`, firstName, lastName, email)
}
