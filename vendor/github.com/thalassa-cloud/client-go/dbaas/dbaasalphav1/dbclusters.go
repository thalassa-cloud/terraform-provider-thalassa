package dbaasalphav1

import (
	"context"
	"fmt"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

const (
	DbClusterEndpoint = "/v1/dbaas/dbclusters"
)

// ListDbClusters lists all dbClusters for a given organisation.
func (c *Client) ListDbClusters(ctx context.Context, listRequest *ListDbClustersRequest) ([]DbCluster, error) {
	dbClusters := []DbCluster{}
	req := c.R().SetResult(&dbClusters)

	if listRequest != nil {
		for _, filter := range listRequest.Filters {
			for k, v := range filter.ToParams() {
				req = req.SetQueryParam(k, v)
			}
		}
	}

	resp, err := c.Do(ctx, req, client.GET, DbClusterEndpoint)
	if err != nil {
		return nil, err
	}

	if err := c.Check(resp); err != nil {
		return dbClusters, err
	}
	return dbClusters, nil
}

// GetDbCluster retrieves a specific dbCluster by its identity.
func (c *Client) GetDbCluster(ctx context.Context, dbClusterIdentity string) (*DbCluster, error) {
	if dbClusterIdentity == "" {
		return nil, fmt.Errorf("identity is required")
	}

	var dbCluster *DbCluster
	req := c.R().SetResult(&dbCluster)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", DbClusterEndpoint, dbClusterIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return dbCluster, err
	}
	return dbCluster, nil
}

// CreateDbCluster creates a new dbCluster.
func (c *Client) CreateDbCluster(ctx context.Context, create CreateDbClusterRequest) (*DbCluster, error) {
	if create.SubnetIdentity == "" {
		return nil, fmt.Errorf("subnet is required")
	}
	if create.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	var dbCluster *DbCluster
	req := c.R().
		SetBody(create).SetResult(&dbCluster)

	resp, err := c.Do(ctx, req, client.POST, DbClusterEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return dbCluster, err
	}
	return dbCluster, nil
}

// UpdateDbCluster updates an existing dbCluster.
func (c *Client) UpdateDbCluster(ctx context.Context, dbClusterIdentity string, update UpdateDbClusterRequest) (*DbCluster, error) {
	if dbClusterIdentity == "" {
		return nil, fmt.Errorf("identity of the dbCluster to update is required")
	}

	var dbCluster *DbCluster
	req := c.R().
		SetBody(update).SetResult(&dbCluster)

	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s", DbClusterEndpoint, dbClusterIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return dbCluster, err
	}
	return dbCluster, nil
}

// DeleteDbCluster deletes a specific dbCluster by its identity.
func (c *Client) DeleteDbCluster(ctx context.Context, dbClusterIdentity string) error {
	if dbClusterIdentity == "" {
		return fmt.Errorf("identity is required")
	}

	req := c.R()

	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s", DbClusterEndpoint, dbClusterIdentity))
	if err != nil {
		return err
	}
	if err := c.Check(resp); err != nil {
		return err
	}
	return nil
}

type ListDbClustersRequest struct {
	Filters []filters.Filter
}
