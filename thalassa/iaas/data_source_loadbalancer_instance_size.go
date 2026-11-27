package iaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func DataSourceLoadbalancerInstanceSize() *schema.Resource {
	return &schema.Resource{
		Description: "Get a load balancer instance size by name",
		ReadContext: dataSourceLoadbalancerInstanceSizeRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity of the load balancer instance size",
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the load balancer instance size to look up (for example lb-gp-small)",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the load balancer instance size is enabled",
			},
			"ingress_bandwidth_rate": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Ingress bandwidth rate of the load balancer instance size",
			},
			"egress_bandwidth_rate": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Egress bandwidth rate of the load balancer instance size",
			},
		},
	}
}

func dataSourceLoadbalancerInstanceSizeRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)
	sizes, err := client.IaaS().ListLoadbalancerInstanceSizes(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, size := range sizes {
		if size.Name != name {
			continue
		}
		d.SetId(size.Identity)
		_ = d.Set("id", size.Identity)
		_ = d.Set("name", size.Name)
		_ = d.Set("enabled", size.Enabled)
		if size.IngressBandwidthRate != nil {
			_ = d.Set("ingress_bandwidth_rate", int(*size.IngressBandwidthRate))
		}
		if size.EgressBandwidthRate != nil {
			_ = d.Set("egress_bandwidth_rate", int(*size.EgressBandwidthRate))
		}
		return nil
	}

	return diag.FromErr(fmt.Errorf("load balancer instance size %q not found", name))
}
