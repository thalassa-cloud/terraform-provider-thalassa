package iaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	iaas "github.com/thalassa-cloud/client-go/iaas"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func DataSourceVpcDefaultRouteTable() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve detailed information about the default route table for a specified Virtual Private Cloud (VPC).",
		ReadContext: dataSourceVpcDefaultRouteTableRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "ID of the route table",
				Computed:    true,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the VPC for which the default route table will be retrieved",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the route table",
			},
			"slug": {
				Type:        schema.TypeString,
				Description: "Slug of the route table",
				Computed:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the route table",
			},
			"labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Labels on the route table",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Annotations on the route table",
			},
			"is_default": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether this is the default route table for the VPC",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last update time",
			},
		},
	}
}

func dataSourceVpcDefaultRouteTableRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	prov := provider.GetProvider(m)
	vpcID := d.Get("vpc_id").(string)

	// Try via VPC reference first
	vpc, err := prov.Client.IaaS().GetVpc(ctx, vpcID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to get VPC %s: %w", vpcID, err))
	}
	if vpc != nil && vpc.RouteTable != nil {
		return setDefaultRouteTable(d, vpc.RouteTable)
	}

	// Fallback: list all route tables and find default for this VPC
	rts, err := prov.Client.IaaS().ListRouteTables(ctx, &iaas.ListRouteTablesRequest{})
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to list route tables: %w", err))
	}
	for i := range rts {
		rt := rts[i]
		if rt.Vpc != nil && rt.Vpc.Identity == vpcID && rt.IsDefault {
			return setDefaultRouteTable(d, &rt)
		}
	}

	return diag.FromErr(fmt.Errorf("default route table not found for VPC %s", vpcID))
}

func setDefaultRouteTable(d *schema.ResourceData, rt *iaas.RouteTable) diag.Diagnostics {
	d.SetId(rt.Identity)
	d.Set("id", rt.Identity)
	d.Set("name", rt.Name)
	d.Set("slug", rt.Slug)
	if rt.Description != nil {
		d.Set("description", *rt.Description)
	}
	if rt.Labels != nil {
		d.Set("labels", rt.Labels)
	}
	if rt.Annotations != nil {
		d.Set("annotations", rt.Annotations)
	}
	d.Set("is_default", rt.IsDefault)
	d.Set("created_at", rt.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
	if rt.UpdatedAt != nil {
		d.Set("updated_at", rt.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))
	}
	return nil
}
