package kms

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	tckms "github.com/thalassa-cloud/client-go/kms"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"

	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func ResourceKmsKey() *schema.Resource {
	return &schema.Resource{
		Description:   "Create and manage a KMS key in Thalassa Cloud",
		CreateContext: resourceKmsKeyCreate,
		ReadContext:   resourceKmsKeyRead,
		UpdateContext: resourceKmsKeyUpdate,
		DeleteContext: resourceKmsKeyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Platform identity of the KMS key.",
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Organisation ID. Defaults to the provider organisation.",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Region slug where the key is stored (e.g. nl-01).",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StringLenBetween(1, 255),
				Description:  "Name of the KMS key.",
			},
			"slug": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL-safe slug of the KMS key.",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.StringLenBetween(0, 1024),
				Description:  "Human-readable description.",
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Labels for the key.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Annotations for the key.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"key_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StringInSlice(kmsKeyTypes, false),
				Description:  "Key algorithm and type.",
			},
			"export_allowed": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Whether the key material may be exported.",
			},
			"key_rotation_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether automatic key rotation is enabled.",
			},
			"rotation_period_in_days": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validate.IntAtLeast(1),
				Description:  "Automatic rotation period in days.",
			},
			"import_key_material": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				ForceNew:    true,
				Description: "Base64-encoded key material for BYOK import.",
			},
			"hash_function": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Hash function used when importing key material.",
			},
			"allow_rotation": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Whether imported keys may be rotated.",
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.StringInSlice(kmsKeyStatuses, false),
				Description:  "Desired key status: active or disabled. pending_deletion is set by the platform after delete is requested.",
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
			"object_version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Platform object version for optimistic concurrency.",
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
			"cancel_scheduled_deletion": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "When true, cancels a pending deletion on the next apply.",
			},
		},
	}
}

func resourceKmsKeyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	region := d.Get("region").(string)

	summary, err := client.KMS().GetSummary(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("checking KMS availability: %w", err))
	}
	if !regionKmsAvailable(summary, region) {
		return diag.Errorf("KMS is not available in region %q for this organisation", region)
	}

	createReq := tckms.CreateKmsKeyRequest{
		Name:               d.Get("name").(string),
		Description:        d.Get("description").(string),
		Labels:             convert.ConvertToMap(d.Get("labels")),
		Annotations:        convert.ConvertToMap(d.Get("annotations")),
		KeyType:            tckms.KmsKeyType(d.Get("key_type").(string)),
		ExportAllowed:      d.Get("export_allowed").(bool),
		KeyRotationEnabled: d.Get("key_rotation_enabled").(bool),
	}

	if v, ok := d.GetOk("rotation_period_in_days"); ok {
		days := v.(int)
		createReq.RotationPeriodInDays = &days
	}
	if v, ok := d.GetOk("import_key_material"); ok {
		createReq.ImportKeyMaterial = v.(string)
	}
	if v, ok := d.GetOk("hash_function"); ok {
		createReq.HashFunction = v.(string)
	}
	if v, ok := d.GetOk("allow_rotation"); ok {
		createReq.AllowRotation = v.(bool)
	}

	key, err := client.KMS().CreateKey(ctx, region, createReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf("creating KMS key: %w", err))
	}

	d.SetId(key.Identity)

	if status := d.Get("status").(string); status == string(tckms.KmsKeyStatusDisabled) {
		key, err = client.KMS().DisableKey(ctx, region, key.Identity)
		if err != nil {
			return diag.FromErr(fmt.Errorf("disabling KMS key after create: %w", err))
		}
	}

	return resourceKmsKeyReadWithKey(ctx, d, m, region, key)
}

func resourceKmsKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	region := d.Get("region").(string)
	identity := d.Id()

	key, err := client.KMS().GetKey(ctx, region, identity)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("reading KMS key: %w", err))
	}

	return resourceKmsKeyReadWithKey(ctx, d, m, region, key)
}

func resourceKmsKeyReadWithKey(ctx context.Context, d *schema.ResourceData, m interface{}, region string, key *tckms.KmsKey) diag.Diagnostics {
	if err := setKmsKeyState(d, key, region); err != nil {
		return diag.FromErr(err)
	}

	desiredStatus := d.Get("status").(string)
	if desiredStatus == "" || key.Status == tckms.KmsKeyStatusPendingDeletion {
		if err := d.Set("status", string(key.Status)); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func resourceKmsKeyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	region := d.Get("region").(string)
	identity := d.Id()

	if d.Get("cancel_scheduled_deletion").(bool) {
		if err := client.KMS().CancelDeletion(ctx, region, identity); err != nil {
			return diag.FromErr(fmt.Errorf("cancelling KMS key deletion: %w", err))
		}
		if err := d.Set("cancel_scheduled_deletion", false); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("key_rotation_enabled") || d.HasChange("rotation_period_in_days") {
		update := tckms.UpdateRotationRequest{}
		if v, ok := d.GetOk("key_rotation_enabled"); ok {
			enabled := v.(bool)
			update.KeyRotationEnabled = &enabled
		}
		if v, ok := d.GetOk("rotation_period_in_days"); ok {
			days := v.(int)
			update.RotationPeriodInDays = &days
		}
		if _, err := client.KMS().UpdateRotation(ctx, region, identity, update); err != nil {
			return diag.FromErr(fmt.Errorf("updating KMS key rotation: %w", err))
		}
	}

	if d.HasChange("status") {
		switch d.Get("status").(string) {
		case string(tckms.KmsKeyStatusActive):
			if _, err := client.KMS().EnableKey(ctx, region, identity); err != nil {
				return diag.FromErr(fmt.Errorf("enabling KMS key: %w", err))
			}
		case string(tckms.KmsKeyStatusDisabled):
			if _, err := client.KMS().DisableKey(ctx, region, identity); err != nil {
				return diag.FromErr(fmt.Errorf("disabling KMS key: %w", err))
			}
		}
	}

	if d.HasChanges("description", "labels", "annotations") {
		// Metadata updates are not exposed by the KMS API; only rotation and status are mutable.
	}

	return resourceKmsKeyRead(ctx, d, m)
}

func resourceKmsKeyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	region := d.Get("region").(string)
	identity := d.Id()

	if err := client.KMS().DeleteKey(ctx, region, identity); err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("scheduling KMS key deletion: %w", err))
	}

	d.SetId("")
	return nil
}
