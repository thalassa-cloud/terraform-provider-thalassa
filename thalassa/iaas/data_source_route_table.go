package iaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	iaas "github.com/thalassa-cloud/client-go/iaas"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func DataSourceRouteTable() *schema.Resource {
	return &schema.Resource{
		Description: "Get a VPC route table",
		ReadContext: dataSourceRouteTableRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"identity": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the route table",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the route table",
			},
			"slug": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the route table",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the VPC for the route table",
			},
			"label_selector": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Match route tables that have all provided labels (exact key=value match)",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"is_default": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Filter by default route table",
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

func dataSourceRouteTableRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	prov := provider.GetProvider(m)

	if id := d.Get("identity").(string); id != "" {
		rt, err := prov.Client.IaaS().GetRouteTable(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}
		if rt == nil {
			return diag.FromErr(fmt.Errorf("route table %s not found", id))
		}
		return setRouteTableData(d, rt)
	}

	rts, err := prov.Client.IaaS().ListRouteTables(ctx, &iaas.ListRouteTablesRequest{})
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)
	vpcID := d.Get("vpc_id").(string)
	labelSelector, hasLabelSelector := d.GetOk("label_selector")
	isDefault, hasIsDefault := d.GetOk("is_default")

	for i := range rts {
		rt := rts[i]
		if name != "" && rt.Name != name {
			continue
		}
		if vpcID != "" && (rt.Vpc == nil || rt.Vpc.Identity != vpcID) {
			continue
		}
		if hasIsDefault && rt.IsDefault != isDefault.(bool) {
			continue
		}
		if hasLabelSelector {
			match := true
			selector := labelSelector.(map[string]interface{})
			for k, v := range selector {
				valStr := v.(string)
				if rt.Labels == nil {
					match = false
					break
				}
				if rtVal, ok := rt.Labels[k]; !ok || rtVal != valStr {
					match = false
					break
				}
			}
			if !match {
				continue
			}
		}
		return setRouteTableData(d, &rt)
	}

	return diag.FromErr(fmt.Errorf("route table not found with the given criteria"))
}

func setRouteTableData(d *schema.ResourceData, rt *iaas.RouteTable) diag.Diagnostics {
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
