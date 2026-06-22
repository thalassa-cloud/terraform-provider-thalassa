package iaas

import (
	"context"
	"fmt"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

const (
	RegionEndpoint = "/v1/regions"
)

type ListRegionsRequest struct {
	Filters []filters.Filter
}

// ListRegions lists all Regions for a given organisation.
func (c *Client) ListRegions(ctx context.Context, listRequest *ListRegionsRequest) ([]Region, error) {
	vpcs := []Region{}
	req := c.R().SetResult(&vpcs)

	if listRequest != nil {
		for _, filter := range listRequest.Filters {
			for k, v := range filter.ToParams() {
				req = req.SetQueryParam(k, v)
			}
		}
	}

	resp, err := c.Do(ctx, req, client.GET, RegionEndpoint)
	if err != nil {
		return nil, err
	}

	if err := c.Check(resp); err != nil {
		return vpcs, err
	}
	return vpcs, nil
}

// GetRegion retrieves a specific Region by its identity.
func (c *Client) GetRegion(ctx context.Context, identity string) (*Region, error) {
	var vpc *Region
	req := c.R().SetResult(&vpc)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", RegionEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return vpc, err
	}
	return vpc, nil
}
