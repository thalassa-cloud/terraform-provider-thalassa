package thalassa

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"

	iaas "github.com/thalassa-cloud/client-go/pkg/iaas"
)

func resourceRouteTable() *schema.Resource {
	return &schema.Resource{
		Description:   "Create an routeTable",
		CreateContext: resourceRouteTableCreate,
		ReadContext:   resourceRouteTableRead,
		UpdateContext: resourceRouteTableUpdate,
		DeleteContext: resourceRouteTableDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the RouteTable. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "VPC of the RouteTable",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StringLenBetween(1, 62),
				Description:  "Name of the RouteTable",
			},
			"slug": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.StringLenBetween(0, 255),
				Description:  "A human readable description about the routeTable",
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Labels for the RouteTable",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Annotations for the RouteTable",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceRouteTableCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := getClient(getProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	createRouteTable := iaas.CreateRouteTable{
		Name:        d.Get("name").(string),
		Description: Ptr(d.Get("description").(string)),
		Labels:      convertToMap(d.Get("labels")),
		Annotations: convertToMap(d.Get("annotations")),
		VpcIdentity: d.Get("vpc_id").(string),
	}

	routeTable, err := client.IaaS().CreateRouteTable(ctx, createRouteTable)

	if err != nil {
		return diag.FromErr(err)
	}
	if routeTable != nil {
		d.SetId(routeTable.Identity)
		d.Set("slug", routeTable.Slug)
		return nil
	}
	return resourceRouteTableRead(ctx, d, m)
}

func resourceRouteTableRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := getClient(getProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("id").(string)
	routeTable, err := client.IaaS().GetRouteTable(ctx, id)
	if err != nil && !tcclient.IsNotFound(err) {
		return diag.FromErr(fmt.Errorf("error getting route table %q: %s", id, err))
	}
	if routeTable == nil {
		return diag.FromErr(fmt.Errorf("route table %q was not found", id))
	}

	d.SetId(routeTable.Identity)
	d.Set("name", routeTable.Name)
	d.Set("slug", routeTable.Slug)
	d.Set("description", routeTable.Description)
	d.Set("labels", routeTable.Labels)
	d.Set("annotations", routeTable.Annotations)
	if routeTable.Vpc != nil {
		d.Set("vpc_id", routeTable.Vpc.Identity)
	}

	return nil
}

func resourceRouteTableUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := getClient(getProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	updateRouteTable := iaas.UpdateRouteTable{
		Name:        Ptr(d.Get("name").(string)),
		Description: Ptr(d.Get("description").(string)),
		Labels:      convertToMap(d.Get("labels")),
		Annotations: convertToMap(d.Get("annotations")),
	}

	id := d.Get("id").(string)

	routeTable, err := client.IaaS().UpdateRouteTable(ctx, id, updateRouteTable)
	if err != nil {
		return diag.FromErr(err)
	}
	if routeTable != nil {
		d.Set("name", routeTable.Name)
		d.Set("description", routeTable.Description)
		d.Set("slug", routeTable.Slug)
		d.Set("labels", routeTable.Labels)
		d.Set("annotations", routeTable.Annotations)
		if routeTable.Vpc != nil {
			d.Set("vpc_id", routeTable.Vpc.Identity)
		}
		return nil
	}

	return resourceRouteTableRead(ctx, d, m)
}

func resourceRouteTableDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := getClient(getProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("id").(string)

	err = client.IaaS().DeleteRouteTable(ctx, id)
	if err != nil && !tcclient.IsNotFound(err) {
		return diag.FromErr(err)
	}

	// wait until the route table is deleted
	for {
		_, err := client.IaaS().GetRouteTable(ctx, id)
		if err != nil && tcclient.IsNotFound(err) {
			break
		}
		time.Sleep(1 * time.Second)
	}

	d.SetId("")

	return nil
}
