package thalassa

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceOrganisations() *schema.Resource {
	return &schema.Resource{
		Description: "Get an organisation",
		ReadContext: dataSourceOrganisationsRead,
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

func dataSourceOrganisationsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := getProvider(m)
	slug := d.Get("slug").(string)

	organisations, err := provider.Client.Me().ListMyOrganisations(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, org := range organisations {
		if slug != "" && org.Slug == slug {
			d.SetId(org.Identity)
			d.Set("id", org.Identity)
			d.Set("name", org.Name)
			d.Set("slug", org.Slug)
			d.Set("description", org.Description)

			// Set labels and annotations directly
			if err := d.Set("labels", org.Labels); err != nil {
				return diag.FromErr(fmt.Errorf("error setting labels: %s", err))
			}

			if err := d.Set("annotations", org.Annotations); err != nil {
				return diag.FromErr(fmt.Errorf("error setting annotations: %s", err))
			}

			return diag.Diagnostics{}
		}
	}
	return diag.FromErr(fmt.Errorf("not found"))
}
