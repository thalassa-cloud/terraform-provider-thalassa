package kms

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	tckms "github.com/thalassa-cloud/client-go/kms"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"

	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func DataSourceKmsKey() *schema.Resource {
	return &schema.Resource{
		Description: "Look up a KMS key by region and identity or name filter",
		ReadContext: dataSourceKmsKeyRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Platform identity of the KMS key.",
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Organisation ID. Defaults to the provider organisation.",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Region slug where the key is stored.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"id", "name"},
				Description:  "Key name (used when id is not set).",
			},
			"slug": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL-safe slug of the KMS key.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Human-readable description.",
			},
			"labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Labels for the key.",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Annotations for the key.",
			},
			"key_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Key algorithm and type.",
			},
			"export_allowed": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the key material may be exported.",
			},
			"key_rotation_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether automatic key rotation is enabled.",
			},
			"rotation_period_in_days": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Automatic rotation period in days.",
			},
			"status": {
				Type:         schema.TypeString,
				Computed:     true,
				ValidateFunc: validate.StringInSlice(append(kmsKeyStatuses, string(tckms.KmsKeyStatusPendingDeletion)), false),
				Description:  "Current key status.",
			},
			"imported": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the key was imported (BYOK).",
			},
			"latest_version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Latest key version number.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation timestamp (RFC3339).",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last update timestamp (RFC3339).",
			},
		},
	}
}

func dataSourceKmsKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	region := d.Get("region").(string)

	var keyIdentity string
	if v, ok := d.GetOk("id"); ok && v.(string) != "" {
		keyIdentity = v.(string)
	} else {
		name := d.Get("name").(string)
		keys, err := client.KMS().ListKeys(ctx, region, nil)
		if err != nil {
			return diag.FromErr(fmt.Errorf("listing KMS keys: %w", err))
		}
		for _, key := range keys {
			if key.Name == name {
				keyIdentity = key.Identity
				break
			}
		}
		if keyIdentity == "" {
			return diag.Errorf("KMS key with name %q not found in region %q", name, region)
		}
	}

	key, err := client.KMS().GetKey(ctx, region, keyIdentity)
	if err != nil {
		if tcclient.IsNotFound(err) {
			return diag.Errorf("KMS key %q not found in region %q", keyIdentity, region)
		}
		return diag.FromErr(fmt.Errorf("reading KMS key: %w", err))
	}

	d.SetId(key.Identity)
	if err := setKmsKeyState(d, key, region); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
