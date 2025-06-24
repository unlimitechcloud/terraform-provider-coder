package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/unlimitechcloud/terraform-provider-coder/coder"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: coder.Provider,
	})
}
