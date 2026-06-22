package iam

import (
	"context"
	"fmt"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

const (
	TeamEndpoint = "/v1/teams"
)

type ListTeamsRequest struct {
	Filters []filters.Filter
}

// ListTeams lists all teams for a given organisation.
func (c *Client) ListTeams(ctx context.Context, request *ListTeamsRequest) ([]Team, error) {
	teams := []Team{}
	req := c.R().SetResult(&teams)
	if request != nil {
		for _, filter := range request.Filters {
			for k, v := range filter.ToParams() {
				req.SetQueryParam(k, v)
			}
		}
	}

	resp, err := c.Do(ctx, req, client.GET, TeamEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return nil, err
	}
	return teams, nil
}

type GetTeamRequest struct {
}

// GetTeam retrieves a specific team by its identity.
func (c *Client) GetTeam(ctx context.Context, identity string, _ *GetTeamRequest) (*Team, error) {
	team := Team{}
	req := c.R().SetResult(&team)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", TeamEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return nil, err
	}
	return &team, nil
}

// CreateTeam creates a new team.
func (c *Client) CreateTeam(ctx context.Context, create CreateTeam) (*Team, error) {
	team := Team{}
	req := c.R().SetBody(create).SetResult(&team)
	resp, err := c.Do(ctx, req, client.POST, TeamEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return nil, err
	}
	return &team, nil
}

// UpdateTeam updates a team.
func (c *Client) UpdateTeam(ctx context.Context, identity string, update UpdateTeam) (*Team, error) {
	team := Team{}
	req := c.R().SetBody(update).SetResult(&team)
	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s", TeamEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return nil, err
	}
	return &team, nil
}

// DeleteTeam deletes a team.
func (c *Client) DeleteTeam(ctx context.Context, identity string) error {
	req := c.R()
	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s", TeamEndpoint, identity))
	if err != nil {
		return err
	}
	if err := c.Check(resp); err != nil {
		return err
	}
	return nil
}
