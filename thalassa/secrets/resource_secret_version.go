package secrets

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	tcsecrets "github.com/thalassa-cloud/client-go/secrets"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"

	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func ResourceSecretVersion() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a secret version in Thalassa Secrets Manager",
		CreateContext: resourceSecretVersionCreate,
		ReadContext:   resourceSecretVersionRead,
		DeleteContext: resourceSecretVersionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSecretVersionImport,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Composite ID: {region}{path}/{version} (e.g. nl-01/app/prod/db/password/3).",
			},
			"organisation_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"path": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateSecretPath,
			},
			"version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Version number assigned by the platform.",
			},
			"secret_string": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{"secret_key_values", "generate_secret"},
				Description:   "Secret string value (not returned on read).",
			},
			"secret_key_values": {
				Type:          schema.TypeMap,
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{"secret_string", "generate_secret"},
				Elem:          &schema.Schema{Type: schema.TypeString},
			},
			"generate_secret": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"secret_string", "secret_key_values"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"length": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"character_set": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceSecretVersionImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	region, path, version, err := parseSecretVersionID(d.Id())
	if err != nil {
		return nil, err
	}
	d.Set("region", region)
	d.Set("path", path)
	d.Set("version", version)
	return []*schema.ResourceData{d}, nil
}

func resourceSecretVersionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	region := d.Get("region").(string)
	path := d.Get("path").(string)

	putReq := tcsecrets.PutSecretValueRequest{}
	if v, ok := d.GetOk("secret_string"); ok {
		putReq.SecretString = tcsecrets.EncodeBytes([]byte(v.(string)))
	}
	if v, ok := d.GetOk("secret_key_values"); ok {
		putReq.SecretKeyValues = convert.ConvertToMap(v)
	}
	if v, ok := d.GetOk("generate_secret"); ok {
		putReq.GenerateSecret = expandGenerateSecret(v.([]interface{}))
	}

	if putReq.SecretString == "" && len(putReq.SecretKeyValues) == 0 && putReq.GenerateSecret == nil {
		return diag.Errorf("one of secret_string, secret_key_values, or generate_secret must be set")
	}

	result, err := client.Secrets().PutSecretValue(ctx, region, path, putReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf("creating secret version: %w", err))
	}

	d.SetId(secretVersionID(region, path, result.Version))
	if err := d.Set("version", result.Version); err != nil {
		return diag.FromErr(err)
	}

	return resourceSecretVersionRead(ctx, d, m)
}

func resourceSecretVersionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	region, path, version, err := parseSecretVersionID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	secret, err := client.Secrets().GetSecret(ctx, region, path, true)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("reading secret metadata: %w", err))
	}

	found := false
	for _, v := range secret.Versions {
		if v.Version == version {
			found = true
			break
		}
	}
	if !found {
		d.SetId("")
		return nil
	}

	if err := d.Set("region", region); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("path", path); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("version", version); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSecretVersionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	region, path, version, err := parseSecretVersionID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err := client.Secrets().DestroySecretVersion(ctx, region, path, version); err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("destroying secret version: %w", err))
	}

	d.SetId("")
	return nil
}
