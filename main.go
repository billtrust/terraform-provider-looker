package main

import (
	"github.com/cambiahealth/terraform-provider-looker/looker"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: looker.Provider,
	})
}
