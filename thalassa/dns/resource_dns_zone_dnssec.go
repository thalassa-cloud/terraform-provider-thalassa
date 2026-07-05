package dns

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	tcdns "github.com/thalassa-cloud/client-go/dns"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"

	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func ResourceDnsZoneDnssec() *schema.Resource {
	return &schema.Resource{
		Description:   "Enable DNSSEC signing for a DNS zone",
		CreateContext: resourceDnsZoneDnssecCreate,
		ReadContext:   resourceDnsZoneDnssecRead,
		DeleteContext: resourceDnsZoneDnssecDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Zone identity (same as zone_id).",
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
				Description: "Platform identity of the DNS zone.",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "KMS region used for DNSSEC signing.",
			},
			"kms_key_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "KMS key ID for signing. Leave blank to auto-provision.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether DNSSEC signing is enabled.",
			},
			"ds_delegated": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether DS records are delegated at the parent.",
			},
			"ds_records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"record":           {Type: schema.TypeString, Computed: true},
						"digest_type_name": {Type: schema.TypeString, Computed: true},
						"key_tag":          {Type: schema.TypeInt, Computed: true},
						"algorithm":        {Type: schema.TypeInt, Computed: true},
						"key_role":         {Type: schema.TypeString, Computed: true},
						"public_key":       {Type: schema.TypeString, Computed: true},
					},
				},
			},
			"last_signed_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_sign_error": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"next_ds_probe_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDnsZoneDnssecCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	zoneID := d.Get("zone_id").(string)
	req := tcdns.SetDnssecRequest{
		Region: d.Get("region").(string),
	}
	if v, ok := d.GetOk("kms_key_id"); ok {
		req.KmsKeyIdentity = v.(string)
	}

	if _, err := client.DNS().SetDnssec(ctx, zoneID, req); err != nil {
		return diag.FromErr(fmt.Errorf("enabling DNSSEC: %w", err))
	}

	d.SetId(zoneID)
	return resourceDnsZoneDnssecRead(ctx, d, m)
}

func resourceDnsZoneDnssecRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	zoneID := d.Id()
	if zoneID == "" {
		zoneID = d.Get("zone_id").(string)
	}

	status, err := client.DNS().GetDnssec(ctx, zoneID)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("reading DNSSEC status: %w", err))
	}

	if !status.Enabled {
		d.SetId("")
		return nil
	}

	d.SetId(zoneID)
	if err := d.Set("zone_id", zoneID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("enabled", status.Enabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ds_delegated", status.DsDelegated); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ds_records", flattenDsRecords(status.DsRecords)); err != nil {
		return diag.FromErr(err)
	}
	if status.Region != "" {
		if err := d.Set("region", status.Region); err != nil {
			return diag.FromErr(err)
		}
	}
	if status.KmsKeyIdentity != "" {
		if err := d.Set("kms_key_id", status.KmsKeyIdentity); err != nil {
			return diag.FromErr(err)
		}
	}
	if status.LastSignedAt != nil {
		if err := d.Set("last_signed_at", status.LastSignedAt.Format(timeFormatRFC3339)); err != nil {
			return diag.FromErr(err)
		}
	}
	if status.LastSignError != nil {
		if err := d.Set("last_sign_error", *status.LastSignError); err != nil {
			return diag.FromErr(err)
		}
	}
	if status.NextDsProbeAt != nil {
		if err := d.Set("next_ds_probe_at", status.NextDsProbeAt.Format(timeFormatRFC3339)); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func resourceDnsZoneDnssecDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := client.DNS().DeleteDnssec(ctx, d.Id()); err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("disabling DNSSEC: %w", err))
	}

	d.SetId("")
	return nil
}
