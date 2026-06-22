package iam

import (
	"context"
	"fmt"
	"time"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/base"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

const (
	ServiceAccountEndpoint = "/v1/service-accounts"
)

// ListServiceAccounts lists all service accounts for the organisation
func (c *Client) ListServiceAccounts(ctx context.Context, request *ListServiceAccountsRequest) ([]ServiceAccount, error) {
	accounts := []ServiceAccount{}
	req := c.R().SetResult(&accounts)
	if request != nil {
		for _, filter := range request.Filters {
			for k, v := range filter.ToParams() {
				req.SetQueryParam(k, v)
			}
		}
	}

	resp, err := c.Do(ctx, req, client.GET, ServiceAccountEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return accounts, err
	}
	return accounts, nil
}

// GetServiceAccount retrieves a specific service account
func (c *Client) GetServiceAccount(ctx context.Context, identity string) (*ServiceAccount, error) {
	var account *ServiceAccount
	req := c.R().SetResult(&account)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", ServiceAccountEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return account, err
	}
	return account, nil
}

// CreateServiceAccount creates a new service account
func (c *Client) CreateServiceAccount(ctx context.Context, create CreateServiceAccountRequest) (*ServiceAccount, error) {
	var account *ServiceAccount
	req := c.R().SetBody(create).SetResult(&account)
	resp, err := c.Do(ctx, req, client.POST, ServiceAccountEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return account, err
	}
	return account, nil
}

// UpdateServiceAccount updates a service account (partial update)
func (c *Client) UpdateServiceAccount(ctx context.Context, identity string, update UpdateServiceAccountRequest) (*ServiceAccount, error) {
	var account *ServiceAccount
	req := c.R().SetBody(update).SetResult(&account)
	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s", ServiceAccountEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return account, err
	}
	return account, nil
}

// DeleteServiceAccount deletes a service account
func (c *Client) DeleteServiceAccount(ctx context.Context, identity string) error {
	req := c.R()
	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s", ServiceAccountEndpoint, identity))
	if err != nil {
		return err
	}
	if err := c.Check(resp); err != nil {
		return err
	}
	return nil
}

// GetServiceAccountAccessCredentials lists access credentials for a service account
func (c *Client) GetServiceAccountAccessCredentials(ctx context.Context, serviceAccountIdentity string) ([]ServiceAccountAccessCredential, error) {
	creds := []ServiceAccountAccessCredential{}
	req := c.R().SetResult(&creds)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s/access-credentials", ServiceAccountEndpoint, serviceAccountIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return creds, err
	}
	return creds, nil
}

