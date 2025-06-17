package iaas

import (
	"context"
	"fmt"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

const (
	NatGatewayEndpoint = "/v1/nat-gateways"
)

type ListNatGatewaysRequest struct {
	Filters []filters.Filter
}

// ListNatGateways lists all NatGateways for a given organisation.
func (c *Client) ListNatGateways(ctx context.Context, listRequest *ListNatGatewaysRequest) ([]VpcNatGateway, error) {
	subnets := []VpcNatGateway{}
	req := c.R().SetResult(&subnets)

	if listRequest != nil {
		for _, filter := range listRequest.Filters {
			for k, v := range filter.ToParams() {
				req = req.SetQueryParam(k, v)
			}
		}
	}

	resp, err := c.Do(ctx, req, client.GET, NatGatewayEndpoint)
	if err != nil {
		return nil, err
	}

	if err := c.Check(resp); err != nil {
		return subnets, err
	}
	return subnets, nil
}

// GetNatGateway retrieves a specific NatGateway by its identity.
func (c *Client) GetNatGateway(ctx context.Context, identity string) (*VpcNatGateway, error) {
	var subnet *VpcNatGateway
	req := c.R().SetResult(&subnet)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", NatGatewayEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return subnet, err
	}
	return subnet, nil
}

// CreateNatGateway creates a new NatGateway.
func (c *Client) CreateNatGateway(ctx context.Context, create CreateVpcNatGateway) (*VpcNatGateway, error) {
	var subnet *VpcNatGateway
	req := c.R().
		SetBody(create).SetResult(&subnet)

	resp, err := c.Do(ctx, req, client.POST, NatGatewayEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return subnet, err
	}
	return subnet, nil
}

// UpdateNatGateway updates an existing NatGateway.
func (c *Client) UpdateNatGateway(ctx context.Context, identity string, update UpdateVpcNatGateway) (*VpcNatGateway, error) {
	var subnet *VpcNatGateway
	req := c.R().
		SetBody(update).SetResult(&subnet)

	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s", NatGatewayEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return subnet, err
	}
	return subnet, nil
}

// DeleteNatGateway deletes a specific NatGateway by its identity.
func (c *Client) DeleteNatGateway(ctx context.Context, identity string) error {
	req := c.R()

	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s", NatGatewayEndpoint, identity))
	if err != nil {
		return err
	}
	if err := c.Check(resp); err != nil {
		return err
	}
	return nil
}
