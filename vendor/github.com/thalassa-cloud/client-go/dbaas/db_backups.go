package dbaas

import (
	"context"
	"fmt"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

const (
	DbBackupEndpoint = "/v1/dbaas/backups"
)

// ListDbBackupsForDbCluster lists all backups for a specific database cluster.
func (c *Client) ListDbBackupsForDbCluster(ctx context.Context, dbClusterIdentity string, listRequest *ListDbBackupsRequest) ([]DbClusterBackup, error) {
	if dbClusterIdentity == "" {
		return nil, fmt.Errorf("database cluster identity is required")
	}

	backups := []DbClusterBackup{}
	req := c.R().SetResult(&backups)

	if listRequest != nil {
		for _, filter := range listRequest.Filters {
			for k, v := range filter.ToParams() {
				req = req.SetQueryParam(k, v)
			}
		}
	}

	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s/backups", DbClusterEndpoint, dbClusterIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return backups, err
	}
	return backups, nil
}

// CreateDbBackup creates a new backup for a database cluster.
func (c *Client) CreateDbBackup(ctx context.Context, dbClusterIdentity string, create CreateDbClusterBackupRequest) (*DbClusterBackup, error) {
	if dbClusterIdentity == "" {
		return nil, fmt.Errorf("database cluster identity is required")
	}
	if create.Name == "" {
		return nil, fmt.Errorf("backup name is required")
	}

	var backup *DbClusterBackup
	req := c.R().SetBody(create).SetResult(&backup)
	resp, err := c.Do(ctx, req, client.POST, fmt.Sprintf("%s/%s/backups", DbClusterEndpoint, dbClusterIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return backup, err
	}
	return backup, nil
}

// ListDbBackupsForOrganisation lists all backups for the organisation.
func (c *Client) ListDbBackupsForOrganisation(ctx context.Context, listRequest *ListDbBackupsRequest) ([]DbClusterBackup, error) {
	backups := []DbClusterBackup{}
	req := c.R().SetResult(&backups)

	if listRequest != nil {
		for _, filter := range listRequest.Filters {
			for k, v := range filter.ToParams() {
				req = req.SetQueryParam(k, v)
			}
		}
	}

	resp, err := c.Do(ctx, req, client.GET, DbBackupEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return backups, err
	}
	return backups, nil
}

// GetDbBackup retrieves a specific backup by its identity.
func (c *Client) GetDbBackup(ctx context.Context, backupIdentity string) (*DbClusterBackup, error) {
	if backupIdentity == "" {
		return nil, fmt.Errorf("backup identity is required")
	}

	var backup *DbClusterBackup
	req := c.R().SetResult(&backup)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", DbBackupEndpoint, backupIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return backup, err
	}
	return backup, nil
}

// DeleteDbBackup deletes a specific backup by its identity.
func (c *Client) DeleteDbBackup(ctx context.Context, backupIdentity string) error {
	if backupIdentity == "" {
		return fmt.Errorf("backup identity is required")
	}

	req := c.R()
	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s", DbBackupEndpoint, backupIdentity))
	if err != nil {
		return err
	}
	return c.Check(resp)
}

// CancelDeleteDbBackup cancels the deletion of a backup.
func (c *Client) CancelDeleteDbBackup(ctx context.Context, backupIdentity string) error {
	if backupIdentity == "" {
		return fmt.Errorf("backup identity is required")
	}

	req := c.R()
	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s/cancel-deletion", DbBackupEndpoint, backupIdentity))
	if err != nil {
		return err
	}
	return c.Check(resp)
}

// ListDbBackupsRequest is the request for listing backups.
type ListDbBackupsRequest struct {
	Filters []filters.Filter
}
