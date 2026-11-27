package iam

import (
	"context"
	"fmt"
	"time"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/base"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

type OrganisationMemberInvite struct {
	Email               string                 `json:"email"`
	Organisation        *base.Organisation     `json:"organisation,omitempty"`
	Role                OrganisationMemberType `json:"role"`
	InvitedByUser       *base.AppUser          `json:"invitedByUser,omitempty"`
	JoinTeamOnAccept    *Team                  `json:"joinTeamOnAccept,omitempty"`
	RolebindingOnAccept *OrganisationRole      `json:"rolebindingOnAccept,omitempty"`
	CreatedAt           time.Time              `json:"createdAt"`
	ExpiresAt           *time.Time             `json:"expiresAt"`
	InviteCode          string                 `json:"inviteCode"`
}

// InviteUserRequest is the request body for creating an organisation member invite.
type InviteUserRequest struct {
	Email               string                 `json:"email"`
	Role                OrganisationMemberType `json:"role"`
	JoinTeamOnAccept    *string                `json:"joinTeamOnAccept,omitempty"`
	RolebindingOnAccept *string                `json:"rolebindingOnAccept,omitempty"`
}

const (
	OrganisationMemberInviteEndpoint = "/v1/invites"
)

type ListOrganisationMemberInvitesRequest struct {
	Filters []filters.Filter
}

// ListOrganisationMemberInvites lists all invites for an organisation
func (c *Client) ListOrganisationMemberInvites(ctx context.Context, request *ListOrganisationMemberInvitesRequest) ([]OrganisationMemberInvite, error) {
	invites := []OrganisationMemberInvite{}
	req := c.R().SetResult(&invites)
	if request != nil {
		for _, filter := range request.Filters {
			for k, v := range filter.ToParams() {
				req = req.SetQueryParam(k, v)
			}
		}
	}
	resp, err := c.Do(ctx, req, client.GET, OrganisationMemberInviteEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return nil, err
	}
	return invites, nil
}

// CreateOrganisationMemberInvite creates an invite for a user to join the organisation.
func (c *Client) CreateOrganisationMemberInvite(ctx context.Context, create InviteUserRequest) (*OrganisationMemberInvite, error) {
	if create.Email == "" {
		return nil, fmt.Errorf("email is required")
	}
	if create.Role == "" {
		return nil, fmt.Errorf("role is required")
	}

	var invite *OrganisationMemberInvite
	req := c.R().SetBody(create).SetResult(&invite)
	resp, err := c.Do(ctx, req, client.POST, OrganisationMemberInviteEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return invite, err
	}
	return invite, nil
}

// AcceptOrganisationMemberInvite accepts the pending invite for the authenticated user.
func (c *Client) AcceptOrganisationMemberInvite(ctx context.Context) error {
	req := c.R()
	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/accept", OrganisationMemberInviteEndpoint))
	if err != nil {
		return err
	}
	return c.Check(resp)
}

// DeleteOrganisationMemberInvite deletes an invite by invite code.
func (c *Client) DeleteOrganisationMemberInvite(ctx context.Context, inviteCode string) error {
	if inviteCode == "" {
		return fmt.Errorf("invite code is required")
	}

	req := c.R()
	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s", OrganisationMemberInviteEndpoint, inviteCode))
	if err != nil {
		return err
	}
	return c.Check(resp)
}