// CreateServiceAccountAccessCredentials creates access credentials for a service account
func (c *Client) CreateServiceAccountAccessCredentials(ctx context.Context, serviceAccountIdentity string, request CreateServiceAccountAccessCredentialRequest) (*ServiceAccountCreatedAccessCredential, error) {
	var created *ServiceAccountCreatedAccessCredential
	req := c.R().SetBody(request).SetResult(&created)
	resp, err := c.Do(ctx, req, client.POST, fmt.Sprintf("%s/%s/access-credentials", ServiceAccountEndpoint, serviceAccountIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return created, err
	}
	return created, nil
}

// DeleteServiceAccountAccessCredentials deletes specific access credentials for a service account
func (c *Client) DeleteServiceAccountAccessCredentials(ctx context.Context, serviceAccountIdentity string, credentialIdentity string) error {
	req := c.R()
	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s/access-credentials/%s", ServiceAccountEndpoint, serviceAccountIdentity, credentialIdentity))
	if err != nil {
		return err
	}
	if err := c.Check(resp); err != nil {
		return err
	}
	return nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Types
// ─────────────────────────────────────────────────────────────────────────────

// AccessCredentialsScope is a type that represents a scope for an access credential
type AccessCredentialsScope string

const (
	// AccessCredentialsScopeAPIRead is a scope that allows read access to the API
	AccessCredentialsScopeAPIRead AccessCredentialsScope = "api:read"
	// AccessCredentialsScopeAPIWrite is a scope that allows write access to the API
	AccessCredentialsScopeAPIWrite AccessCredentialsScope = "api:write"
	// AccessCredentialsScopeKubernetes is a scope that allows access to the Kubernetes API
	AccessCredentialsScopeKubernetes AccessCredentialsScope = "kubernetes"
	// AccessCredentialsScopeObjectStorage is a scope that allows access to the Object Storage API
	AccessCredentialsScopeObjectStorage AccessCredentialsScope = "objectStorage"
)

// ServiceAccount is the response for listing service accounts
type ServiceAccount struct {
	// Identity is the unique identifier for the service account
	Identity string `json:"identity"`
	// Name is the name of the service account
	Name string `json:"name"`
	// Slug is a human-readable unique identifier for the service account
	Slug string `json:"slug"`
	// Description is the description of the service account
	Description *string `json:"description,omitempty"`
	// Annotations is a map of key-value pairs used for storing additional information
	Annotations map[string]string `json:"annotations,omitempty"`
	// Labels is a map of key-value pairs used for filtering and grouping service accounts
	Labels map[string]string `json:"labels,omitempty"`
	// CreatedAt is the timestamp when the service account was created
	CreatedAt time.Time `json:"createdAt"`
	// UpdatedAt is the timestamp when the service account was last updated
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	// DeletedAt is the timestamp when the service account was deleted
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
	// ObjectVersion is the version of the service account
	ObjectVersion int `json:"objectVersion"`
	// Organisation is the organisation that the service account belongs to
	Organisation *base.Organisation `json:"organisation,omitempty"`
	// CreatedBy is the user that created the service account
	CreatedBy *base.AppUser `json:"createdBy,omitempty"`
	// RoleBindings is a list of role bindings for the service account
	RoleBindings []OrganisationRoleBinding `json:"roleBindings,omitempty"`
}

// CreateServiceAccountRequest is the request for creating a service account
type CreateServiceAccountRequest struct {
	// Name is the name of the service account
	Name string `json:"name"`
	// Description is the description of the service account
	Description *string `json:"description,omitempty"`
	// Annotations is a map of key-value pairs used for storing additional information
	Annotations map[string]string `json:"annotations,omitempty"`
	// Labels is a map of key-value pairs used for filtering and grouping service accounts
	Labels map[string]string `json:"labels,omitempty"`
}

// UpdateServiceAccountRequest is the request for updating a service account
type UpdateServiceAccountRequest struct {
	// Name is the name of the service account
	Name *string `json:"name,omitempty"`
	// Description is the description of the service account
	Description *string `json:"description,omitempty"`
	// Annotations is a map of key-value pairs used for storing additional information
	Annotations map[string]string `json:"annotations,omitempty"`
	// Labels is a map of key-value pairs used for filtering and grouping service accounts
	Labels map[string]string `json:"labels,omitempty"`
}

// ServiceAccountAccessCredential is the response for listing access credentials for a service account
type ServiceAccountAccessCredential struct {
	// Identity is the unique identifier for the access credential
	Identity string `json:"identity"`
	// Name is the name of the access credential
	Name string `json:"name"`
	// Description is the description of the access credential
	Description *string `json:"description,omitempty"`
	// CreatedAt is the timestamp when the access credential was created
	CreatedAt time.Time `json:"createdAt"`
	// LastUsedAt is the timestamp when the access credential was last used
	LastUsedAt *time.Time `json:"lastUsedAt,omitempty"`
	// ExpiresAt is the timestamp when the access credential expires
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
	// AccessKey is the access key for the access credential
	AccessKey string `json:"accessKey"`
}

// ServiceAccountCreatedAccessCredential is the response for creating an access credential for a service account
type ServiceAccountCreatedAccessCredential struct {
	// AccessKey is the access key for the access credential
	AccessKey string `json:"accessKey"`
	// AccessSecret is the access secret for the access credential
	AccessSecret string `json:"accessSecret"`
	// Identity is the unique identifier for the access credential
	Identity string `json:"identity"`
	// Scopes is a list of scopes for the access credential
	Scopes []AccessCredentialsScope `json:"scopes,omitempty"`
}

// CreateServiceAccountAccessCredentialRequest is the request for creating an access credential for a service account
type CreateServiceAccountAccessCredentialRequest struct {
	// Name is the name of the access credential
	Name string `json:"name"`
	// Description is the description of the access credential
	Description *string `json:"description,omitempty"`
	// ExpiresAt is the timestamp when the access credential expires
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
	// Scopes is a list of scopes for the access credential
	Scopes []AccessCredentialsScope `json:"scopes,omitempty"`
}

// ListServiceAccountsRequest is the request for listing service accounts
type ListServiceAccountsRequest struct {
	// Filters is a list of filters to apply to the request
	Filters []filters.Filter
}
