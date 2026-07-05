package kms

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func DataSourceKmsSummary() *schema.Resource {
	return &schema.Resource{
		Description: "KMS and Secrets Manager regional availability for the organisation",
		ReadContext: dataSourceKmsSummaryRead,
		Schema: map[string]*schema.Schema{
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Organisation ID. Defaults to the provider organisation.",
			},
			"feature_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether KMS is enabled for the organisation.",
			},
			"regions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Per-region KMS and Secrets availability.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"identity": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"slug": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"kms_available": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceKmsSummaryRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	summary, err := client.KMS().GetSummary(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("kms-summary")
	_ = d.Set("feature_enabled", summary.FeatureEnabled)

	regions := make([]map[string]any, 0, len(summary.Regions))
	for _, region := range summary.Regions {
		regions = append(regions, map[string]any{
			"identity":      region.Identity,
			"name":          region.Name,
			"slug":          region.Slug,
			"kms_available": region.KmsAvailable,
		})
	}
	_ = d.Set("regions", regions)

	return nil
}
