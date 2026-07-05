package dns

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	tcdns "github.com/thalassa-cloud/client-go/dns"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"

	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func ResourceDnsRecord() *schema.Resource {
	return &schema.Resource{
		Description:   "Create and manage a DNS record in a Thalassa DNS zone",
		CreateContext: resourceDnsRecordCreate,
		ReadContext:   resourceDnsRecordRead,
		UpdateContext: resourceDnsRecordUpdate,
		DeleteContext: resourceDnsRecordDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Platform identity of the DNS record (dnsr-…).",
			},
			"organisation_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Platform identity of the DNS zone (dnsz-…).",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Record name relative to the zone (@ for apex, www, *, etc.).",
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StringInSlice(dnsRecordTypes, false),
			},
			"ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     3600,
				Description: "Time to live in seconds.",
			},
			"values": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				Description: "Record values. Format depends on record type.",
				Elem:        &schema.Schema{Type: schema.TypeString},
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

func resourceDnsRecordCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	zoneID := d.Get("zone_id").(string)
	record, err := client.DNS().CreateRecord(ctx, zoneID, tcdns.CreateDnsRecordRequest{
		Name:   d.Get("name").(string),
		Type:   tcdns.DnsRecordType(d.Get("type").(string)),
		TTL:    d.Get("ttl").(int),
		Values: convert.ConvertToStringSlice(d.Get("values")),
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("creating DNS record: %w", err))
	}

	d.SetId(record.Identity)
	return resourceDnsRecordReadWithRecord(ctx, d, zoneID, record)
}

func resourceDnsRecordRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	zoneID := d.Get("zone_id").(string)

	record, err := client.DNS().GetRecord(ctx, zoneID, d.Id())
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("reading DNS record: %w", err))
	}

	return resourceDnsRecordReadWithRecord(ctx, d, zoneID, record)
}

func resourceDnsRecordReadWithRecord(_ context.Context, d *schema.ResourceData, zoneID string, record *tcdns.DnsRecord) diag.Diagnostics {
	if err := setDnsRecordState(d, record, zoneID); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceDnsRecordUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	zoneID := d.Get("zone_id").(string)
	recordID := d.Id()

	_, err = client.DNS().UpdateRecord(ctx, zoneID, recordID, tcdns.UpdateDnsRecordRequest{
		TTL:    d.Get("ttl").(int),
		Values: convert.ConvertToStringSlice(d.Get("values")),
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("updating DNS record: %w", err))
	}

	return resourceDnsRecordRead(ctx, d, m)
}

func resourceDnsRecordDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	zoneID := d.Get("zone_id").(string)
	recordID := d.Id()

	if err := client.DNS().DeleteRecord(ctx, zoneID, recordID); err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("deleting DNS record: %w", err))
	}

	d.SetId("")
	return nil
}
