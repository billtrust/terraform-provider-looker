package looker

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_ModelSet(t *testing.T) {
	name1 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	name2 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: modelSetConfig(name1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("looker_model_set.test", "name", name1),
					resource.TestCheckResourceAttr("looker_model_set.test", "models.#", "1"),
				),
			},
			{
				Config: modelSetConfig(name2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("looker_model_set.test", "name", name2),
					resource.TestCheckResourceAttr("looker_model_set.test", "models.#", "1"),
				),
			},
			{
				ResourceName:      "looker_model_set.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func modelSetConfig(name string) string {
	return fmt.Sprintf(`
	resource "looker_model_set" "test" {
		name = "%s"
		models = ["test"]
	}
	`, name)
}
