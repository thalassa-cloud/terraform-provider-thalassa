package iaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	iaas "github.com/thalassa-cloud/client-go/iaas"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func DataSourceRegion() *schema.Resource {
	return &schema.Resource{
		Description: "Get an region",
		ReadContext: dataSourceRegionRead,
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
				Computed:    true,
				Description: "Reference to the Organisation of the Region. If not provided, the organisation of the (Terraform) provider will be used.",
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

func dataSourceRegionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}
	slug := d.Get("slug").(string)

	regions, err := client.IaaS().ListRegions(ctx, &iaas.ListRegionsRequest{})
	if err != nil {
		return diag.FromErr(err)
	}

	for _, region := range regions {
		if slug != "" && region.Slug == slug {
			d.SetId(region.Identity)
			d.Set("id", region.Identity)
			d.Set("name", region.Name)
			d.Set("slug", region.Slug)
			d.Set("description", region.Description)

			// Set labels and annotations directly
			if err := d.Set("labels", region.Labels); err != nil {
				return diag.FromErr(fmt.Errorf("error setting labels: %s", err))
			}

			if err := d.Set("annotations", region.Annotations); err != nil {
				return diag.FromErr(fmt.Errorf("error setting annotations: %s", err))
			}

			return diag.Diagnostics{}
		}
	}
	return diag.FromErr(fmt.Errorf("not found"))
}
