package dns

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	tcdns "github.com/thalassa-cloud/client-go/dns"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"

	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func ResourceDnsZone() *schema.Resource {
	return &schema.Resource{
		Description:   "Create and manage a DNS zone in Thalassa Cloud",
		CreateContext: resourceDnsZoneCreate,
		ReadContext:   resourceDnsZoneRead,
		UpdateContext: resourceDnsZoneUpdate,
		DeleteContext: resourceDnsZoneDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Platform identity of the DNS zone (dnsz-…).",
			},
			"organisation_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"zone_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "DNS zone name (e.g. example.com). Renaming requires recreation.",
			},
			"slug": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"annotations": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"object_version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDnsZoneCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	zone, err := client.DNS().CreateZone(ctx, tcdns.CreateDnsZoneRequest{
		ZoneName:    d.Get("zone_name").(string),
		Description: d.Get("description").(string),
		Labels:      convert.ConvertToMap(d.Get("labels")),
		Annotations: convert.ConvertToMap(d.Get("annotations")),
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("creating DNS zone: %w", err))
	}

	d.SetId(zone.Identity)
	return resourceDnsZoneReadWithZone(ctx, d, zone)
}

func resourceDnsZoneRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	zone, err := client.DNS().GetZone(ctx, d.Id())
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("reading DNS zone: %w", err))
	}

	return resourceDnsZoneReadWithZone(ctx, d, zone)
}

func resourceDnsZoneReadWithZone(_ context.Context, d *schema.ResourceData, zone *tcdns.DnsZone) diag.Diagnostics {
	if err := setDnsZoneState(d, zone); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceDnsZoneUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.DNS().UpdateZone(ctx, d.Id(), tcdns.UpdateDnsZoneRequest{
		Description: d.Get("description").(string),
		Labels:      convert.ConvertToMap(d.Get("labels")),
		Annotations: convert.ConvertToMap(d.Get("annotations")),
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("updating DNS zone: %w", err))
	}

	return resourceDnsZoneRead(ctx, d, m)
}

func resourceDnsZoneDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := client.DNS().DeleteZone(ctx, d.Id()); err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("deleting DNS zone: %w", err))
	}

	d.SetId("")
	return nil
}
