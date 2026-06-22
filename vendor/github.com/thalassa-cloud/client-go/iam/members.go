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
	OrganisationMemberEndpoint = "/v1/memberships"
)

// ListOrganisationMembers lists all members of an organisation
func (c *Client) ListOrganisationMembers(ctx context.Context, request *ListMembersRequest) ([]OrganisationMember, error) {
	members := []OrganisationMember{}
	req := c.R().SetResult(&members)
	resp, err := c.Do(ctx, req, client.GET, OrganisationMemberEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return nil, err
	}
	return members, nil
}

// DeleteOrganisationMember deletes a member from an organisation
func (c *Client) DeleteOrganisationMember(ctx context.Context, identity string) error {
	req := c.R()
	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s", OrganisationMemberEndpoint, identity))
	if err != nil {
		return err
	}
	if err := c.Check(resp); err != nil {
		return err
	}
	return nil
}

// UpdateOrganisationMember updates a member of an organisation
func (c *Client) UpdateOrganisationMember(ctx context.Context, identity string, request UpdateOrganisationMemberRequest) error {
	req := c.R().SetBody(request)
	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s", OrganisationMemberEndpoint, identity))
	if err != nil {
		return err
	}
	if err := c.Check(resp); err != nil {
		return err
	}
	return nil
}

type ListMembersRequest struct {
	Filters []filters.Filter
}

type OrganisationMember struct {
	Identity   string                 `json:"identity"`
	CreatedAt  time.Time              `json:"createdAt"`
	User       *base.AppUser          `json:"user"`
	MemberType OrganisationMemberType `json:"role"`
}

// OrganisationMemberType is a type that represents a role of a member in an organisation
type OrganisationMemberType string

const (
	// OrganisationMemberTypeOwner is a role that indicates that the user is an owner of the organisation
	OrganisationMemberTypeOwner OrganisationMemberType = "OWNER"
	// OrganisationMemberTypeMember is a role that indicates that the user is a member of the organisation
	OrganisationMemberTypeMember OrganisationMemberType = "MEMBER"
)

type UpdateOrganisationMemberRequest struct {
	MemberType OrganisationMemberType `json:"role"`
}
