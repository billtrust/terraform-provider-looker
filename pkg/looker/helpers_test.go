package looker

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func providers() map[string]*schema.Provider {
	p := Provider()
	return map[string]*schema.Provider{
		"looker": p,
	}
}
