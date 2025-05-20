package thalassa

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	iaas "github.com/thalassa-cloud/client-go/iaas"
)

func dataSourceMachineType() *schema.Resource {
	return &schema.Resource{
		Description: "Get an machine type",
		ReadContext: dataSourceMachineTypeRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"slug": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cpu_cores": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ram_mb": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceMachineTypeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := getProvider(m)
	slug := d.Get("slug").(string)

	machineTypes, err := provider.Client.IaaS().ListMachineTypes(ctx, &iaas.ListMachineTypesRequest{})
	if err != nil {
		return diag.FromErr(err)
	}

	for _, machineType := range machineTypes {
		if slug != "" && machineType.Slug == slug {
			d.SetId(machineType.Identity)
			d.Set("id", machineType.Identity)
			d.Set("name", machineType.Name)
			d.Set("slug", machineType.Slug)
			d.Set("description", machineType.Description)
			d.Set("cpu_cores", machineType.Vcpus)
			d.Set("ram_mb", machineType.RamMb)
			return diag.Diagnostics{}
		}
	}
	return diag.FromErr(fmt.Errorf("not found"))
}
