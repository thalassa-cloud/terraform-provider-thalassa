package iam

import (
	"context"
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
	resp, err := c.Do(ctx, req, client.GET, OrganisationMemberInviteEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return nil, err
	}
	return invites, nil
}
