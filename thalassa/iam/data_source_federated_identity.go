package iam

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	iam "github.com/thalassa-cloud/client-go/iam"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func DataSourceFederatedIdentity() *schema.Resource {
	return &schema.Resource{
		Description: "Get a federated identity by identity or name",
		ReadContext: dataSourceFederatedIdentityRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
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
				Optional:     true,
				ValidateFunc: validate.StringLenBetween(1, 255),
				Description:  "Name of the federated identity",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Human-readable description of the federated identity",
			},
			"labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Labels for the federated identity",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"annotations": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Annotations for the federated identity",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"service_account_identity": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity of the bound service account",
			},
			"provider_identity": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity of the federated identity provider",
			},
			"provider_subject": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Subject identifier from the OIDC provider",
			},
			"trusted_audiences": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of trusted audiences",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"audience_match_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "How audience matching is performed",
			},
			"allowed_scopes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Scopes the federated identity is allowed to access",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the federated identity",
			},
			"expires_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Expiration timestamp (RFC3339)",
			},
			"conditions": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Conditions/claims matcher rules as JSON",
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

func dataSourceFederatedIdentityRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	identity := d.Get("id").(string)
	name := d.Get("name").(string)

	if identity == "" && name == "" {
		return diag.FromErr(fmt.Errorf("either 'id' or 'name' must be provided to look up a federated identity"))
	}

	if identity != "" {
		fi, err := client.IAM().GetFederatedIdentity(ctx, identity)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error getting federated identity by id: %w", err))
		}
		if fi == nil {
			return diag.FromErr(fmt.Errorf("no federated identity found with id '%s'", identity))
		}
		return setFederatedIdentityState(d, fi)
	}

	identities, err := client.IAM().ListFederatedIdentities(ctx, &iam.ListFederatedIdentitiesRequest{})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error listing federated identities: %w", err))
	}

	var matches []iam.FederatedIdentity
	for _, fi := range identities {
		if fi.Name == name {
			matches = append(matches, fi)
		}
	}
	if len(matches) == 0 {
		return diag.FromErr(fmt.Errorf("no federated identity found with name '%s'", name))
	}
	if len(matches) > 1 {
		var ids []string
		for _, fi := range matches {
			ids = append(ids, fi.Identity)
		}
		return diag.FromErr(fmt.Errorf("multiple federated identities found with name '%s', please specify one of these ids: %v", name, ids))
	}

	return setFederatedIdentityState(d, &matches[0])
}
