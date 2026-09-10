package iam

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	iam "github.com/thalassa-cloud/client-go/iam"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

var (
	federatedIdentityAudienceMatchModes = []string{
		string(iam.AudienceMatchModeExact),
		string(iam.AudienceMatchModeAny),
		string(iam.AudienceMatchModeAll),
	}
	federatedIdentityAllowedScopes = []string{
		string(iam.AccessCredentialsScopeAPIRead),
		string(iam.AccessCredentialsScopeAPIWrite),
		string(iam.AccessCredentialsScopeKubernetes),
		string(iam.AccessCredentialsScopeObjectStorage),
	}
	federatedIdentityStatuses = []string{
		string(iam.FederatedIdentityStatusActive),
		string(iam.FederatedIdentityStatusInactive),
		string(iam.FederatedIdentityStatusExpired),
		string(iam.FederatedIdentityStatusRevoked),
	}
)

func ResourceFederatedIdentity() *schema.Resource {
	return &schema.Resource{
		Description:   "Create and manage a federated identity in Thalassa Cloud",
		CreateContext: resourceFederatedIdentityCreate,
		ReadContext:   resourceFederatedIdentityRead,
		UpdateContext: resourceFederatedIdentityUpdate,
		DeleteContext: resourceFederatedIdentityDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity of the federated identity",
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Organisation of the federated identity. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.StringLenBetween(1, 255),
				Description:  "Name of the federated identity",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				ValidateFunc: validate.StringLenBetween(0, 1024),
				Description:  "Human-readable description of the federated identity",
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Labels for the federated identity",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Annotations for the federated identity",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"service_account_identity": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Identity of the service account to bind",
			},
			"provider_identity": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Identity of the federated identity provider",
			},
			"provider_subject": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StringLenBetween(1, 2048),
				Description:  "Subject identifier from the OIDC provider (sub claim)",
			},
			"trusted_audiences": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of trusted audiences (may be empty)",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"audience_match_mode": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.StringInSlice(federatedIdentityAudienceMatchModes, false),
				Description:  "How audience matching is performed: exact, any, or all",
			},
			"allowed_scopes": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Scopes the federated identity is allowed to access",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validate.StringInSlice(federatedIdentityAllowedScopes, false),
				},
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.StringInSlice(federatedIdentityStatuses, false),
				Description:  "Status of the federated identity",
			},
			"expires_at": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Expiration timestamp (RFC3339). Omit for no expiration.",
			},
			"conditions": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional conditions/claims matcher rules as a JSON object",
			},
			"last_used_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last used timestamp (RFC3339)",
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
				Description: "Object version of the federated identity",
			},
		},
	}
}

func expandAllowedScopes(v any) []iam.AccessCredentialsScope {
	raw := convert.ConvertToStringSlice(v)
	scopes := make([]iam.AccessCredentialsScope, 0, len(raw))
	for _, s := range raw {
		scopes = append(scopes, iam.AccessCredentialsScope(s))
	}
	return scopes
}

func flattenAllowedScopes(scopes []iam.AccessCredentialsScope) []string {
	result := make([]string, 0, len(scopes))
	for _, s := range scopes {
		result = append(result, string(s))
	}
	return result
}

func parseConditions(raw string) (map[string]any, error) {
	if raw == "" {
		return nil, nil
	}
	var conditions map[string]any
	if err := json.Unmarshal([]byte(raw), &conditions); err != nil {
		return nil, fmt.Errorf("invalid conditions JSON: %w", err)
	}
	return conditions, nil
}

func flattenConditions(conditions map[string]any) (string, error) {
	if conditions == nil {
		return "", nil
	}
	b, err := json.Marshal(conditions)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func parseExpiresAt(raw string) (*time.Time, error) {
	if raw == "" {
		return nil, nil
	}
	parsed, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return nil, fmt.Errorf("invalid expires_at, expected RFC3339: %w", err)
	}
	return &parsed, nil
}

func resourceFederatedIdentityCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	expiresAt, err := parseExpiresAt(d.Get("expires_at").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	conditions, err := parseConditions(d.Get("conditions").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	createReq := iam.CreateFederatedIdentityRequest{
		Name:                   d.Get("name").(string),
		Description:            d.Get("description").(string),
		Labels:                 convert.ConvertToMap(d.Get("labels")),
		Annotations:            convert.ConvertToMap(d.Get("annotations")),
		ServiceAccountIdentity: d.Get("service_account_identity").(string),
		ProviderIdentity:       d.Get("provider_identity").(string),
		ProviderSubject:        d.Get("provider_subject").(string),
		TrustedAudiences:       convert.ConvertToStringSlice(d.Get("trusted_audiences")),
		AudienceMatchMode:      iam.AudienceMatchMode(d.Get("audience_match_mode").(string)),
		AllowedScopes:          expandAllowedScopes(d.Get("allowed_scopes")),
		ExpiresAt:              expiresAt,
		Conditions:             conditions,
	}

	identity, err := client.IAM().CreateFederatedIdentity(ctx, createReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating federated identity: %w", err))
	}
	if identity == nil {
		return diag.FromErr(fmt.Errorf("error creating federated identity: empty response"))
	}

	d.SetId(identity.Identity)

	if status := d.Get("status").(string); status != "" && status != string(identity.Status) {
		updateReq := iam.UpdateFederatedIdentityRequest{
			Status: iam.FederatedIdentityStatus(status),
		}
		if _, err := client.IAM().UpdateFederatedIdentity(ctx, identity.Identity, updateReq); err != nil {
			return diag.FromErr(fmt.Errorf("error setting federated identity status after create: %w", err))
		}
	}

	return resourceFederatedIdentityRead(ctx, d, m)
}

func resourceFederatedIdentityRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	identity, err := client.IAM().GetFederatedIdentity(ctx, d.Id())
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error getting federated identity: %w", err))
	}
	if identity == nil {
		d.SetId("")
		return nil
	}

	return setFederatedIdentityState(d, identity)
}

func setFederatedIdentityState(d *schema.ResourceData, identity *iam.FederatedIdentity) diag.Diagnostics {
	d.SetId(identity.Identity)
	_ = d.Set("name", identity.Name)
	_ = d.Set("description", identity.Description)
	_ = d.Set("labels", identity.Labels)
	_ = d.Set("annotations", identity.Annotations)
	_ = d.Set("provider_subject", identity.ProviderSubject)
	trustedAudiences := identity.TrustedAudiences
	if trustedAudiences == nil {
		trustedAudiences = []string{}
	}
	_ = d.Set("trusted_audiences", trustedAudiences)
	_ = d.Set("audience_match_mode", string(identity.AudienceMatchMode))
	_ = d.Set("allowed_scopes", flattenAllowedScopes(identity.AllowedScopes))
	_ = d.Set("status", string(identity.Status))
	_ = d.Set("created_at", identity.CreatedAt.Format(TimeFormatRFC3339))
	_ = d.Set("object_version", identity.ObjectVersion)

	if identity.ServiceAccount != nil && identity.ServiceAccount.Identity != "" {
		_ = d.Set("service_account_identity", identity.ServiceAccount.Identity)
	}
	if identity.Provider != nil && identity.Provider.Identity != "" {
		_ = d.Set("provider_identity", identity.Provider.Identity)
	}
	if identity.UpdatedAt != nil {
		_ = d.Set("updated_at", identity.UpdatedAt.Format(TimeFormatRFC3339))
	}
	if identity.LastUsedAt != nil {
		_ = d.Set("last_used_at", identity.LastUsedAt.Format(TimeFormatRFC3339))
	} else {
		_ = d.Set("last_used_at", "")
	}
	if identity.ExpiresAt != nil {
		_ = d.Set("expires_at", identity.ExpiresAt.Format(TimeFormatRFC3339))
	} else {
		_ = d.Set("expires_at", "")
	}

	conditionsJSON, err := flattenConditions(identity.Conditions)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error serializing conditions: %w", err))
	}
	_ = d.Set("conditions", conditionsJSON)

	return nil
}

func resourceFederatedIdentityUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	expiresAt, err := parseExpiresAt(d.Get("expires_at").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	conditions, err := parseConditions(d.Get("conditions").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	updateReq := iam.UpdateFederatedIdentityRequest{
		Name:              d.Get("name").(string),
		Description:       d.Get("description").(string),
		Labels:            convert.ConvertToMap(d.Get("labels")),
		Annotations:       convert.ConvertToMap(d.Get("annotations")),
		TrustedAudiences:  convert.ConvertToStringSlice(d.Get("trusted_audiences")),
		AudienceMatchMode: iam.AudienceMatchMode(d.Get("audience_match_mode").(string)),
		AllowedScopes:     expandAllowedScopes(d.Get("allowed_scopes")),
		ExpiresAt:         expiresAt,
		Conditions:        conditions,
	}
	if status := d.Get("status").(string); status != "" {
		updateReq.Status = iam.FederatedIdentityStatus(status)
	}

	identity, err := client.IAM().UpdateFederatedIdentity(ctx, d.Id(), updateReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating federated identity: %w", err))
	}
	if identity != nil {
		if diags := setFederatedIdentityState(d, identity); diags != nil {
			return diags
		}
	}

	return resourceFederatedIdentityRead(ctx, d, m)
}

func resourceFederatedIdentityDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := client.IAM().DeleteFederatedIdentity(ctx, d.Id()); err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error deleting federated identity: %w", err))
	}

	d.SetId("")
	return nil
}
