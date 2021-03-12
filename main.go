package main

import (
	"github.com/Foxtel-DnA/terraform-provider-looker/looker"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: looker.Provider,
	})
}
