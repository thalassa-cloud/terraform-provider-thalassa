package iam

import (
	"time"

	"github.com/thalassa-cloud/client-go/pkg/base"
)

type Team struct {
	Identity    string            `json:"identity"`
	Name        string            `json:"name"`
	Slug        string            `json:"slug"`
	Description string            `json:"description"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   *time.Time        `json:"updatedAt,omitempty"`
	Members     []TeamMember      `json:"members"`
}

type CreateTeam struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
}

type UpdateTeam struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
}

type TeamMember struct {
	Identity  string       `json:"identity"`
	Role      string       `json:"role"`
	User      base.AppUser `json:"user"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt *time.Time   `json:"updatedAt,omitempty"`
}

// LocalJWKS represents a locally stored JWKS (JSON Web Key Set)
// This allows offline verification without fetching from the issuer
type LocalJWKS struct {
	Keys []map[string]interface{} `json:"keys"`
}

// FederatedIdentityProviderStatus represents the status of a federated identity provider
type FederatedIdentityProviderStatus string

const (
	FederatedIdentityProviderStatusActive   FederatedIdentityProviderStatus = "active"
	FederatedIdentityProviderStatusInactive FederatedIdentityProviderStatus = "inactive"
)

// FederatedIdentityProvider represents an OIDC provider that can be used for token exchange
type FederatedIdentityProvider struct {
	// Identity is a unique identifier for the federated identity provider
	Identity string `json:"identity"`
	// Name is a human-readable name of the federated identity provider
	Name string `json:"name"`
	// Description is a human-readable description of the federated identity provider
	Description string `json:"description,omitempty"`
	// Annotations is a map of key-value pairs used for storing additional information
	Annotations map[string]string `json:"annotations,omitempty"`
	// Labels is a map of key-value pairs used for filtering and grouping objects
	Labels map[string]string `json:"labels,omitempty"`

	// CreatedAt is the timestamp when the object was created
	CreatedAt time.Time `json:"createdAt"`
	// UpdatedAt is the timestamp when the object was last updated
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	// DeletedAt is the timestamp when the object was deleted
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
	// ObjectVersion represents the version of the object
	ObjectVersion int64 `json:"objectVersion"`

	// ProviderIssuer is the issuer URL of the OIDC provider
	// This must be unique per organisation
	// This is the 'iss' claim from the JWT token
	ProviderIssuer string `json:"providerIssuer"`

	// ProviderJwksURI is an optional JWKS URI to override the discovered one from the issuer's openid-configuration
	// This is useful for caching or when the discovered URI is not accessible
	ProviderJwksURI *string `json:"providerJwksUri,omitempty"`

	// LocalJWKS is an optional locally stored JWKS for offline verification
	// If provided, this JWKS will be used instead of fetching from ProviderJwksURI or issuer
	// This is useful for air-gapped environments or when you want to pin specific keys
	LocalJWKS *LocalJWKS `json:"localJwks,omitempty"`

	// Status is the current status of the provider
	Status FederatedIdentityProviderStatus `json:"status"`

	// CreatedBy is the user that created the provider
	CreatedBy *base.AppUser `json:"createdBy,omitempty"`

	// ParentResourceIdentity is the identity of the parent resource that the provider is linked to
	// This is used to link a federated identity provider to a parent resource, i.e. a kubernetes cluster
	ParentResourceIdentity *string `json:"parentResourceIdentity,omitempty"`
	// ParentResourceType is the type of the parent resource that the provider is linked to
	// This is used to link a federated identity provider to a parent resource, i.e. a kubernetes cluster
	ParentResourceType *string `json:"parentResourceType,omitempty"`
}

// FederatedIdentityStatus represents the status of a federated identity
type FederatedIdentityStatus string

const (
	FederatedIdentityStatusActive   FederatedIdentityStatus = "active"
	FederatedIdentityStatusInactive FederatedIdentityStatus = "inactive"
	FederatedIdentityStatusExpired  FederatedIdentityStatus = "expired"
	FederatedIdentityStatusRevoked  FederatedIdentityStatus = "revoked"
)

// AudienceMatchMode represents how audience matching should be performed
type AudienceMatchMode string

const (
	AudienceMatchModeExact AudienceMatchMode = "exact" // Must match exactly one of the trusted audiences
	AudienceMatchModeAny   AudienceMatchMode = "any"   // Must match any of the trusted audiences (default)
	AudienceMatchModeAll   AudienceMatchMode = "all"   // Must match all trusted audiences
)

// FederatedIdentity represents a federated identity that can be used for OIDC token provisioning
type FederatedIdentity struct {
	// Identity is a unique identifier for the federated identity
	Identity string `json:"identity"`
	// Name is a human-readable name of the federated identity
	Name string `json:"name"`
	// Description is a human-readable description of the federated identity
	Description string `json:"description,omitempty"`
	// Annotations is a map of key-value pairs used for storing additional information
	Annotations map[string]string `json:"annotations,omitempty"`
	// Labels is a map of key-value pairs used for filtering and grouping objects
	Labels map[string]string `json:"labels,omitempty"`

	// CreatedAt is the timestamp when the object was created
	CreatedAt time.Time `json:"createdAt"`
	// UpdatedAt is the timestamp when the object was last updated
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	// DeletedAt is the timestamp when the object was deleted
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
	// ObjectVersion represents the version of the object
	ObjectVersion int64 `json:"objectVersion"`

	// ServiceAccount is the service account that this federated identity is bound to
	ServiceAccount *ServiceAccount `json:"serviceAccount,omitempty"`

	// Provider is the federated identity provider
	Provider *FederatedIdentityProvider `json:"provider,omitempty"`

	// ProviderSubject is the subject identifier from the OIDC provider
	// This is the 'sub' claim from the JWT token
	ProviderSubject string `json:"providerSubject"`

	// TrustedAudiences is a list of trusted audiences for the federated identity
	// These are the audiences that the JWT token must contain to be considered valid
	TrustedAudiences []string `json:"trustedAudiences"`

	// AudienceMatchMode defines how audience matching should be performed
	AudienceMatchMode AudienceMatchMode `json:"audienceMatchMode"`

	// AllowedScopes is a list of scopes that the federated identity is allowed to access
	AllowedScopes []AccessCredentialsScope `json:"allowedScopes"`

	// Status is the current status of the federated identity
	Status FederatedIdentityStatus `json:"status"`

	// ExpiresAt is the timestamp when the federated identity will expire
	// If not set, the federated identity will never expire
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`

	// LastUsedAt is the timestamp when the federated identity was last used
	LastUsedAt *time.Time `json:"lastUsedAt,omitempty"`

	// Conditions is a JSONB field containing conditions/claims matcher rules
	// This allows locking to specific branches, environments, PRs, ref_protected, workflow_ref, etc.
	// Example: {"branch": "main", "environment": "production", "ref_protected": true}
	Conditions map[string]interface{} `json:"conditions,omitempty"`

	// CreatedBy is the user that created the federated identity
	CreatedBy *base.AppUser `json:"createdBy,omitempty"`
}
