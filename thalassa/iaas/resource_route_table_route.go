package iaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"

	iaas "github.com/thalassa-cloud/client-go/iaas"
)

func resourceRouteTableRoute() *schema.Resource {
	return &schema.Resource{
		Description:   "Create an route table route with a destination cidr block, target gateway, target nat gateway and gateway address within a route table.",
		CreateContext: resourceRouteTableRouteCreate,
		ReadContext:   resourceRouteTableRouteRead,
		UpdateContext: resourceRouteTableRouteUpdate,
		DeleteContext: resourceRouteTableRouteDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Organisation of the RouteTable",
			},
			"route_table_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "RouteTable of the Route",
			},
			"destination_cidr": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.IsCIDR,
				Description:  "Destination CIDR of the Route",
			},
			"notes": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.StringLenBetween(0, 255),
				Description:  "Notes for the Route",
			},
			"target_gateway": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Target Gateway of the Route",
			},
			"target_natgateway": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Target NAT Gateway of the Route",
			},
			"gateway_address": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.IsIPAddress,
				Description:  "Gateway Address of the Route",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceRouteTableRouteCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	createRouteTableRoute := iaas.CreateRouteTableRoute{
		DestinationCidrBlock:     d.Get("destination_cidr").(string),
		TargetGatewayIdentity:    d.Get("target_gateway").(string),
		TargetNatGatewayIdentity: d.Get("target_natgateway").(string),
		GatewayAddress:           d.Get("gateway_address").(string),
	}
	route, err := client.IaaS().CreateRouteTableRoute(ctx, d.Get("route_table_id").(string), createRouteTableRoute)

	if err != nil {
		return diag.FromErr(err)
	}
	if route != nil {
		d.SetId(route.Identity)
		return nil
	}
	return resourceRouteTableRead(ctx, d, m)
}

func resourceRouteTableRouteRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("id").(string)
	routeTable := d.Get("route_table_id").(string)
	route, err := client.IaaS().GetRouteTableRoute(ctx, routeTable, id)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error getting route: %s", err))
	}
	if route == nil {
		return diag.FromErr(fmt.Errorf("route was not found"))
	}

	d.SetId(route.Identity)
	d.Set("destination_cidr", route.DestinationCidrBlock)
	if route.TargetGateway != nil {
		d.Set("target_gateway", route.TargetGateway.Identity)
	}
	if route.TargetNatGateway != nil {
		d.Set("target_natgateway", route.TargetNatGateway.Identity)
	}
	if route.GatewayAddress != nil {
		d.Set("gateway_address", route.GatewayAddress)
	}
	return nil
}

func resourceRouteTableRouteUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	updateRouteTableRoute := iaas.UpdateRouteTableRoute{
		DestinationCidrBlock:     d.Get("destination_cidr").(string),
		TargetGatewayIdentity:    d.Get("target_gateway").(string),
		TargetNatGatewayIdentity: d.Get("target_natgateway").(string),
		GatewayAddress:           d.Get("gateway_address").(string),
	}

	id := d.Get("id").(string)
	routeTableIdentity := d.Get("route_table_id").(string)

	// get the route table
	rt, err := client.IaaS().GetRouteTable(ctx, routeTableIdentity)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting route table %q: %s", routeTableIdentity, err))
	}
	if rt == nil {
		return diag.FromErr(fmt.Errorf("route table %q was not found", routeTableIdentity))
	}

	routeTable, err := client.IaaS().UpdateRouteTableRoute(ctx, rt.Identity, id, updateRouteTableRoute)
	if err != nil {
		return diag.FromErr(err)
	}
	if routeTable != nil {
		return nil
	}
	return resourceRouteTableRead(ctx, d, m)
}

func resourceRouteTableRouteDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}
	id := d.Get("id").(string)
	routeTableIdentity := d.Get("route_table_id").(string)

	err = client.IaaS().DeleteRouteTableRoute(ctx, routeTableIdentity, id)
	if err != nil {
		if !tcclient.IsNotFound(err) {
			return diag.FromErr(fmt.Errorf("error deleting routeTable: %s", err))
		}
	}
	d.SetId("")
	return nil
}
