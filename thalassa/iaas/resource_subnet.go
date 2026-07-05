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
		Description: "Create a subnet in a VPC. Subnets provide network segments for resources. " +
			"A VPC can have multiple subnets with unique CIDR blocks. IPv4, IPv6, and dual-stack are supported. " +
			"The CIDR cannot be changed after creation.",
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
				Optional:    true,
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
				Default:      "",
				Optional:     true,
				ValidateFunc: validate.StringLenBetween(0, 255),
				Description:  "A human readable description about the subnet",
			},
			"labels": {
				Type:        schema.TypeMap,
				Default:     make(map[string]string),
				Optional:    true,
				Description: "Labels for the Subnet",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Default:     make(map[string]string),
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

func resourceSubnetCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
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
		_ = d.Set("slug", subnet.Slug)
		_ = d.Set("type", subnet.Type)
		_ = d.Set("status", subnet.Status)

		// wait until the subnet is ready
		ctxWithTimeout, cancel := context.WithTimeout(ctx, 20*time.Minute)
		defer cancel()

		if subnet, err = client.IaaS().WaitUntilSubnetReady(ctxWithTimeout, subnet.Identity); err != nil {
			return diag.FromErr(fmt.Errorf("error waiting for subnet to be ready: %w", err))
		}

		_ = d.Set("status", subnet.Status)
		_ = d.Set("ipv4_addresses_used", subnet.V4usingIPs)
		_ = d.Set("ipv4_addresses_available", subnet.V4availableIPs)
		_ = d.Set("ipv6_addresses_used", subnet.V6usingIPs)
		_ = d.Set("ipv6_addresses_available", subnet.V6availableIPs)
	}
	return resourceSubnetRead(ctx, d, m)
}

func resourceSubnetRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
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
		return diag.FromErr(fmt.Errorf("error getting subnet: %w", err))
	}
	if subnet == nil {
		d.SetId("")
		return nil
	}

	d.SetId(subnet.Identity)
	_ = d.Set("name", subnet.Name)
	_ = d.Set("slug", subnet.Slug)
	_ = d.Set("description", subnet.Description)
	_ = d.Set("labels", subnet.Labels)
	_ = d.Set("annotations", subnet.Annotations)
	_ = d.Set("cidr", subnet.Cidr)
	if subnet.Vpc != nil {
		_ = d.Set("vpc_id", subnet.Vpc.Identity)
	} else if subnet.VpcIdentity != "" {
		_ = d.Set("vpc_id", subnet.VpcIdentity)
	}
	if subnet.RouteTable != nil {
		_ = d.Set("route_table_id", subnet.RouteTable.Identity)
	}
	_ = d.Set("status", subnet.Status)
	_ = d.Set("type", subnet.Type)

	_ = d.Set("ipv4_addresses_used", subnet.V4usingIPs)
	_ = d.Set("ipv4_addresses_available", subnet.V4availableIPs)
	_ = d.Set("ipv6_addresses_used", subnet.V6usingIPs)
	_ = d.Set("ipv6_addresses_available", subnet.V6availableIPs)

	return nil
}

func resourceSubnetUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating client: %w", err))
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
		return diag.FromErr(fmt.Errorf("error updating subnet: %w", err))
	}
	if subnet != nil {
		_ = d.Set("name", subnet.Name)
		_ = d.Set("description", subnet.Description)
		_ = d.Set("slug", subnet.Slug)
		_ = d.Set("labels", subnet.Labels)
		_ = d.Set("annotations", subnet.Annotations)
		_ = d.Set("ipv4_addresses_used", subnet.V4usingIPs)
		_ = d.Set("ipv4_addresses_available", subnet.V4availableIPs)
		_ = d.Set("ipv6_addresses_used", subnet.V6usingIPs)
		_ = d.Set("ipv6_addresses_available", subnet.V6availableIPs)

		if subnet.RouteTable != nil {
			_ = d.Set("route_table_id", subnet.RouteTable.Identity)
		}

		ctxWithTimeout, cancel := context.WithTimeout(ctx, 20*time.Minute)
		defer cancel()

		if subnet, err = client.IaaS().WaitUntilSubnetReady(ctxWithTimeout, subnet.Identity); err != nil {
			return diag.FromErr(fmt.Errorf("error waiting for subnet to be ready: %w", err))
		}

		_ = d.Set("status", subnet.Status)
		_ = d.Set("ipv4_addresses_used", subnet.V4usingIPs)
		_ = d.Set("ipv4_addresses_available", subnet.V4availableIPs)
		_ = d.Set("ipv6_addresses_used", subnet.V6usingIPs)
		_ = d.Set("ipv6_addresses_available", subnet.V6availableIPs)
	}

	return resourceSubnetRead(ctx, d, m)
}

func resourceSubnetDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	if id == "" {
		return nil
	}

	err = client.IaaS().DeleteSubnet(ctx, id)
	if err != nil && !tcclient.IsNotFound(err) {
		return diag.FromErr(err)
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 20*time.Minute)
	defer cancel()

	for {
		select {
		case <-ctxWithTimeout.Done():
			return diag.FromErr(fmt.Errorf("timeout waiting for subnet %s to be deleted: %w", id, ctxWithTimeout.Err()))
		default:
		}

		subnet, err := client.IaaS().GetSubnet(ctx, id)
		if err != nil {
			if tcclient.IsNotFound(err) {
				d.SetId("")
				return nil
			}
			return diag.FromErr(err)
		}

		switch subnet.Status {
		case iaas.SubnetStatusDeleted:
			d.SetId("")
			return nil
		case iaas.SubnetStatusDeleting:
			// Deletion in progress; keep polling.
		case iaas.SubnetStatusReady, iaas.SubnetStatusActive:
			// The delete API can return before the subnet transitions to deleting.
			if err := client.IaaS().DeleteSubnet(ctx, id); err != nil && !tcclient.IsNotFound(err) {
				return diag.FromErr(err)
			}
		case iaas.SubnetStatusFailed:
			return diag.FromErr(fmt.Errorf("subnet %s failed to delete (status: %s)", id, subnet.Status))
		default:
			return diag.FromErr(fmt.Errorf("subnet %s unexpected status during delete: %s", id, subnet.Status))
		}

		time.Sleep(iaas.DefaultPollIntervalForWaiting)
	}
}
