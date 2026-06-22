package dbaas

import (
	"context"
	"fmt"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

type DatabaseInstanceTypeCategory struct {
	// Name is the name of the database instance type category.
	Name string `json:"name,omitempty"`
	// Description is the description of the database instance type category.
	Description string `json:"description,omitempty"`
	// InstanceTypes is the list of database instance types in the category.
	InstanceTypes []DatabaseInstanceType `json:"instanceTypes,omitempty"`
}

type DatabaseInstanceType struct {
	// Name is the name of the database instance type.
	Name string `json:"name,omitempty"`
	// Slug is the slug of the database instance type.
	Slug string `json:"slug,omitempty"`
	// Identity is the identity of the database instance type.
	Identity string `json:"identity,omitempty"`
	// Description is the description of the database instance type.
	Description string `json:"description,omitempty"`
	// CategorySlug is the slug of the database instance type category.
	CategorySlug string `json:"categorySlug,omitempty"`
	// Cpus is the number of CPUs of the database instance type.
	Cpus int `json:"cpus,omitempty"`
	// Memory is the memory of the database instance type. In GB.
	Memory int `json:"memory,omitempty"`
	// MaxStorage is the maximum storage of the database instance type. In GB.
	MaxStorage int `json:"maxStorage,omitempty"`
	// Bandwidth is the bandwidth of the database instance type. In MB/s.
	Bandwidth int `json:"bandwidth,omitempty"`
	// Architecture is the architecture of the database instance type.
	Architecture string `json:"architecture,omitempty"`
}

const (
	DatabaseInstanceTypeEndpoint = "/v1/dbaas/instance-types"
)

type ListDatabaseInstanceTypesRequest struct {
	Filters []filters.Filter
}

// ListDatabaseInstanceTypes lists all DatabaseInstanceTypes for the current organisation.
// The current organisation is determined by the client's organisation identity.
func (c *Client) ListDatabaseInstanceTypes(ctx context.Context, listRequest *ListDatabaseInstanceTypesRequest) ([]DatabaseInstanceType, error) {
	databaseInstanceTypes := []DatabaseInstanceType{}
	req := c.R().SetResult(&databaseInstanceTypes)

	if listRequest != nil {
		for _, filter := range listRequest.Filters {
			for k, v := range filter.ToParams() {
				req = req.SetQueryParam(k, v)
			}
		}
	}

	resp, err := c.Do(ctx, req, client.GET, DatabaseInstanceTypeEndpoint)
	if err != nil {
		return nil, err
	}

	if err := c.Check(resp); err != nil {
		return databaseInstanceTypes, err
	}
	return databaseInstanceTypes, nil
}

// GetDatabaseInstanceType retrieves a specific DatabaseInstanceType by its identity.
// The identity is the unique identifier for the DatabaseInstanceType.
func (c *Client) GetDatabaseInstanceType(ctx context.Context, identity string) (*DatabaseInstanceType, error) {
	var databaseInstanceType *DatabaseInstanceType
	req := c.R().SetResult(&databaseInstanceType)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", DatabaseInstanceTypeEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return databaseInstanceType, err
	}
	return databaseInstanceType, nil
}

// ListDatabaseInstanceTypeCategories lists all DatabaseInstanceTypeCategories for the current organisation.
// The current organisation is determined by the client's organisation identity.
func (c *Client) ListDatabaseInstanceTypeCategories(ctx context.Context) ([]DatabaseInstanceTypeCategory, error) {
	databaseInstanceTypeCategories := []DatabaseInstanceTypeCategory{}
	req := c.R().SetResult(&databaseInstanceTypeCategories)

	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/by-categories", DatabaseInstanceTypeEndpoint))
	if err != nil {
		return nil, err
	}

	if err := c.Check(resp); err != nil {
		return databaseInstanceTypeCategories, err
	}

	return databaseInstanceTypeCategories, nil
}
