package secrets

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
	tcsecrets "github.com/thalassa-cloud/client-go/secrets"

	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func ResourceSecretAccessPolicy() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage access policy for a secret in Thalassa Secrets Manager",
		CreateContext: resourceSecretAccessPolicyCreate,
		ReadContext:   resourceSecretAccessPolicyRead,
		UpdateContext: resourceSecretAccessPolicyUpdate,
		DeleteContext: resourceSecretAccessPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSecretImport,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Composite ID: {region}{path} (same as thalassa_secret).",
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
			"statement": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Access policy statements.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"effect": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.StringInSlice([]string{"Allow", "Deny"}, false),
						},
						"actions": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"principals": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func resourceSecretAccessPolicyCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	return resourceSecretAccessPolicyUpdate(ctx, d, m)
}

func resourceSecretAccessPolicyRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
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
		return diag.FromErr(fmt.Errorf("reading secret access policy: %w", err))
	}

	_ = d.Set("region", region)
	_ = d.Set("path", path)

	if secret.AccessPolicy != nil {
		_ = d.Set("statement", flattenAccessPolicyStatements(secret.AccessPolicy.Statements))
	}

	return nil
}

func resourceSecretAccessPolicyUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	region := d.Get("region").(string)
	path := d.Get("path").(string)

	statements := expandAccessPolicyStatements(d.Get("statement").([]any))
	_, err = client.Secrets().UpdateAccessPolicy(ctx, region, path, tcsecrets.UpdateAccessPolicyRequest{
		AccessPolicy: tcsecrets.SecretPolicy{Statements: statements},
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("updating secret access policy: %w", err))
	}

	d.SetId(secretID(region, path))
	return resourceSecretAccessPolicyRead(ctx, d, m)
}

func resourceSecretAccessPolicyDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	region, path, err := parseSecretID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.Secrets().UpdateAccessPolicy(ctx, region, path, tcsecrets.UpdateAccessPolicyRequest{
		AccessPolicy: tcsecrets.SecretPolicy{Statements: []tcsecrets.SecretPolicyStatement{}},
	})
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("clearing secret access policy: %w", err))
	}

	d.SetId("")
	return nil
}
