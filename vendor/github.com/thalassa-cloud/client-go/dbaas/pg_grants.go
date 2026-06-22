package dbaas

import (
	"context"
	"fmt"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

// DBaaS Cluster Grant Operations

type ListDbGrantsRequest struct {
	Filters []filters.Filter
}

// ListDbGrants lists all DBaaS Cluster grants for a database cluster.
func (c *Client) ListDbGrants(ctx context.Context, dbClusterIdentity string, listRequest *ListDbGrantsRequest) ([]DbClusterPostgresGrant, error) {
	if dbClusterIdentity == "" {
		return nil, fmt.Errorf("database cluster identity is required")
	}

	grants := []DbClusterPostgresGrant{}
	req := c.R().SetResult(&grants)
	if listRequest != nil {
		for _, filter := range listRequest.Filters {
			for k, v := range filter.ToParams() {
				req = req.SetQueryParam(k, v)
			}
		}
	}

	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s/postgres-grants", DbClusterEndpoint, dbClusterIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return nil, err
	}
	return grants, nil
}

// CreatePgGrant creates a new PostgreSQL grant for a role on a database in a database cluster.
func (c *Client) CreatePgGrant(ctx context.Context, dbClusterIdentity string, create CreatePgGrantRequest) (*DbClusterPostgresGrant, error) {
	if dbClusterIdentity == "" {
		return nil, fmt.Errorf("database cluster identity is required")
	}
	if create.Name == "" {
		return nil, fmt.Errorf("grant name is required")
	}
	if create.RoleName == "" {
		return nil, fmt.Errorf("role name is required")
	}
	if create.DatabaseName == "" {
		return nil, fmt.Errorf("database name is required")
	}

	var grant *DbClusterPostgresGrant
	req := c.R().SetBody(create).SetResult(&grant)
	resp, err := c.Do(ctx, req, client.POST, fmt.Sprintf("%s/%s/postgres-grants", DbClusterEndpoint, dbClusterIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return grant, err
	}
	return grant, nil
}

// UpdatePgGrant updates an existing PostgreSQL grant for a role on a database in a database cluster.
func (c *Client) UpdatePgGrant(ctx context.Context, dbClusterIdentity string, grantIdentity string, update UpdatePgGrantRequest) (*DbClusterPostgresGrant, error) {
	if dbClusterIdentity == "" {
		return nil, fmt.Errorf("database cluster identity is required")
	}
	if grantIdentity == "" {
		return nil, fmt.Errorf("postgres grant identity is required")
	}

	var grant *DbClusterPostgresGrant
	req := c.R().SetBody(update).SetResult(&grant)
	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s/postgres-grants/%s", DbClusterEndpoint, dbClusterIdentity, grantIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return grant, err
	}
	return grant, nil
}

// DeletePgGrant deletes a PostgreSQL grant from a database cluster.
func (c *Client) DeletePgGrant(ctx context.Context, dbClusterIdentity string, grantIdentity string) error {
	if dbClusterIdentity == "" {
		return fmt.Errorf("database cluster identity is required")
	}
	if grantIdentity == "" {
		return fmt.Errorf("postgres grant identity is required")
	}

	req := c.R()
	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s/postgres-grants/%s", DbClusterEndpoint, dbClusterIdentity, grantIdentity))
	if err != nil {
		return err
	}
	return c.Check(resp)
}
