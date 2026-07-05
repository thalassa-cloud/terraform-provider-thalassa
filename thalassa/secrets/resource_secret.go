package secrets

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
	tcsecrets "github.com/thalassa-cloud/client-go/secrets"

	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func ResourceSecret() *schema.Resource {
	return &schema.Resource{
		Description:   "Create and manage a secret in Thalassa Secrets Manager",
		CreateContext: resourceSecretCreate,
		ReadContext:   resourceSecretRead,
		UpdateContext: resourceSecretUpdate,
		DeleteContext: resourceSecretDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSecretImport,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Composite ID: {region}{path} (e.g. nl-01/app/prod/db/password).",
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
				Description: "Region slug where the secret is stored.",
			},
			"path": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     validateSecretPath,
				Description:      "Absolute secret path (must start with /).",
				DiffSuppressFunc: func(_, old, new string, _ *schema.ResourceData) bool { return old == new },
			},
			"kms_key_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "KMS key identity used to encrypt the secret.",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.StringLenBetween(0, 1024),
				Description:  "Human-readable description.",
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
			"generate_secret": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ForceNew:      true,
				ConflictsWith: []string{"secret_string", "secret_key_values"},
				Description:   "Generate a random secret value on create. Mutually exclusive with secret_string and secret_key_values.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"byte_length": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Byte length of the generated secret.",
						},
					},
				},
			},
			"secret_string": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				ForceNew:      true,
				ConflictsWith: []string{"generate_secret", "secret_key_values"},
				Description:   "Initial secret string value (create only; not returned on read).",
			},
			"secret_key_values": {
				Type:          schema.TypeMap,
				Optional:      true,
				Sensitive:     true,
				ForceNew:      true,
				ConflictsWith: []string{"generate_secret", "secret_string"},
				Elem:          &schema.Schema{Type: schema.TypeString},
				Description:   "Initial key-value secret payload (create only; not returned on read).",
			},
			"current_version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Current active secret version.",
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_accessed_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSecretImport(ctx context.Context, d *schema.ResourceData, m any) ([]*schema.ResourceData, error) {
	region, path, err := parseSecretID(d.Id())
	if err != nil {
		return nil, err
	}
	d.Set("region", region)
	d.Set("path", path)
	return []*schema.ResourceData{d}, nil
}

func resourceSecretCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	region := d.Get("region").(string)
	path := d.Get("path").(string)

	createReq := tcsecrets.CreateSecretRequest{
		Path:           path,
		Description:    d.Get("description").(string),
		Labels:         convert.ConvertToMap(d.Get("labels")),
		Annotations:    convert.ConvertToMap(d.Get("annotations")),
		KmsKeyIdentity: d.Get("kms_key_id").(string),
	}

	if v, ok := d.GetOk("secret_string"); ok {
		createReq.SecretString = tcsecrets.EncodeBytes([]byte(v.(string)))
	}
	if v, ok := d.GetOk("secret_key_values"); ok {
		createReq.SecretKeyValues = convert.ConvertToMap(v)
	}
	if v, ok := d.GetOk("generate_secret"); ok {
		createReq.GenerateSecret = expandGenerateSecret(v.([]any))
	}

	secret, err := client.Secrets().CreateSecret(ctx, region, createReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf("creating secret: %w", err))
	}

	d.SetId(secretID(region, secret.Path))
	return resourceSecretRead(ctx, d, m)
}

func resourceSecretRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	region, path, err := parseSecretID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	secret, err := client.Secrets().GetSecret(ctx, region, path, false)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("reading secret: %w", err))
	}

	if err := setSecretState(d, secret, region); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSecretUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	if !d.HasChanges("description", "labels", "annotations") {
		return resourceSecretRead(ctx, d, m)
	}

	return diag.Diagnostics{{
		Severity: diag.Warning,
		Summary:  "Secret metadata update not supported",
		Detail:   "Only description, labels, and annotations were changed but the Secrets Manager API does not support metadata updates. No changes were applied.",
	}}
}

func resourceSecretDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	region, path, err := parseSecretID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err := client.Secrets().DeleteSecret(ctx, region, path); err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("deleting secret: %w", err))
	}

	d.SetId("")
	return nil
}
