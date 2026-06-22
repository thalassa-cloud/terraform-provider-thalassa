package iaas

import (
	"context"
	"fmt"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

const (
	VolumeTypeEndpoint = "/v1/volume-types"
)

type ListVolumeTypesRequest struct {
	Filters []filters.Filter
}

// ListVolumeTypes lists all volume types.
func (c *Client) ListVolumeTypes(ctx context.Context, listRequest *ListVolumeTypesRequest) ([]VolumeType, error) {
	var volumeTypes []VolumeType
	req := c.R().SetResult(&volumeTypes)

	if listRequest != nil {
		for _, filter := range listRequest.Filters {
			for k, v := range filter.ToParams() {
				req = req.SetQueryParam(k, v)
			}
		}
	}

	resp, err := c.Do(ctx, req, client.GET, VolumeTypeEndpoint)
	if err != nil {
		return nil, err
	}

	if err := c.Check(resp); err != nil {
		return nil, err
	}

	return volumeTypes, nil
}

// GetVolumeType gets a volume type by its identity.
func (c *Client) GetVolumeType(ctx context.Context, identity string) (*VolumeType, error) {
	var volumeType *VolumeType
	req := c.R().SetResult(&volumeType)

	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", VolumeTypeEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return nil, err
	}
	return volumeType, nil
}
