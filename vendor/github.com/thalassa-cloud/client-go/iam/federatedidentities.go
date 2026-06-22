package iam

import (
	"context"
	"fmt"
	"time"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

const (
	FederatedIdentityEndpoint = "/v1/federated-identities"
)

// ListFederatedIdentities lists all federated identities for the organisation
func (c *Client) ListFederatedIdentities(ctx context.Context, request *ListFederatedIdentitiesRequest) ([]FederatedIdentity, error) {
	identities := []FederatedIdentity{}
	req := c.R().SetResult(&identities)
	if request != nil {
		for _, filter := range request.Filters {
			for k, v := range filter.ToParams() {
				req.SetQueryParam(k, v)
			}
		}
	}

	resp, err := c.Do(ctx, req, client.GET, FederatedIdentityEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return identities, err
	}
	return identities, nil
}

// GetFederatedIdentity retrieves a specific federated identity
func (c *Client) GetFederatedIdentity(ctx context.Context, identity string) (*FederatedIdentity, error) {
	var federatedIdentity *FederatedIdentity
	req := c.R().SetResult(&federatedIdentity)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", FederatedIdentityEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return federatedIdentity, err
	}
	return federatedIdentity, nil
}

// CreateFederatedIdentity creates a new federated identity
func (c *Client) CreateFederatedIdentity(ctx context.Context, create CreateFederatedIdentityRequest) (*FederatedIdentity, error) {
	var federatedIdentity *FederatedIdentity
	req := c.R().SetBody(create).SetResult(&federatedIdentity)
	resp, err := c.Do(ctx, req, client.POST, FederatedIdentityEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return federatedIdentity, err
	}
	return federatedIdentity, nil
}

// UpdateFederatedIdentity updates a federated identity (partial update)
func (c *Client) UpdateFederatedIdentity(ctx context.Context, identity string, update UpdateFederatedIdentityRequest) (*FederatedIdentity, error) {
	var federatedIdentity *FederatedIdentity
	req := c.R().SetBody(update).SetResult(&federatedIdentity)
	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s", FederatedIdentityEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return federatedIdentity, err
	}
	return federatedIdentity, nil
}

// DeleteFederatedIdentity deletes a federated identity
func (c *Client) DeleteFederatedIdentity(ctx context.Context, identity string) error {
	req := c.R()
	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s", FederatedIdentityEndpoint, identity))
	if err != nil {
		return err
	}
	if err := c.Check(resp); err != nil {
		return err
	}
	return nil
}

// ListFederatedIdentitiesRequest is the request for listing federated identities
type ListFederatedIdentitiesRequest struct {
	// Filters is a list of filters to apply to the request
	Filters []filters.Filter
}

// CreateFederatedIdentityRequest is the request for creating a federated identity
type CreateFederatedIdentityRequest struct {
	// Name of the federated identity
	Name string `json:"name"`

	// Description of the federated identity
	Description string `json:"description,omitempty"`

	// Annotations for the federated identity
	Annotations map[string]string `json:"annotations,omitempty"`

	// Labels for the federated identity
	Labels map[string]string `json:"labels,omitempty"`

	// ServiceAccountIdentity is the identity of the service account to bind to
	ServiceAccountIdentity string `json:"serviceAccountIdentity"`

	// ProviderIdentity is the identity of the federated identity provider
	ProviderIdentity string `json:"providerIdentity"`

	// ProviderSubject is the subject identifier from the OIDC provider
	// This is the 'sub' claim from the JWT token
	ProviderSubject string `json:"providerSubject"`

	// TrustedAudiences is a list of trusted audiences for the federated identity
	TrustedAudiences []string `json:"trustedAudiences"`

	// AudienceMatchMode defines how audience matching should be performed
	AudienceMatchMode AudienceMatchMode `json:"audienceMatchMode"`

	// AllowedScopes is a list of scopes that the federated identity is allowed to access
	AllowedScopes []AccessCredentialsScope `json:"allowedScopes"`

	// ExpiresAt is the timestamp when the federated identity will expire
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`

	// Conditions is a JSONB field containing conditions/claims matcher rules
	Conditions map[string]interface{} `json:"conditions,omitempty"`
}

// UpdateFederatedIdentityRequest is the request for updating a federated identity
type UpdateFederatedIdentityRequest struct {
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`

	// TrustedAudiences is a list of trusted audiences for the federated identity
	TrustedAudiences []string `json:"trustedAudiences,omitempty"`

	// AudienceMatchMode defines how audience matching should be performed
	AudienceMatchMode AudienceMatchMode `json:"audienceMatchMode,omitempty"`

	// AllowedScopes is a list of scopes that the federated identity is allowed to access
	AllowedScopes []AccessCredentialsScope `json:"allowedScopes,omitempty"`

	// Status is the current status of the federated identity
	Status FederatedIdentityStatus `json:"status,omitempty"`

	// ExpiresAt is the timestamp when the federated identity will expire
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`

	// Conditions is a JSONB field containing conditions/claims matcher rules
	Conditions map[string]interface{} `json:"conditions,omitempty"`
}
