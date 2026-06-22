package iaas

import (
	"context"
	"fmt"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

const ReservedIPEndpoint = "/v1/reserved-ips"

// ListReservedIPs lists reserved IPs for the organisation (auth context).
func (c *Client) ListReservedIPs(ctx context.Context, listRequest *ListReservedIPsRequest) ([]ReservedIP, error) {
	var out []ReservedIP
	req := c.R().SetResult(&out)
	if listRequest != nil {
		for _, filter := range listRequest.Filters {
			for k, v := range filter.ToParams() {
				req = req.SetQueryParam(k, v)
			}
		}
	}
	resp, err := c.Do(ctx, req, client.GET, ReservedIPEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return out, err
	}
	return out, nil
}

// CreateReservedIP creates a reserved IP (201).
func (c *Client) CreateReservedIP(ctx context.Context, create CreateReservedIpRequest) (*ReservedIP, error) {
	if create.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if create.Region == "" {
		return nil, fmt.Errorf("region is required")
	}
	var fip *ReservedIP
	req := c.R().SetBody(create).SetResult(&fip)
	resp, err := c.Do(ctx, req, client.POST, ReservedIPEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return fip, err
	}
	return fip, nil
}

// GetReservedIP returns a reserved IP by identity.
func (c *Client) GetReservedIP(ctx context.Context, identity string) (*ReservedIP, error) {
	if identity == "" {
		return nil, fmt.Errorf("identity is required")
	}
	var fip *ReservedIP
	req := c.R().SetResult(&fip)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", ReservedIPEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return fip, err
	}
	return fip, nil
}

// UpdateReservedIP updates name, description, labels, and annotations.
func (c *Client) UpdateReservedIP(ctx context.Context, identity string, update UpdateReservedIpRequest) (*ReservedIP, error) {
	if identity == "" {
		return nil, fmt.Errorf("identity is required")
	}
	if update.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	var fip *ReservedIP
	req := c.R().SetBody(update).SetResult(&fip)
	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s", ReservedIPEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return fip, err
	}
	return fip, nil
}

// DeleteReservedIP deletes a reserved IP. If attached, the API disassociates first (204).
func (c *Client) DeleteReservedIP(ctx context.Context, identity string) error {
	if identity == "" {
		return fmt.Errorf("identity is required")
	}
	resp, err := c.Do(ctx, c.R(), client.DELETE, fmt.Sprintf("%s/%s", ReservedIPEndpoint, identity))
	if err != nil {
		return err
	}
	return c.Check(resp)
}

// AssociateReservedIP attaches the reserved IP to a load balancer or NAT gateway.
func (c *Client) AssociateReservedIP(ctx context.Context, identity string, body AssociateReservedIpRequest) (*ReservedIP, error) {
	if identity == "" {
		return nil, fmt.Errorf("identity is required")
	}
	var fip *ReservedIP
	req := c.R().SetBody(body).SetResult(&fip)
	resp, err := c.Do(ctx, req, client.POST, fmt.Sprintf("%s/%s/associate", ReservedIPEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return fip, err
	}
	return fip, nil
}

// DisassociateReservedIP detaches the reserved IP from its current target.
func (c *Client) DisassociateReservedIP(ctx context.Context, identity string) (*ReservedIP, error) {
	if identity == "" {
		return nil, fmt.Errorf("identity is required")
	}
	var fip *ReservedIP
	req := c.R().SetResult(&fip)
	resp, err := c.Do(ctx, req, client.POST, fmt.Sprintf("%s/%s/disassociate", ReservedIPEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return fip, err
	}
	return fip, nil
}

// ListReservedIPsRequest holds query filters for ListReservedIPs.
// Use filters.FilterKeyValue with filters.FilterReservedIp, FilterName, FilterIdentity,
// FilterSlug, FilterStatus, FilterRegion, or LabelFilter as for other IaaS list APIs.
type ListReservedIPsRequest struct {
	Filters []filters.Filter
}
