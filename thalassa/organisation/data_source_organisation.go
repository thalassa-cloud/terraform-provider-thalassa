package organisation

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func DataSourceOrganisations() *schema.Resource {
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

func dataSourceOrganisationsRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	provider := provider.GetProvider(m)
	slug := d.Get("slug").(string)

	organisations, err := provider.Client.Me().ListMyOrganisations(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, org := range organisations {
		if slug != "" && org.Slug == slug {
			d.SetId(org.Identity)
			_ = d.Set("id", org.Identity)
			_ = d.Set("name", org.Name)
			_ = d.Set("slug", org.Slug)
			_ = d.Set("description", org.Description)

			// Set labels and annotations directly
			_ = d.Set("labels", org.Labels)

			_ = d.Set("annotations", org.Annotations)

			return diag.Diagnostics{}
		}
	}
	return diag.FromErr(fmt.Errorf("not found"))
}
