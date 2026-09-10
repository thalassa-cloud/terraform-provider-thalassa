package iam

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	iam "github.com/thalassa-cloud/client-go/iam"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func ResourceFederatedIdentityProvider() *schema.Resource {
	return &schema.Resource{
		Description:   "Create and manage a federated identity provider in Thalassa Cloud",
		CreateContext: resourceFederatedIdentityProviderCreate,
		ReadContext:   resourceFederatedIdentityProviderRead,
		UpdateContext: resourceFederatedIdentityProviderUpdate,
		DeleteContext: resourceFederatedIdentityProviderDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity of the federated identity provider",
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Organisation of the provider. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.StringLenBetween(1, 255),
				Description:  "Name of the federated identity provider",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				ValidateFunc: validate.StringLenBetween(0, 1024),
				Description:  "Human-readable description of the provider",
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Labels for the provider",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Annotations for the provider",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"provider_issuer": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StringLenBetween(1, 2048),
				Description:  "OIDC issuer URL of the identity provider (iss claim)",
			},
			"provider_jwks_uri": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional JWKS URI override for the provider",
			},
			"local_jwks": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				ForceNew:    true,
				Description: "Optional locally stored JWKS as JSON object with a keys array, e.g. {\"keys\":[...]}",
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.StringInSlice([]string{string(iam.FederatedIdentityProviderStatusActive), string(iam.FederatedIdentityProviderStatusInactive)}, false),
				Description:  "Status of the provider: active or inactive",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation timestamp (RFC3339)",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last update timestamp (RFC3339)",
			},
			"object_version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Object version of the provider",
			},
		},
	}
}

func parseLocalJWKS(raw string) (*iam.LocalJWKS, error) {
	if raw == "" {
		return nil, nil
	}
	var jwks iam.LocalJWKS
	if err := json.Unmarshal([]byte(raw), &jwks); err != nil {
		return nil, fmt.Errorf("invalid local_jwks JSON: %w", err)
	}
	if jwks.Keys == nil {
		return nil, fmt.Errorf("local_jwks must contain a keys array")
	}
	return &jwks, nil
}

func resourceFederatedIdentityProviderCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	createReq := iam.CreateFederatedIdentityProviderRequest{
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		Labels:         convert.ConvertToMap(d.Get("labels")),
		Annotations:    convert.ConvertToMap(d.Get("annotations")),
		ProviderIssuer: d.Get("provider_issuer").(string),
	}

	if v, ok := d.GetOk("provider_jwks_uri"); ok {
		uri := v.(string)
		createReq.ProviderJwksURI = &uri
	}

	localJWKS, err := parseLocalJWKS(d.Get("local_jwks").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	createReq.LocalJWKS = localJWKS

	status := d.Get("status").(string)
	if status == "" {
		status = string(iam.FederatedIdentityProviderStatusActive)
	}
	createReq.Status = iam.FederatedIdentityProviderStatus(status)

	providerObj, err := client.IAM().CreateFederatedIdentityProvider(ctx, createReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating federated identity provider: %w", err))
	}
	if providerObj == nil {
		return diag.FromErr(fmt.Errorf("error creating federated identity provider: empty response"))
	}

	d.SetId(providerObj.Identity)
	return resourceFederatedIdentityProviderRead(ctx, d, m)
}

func resourceFederatedIdentityProviderRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	providerObj, err := client.IAM().GetFederatedIdentityProvider(ctx, d.Id())
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error getting federated identity provider: %w", err))
	}
	if providerObj == nil {
		d.SetId("")
		return nil
	}

	return setFederatedIdentityProviderState(d, providerObj)
}

func setFederatedIdentityProviderState(d *schema.ResourceData, providerObj *iam.FederatedIdentityProvider) diag.Diagnostics {
	d.SetId(providerObj.Identity)
	_ = d.Set("name", providerObj.Name)
	_ = d.Set("description", providerObj.Description)
	_ = d.Set("labels", providerObj.Labels)
	_ = d.Set("annotations", providerObj.Annotations)
	_ = d.Set("provider_issuer", providerObj.ProviderIssuer)
	_ = d.Set("status", string(providerObj.Status))
	_ = d.Set("created_at", providerObj.CreatedAt.Format(TimeFormatRFC3339))
	_ = d.Set("object_version", providerObj.ObjectVersion)
	if providerObj.UpdatedAt != nil {
		_ = d.Set("updated_at", providerObj.UpdatedAt.Format(TimeFormatRFC3339))
	}
	if providerObj.ProviderJwksURI != nil {
		_ = d.Set("provider_jwks_uri", *providerObj.ProviderJwksURI)
	} else {
		_ = d.Set("provider_jwks_uri", "")
	}
	return nil
}

func resourceFederatedIdentityProviderUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	updateReq := iam.UpdateFederatedIdentityProviderRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Labels:      convert.ConvertToMap(d.Get("labels")),
		Annotations: convert.ConvertToMap(d.Get("annotations")),
	}

	if d.HasChange("provider_jwks_uri") {
		uri := d.Get("provider_jwks_uri").(string)
		updateReq.ProviderJwksURI = &uri
	}

	if status := d.Get("status").(string); status != "" {
		updateReq.Status = iam.FederatedIdentityProviderStatus(status)
	}

	providerObj, err := client.IAM().UpdateFederatedIdentityProvider(ctx, d.Id(), updateReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating federated identity provider: %w", err))
	}
	if providerObj != nil {
		if diags := setFederatedIdentityProviderState(d, providerObj); diags != nil {
			return diags
		}
	}

	return resourceFederatedIdentityProviderRead(ctx, d, m)
}

func resourceFederatedIdentityProviderDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := client.IAM().DeleteFederatedIdentityProvider(ctx, d.Id()); err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error deleting federated identity provider: %w", err))
	}

	d.SetId("")
	return nil
}
