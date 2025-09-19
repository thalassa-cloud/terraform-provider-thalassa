package iam

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"

	iam "github.com/thalassa-cloud/client-go/iam"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

func ResourceServiceAccountAccessCredential() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage access credentials for a service account in Thalassa Cloud",
		CreateContext: resourceServiceAccountAccessCredentialCreate,
		ReadContext:   resourceServiceAccountAccessCredentialRead,
		DeleteContext: resourceServiceAccountAccessCredentialDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_account_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Identity of the service account",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StringLenBetween(1, 255),
				Description:  "Name of the access credential",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validate.StringLenBetween(0, 255),
				Description:  "Description of the access credential",
			},
			"expires_at": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Expiration timestamp of the access credential (RFC3339 format)",
			},
			"scopes": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "List of scopes for the access credential",
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validate.StringInSlice([]string{
						"api:read",
						"api:write",
						"kubernetes",
						"objectStorage",
					}, false),
				},
			},
			"access_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Access key for the credential",
			},
			"access_secret": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Access secret for the credential",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation timestamp of the access credential",
			},
			"last_used_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last used timestamp of the access credential",
			},
		},
	}
}

func resourceServiceAccountAccessCredentialCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	serviceAccountID := d.Get("service_account_id").(string)

	// Convert scopes list
	scopes := make([]iam.AccessCredentialsScope, 0)
	if s, ok := d.Get("scopes").([]interface{}); ok {
		for _, scope := range s {
			scopes = append(scopes, iam.AccessCredentialsScope(scope.(string)))
		}
	}

	// Parse expires_at if provided
	var expiresAt *time.Time
	if expiresAtStr, ok := d.Get("expires_at").(string); ok && expiresAtStr != "" {
		parsed, err := time.Parse(time.RFC3339, expiresAtStr)
		if err != nil {
			return diag.FromErr(fmt.Errorf("invalid expires_at format: %s", err))
		}
		expiresAt = &parsed
	}

	createReq := iam.CreateServiceAccountAccessCredentialRequest{
		Name:        d.Get("name").(string),
		Description: convert.Ptr(d.Get("description").(string)),
		ExpiresAt:   expiresAt,
		Scopes:      scopes,
	}

	credential, err := client.IAM().CreateServiceAccountAccessCredentials(ctx, serviceAccountID, createReq)
	if err != nil {
		return diag.FromErr(err)
	}

	if credential != nil {
		d.SetId(credential.Identity)
		d.Set("access_key", credential.AccessKey)
		d.Set("access_secret", credential.AccessSecret)
		return resourceServiceAccountAccessCredentialRead(ctx, d, m)
	}

	return diag.FromErr(fmt.Errorf("failed to create service account access credential"))
}

func resourceServiceAccountAccessCredentialRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	serviceAccountID := d.Get("service_account_id").(string)
	credentialID := d.Get("id").(string)

	// Get all access credentials for the service account
	credentials, err := client.IAM().GetServiceAccountAccessCredentials(ctx, serviceAccountID)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	// Find the specific credential
	var credential *iam.ServiceAccountAccessCredential
	for _, c := range credentials {
		if c.Identity == credentialID {
			credential = &c
			break
		}
	}

	if credential == nil {
		d.SetId("")
		return nil
	}

	d.SetId(credential.Identity)
	d.Set("name", credential.Name)
	d.Set("description", credential.Description)
	d.Set("access_key", credential.AccessKey)
	d.Set("created_at", credential.CreatedAt.Format(TimeFormatRFC3339))

	if credential.LastUsedAt != nil {
		d.Set("last_used_at", credential.LastUsedAt.Format(TimeFormatRFC3339))
	}
	if credential.ExpiresAt != nil {
		d.Set("expires_at", credential.ExpiresAt.Format(TimeFormatRFC3339))
	}

	return nil
}

func resourceServiceAccountAccessCredentialDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	serviceAccountID := d.Get("service_account_id").(string)
	credentialID := d.Get("id").(string)

	err = client.IAM().DeleteServiceAccountAccessCredentials(ctx, serviceAccountID, credentialID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
