package iam

import (
	"context"
	"fmt"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

const (
	FederatedIdentityProviderEndpoint = "/v1/federated-identity-providers"
)

// ListFederatedIdentityProviders lists all federated identity providers for the organisation
func (c *Client) ListFederatedIdentityProviders(ctx context.Context, request *ListFederatedIdentityProvidersRequest) ([]FederatedIdentityProvider, error) {
	providers := []FederatedIdentityProvider{}
	req := c.R().SetResult(&providers)
	if request != nil {
		for _, filter := range request.Filters {
			for k, v := range filter.ToParams() {
				req = req.SetQueryParam(k, v)
			}
		}
	}

	resp, err := c.Do(ctx, req, client.GET, FederatedIdentityProviderEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return providers, err
	}
	return providers, nil
}

// GetFederatedIdentityProvider retrieves a specific federated identity provider
func (c *Client) GetFederatedIdentityProvider(ctx context.Context, identity string) (*FederatedIdentityProvider, error) {
	var provider *FederatedIdentityProvider
	req := c.R().SetResult(&provider)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", FederatedIdentityProviderEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return provider, err
	}
	return provider, nil
}

// CreateFederatedIdentityProvider creates a new federated identity provider
func (c *Client) CreateFederatedIdentityProvider(ctx context.Context, create CreateFederatedIdentityProviderRequest) (*FederatedIdentityProvider, error) {
	var provider *FederatedIdentityProvider
	req := c.R().SetBody(create).SetResult(&provider)
	resp, err := c.Do(ctx, req, client.POST, FederatedIdentityProviderEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return provider, err
	}
	return provider, nil
}

// UpdateFederatedIdentityProvider updates a federated identity provider (partial update)
func (c *Client) UpdateFederatedIdentityProvider(ctx context.Context, identity string, update UpdateFederatedIdentityProviderRequest) (*FederatedIdentityProvider, error) {
	var provider *FederatedIdentityProvider
	req := c.R().SetBody(update).SetResult(&provider)
	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s", FederatedIdentityProviderEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return provider, err
	}
	return provider, nil
}

// DeleteFederatedIdentityProvider deletes a federated identity provider
func (c *Client) DeleteFederatedIdentityProvider(ctx context.Context, identity string) error {
	req := c.R()
	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s", FederatedIdentityProviderEndpoint, identity))
	if err != nil {
		return err
	}
	if err := c.Check(resp); err != nil {
		return err
	}
	return nil
}

// ListFederatedIdentityProvidersRequest is the request for listing federated identity providers
type ListFederatedIdentityProvidersRequest struct {
	// Filters is a list of filters to apply to the request
	Filters []filters.Filter
}

// CreateFederatedIdentityProviderRequest is the request for creating a federated identity provider
type CreateFederatedIdentityProviderRequest struct {
	// Name of the federated identity provider
	Name string `json:"name"`

	// Description of the federated identity provider
	Description string `json:"description,omitempty"`

	// Annotations for the federated identity provider
	Annotations map[string]string `json:"annotations,omitempty"`

	// Labels for the federated identity provider
	Labels map[string]string `json:"labels,omitempty"`

	// ProviderIssuer is the issuer URL of the OIDC provider
	// This must be unique per organisation
	ProviderIssuer string `json:"providerIssuer"`

	// ProviderJwksURI is an optional JWKS URI to override the discovered one from the issuer's openid-configuration
	ProviderJwksURI *string `json:"providerJwksUri,omitempty"`

	// LocalJWKS is an optional locally stored JWKS for offline verification
	LocalJWKS *LocalJWKS `json:"localJwks,omitempty"`

	// Status is the current status of the provider
	Status FederatedIdentityProviderStatus `json:"status,omitempty"`
}

// UpdateFederatedIdentityProviderRequest is the request for updating a federated identity provider
type UpdateFederatedIdentityProviderRequest struct {
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`

	// ProviderJwksURI is an optional JWKS URI to override the discovered one
	ProviderJwksURI *string `json:"providerJwksUri,omitempty"`

	// Status is the current status of the provider
	Status FederatedIdentityProviderStatus `json:"status,omitempty"`
}
