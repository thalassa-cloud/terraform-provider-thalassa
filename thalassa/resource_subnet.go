package thalassa

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"

	iaas "github.com/thalassa-cloud/client-go/pkg/iaas"
)

func resourceSubnet() *schema.Resource {
	return &schema.Resource{
		Description:   "Create an subnet",
		CreateContext: resourceSubnetCreate,
		ReadContext:   resourceSubnetRead,
		UpdateContext: resourceSubnetUpdate,
		DeleteContext: resourceSubnetDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organisation": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Organisation of the Subnet",
			},
			"vpc": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "VPC of the Subnet",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the Subnet",
			},
			"slug": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cidr": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "CIDR of the Subnet",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A human readable description about the subnet",
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Labels for the Subnet",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Annotations for the Subnet",
			},
			"zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Zone of the Subnet",
			},
			"route_table": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Route Table of the Subnet",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceSubnetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := getClient(getProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	createSubnet := iaas.CreateSubnet{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Labels:      convertToMap(d.Get("labels")),
		Annotations: convertToMap(d.Get("annotations")),
		VpcIdentity: d.Get("vpc").(string),
		CloudZone:   d.Get("zone").(string),
		Cidr:        d.Get("cidr").(string),
	}

	if routeTable, ok := d.GetOk("route_table"); ok {
		createSubnet.AssociatedRouteTableIdentity = Ptr(routeTable.(string))
	}

	subnet, err := client.IaaS().CreateSubnet(ctx, createSubnet)

	if err != nil {
		return diag.FromErr(err)
	}
	if subnet != nil {
		d.SetId(subnet.Identity)
		d.Set("slug", subnet.Slug)
		return nil
	}
	return resourceSubnetRead(ctx, d, m)
}

func resourceSubnetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := getClient(getProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	slug := d.Get("slug").(string)
	subnet, err := client.IaaS().GetSubnet(ctx, slug)
	if err != nil && !tcclient.IsNotFound(err) {
		return diag.FromErr(fmt.Errorf("error getting subnet: %s", err))
	}
	if subnet == nil {
		return diag.FromErr(fmt.Errorf("subnet was not found"))
	}

	d.SetId(subnet.Identity)
	d.Set("name", subnet.Name)
	d.Set("slug", subnet.Slug)
	d.Set("description", subnet.Description)
	d.Set("labels", subnet.Labels)
	d.Set("annotations", subnet.Annotations)
	d.Set("zone", subnet.CloudZone)
	if subnet.RouteTable != nil {
		d.Set("route_table", subnet.RouteTable.Slug)
	}

	return nil
}

func resourceSubnetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := getClient(getProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	updateSubnet := iaas.UpdateSubnet{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Labels:      convertToMap(d.Get("labels")),
		Annotations: convertToMap(d.Get("annotations")),
	}

	slug := d.Get("slug").(string)

	subnet, err := client.IaaS().UpdateSubnet(ctx, slug, updateSubnet)
	if err != nil {
		return diag.FromErr(err)
	}
	if subnet != nil {
		d.Set("name", subnet.Name)
		d.Set("description", subnet.Description)
		d.Set("slug", subnet.Slug)
		d.Set("labels", subnet.Labels)
		d.Set("annotations", subnet.Annotations)
		if subnet.RouteTable != nil {
			d.Set("route_table", subnet.RouteTable.Slug)
		}
		return nil
	}

	return resourceSubnetRead(ctx, d, m)
}

func resourceSubnetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := getClient(getProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("id").(string)

	err = client.IaaS().DeleteSubnet(ctx, id)
	if err != nil && !tcclient.IsNotFound(err) {
		return diag.FromErr(err)
	}

	// wait until the subnet is deleted
	for {
		_, err := client.IaaS().GetSubnet(ctx, id)
		if err != nil && tcclient.IsNotFound(err) {
			break
		}
		time.Sleep(1 * time.Second)
	}

	d.SetId("")
	return nil
}
