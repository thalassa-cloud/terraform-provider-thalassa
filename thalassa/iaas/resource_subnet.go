package iaas

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"

	iaas "github.com/thalassa-cloud/client-go/iaas"
)

func resourceSubnet() *schema.Resource {
	return &schema.Resource{
		Description:   "Create an subnet in a VPC. Subnets are used to create a network for your resources. A VPC can have multiple subnets, and each subnet must have a different CIDR block. IPv4, IPv6 and Dual-stack subnets are supported. After creationg the CIDR cannot be changed.",
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
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the Subnet",
			},
			"ipv4_addresses_used": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of IPv4 addresses used in the Subnet",
			},
			"ipv4_addresses_available": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of IPv4 addresses available in the Subnet",
			},
			"ipv6_addresses_used": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of IPv6 addresses used in the Subnet",
			},
			"ipv6_addresses_available": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of IPv6 addresses available in the Subnet",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceSubnetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	createSubnet := iaas.CreateSubnet{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Labels:      convert.ConvertToMap(d.Get("labels")),
		Annotations: convert.ConvertToMap(d.Get("annotations")),
		VpcIdentity: d.Get("vpc_id").(string),
		Cidr:        d.Get("cidr").(string),
	}

	if routeTable, ok := d.GetOk("route_table_id"); ok {
		createSubnet.AssociatedRouteTableIdentity = convert.Ptr(routeTable.(string))
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
		ctxWithTimeout, cancel := context.WithTimeout(ctx, 20*time.Minute)
		defer cancel()
		for {
			select {
			case <-ctxWithTimeout.Done():
				if subnet != nil {
					return diag.FromErr(fmt.Errorf("timeout while waiting for subnet to be ready. Current status: %s", subnet.Status))
				}
				return diag.FromErr(fmt.Errorf("timeout while waiting for subnet to be ready"))
			case <-time.After(1 * time.Second):
				// continue
				subnet, err = client.IaaS().GetSubnet(ctxWithTimeout, subnet.Identity)
				if err != nil {
					if tcclient.IsNotFound(err) {
						return diag.FromErr(fmt.Errorf("subnet %s was not found after creation", subnet.Identity))
					}
					return diag.FromErr(err)
				}

				if subnet.Status == iaas.SubnetStatusReady {
					d.Set("status", subnet.Status)
					d.Set("ipv4_addresses_used", subnet.V4usingIPs)
					d.Set("ipv4_addresses_available", subnet.V4availableIPs)
					d.Set("ipv6_addresses_used", subnet.V6usingIPs)
					d.Set("ipv6_addresses_available", subnet.V6availableIPs)

					return nil
				} else if subnet.Status == iaas.SubnetStatusFailed {
					return diag.FromErr(fmt.Errorf("subnet is in failed state: %s", subnet.Status))
				}
				d.Set("status", subnet.Status)
			}
		}
	}
	return resourceSubnetRead(ctx, d, m)
}

func resourceSubnetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	identity := d.Get("id").(string)
	subnet, err := client.IaaS().GetSubnet(ctx, identity)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error getting subnet: %s", err))
	}
	if subnet == nil {
		d.SetId("")
		return nil
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

	d.Set("ipv4_addresses_used", subnet.V4usingIPs)
	d.Set("ipv4_addresses_available", subnet.V4availableIPs)
	d.Set("ipv6_addresses_used", subnet.V6usingIPs)
	d.Set("ipv6_addresses_available", subnet.V6availableIPs)

	return nil
}

func resourceSubnetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	updateSubnet := iaas.UpdateSubnet{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Labels:      convert.ConvertToMap(d.Get("labels")),
		Annotations: convert.ConvertToMap(d.Get("annotations")),
	}

	if routeTable, ok := d.GetOk("route_table_id"); ok {
		updateSubnet.AssociatedRouteTableIdentity = convert.Ptr(routeTable.(string))
	}

	identity := d.Get("id").(string)

	subnet, err := client.IaaS().UpdateSubnet(ctx, identity, updateSubnet)
	if err != nil {
		return diag.FromErr(err)
	}
	if subnet != nil {
		d.Set("name", subnet.Name)
		d.Set("description", subnet.Description)
		d.Set("slug", subnet.Slug)
		d.Set("labels", subnet.Labels)
		d.Set("annotations", subnet.Annotations)
		d.Set("ipv4_addresses_used", subnet.V4usingIPs)
		d.Set("ipv4_addresses_available", subnet.V4availableIPs)
		d.Set("ipv6_addresses_used", subnet.V6usingIPs)
		d.Set("ipv6_addresses_available", subnet.V6availableIPs)

		if subnet.RouteTable != nil {
			d.Set("route_table_id", subnet.RouteTable.Identity)
		}

		ctxWithTimeout, cancel := context.WithTimeout(ctx, 20*time.Minute)
		defer cancel()
		// wait until the subnet is ready
		for {
			select {
			case <-ctxWithTimeout.Done():
				return diag.FromErr(fmt.Errorf("timeout while waiting for subnet to be ready. Current status: %s", subnet.Status))
			case <-time.After(1 * time.Second):
				// continue
				subnet, err = client.IaaS().GetSubnet(ctxWithTimeout, subnet.Identity)
				if err != nil {
					if tcclient.IsNotFound(err) {
						return diag.FromErr(fmt.Errorf("subnet %s was not found after update", subnet.Identity))
					}
					return diag.FromErr(err)
				}
				if subnet.Status == iaas.SubnetStatusReady {
					d.Set("status", subnet.Status)
					return nil
				} else if subnet.Status == iaas.SubnetStatusFailed {
					return diag.FromErr(fmt.Errorf("subnet is in failed state: %s", subnet.Status))
				}
			}
		}
	}

	return resourceSubnetRead(ctx, d, m)
}

func resourceSubnetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("id").(string)

	err = client.IaaS().DeleteSubnet(ctx, id)
	if err != nil && !tcclient.IsNotFound(err) {
		return diag.FromErr(err)
	}

	// wait until the subnet is deleted
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 20*time.Minute)
	defer cancel()
	for {
		select {
		case <-ctxWithTimeout.Done():
			return diag.FromErr(fmt.Errorf("timeout while waiting for subnet to be deleted"))
		case <-time.After(1 * time.Second):
			// continue
			_, err := client.IaaS().GetSubnet(ctxWithTimeout, id)
			if err != nil && tcclient.IsNotFound(err) {
				d.SetId("")
				return nil
			}
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
}
