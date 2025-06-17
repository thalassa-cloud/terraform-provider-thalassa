package iaas

import (
	"context"
	"fmt"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

const (
	VpcEndpoint = "/v1/vpcs"
)

type ListVpcsRequest struct {
	Filters []filters.Filter
}

// ListVpcs lists all VPCs for a given organisation.
func (c *Client) ListVpcs(ctx context.Context, request *ListVpcsRequest) ([]Vpc, error) {
	vpcs := []Vpc{}
	req := c.R().SetResult(&vpcs)
	if request != nil {
		for _, filter := range request.Filters {
			for k, v := range filter.ToParams() {
				req.SetQueryParam(k, v)
			}
		}
	}

	resp, err := c.Do(ctx, req, client.GET, VpcEndpoint)
	if err != nil {
		return nil, err
	}

	if err := c.Check(resp); err != nil {
		return vpcs, err
	}
	return vpcs, nil
}

// GetVpc retrieves a specific VPC by its identity.
func (c *Client) GetVpc(ctx context.Context, identity string) (*Vpc, error) {
	var vpc *Vpc
	req := c.R().SetResult(&vpc)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", VpcEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return vpc, err
	}
	return vpc, nil
}

// CreateVpc creates a new VPC.
func (c *Client) CreateVpc(ctx context.Context, create CreateVpc) (*Vpc, error) {
	var vpc *Vpc
	req := c.R().
		SetBody(create).SetResult(&vpc)

	resp, err := c.Do(ctx, req, client.POST, VpcEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return vpc, err
	}
	return vpc, nil
}

// UpdateVpc updates an existing VPC.
func (c *Client) UpdateVpc(ctx context.Context, identity string, update UpdateVpc) (*Vpc, error) {
	var vpc *Vpc
	req := c.R().
		SetBody(update).SetResult(&vpc)

	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s", VpcEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return vpc, err
	}
	return vpc, nil
}

// DeleteVpc deletes a specific VPC by its identity.
func (c *Client) DeleteVpc(ctx context.Context, identity string) error {
	req := c.R()

	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s", VpcEndpoint, identity))
	if err != nil {
		return err
	}
	if err := c.Check(resp); err != nil {
		return err
	}
	return nil
}
