package iaas

import (
	"context"
	"fmt"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

const (
	MachineTypeEndpoint = "/v1/machine-types"
)

type ListMachineTypesRequest struct {
	Filters []filters.Filter
}

// ListMachineTypes lists all MachineTypes for the current organisation.
// The current organisation is determined by the client's organisation identity.
func (c *Client) ListMachineTypes(ctx context.Context, listRequest *ListMachineTypesRequest) ([]MachineType, error) {
	machineTypes := []MachineType{}
	req := c.R().SetResult(&machineTypes)

	if listRequest != nil {
		for _, filter := range listRequest.Filters {
			for k, v := range filter.ToParams() {
				req = req.SetQueryParam(k, v)
			}
		}
	}

	resp, err := c.Do(ctx, req, client.GET, MachineTypeEndpoint)
	if err != nil {
		return nil, err
	}

	if err := c.Check(resp); err != nil {
		return machineTypes, err
	}
	return machineTypes, nil
}

// GetMachineType retrieves a specific MachineType by its identity.
// The identity is the unique identifier for the MachineType.
func (c *Client) GetMachineType(ctx context.Context, identity string) (*MachineType, error) {
	var machineType *MachineType
	req := c.R().SetResult(&machineType)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", MachineTypeEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return machineType, err
	}
	return machineType, nil
}

// ListMachineTypeCategories lists all MachineTypeCategories for the current organisation.
// The current organisation is determined by the client's organisation identity.
func (c *Client) ListMachineTypeCategories(ctx context.Context) ([]MachineTypeCategory, error) {
	machineTypeCategories := []MachineTypeCategory{}
	req := c.R().SetResult(&machineTypeCategories)

	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/by-categories", MachineTypeEndpoint))
	if err != nil {
		return nil, err
	}

	if err := c.Check(resp); err != nil {
		return machineTypeCategories, err
	}

	return machineTypeCategories, nil
}
