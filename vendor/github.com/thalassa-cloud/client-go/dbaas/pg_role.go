package dbaas

import (
	"context"
	"fmt"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

// PostgreSQL Role Operations

type ListPgRolesRequest struct {
	Filters []filters.Filter
}

// ListPgRoles lists all PostgreSQL roles for a database cluster.
func (c *Client) ListPgRoles(ctx context.Context, dbClusterIdentity string, listRequest *ListPgRolesRequest) ([]DbClusterPostgresRole, error) {
	if dbClusterIdentity == "" {
		return nil, fmt.Errorf("database cluster identity is required")
	}

	roles := []DbClusterPostgresRole{}
	req := c.R().SetResult(&roles)
	if listRequest != nil {
		for _, filter := range listRequest.Filters {
			for k, v := range filter.ToParams() {
				req = req.SetQueryParam(k, v)
			}
		}
	}
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s/postgres-roles", DbClusterEndpoint, dbClusterIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return nil, err
	}
	return roles, nil
}

// CreatePgRole creates a new PostgreSQL role in a database cluster.
func (c *Client) CreatePgRole(ctx context.Context, dbClusterIdentity string, create CreatePgRoleRequest) (*DbClusterPostgresRole, error) {
	if dbClusterIdentity == "" {
		return nil, fmt.Errorf("database cluster identity is required")
	}
	if create.Name == "" {
		return nil, fmt.Errorf("role name is required")
	}

	var role *DbClusterPostgresRole
	req := c.R().SetBody(create).SetResult(&role)
	resp, err := c.Do(ctx, req, client.POST, fmt.Sprintf("%s/%s/postgres-roles", DbClusterEndpoint, dbClusterIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return role, err
	}
	return role, nil
}

// UpdatePgRole updates an existing PostgreSQL role in a database cluster.
func (c *Client) UpdatePgRole(ctx context.Context, dbClusterIdentity string, postgresRoleIdentity string, update UpdatePgRoleRequest) (*DbClusterPostgresRole, error) {
	if dbClusterIdentity == "" {
		return nil, fmt.Errorf("database cluster identity is required")
	}
	if postgresRoleIdentity == "" {
		return nil, fmt.Errorf("postgres role identity is required")
	}

	var role *DbClusterPostgresRole
	req := c.R().SetBody(update).SetResult(&role)
	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s/postgres-roles/%s", DbClusterEndpoint, dbClusterIdentity, postgresRoleIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return role, err
	}
	return role, nil
}

// DeletePgRole deletes a PostgreSQL role from a database cluster.
func (c *Client) DeletePgRole(ctx context.Context, dbClusterIdentity string, postgresRoleIdentity string) error {
	if dbClusterIdentity == "" {
		return fmt.Errorf("database cluster identity is required")
	}
	if postgresRoleIdentity == "" {
		return fmt.Errorf("postgres role identity is required")
	}

	req := c.R()
	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s/postgres-roles/%s", DbClusterEndpoint, dbClusterIdentity, postgresRoleIdentity))
	if err != nil {
		return err
	}
	return c.Check(resp)
}

// CancelDeletePgRole cancels the deletion of a PostgreSQL role from a database cluster.
func (c *Client) CancelDeletePgRole(ctx context.Context, dbClusterIdentity string, postgresRoleIdentity string) error {
	if dbClusterIdentity == "" {
		return fmt.Errorf("database cluster identity is required")
	}
	if postgresRoleIdentity == "" {
		return fmt.Errorf("postgres role identity is required")
	}

	req := c.R()
	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s/postgres-roles/%s/cancel-deletion", DbClusterEndpoint, dbClusterIdentity, postgresRoleIdentity))
	if err != nil {
		return err
	}
	return c.Check(resp)
}
