package thalassa

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRegions() *schema.Resource {
	return &schema.Resource{
		Description: "Get a list of regions",
		ReadContext: dataSourceRegionsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organisation": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The organisation to get the regions for. If not provided, the current organisation will be used.",
			},
			"regions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The identity of the region.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the region.",
						},
						"slug": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The slug of the region.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the region.",
						},
						"labels": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "The labels of the region.",
						},
						"annotations": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "The annotations of the region.",
						},
					},
				},
			},
		},
	}
}

func dataSourceRegionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := getProvider(m)

	regions, err := provider.Client.IaaS().ListRegions(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("regions")

	regionsList := []map[string]interface{}{}

	for _, region := range regions {
		regionsList = append(regionsList, map[string]interface{}{
			"id":          region.Identity,
			"name":        region.Name,
			"slug":        region.Slug,
			"description": region.Description,
			"labels":      region.Labels,
			"annotations": region.Annotations,
		})
	}

	d.Set("regions", regionsList)

	return nil
}
