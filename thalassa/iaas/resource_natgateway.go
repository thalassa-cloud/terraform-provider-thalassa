package iaas

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"

	iaas "github.com/thalassa-cloud/client-go/iaas"
)

func resourceNatGateway() *schema.Resource {
	return &schema.Resource{
		Description:   "Create an NAT Gateway within a VPC",
		CreateContext: resourceNatGatewayCreate,
		ReadContext:   resourceNatGatewayRead,
		UpdateContext: resourceNatGatewayUpdate,
		DeleteContext: resourceNatGatewayDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the NatGateway. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "VPC of the NatGateway",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Subnet of the NatGateway",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.StringLenBetween(1, 62),
				ForceNew:     true,
				Description:  "Name of the NatGateway",
			},
			"slug": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.StringLenBetween(0, 255),
				Description:  "A human readable description about the natGateway",
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Labels for the NatGateway",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Annotations for the NatGateway",
			},
			"endpoint_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Endpoint IP of the NatGateway",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the NatGateway",
			},
			"v4_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "V4 IP of the NatGateway",
			},
			"v6_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "V6 IP of the NatGateway",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceNatGatewayCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	createNatGateway := iaas.CreateVpcNatGateway{
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		Labels:         convert.ConvertToMap(d.Get("labels")),
		Annotations:    convert.ConvertToMap(d.Get("annotations")),
		SubnetIdentity: d.Get("subnet_id").(string),
	}

	natGateway, err := client.IaaS().CreateNatGateway(ctx, createNatGateway)

	if err != nil {
		return diag.FromErr(err)
	}
	if natGateway != nil {
		d.SetId(natGateway.Identity)
		d.Set("slug", natGateway.Slug)

		// wait until the natGateway is ready and has an endpoint IP
		ctxWithTimeout, cancel := context.WithTimeout(ctx, 20*time.Minute)
		defer cancel()
		for {
			select {
			case <-ctxWithTimeout.Done():
				return diag.FromErr(fmt.Errorf("timeout while waiting for natGateway to be ready"))
			case <-time.After(1 * time.Second):
				// continue
				natGateway, err = client.IaaS().GetNatGateway(ctxWithTimeout, natGateway.Identity)
				if err != nil {
					if tcclient.IsNotFound(err) {
						return diag.FromErr(fmt.Errorf("natGateway %s was not found after creation", natGateway.Identity))
					}
					return diag.FromErr(err)
				}
				if natGateway == nil {
					return diag.FromErr(fmt.Errorf("natGateway %s was not found after creation", natGateway.Identity))
				}
				if strings.TrimSpace(natGateway.EndpointIP) != "" {
					d.Set("status", natGateway.Status)
					d.Set("endpoint_ip", natGateway.EndpointIP)
					d.Set("v4_ip", natGateway.V4IP)
					d.Set("v6_ip", natGateway.V6IP)
					return nil
				}
			}
		}
	}
	return resourceNatGatewayRead(ctx, d, m)
}

func resourceNatGatewayRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("id").(string)
	natGateway, err := client.IaaS().GetNatGateway(ctx, id)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error getting natGateway: %s", err))
	}
	if natGateway == nil {
		return diag.FromErr(fmt.Errorf("natGateway was not found"))
	}

	d.SetId(natGateway.Identity)
	d.Set("name", natGateway.Name)
	d.Set("slug", natGateway.Slug)
	d.Set("description", natGateway.Description)
	d.Set("labels", natGateway.Labels)
	d.Set("annotations", natGateway.Annotations)
	d.Set("subnet_id", natGateway.Subnet.Identity)
	d.Set("vpc_id", natGateway.Vpc.Identity)
	d.Set("endpoint_ip", natGateway.EndpointIP)
	d.Set("status", natGateway.Status)
	d.Set("v4_ip", natGateway.V4IP)
	d.Set("v6_ip", natGateway.V6IP)
	return nil
}

func resourceNatGatewayUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	updateNatGateway := iaas.UpdateVpcNatGateway{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Labels:      convert.ConvertToMap(d.Get("labels")),
		Annotations: convert.ConvertToMap(d.Get("annotations")),
	}

	slug := d.Get("slug").(string)

	natGateway, err := client.IaaS().UpdateNatGateway(ctx, slug, updateNatGateway)
	if err != nil {
		return diag.FromErr(err)
	}
	if natGateway != nil {
		d.Set("name", natGateway.Name)
		d.Set("description", natGateway.Description)
		d.Set("slug", natGateway.Slug)
		d.Set("labels", natGateway.Labels)
		d.Set("annotations", natGateway.Annotations)
		return nil
	}

	return resourceNatGatewayRead(ctx, d, m)
}

func resourceNatGatewayDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("id").(string)

	err = client.IaaS().DeleteNatGateway(ctx, id)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	// wait until the natGateway is deleted
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 20*time.Minute)
	defer cancel()
	if err := client.IaaS().WaitUntilNatGatewayDeleted(ctxWithTimeout, id); err != nil {
		return diag.FromErr(fmt.Errorf("error waiting for natGateway to be deleted: %w", err))
	}
	d.SetId("")
	return nil
}
