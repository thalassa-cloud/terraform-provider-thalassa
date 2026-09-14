package dbaas

import (
	"context"
	"fmt"

	"github.com/thalassa-cloud/client-go/pkg/client"
)

// ListDbClusterRevisions lists revisions for a database cluster.
func (c *Client) ListDbClusterRevisions(ctx context.Context, dbClusterIdentity string) ([]DbClusterRevision, error) {
	if dbClusterIdentity == "" {
		return nil, fmt.Errorf("database cluster identity is required")
	}

	revisions := []DbClusterRevision{}
	req := c.R().SetResult(&revisions)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s/revisions", DbClusterEndpoint, dbClusterIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return revisions, err
	}
	return revisions, nil
}

// ListScheduledMaintenances lists scheduled maintenances for a database cluster.
func (c *Client) ListScheduledMaintenances(ctx context.Context, dbClusterIdentity string) ([]DbClusterScheduledMaintenance, error) {
	if dbClusterIdentity == "" {
		return nil, fmt.Errorf("database cluster identity is required")
	}

	maintenances := []DbClusterScheduledMaintenance{}
	req := c.R().SetResult(&maintenances)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s/scheduled-maintenances", DbClusterEndpoint, dbClusterIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return maintenances, err
	}
	return maintenances, nil
}

// StartScheduledMaintenance starts a scheduled maintenance immediately.
func (c *Client) StartScheduledMaintenance(ctx context.Context, dbClusterIdentity, maintenanceIdentity string) (*DbClusterScheduledMaintenance, error) {
	if dbClusterIdentity == "" {
		return nil, fmt.Errorf("database cluster identity is required")
	}
	if maintenanceIdentity == "" {
		return nil, fmt.Errorf("maintenance identity is required")
	}

	var maintenance *DbClusterScheduledMaintenance
	req := c.R().SetResult(&maintenance)
	resp, err := c.Do(ctx, req, client.POST, fmt.Sprintf("%s/%s/scheduled-maintenances/%s/start", DbClusterEndpoint, dbClusterIdentity, maintenanceIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return maintenance, err
	}
	return maintenance, nil
}

// PostponeScheduledMaintenance postpones a scheduled maintenance to the next maintenance window.
func (c *Client) PostponeScheduledMaintenance(ctx context.Context, dbClusterIdentity, maintenanceIdentity string) (*DbClusterScheduledMaintenance, error) {
	if dbClusterIdentity == "" {
		return nil, fmt.Errorf("database cluster identity is required")
	}
	if maintenanceIdentity == "" {
		return nil, fmt.Errorf("maintenance identity is required")
	}

	var maintenance *DbClusterScheduledMaintenance
	req := c.R().SetResult(&maintenance)
	resp, err := c.Do(ctx, req, client.POST, fmt.Sprintf("%s/%s/scheduled-maintenances/%s/postpone", DbClusterEndpoint, dbClusterIdentity, maintenanceIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return maintenance, err
	}
	return maintenance, nil
}
