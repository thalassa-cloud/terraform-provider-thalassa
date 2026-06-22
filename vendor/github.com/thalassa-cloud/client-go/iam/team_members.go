package iam

import (
	"context"
	"fmt"

	"github.com/thalassa-cloud/client-go/pkg/client"
)

type AddTeamMemberRequest struct {
	UserIdentity string `json:"userIdentity"`
	Role         string `json:"role"`
}

// AddTeamMember adds a member to a team.
func (c *Client) AddTeamMember(ctx context.Context, teamID string, request AddTeamMemberRequest) error {
	req := c.R().SetBody(request)
	resp, err := c.Do(ctx, req, client.POST, fmt.Sprintf("%s/%s/members", TeamEndpoint, teamID))
	if err != nil {
		return err
	}
	if err := c.Check(resp); err != nil {
		return err
	}
	return nil
}

// RemoveTeamMember removes a member from a team.
func (c *Client) RemoveTeamMember(ctx context.Context, teamID string, memberIdentity string) error {
	req := c.R()
	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s/members/%s", TeamEndpoint, teamID, memberIdentity))
	if err != nil {
		return err
	}
	if err := c.Check(resp); err != nil {
		return err
	}
	return nil
}
