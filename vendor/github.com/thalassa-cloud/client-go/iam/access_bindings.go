package iam

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/thalassa-cloud/client-go/pkg/client"
)

const (
	IamAccessBindingsEndpoint = "/v1/projects/iam/access-bindings"

	ResourceAccessBindingPrincipalUser           ResourceAccessBindingPrincipalType = "user"
	ResourceAccessBindingPrincipalServiceAccount ResourceAccessBindingPrincipalType = "service_account"
)

// ResourceAccessBindingPrincipalType identifies the principal holding a resource access binding.
type ResourceAccessBindingPrincipalType string

// ResourceAccessBindingItem is one principal binding for a resource.
type ResourceAccessBindingItem struct {
	ResourceType         string                             `json:"resourceType"`
	BindingType          string                             `json:"bindingType"`
	BindingIdentity      string                             `json:"bindingIdentity"`
	PolicyIdentity       string                             `json:"policyIdentity,omitempty"`
	PolicyName           string                             `json:"policyName,omitempty"`
	PrincipalType        ResourceAccessBindingPrincipalType `json:"principalType,omitempty"`
	PrincipalIdentity    string                             `json:"principalIdentity,omitempty"`
	PrincipalDisplayName string                             `json:"principalDisplayName,omitempty"`
	Permissions          []string                           `json:"permissions"`
	ExpiresAt            *time.Time                         `json:"expiresAt,omitempty"`
}

// ResourceAccessBindingsResponse is the response for GET /v1/projects/iam/access-bindings.
type ResourceAccessBindingsResponse struct {
	ResourceIdentity string                      `json:"resourceIdentity"`
	Bindings         []ResourceAccessBindingItem `json:"bindings"`
}

// ListResourceAccessBindingsRequest configures GET /v1/projects/iam/access-bindings.
type ListResourceAccessBindingsRequest struct {
	ResourceTypes    []string
	ResourceIdentity string
	Permissions      []PermissionType
	BindingTypes     []string
}

// ListResourceAccessBindings lists principals with IAM access to a resource identity.
func (c *Client) ListResourceAccessBindings(ctx context.Context, request ListResourceAccessBindingsRequest) (*ResourceAccessBindingsResponse, error) {
	if strings.TrimSpace(request.ResourceIdentity) == "" {
		return nil, fmt.Errorf("resourceIdentity is required")
	}
	if len(request.ResourceTypes) == 0 {
		return nil, fmt.Errorf("resourceType is required")
	}
	for _, resourceType := range request.ResourceTypes {
		if strings.TrimSpace(resourceType) == "" {
			return nil, fmt.Errorf("resourceType must not be empty")
		}
	}

	result := ResourceAccessBindingsResponse{}
	req := c.R().SetResult(&result)
	for _, resourceType := range request.ResourceTypes {
		req.QueryParam.Add("resourceType", resourceType)
	}
	req.SetQueryParam("resourceIdentity", strings.TrimSpace(request.ResourceIdentity))
	for _, permission := range request.Permissions {
		req.QueryParam.Add("permissions", string(permission))
	}
	for _, bindingType := range request.BindingTypes {
		req.QueryParam.Add("bindingTypes", bindingType)
	}

	resp, err := c.Do(ctx, req, client.GET, IamAccessBindingsEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return nil, err
	}
	return &result, nil
}
