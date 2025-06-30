package dbaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/thalassa-cloud/client-go/dbaas/dbaasalphav1"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func DataSourceDbCluster() *schema.Resource {
	return &schema.Resource{
		Description: "Get an Kubernetes cluster",
		ReadContext: dataSourceDbClusterRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Kubernetes version.",
			},
			"slug": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The slug of the Kubernetes version.",
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Organisation of the Kubernetes Cluster",
			},
		},
	}
}

func dataSourceDbClusterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := provider.GetProvider(m)
	slug := d.Get("slug").(string)

	clusters, err := provider.Client.DbaaSAlphaV1().ListDbClusters(ctx, &dbaasalphav1.ListDbClustersRequest{})
	if err != nil {
		return diag.FromErr(err)
	}

	for _, cluster := range clusters {
		if slug != "" && cluster.Identity == slug {

			return diag.Diagnostics{}
		}
	}
	return diag.FromErr(fmt.Errorf("not found"))
}
