package iaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"

	iaas "github.com/thalassa-cloud/client-go/iaas"
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
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
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
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Slug of the Vpc",
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
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	// ensure the region exists
	region, err := client.IaaS().GetRegion(ctx, d.Get("region").(string))
	if err != nil {
		if !tcclient.IsNotFound(err) {
			return diag.FromErr(err)
		}
		// check if we can find the region using list
		regions, err := client.IaaS().ListRegions(ctx, &iaas.ListRegionsRequest{})
		if err != nil {
			return diag.FromErr(err)
		}
		for _, r := range regions {
			if r.Identity == d.Get("region").(string) || r.Slug == d.Get("region").(string) {
				region = &r
				break
			}
		}
	}
	if region == nil {
		return diag.FromErr(fmt.Errorf("region not found"))
	}

	createVpc := iaas.CreateVpc{
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		Labels:              convert.ConvertToMap(d.Get("labels")),
		Annotations:         convert.ConvertToMap(d.Get("annotations")),
		CloudRegionIdentity: d.Get("region").(string),
		VpcCidrs:            convert.ConvertToStringSlice(d.Get("cidrs")),
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
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	identity := d.Get("id").(string)
	vpc, err := client.IaaS().GetVpc(ctx, identity)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error getting vpc: %s", err))
	}
	if vpc == nil {
		d.SetId("")
		return nil
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
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	updateVpc := iaas.UpdateVpc{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Labels:      convert.ConvertToMap(d.Get("labels")),
		Annotations: convert.ConvertToMap(d.Get("annotations")),
		VpcCidrs:    convert.ConvertToStringSlice(d.Get("cidrs")),
	}

	identity := d.Get("id").(string)

	vpc, err := client.IaaS().UpdateVpc(ctx, identity, updateVpc)
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
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	identity := d.Get("id").(string)
	error := client.IaaS().DeleteVpc(ctx, identity)
	if error != nil {
		return diag.FromErr(error)
	}

	d.SetId("")

	return nil
}
