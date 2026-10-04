package iam

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/thalassa-cloud/client-go/pkg/base"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

const (
	MyAccessElevationsEndpoint = "/v1/projects/iam/my-access-elevations"
	AccessElevationsEndpoint   = "/v1/projects/iam/access-elevations"
)

// IamAccessElevationStatus is the lifecycle state of a temporary access elevation request.
type IamAccessElevationStatus string

const (
	IamAccessElevationStatusPending   IamAccessElevationStatus = "pending"
	IamAccessElevationStatusApproved  IamAccessElevationStatus = "approved"
	IamAccessElevationStatusRejected  IamAccessElevationStatus = "rejected"
	IamAccessElevationStatusCancelled IamAccessElevationStatus = "cancelled"
	IamAccessElevationStatusExpired   IamAccessElevationStatus = "expired"
	IamAccessElevationStatusRevoked   IamAccessElevationStatus = "revoked"
)

// IamAccessElevationRequest is an organisation- or project-scoped temporary IAM access request.
type IamAccessElevationRequest struct {
	Identity      string            `json:"identity"`
	Annotations   map[string]string `json:"annotations,omitempty"`
	Labels        map[string]string `json:"labels,omitempty"`
	CreatedAt     time.Time         `json:"createdAt"`
	UpdatedAt     *time.Time        `json:"updatedAt,omitempty"`
	ObjectVersion int64             `json:"objectVersion"`

	Organisation *base.Organisation `json:"organisation,omitempty"`
	ProjectId    string             `json:"projectId,omitempty"`
	Project      *IamPolicyProject  `json:"project,omitempty"`

	IamPolicy *IamPolicy `json:"iamPolicy,omitempty"`

	Requester *base.AppUser `json:"requester,omitempty"`

	Status             IamAccessElevationStatus `json:"status"`
	Reason             string                   `json:"reason"`
	RequestedExpiresAt time.Time                `json:"requestedExpiresAt"`

	ReviewedAt *time.Time    `json:"reviewedAt,omitempty"`
	ReviewedBy *base.AppUser `json:"reviewedBy,omitempty"`
	ReviewNote *string       `json:"reviewNote,omitempty"`

	IamPolicyBinding *IamPolicyBinding `json:"iamPolicyBinding,omitempty"`
}

// CreateAccessElevationRequest is the body for POST /v1/projects/iam/my-access-elevations.
type CreateAccessElevationRequest struct {
	Annotations    map[string]string `json:"annotations"`
	Labels         map[string]string `json:"labels"`
	PolicyIdentity string            `json:"policyIdentity"`
	Reason         string            `json:"reason"`
	Duration       string            `json:"duration"`
	ExpiresAt      *time.Time        `json:"expiresAt,omitempty"`
}

// ReviewAccessElevationRequest is the body for approve/reject/revoke.
type ReviewAccessElevationRequest struct {
	ReviewNote string `json:"reviewNote"`
}

// ListMyAccessElevations lists access elevation requests created by the authenticated user.
func (c *Client) ListMyAccessElevations(ctx context.Context) ([]IamAccessElevationRequest, error) {
	requests := []IamAccessElevationRequest{}
	req := c.R().SetResult(&requests)
	resp, err := c.Do(ctx, req, client.GET, MyAccessElevationsEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return nil, err
	}
	return requests, nil
}

// GetMyAccessElevation retrieves one of the authenticated user's access elevation requests.
func (c *Client) GetMyAccessElevation(ctx context.Context, identity string) (*IamAccessElevationRequest, error) {
	if strings.TrimSpace(identity) == "" {
		return nil, fmt.Errorf("identity is required")
	}

	request := IamAccessElevationRequest{}
	req := c.R().SetResult(&request)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", MyAccessElevationsEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return nil, err
	}
	return &request, nil
}

// CreateAccessElevation creates a temporary access elevation request for a policy.
func (c *Client) CreateAccessElevation(ctx context.Context, create CreateAccessElevationRequest) (*IamAccessElevationRequest, error) {
	if strings.TrimSpace(create.PolicyIdentity) == "" {
		return nil, fmt.Errorf("policyIdentity is required")
	}
	if strings.TrimSpace(create.Reason) == "" {
		return nil, fmt.Errorf("reason is required")
	}
	if create.ExpiresAt == nil && strings.TrimSpace(create.Duration) == "" {
		return nil, fmt.Errorf("duration or expiresAt is required")
	}

	request := IamAccessElevationRequest{}
	req := c.R().SetBody(create).SetResult(&request)
	resp, err := c.Do(ctx, req, client.POST, MyAccessElevationsEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return nil, err
	}
	return &request, nil
}

// CancelAccessElevation cancels a pending access elevation request owned by the caller.
func (c *Client) CancelAccessElevation(ctx context.Context, identity string) (*IamAccessElevationRequest, error) {
	if strings.TrimSpace(identity) == "" {
		return nil, fmt.Errorf("identity is required")
	}

	request := IamAccessElevationRequest{}
	req := c.R().SetResult(&request)
	resp, err := c.Do(ctx, req, client.POST, fmt.Sprintf("%s/%s/cancel", MyAccessElevationsEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return nil, err
	}
	return &request, nil
}

// ListAccessElevations lists access elevation requests visible to approvers in the current scope.
func (c *Client) ListAccessElevations(ctx context.Context) ([]IamAccessElevationRequest, error) {
	requests := []IamAccessElevationRequest{}
	req := c.R().SetResult(&requests)
	resp, err := c.Do(ctx, req, client.GET, AccessElevationsEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return nil, err
	}
	return requests, nil
}

// GetAccessElevation retrieves an access elevation request for approvers.
func (c *Client) GetAccessElevation(ctx context.Context, identity string) (*IamAccessElevationRequest, error) {
	if strings.TrimSpace(identity) == "" {
		return nil, fmt.Errorf("identity is required")
	}

	request := IamAccessElevationRequest{}
	req := c.R().SetResult(&request)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", AccessElevationsEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return nil, err
	}
	return &request, nil
}

// ApproveAccessElevation approves a pending access elevation request.
func (c *Client) ApproveAccessElevation(ctx context.Context, identity string, review ReviewAccessElevationRequest) (*IamAccessElevationRequest, error) {
	return c.reviewAccessElevation(ctx, identity, "approve", review)
}

// RejectAccessElevation rejects a pending access elevation request.
func (c *Client) RejectAccessElevation(ctx context.Context, identity string, review ReviewAccessElevationRequest) (*IamAccessElevationRequest, error) {
	return c.reviewAccessElevation(ctx, identity, "reject", review)
}

// RevokeAccessElevation revokes an approved access elevation request.
func (c *Client) RevokeAccessElevation(ctx context.Context, identity string, review ReviewAccessElevationRequest) (*IamAccessElevationRequest, error) {
	return c.reviewAccessElevation(ctx, identity, "revoke", review)
}

func (c *Client) reviewAccessElevation(ctx context.Context, identity, action string, review ReviewAccessElevationRequest) (*IamAccessElevationRequest, error) {
	if strings.TrimSpace(identity) == "" {
		return nil, fmt.Errorf("identity is required")
	}

	request := IamAccessElevationRequest{}
	req := c.R().SetBody(review).SetResult(&request)
	resp, err := c.Do(ctx, req, client.POST, fmt.Sprintf("%s/%s/%s", AccessElevationsEndpoint, identity, action))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return nil, err
	}
	return &request, nil
}
