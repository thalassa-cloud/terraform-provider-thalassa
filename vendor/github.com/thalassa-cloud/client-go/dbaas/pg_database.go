package dbaas

import (
	"context"
	"fmt"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

// PostgreSQL Database Operations
func (c *Client) ListPgDatabases(ctx context.Context, dbClusterIdentity string, listRequest *ListPgDatabasesRequest) ([]DbClusterPostgresDatabase, error) {
	if dbClusterIdentity == "" {
		return nil, fmt.Errorf("database cluster identity is required")
	}

	databases := []DbClusterPostgresDatabase{}
	req := c.R().SetResult(&databases)
	if listRequest != nil {
		for _, filter := range listRequest.Filters {
			for k, v := range filter.ToParams() {
				req = req.SetQueryParam(k, v)
			}
		}
	}
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s/postgres-databases", DbClusterEndpoint, dbClusterIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return nil, err
	}
	return databases, nil
}

type ListPgDatabasesRequest struct {
	Filters []filters.Filter
}

// CreatePgDatabase creates a new PostgreSQL database in a database cluster.
func (c *Client) CreatePgDatabase(ctx context.Context, dbClusterIdentity string, create CreatePgDatabaseRequest) (*DbClusterPostgresDatabase, error) {
	if dbClusterIdentity == "" {
		return nil, fmt.Errorf("database cluster identity is required")
	}
	if create.Name == "" {
		return nil, fmt.Errorf("database name is required")
	}
	if create.Owner == "" {
		return nil, fmt.Errorf("database owner is required")
	}

	var database *DbClusterPostgresDatabase
	req := c.R().SetBody(create).SetResult(&database)
	resp, err := c.Do(ctx, req, client.POST, fmt.Sprintf("%s/%s/postgres-databases", DbClusterEndpoint, dbClusterIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return database, err
	}
	return database, nil
}

// UpdatePgDatabase updates an existing PostgreSQL database in a database cluster.
func (c *Client) UpdatePgDatabase(ctx context.Context, dbClusterIdentity string, postgresDatabaseIdentity string, update UpdatePgDatabaseRequest) (*DbClusterPostgresDatabase, error) {
	if dbClusterIdentity == "" {
		return nil, fmt.Errorf("database cluster identity is required")
	}
	if postgresDatabaseIdentity == "" {
		return nil, fmt.Errorf("postgres database identity is required")
	}

	var database *DbClusterPostgresDatabase
	req := c.R().SetBody(update).SetResult(&database)
	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s/postgres-databases/%s", DbClusterEndpoint, dbClusterIdentity, postgresDatabaseIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return database, err
	}
	return database, nil
}

// DeletePgDatabase deletes a PostgreSQL database from a database cluster.
func (c *Client) DeletePgDatabase(ctx context.Context, dbClusterIdentity string, postgresDatabaseIdentity string, immediate bool) error {
	if dbClusterIdentity == "" {
		return fmt.Errorf("database cluster identity is required")
	}
	if postgresDatabaseIdentity == "" {
		return fmt.Errorf("postgres database identity is required")
	}

	req := c.R()
	if immediate {
		req.SetQueryParam("deleteImmediately", "true") // If true, the database will be deleted immediately. If false, the database will be deleted after a soft deletion grace period.
	}
	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s/postgres-databases/%s", DbClusterEndpoint, dbClusterIdentity, postgresDatabaseIdentity))
	if err != nil {
		return err
	}
	return c.Check(resp)
}

// CancelDeletePgDatabase cancels the deletion of a PostgreSQL database from a database cluster.
func (c *Client) CancelDeletePgDatabase(ctx context.Context, dbClusterIdentity string, postgresDatabaseIdentity string) error {
	if dbClusterIdentity == "" {
		return fmt.Errorf("database cluster identity is required")
	}
	if postgresDatabaseIdentity == "" {
		return fmt.Errorf("postgres database identity is required")
	}

	req := c.R()
	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s/postgres-databases/%s/cancel-deletion", DbClusterEndpoint, dbClusterIdentity, postgresDatabaseIdentity))
	if err != nil {
		return err
	}
	return c.Check(resp)
}
