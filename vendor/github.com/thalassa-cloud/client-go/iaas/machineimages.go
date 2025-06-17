package iaas

import (
	"context"
	"fmt"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

const (
	MachineImageEndpoint = "/v1/images"
)

type ListMachineImagesRequest struct {
	Filters []filters.Filter
}

// ListMachineImages lists all MachineImages for the current organisation.
// The current organisation is determined by the client's organisation identity.
func (c *Client) ListMachineImages(ctx context.Context, listRequest *ListMachineImagesRequest) ([]MachineImage, error) {
	machineImages := []MachineImage{}
	req := c.R().SetResult(&machineImages)

	if listRequest != nil {
		for _, filter := range listRequest.Filters {
			for k, v := range filter.ToParams() {
				req = req.SetQueryParam(k, v)
			}
		}
	}

	resp, err := c.Do(ctx, req, client.GET, MachineImageEndpoint)
	if err != nil {
		return nil, err
	}

	if err := c.Check(resp); err != nil {
		return machineImages, err
	}
	return machineImages, nil
}

// GetMachineImage retrieves a specific MachineImage by its identity.
// The identity is the unique identifier for the MachineImage.
func (c *Client) GetMachineImage(ctx context.Context, identity string) (*MachineImage, error) {
	var machineImage *MachineImage
	req := c.R().SetResult(&machineImage)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", MachineImageEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return machineImage, err
	}
	return machineImage, nil
}
