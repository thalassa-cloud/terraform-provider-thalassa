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
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the Machine Type. If not provided, the organisation of the (Terraform) provider will be used.",
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
