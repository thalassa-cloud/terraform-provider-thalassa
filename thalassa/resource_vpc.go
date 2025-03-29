package thalassa

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"

	iaas "github.com/thalassa-cloud/client-go/pkg/iaas"
)

func resourceVpc() *schema.Resource {
	return &schema.Resource{
		Description:   "Create an vpc",
		CreateContext: resourceVpcCreate,
		ReadContext:   resourceVpcRead,
		UpdateContext: resourceVpcUpdate,
		DeleteContext: resourceVpcDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organisation": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the Vpc. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StringLenBetween(1, 62),
				Description:  "Name of the Vpc",
			},
			"slug": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cidrs": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of CIDRs for the Vpc",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validate.IsCIDR,
				},
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.StringLenBetween(0, 255),
				Description:  "A human readable description about the vpc",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Region of the Vpc. Provide the identity of the region. Can only be set on creation.",
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Labels for the Vpc",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Annotations for the Vpc",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the Vpc",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceVpcCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := getClient(getProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	createVpc := iaas.CreateVpc{
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		Labels:              convertToMap(d.Get("labels")),
		Annotations:         convertToMap(d.Get("annotations")),
		CloudRegionIdentity: d.Get("region").(string),
		VpcCidrs:            convertToStringSlice(d.Get("cidrs")),
	}

	vpc, err := client.IaaS().CreateVpc(ctx, createVpc)

	if err != nil {
		return diag.FromErr(err)
	}
	if vpc != nil {
		d.SetId(vpc.Identity)
		d.Set("slug", vpc.Slug)
		return nil
	}
	return resourceVpcRead(ctx, d, m)
}

func resourceVpcRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := getClient(getProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	slug := d.Get("slug").(string)
	vpc, err := client.IaaS().GetVpc(ctx, slug)
	if err != nil && !tcclient.IsNotFound(err) {
		return diag.FromErr(fmt.Errorf("error getting vpc: %s", err))
	}
	if vpc == nil {
		return diag.FromErr(fmt.Errorf("vpc was not found"))
	}

	d.SetId(vpc.Identity)
	d.Set("name", vpc.Name)
	d.Set("slug", vpc.Slug)
	d.Set("description", vpc.Description)
	d.Set("labels", vpc.Labels)
	d.Set("annotations", vpc.Annotations)
	d.Set("region", vpc.CloudRegion.Slug)
	d.Set("cidrs", vpc.CIDRs)
	d.Set("status", vpc.Status)
	return nil
}

func resourceVpcUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := getClient(getProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	updateVpc := iaas.UpdateVpc{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Labels:      convertToMap(d.Get("labels")),
		Annotations: convertToMap(d.Get("annotations")),
		VpcCidrs:    convertToStringSlice(d.Get("cidrs")),
	}

	id := d.Get("id").(string)

	vpc, err := client.IaaS().UpdateVpc(ctx, id, updateVpc)
	if err != nil {
		return diag.FromErr(err)
	}
	if vpc != nil {
		d.Set("name", vpc.Name)
		d.Set("description", vpc.Description)
		d.Set("slug", vpc.Slug)
		d.Set("labels", vpc.Labels)
		d.Set("annotations", vpc.Annotations)
		d.Set("cidrs", vpc.CIDRs)
		d.Set("status", vpc.Status)
		return nil
	}

	return resourceVpcRead(ctx, d, m)
}

func resourceVpcDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := getClient(getProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("id").(string)
	error := client.IaaS().DeleteVpc(ctx, id)
	if error != nil {
		return diag.FromErr(error)
	}

	d.SetId("")

	return nil
}
