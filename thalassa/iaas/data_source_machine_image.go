package iaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	iaas "github.com/thalassa-cloud/client-go/iaas"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func DataSourceMachineImage() *schema.Resource {
	return &schema.Resource{
		Description: "Get an machine image",
		ReadContext: dataSourceMachineImageRead,
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
			"labels": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"annotations": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func dataSourceMachineImageRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}
	slug := d.Get("slug").(string)

	machineImages, err := client.IaaS().ListMachineImages(ctx, &iaas.ListMachineImagesRequest{})
	if err != nil {
		return diag.FromErr(err)
	}

	for _, machineImage := range machineImages {
		if slug != "" && machineImage.Slug == slug {
			d.SetId(machineImage.Identity)
			d.Set("id", machineImage.Identity)
			d.Set("name", machineImage.Name)
			d.Set("slug", machineImage.Slug)
			d.Set("description", machineImage.Description)

			// Set labels and annotations directly
			if err := d.Set("labels", machineImage.Labels); err != nil {
				return diag.FromErr(fmt.Errorf("error setting labels: %s", err))
			}
			return diag.Diagnostics{}
		}
	}
	return diag.FromErr(fmt.Errorf("not found"))
}
