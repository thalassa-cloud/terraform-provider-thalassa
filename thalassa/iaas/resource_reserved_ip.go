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

func resourceReservedIP() *schema.Resource {
	return &schema.Resource{
		Description:   "Reserve a public IPv4/IPv6 address. The IP address will be reserved in the region specified and can be attached to a load balancer or NAT gateway during creation in the same region.",
		CreateContext: resourceReservedIPCreate,
		ReadContext:   resourceReservedIPRead,
		UpdateContext: resourceReservedIPUpdate,
		DeleteContext: resourceReservedIPDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the reserved IP. If not provided, the organisation configured in the Terraform provider will be used.",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Region of the reserved IP. Provide the identity of the region. Can only be set on creation.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.StringLenBetween(1, 62),
				Description:  "Display name of the reserved IP.",
			},
			"slug": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Slug of the reserved IP.",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.StringLenBetween(0, 255),
				Description:  "Human-readable description of the reserved IP.",
			},
			"labels": {
				Type:        schema.TypeMap,
				Default:     make(map[string]string),
				Optional:    true,
				Description: "Labels for the reserved IP.",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Default:     make(map[string]string),
				Optional:    true,
				Description: "Annotations for the reserved IP.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Provisioning and attachment status of the reserved IP.",
			},
			"ipv4_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Allocated public IPv4 address, when available.",
			},
			"ipv6_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Allocated public IPv6 address, when available.",
			},
			"attached_to_resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the resource this IP is attached to, if any.",
			},
			"attached_to_resource_identity": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity of the resource this IP is attached to, if any.",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CustomizeDiff: func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
			_, new := diff.GetChange("description")
			if new == nil {
				return diff.SetNew("description", "")
			}
			return nil
		},
	}
}

func resourceReservedIPCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	regionInput := d.Get("region").(string)
	var region *iaas.Region
	region, err = client.IaaS().GetRegion(ctx, regionInput)
	if err != nil {
		if !tcclient.IsNotFound(err) {
			return diag.FromErr(err)
		}
		regions, err := client.IaaS().ListRegions(ctx, &iaas.ListRegionsRequest{})
		if err != nil {
			return diag.FromErr(err)
		}
		for _, r := range regions {
			if r.Identity == regionInput || r.Slug == regionInput {
				region = &r
				break
			}
		}
	}
	if region == nil {
		return diag.FromErr(fmt.Errorf("region not found"))
	}

	create := iaas.CreateReservedIpRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Labels:      convert.ConvertToMap(d.Get("labels")),
		Annotations: convert.ConvertToMap(d.Get("annotations")),
		Region:      regionInput,
	}

	fip, err := client.IaaS().CreateReservedIP(ctx, create)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create reserved IP: %w", err))
	}
	if fip == nil {
		return diag.FromErr(fmt.Errorf("reserved IP was not returned after creation"))
	}

	d.SetId(fip.Identity)

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 20*time.Minute)
	defer cancel()
	for {
		select {
		case <-ctxWithTimeout.Done():
			return diag.FromErr(fmt.Errorf("timeout while waiting for reserved IP to become available"))
		case <-time.After(1 * time.Second):
		}
		fip, err = client.IaaS().GetReservedIP(ctxWithTimeout, fip.Identity)
		if err != nil {
			if tcclient.IsNotFound(err) {
				return diag.FromErr(fmt.Errorf("reserved IP %s was not found after creation", d.Id()))
			}
			return diag.FromErr(err)
		}
		if fip == nil {
			return diag.FromErr(fmt.Errorf("reserved IP %s was not found after creation", d.Id()))
		}
		switch fip.Status {
		case iaas.ReservedIpStatusAvailable, iaas.ReservedIpStatusAttached:
			return resourceReservedIPRead(ctx, d, m)
		case iaas.ReservedIpStatusFailed:
			return diag.FromErr(fmt.Errorf("reserved IP provisioning failed (status: %s)", fip.Status))
		case iaas.ReservedIpStatusDeleted, iaas.ReservedIpStatusDeleting:
			return diag.FromErr(fmt.Errorf("reserved IP entered unexpected status %s", fip.Status))
		}
	}
}

func resourceReservedIPRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to get client: %w", err))
	}

	id := d.Get("id").(string)
	fip, err := client.IaaS().GetReservedIP(ctx, id)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error getting reserved IP: %w", err))
	}
	if fip == nil {
		d.SetId("")
		return nil
	}

	d.SetId(fip.Identity)
	d.Set("name", fip.Name)
	d.Set("slug", fip.Slug)
	d.Set("description", fip.Description)
	d.Set("labels", fip.Labels)
	d.Set("annotations", fip.Annotations)
	d.Set("status", string(fip.Status))
	d.Set("ipv4_address", fip.IPv4Address)
	d.Set("ipv6_address", fip.IPv6Address)
	d.Set("attached_to_resource_type", string(fip.AttachedToResourceType))
	d.Set("attached_to_resource_identity", fip.AttachedToResourceIdentity)

	if fip.Region != nil {
		d.Set("region", fip.Region.Slug)
	}

	return nil
}

func resourceReservedIPUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to get client: %w", err))
	}

	update := iaas.UpdateReservedIpRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Labels:      convert.ConvertToMap(d.Get("labels")),
		Annotations: convert.ConvertToMap(d.Get("annotations")),
	}

	id := d.Get("id").(string)
	fip, err := client.IaaS().UpdateReservedIP(ctx, id, update)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("failed to update reserved IP: %w", err))
	}
	if fip != nil {
		d.Set("name", fip.Name)
		d.Set("description", fip.Description)
		d.Set("slug", fip.Slug)
		d.Set("labels", fip.Labels)
		d.Set("annotations", fip.Annotations)
	}
	return resourceReservedIPRead(ctx, d, m)
}

func resourceReservedIPDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to get client: %w", err))
	}

	id := d.Get("id").(string)
	err = client.IaaS().DeleteReservedIP(ctx, id)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("failed to delete reserved IP: %w", err))
	}

	d.SetId("")
	return nil
}
