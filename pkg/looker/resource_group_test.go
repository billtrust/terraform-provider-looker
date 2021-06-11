package looker

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_Group(t *testing.T) {
	name1 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	name2 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: groupConfig(name1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("looker_group.test", "name", name1),
				),
			},
			{
				Config: groupConfig(name2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("looker_group.test", "name", name2),
				),
			},
			{
				ResourceName:      "looker_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func groupConfig(name string) string {
	return fmt.Sprintf(`
	resource "looker_group" "test" {
		name = "%s"
	}
	`, name)
}
