package iaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	iaas "github.com/thalassa-cloud/client-go/iaas"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func DataSourceMachineType() *schema.Resource {
	return &schema.Resource{
		Description: "Get an machine type",
		ReadContext: dataSourceMachineTypeRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the machine type",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the machine type",
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the Machine Type. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"slug": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The slug of the machine type to look up",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A description of the machine type and its specifications",
			},
			"cpu_cores": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of CPU cores available in this machine type",
			},
			"ram_mb": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The amount of RAM in megabytes available in this machine type",
			},
		},
	}
}

func dataSourceMachineTypeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}
	slug := d.Get("slug").(string)

	machineTypes, err := client.IaaS().ListMachineTypes(ctx, &iaas.ListMachineTypesRequest{})
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
