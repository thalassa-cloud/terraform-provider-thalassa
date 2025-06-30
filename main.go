package main

import (
	"flag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa"
)

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "Enable debug mode")
	flag.Parse()

	plugin.Serve(&plugin.ServeOpts{
		Debug: debug,
		ProviderFunc: func() *schema.Provider {
			return thalassa.Provider()
		},
		ProviderAddr: "terraform.local/local/thalassa",
	})
}
