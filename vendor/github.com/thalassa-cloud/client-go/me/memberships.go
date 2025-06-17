package me

import (
	"context"

	"github.com/thalassa-cloud/client-go/pkg/base"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

const (
	OrganisationMemberEndpoint = "/v1/me/organisation-memberships"
)

// ListMyMemberships lists all memberships for the current user.
func (c *Client) ListMyMemberships(ctx context.Context) ([]base.OrganisationMember, error) {
	memberships := []base.OrganisationMember{}
	req := c.R().SetResult(&memberships)
	resp, err := c.Do(ctx, req, client.GET, OrganisationMemberEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return memberships, err
	}
	return memberships, nil
}

// ListMyOrganisations lists all organisations for the current user.
func (c *Client) ListMyOrganisations(ctx context.Context) ([]base.Organisation, error) {
	memberships, err := c.ListMyMemberships(ctx)
	if err != nil {
		return nil, err
	}
	organisations := []base.Organisation{}
	for _, membership := range memberships {
		if membership.Organisation != nil {
			organisations = append(organisations, *membership.Organisation)
		}
	}
	return organisations, nil
}
