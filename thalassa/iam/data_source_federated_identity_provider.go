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

func DataSourceFederatedIdentityProvider() *schema.Resource {
	return &schema.Resource{
		Description: "Get a federated identity provider by identity or name",
		ReadContext: dataSourceFederatedIdentityProviderRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
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
				Optional:     true,
				ValidateFunc: validate.StringLenBetween(1, 255),
				Description:  "Name of the federated identity provider",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Human-readable description of the provider",
			},
			"labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Labels for the provider",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"annotations": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Annotations for the provider",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"provider_issuer": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "OIDC issuer URL of the identity provider",
			},
			"provider_jwks_uri": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Optional JWKS URI override for the provider",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the provider",
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

func dataSourceFederatedIdentityProviderRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	identity := d.Get("id").(string)
	name := d.Get("name").(string)

	if identity == "" && name == "" {
		return diag.FromErr(fmt.Errorf("either 'id' or 'name' must be provided to look up a federated identity provider"))
	}

	if identity != "" {
		providerObj, err := client.IAM().GetFederatedIdentityProvider(ctx, identity)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error getting federated identity provider by id: %w", err))
		}
		if providerObj == nil {
			return diag.FromErr(fmt.Errorf("no federated identity provider found with id '%s'", identity))
		}
		return setFederatedIdentityProviderState(d, providerObj)
	}

	providers, err := client.IAM().ListFederatedIdentityProviders(ctx, &iam.ListFederatedIdentityProvidersRequest{})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error listing federated identity providers: %w", err))
	}

	var matches []iam.FederatedIdentityProvider
	for _, p := range providers {
		if p.Name == name {
			matches = append(matches, p)
		}
	}
	if len(matches) == 0 {
		return diag.FromErr(fmt.Errorf("no federated identity provider found with name '%s'", name))
	}
	if len(matches) > 1 {
		var identities []string
		for _, p := range matches {
			identities = append(identities, p.Identity)
		}
		return diag.FromErr(fmt.Errorf("multiple federated identity providers found with name '%s', please specify one of these ids: %v", name, identities))
	}

	return setFederatedIdentityProviderState(d, &matches[0])
}
