package dbaasalphav1

import (
	"context"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

const (
	DbEngineEndpoint = "/v1/dbaas/engine-versions"
)

// ListDatabaseEngines lists all Database engines available for the organisation.
func (c *Client) ListDatabaseEngines(ctx context.Context, listRequest *ListDatabaseEnginesRequest) (*ListDatabaseEnginesResponse, error) {
	dbEngines := ListDatabaseEnginesResponse{}
	req := c.R().SetResult(&dbEngines)

	if listRequest != nil {
		for _, filter := range listRequest.Filters {
			for k, v := range filter.ToParams() {
				req = req.SetQueryParam(k, v)
			}
		}
	}

	resp, err := c.Do(ctx, req, client.GET, DbEngineEndpoint)
	if err != nil {
		return nil, err
	}

	if err := c.Check(resp); err != nil {
		return &dbEngines, err
	}
	return &dbEngines, nil
}

// ListDatabaseEnginesRequest is the request for the ListDatabaseEngines function.
type ListDatabaseEnginesRequest struct {
	Filters []filters.Filter
}

type ListDatabaseEnginesResponse struct {
	Engines map[DbClusterDatabaseEngine][]DbClusterEngineVersion `json:"engines"`
}
