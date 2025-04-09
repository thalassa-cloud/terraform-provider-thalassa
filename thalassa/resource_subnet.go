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
			"organisation_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the Subnet. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "VPC of the Subnet",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.StringLenBetween(1, 62),
				ForceNew:     true,
				Description:  "Name of the Subnet",
			},
			"slug": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Slug of the Subnet",
			},
			"cidr": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IsCIDR,
				Description:  "CIDR of the Subnet",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.StringLenBetween(0, 255),
				Description:  "A human readable description about the subnet",
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
			"route_table_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Route Table of the Subnet",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the Subnet",
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
		VpcIdentity: d.Get("vpc_id").(string),
		Cidr:        d.Get("cidr").(string),
	}

	if routeTable, ok := d.GetOk("route_table_id"); ok {
		createSubnet.AssociatedRouteTableIdentity = Ptr(routeTable.(string))
	}

	subnet, err := client.IaaS().CreateSubnet(ctx, createSubnet)

	if err != nil {
		return diag.FromErr(err)
	}
	if subnet != nil {
		d.SetId(subnet.Identity)
		d.Set("slug", subnet.Slug)
		d.Set("type", subnet.Type)
		d.Set("status", subnet.Status)

		// wait until the subnet is ready
		for {
			select {
			case <-ctx.Done():
				if subnet != nil {
					return diag.FromErr(fmt.Errorf("timeout while waiting for subnet to be ready. Current status: %s", subnet.Status))
				}
				return diag.FromErr(fmt.Errorf("timeout while waiting for subnet to be ready"))
			case <-time.After(1 * time.Second):
				// continue
				subnet, err := client.IaaS().GetSubnet(ctx, subnet.Identity)
				if err != nil {
					if tcclient.IsNotFound(err) {
						return diag.FromErr(fmt.Errorf("subnet was not found after creation"))
					}
					return diag.FromErr(err)
				}

				if subnet.Status == iaas.SubnetStatusReady {
					d.Set("status", subnet.Status)
					return nil
				}
				d.Set("status", subnet.Status)
			}
		}
	}
	return resourceSubnetRead(ctx, d, m)
}

func resourceSubnetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := getClient(getProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	identity := d.Get("id").(string)
	subnet, err := client.IaaS().GetSubnet(ctx, identity)
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
	if subnet.RouteTable != nil {
		d.Set("route_table_id", subnet.RouteTable.Identity)
	}
	d.Set("status", subnet.Status)
	d.Set("type", subnet.Type)
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
			d.Set("route_table_id", subnet.RouteTable.Identity)
		}

		// wait until the subnet is ready
		for {
			select {
			case <-ctx.Done():
				return diag.FromErr(fmt.Errorf("timeout while waiting for subnet to be ready"))
			case <-time.After(1 * time.Second):
				// continue
				subnet, err := client.IaaS().GetSubnet(ctx, subnet.Identity)
				if err != nil {
					return diag.FromErr(err)
				}
				if subnet.Status == iaas.SubnetStatusReady {
					d.Set("status", subnet.Status)
					return nil
				}
			}
		}
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
