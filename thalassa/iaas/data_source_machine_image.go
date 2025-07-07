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
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity of the machine image",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the machine image",
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the Machine Image. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"slug": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Slug of the machine image",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the machine image",
			},
			"labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Labels of the machine image",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Annotations of the machine image",
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

	name := d.Get("name").(string)

	for _, machineImage := range machineImages {
		if slug != "" && machineImage.Slug == slug || name != "" && machineImage.Name == name {
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
