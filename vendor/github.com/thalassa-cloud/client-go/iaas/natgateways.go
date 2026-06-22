package iaas

import (
	"context"
	"errors"
	"fmt"
	"time"

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

// WaitUntilNatGatewayHasEndpoint waits until the nat gateway has an endpoint.
// It returns the nat gateway when it has an endpoint or an error if the nat gateway fails to get an endpoint.
// You are responsible for providing a context that can be cancelled, and for handling the error case.
// Example: ctxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// natGateway, err := c.WaitUntilNatGatewayHasEndpoint(ctxt, "nat-gateway-identity1234")
//
//	if err != nil {
//		log.Fatalf("Failed to wait for nat gateway to have an endpoint: %v", err)
//	}
//	defer cancel()
func (c *Client) WaitUntilNatGatewayHasEndpoint(ctx context.Context, identity string) (*VpcNatGateway, error) {
	natGateway, err := c.GetNatGateway(ctx, identity)
	if err != nil {
		return nil, err
	}
	if natGateway.EndpointIP != "" {
		return natGateway, nil
	}
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(DefaultPollIntervalForWaiting):
			natGateway, err := c.GetNatGateway(ctx, identity)
			if err != nil {
				return nil, err
			}
			if natGateway.EndpointIP != "" {
				return natGateway, nil
			}
		}
	}
}

// WaitUntilNatGatewayDeleted waits until the nat gateway is deleted.
// It returns an error if the nat gateway fails to delete.
// You are responsible for providing a context that can be cancelled, and for handling the error case.
// Example: ctxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// err := c.WaitUntilNatGatewayDeleted(ctxt, "nat-gateway-identity1234")
//
//	if err != nil {
//		log.Fatalf("Failed to wait for nat gateway to be deleted: %v", err)
//	}
//	defer cancel()
func (c *Client) WaitUntilNatGatewayDeleted(ctx context.Context, identity string) error {
	natGateway, err := c.GetNatGateway(ctx, identity)
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			return nil
		}
		return err
	}
	if natGateway.Status == "deleted" {
		return nil
	}
	if natGateway.Status != "deleting" {
		return fmt.Errorf("nat gateway %s is not being deleted", identity)
	}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(DefaultPollIntervalForWaiting):
			natGateway, err := c.GetNatGateway(ctx, identity)
			if err != nil {
				if errors.Is(err, client.ErrNotFound) {
					return nil
				}
				return err
			}
			if natGateway.Status == "deleted" {
				return nil
			}
			if natGateway.Status != "deleting" {
				return fmt.Errorf("nat gateway %s is not being deleted", identity)
			}
		}
	}
}
