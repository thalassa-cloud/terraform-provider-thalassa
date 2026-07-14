package dbaas

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/thalassa-cloud/client-go/dbaas"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/client-go/thalassa"
)

const (
	dbClusterReadyPollInterval  = time.Second
	dbClusterDeletePollInterval = time.Second
)

func waitForReadyDbCluster(ctx context.Context, client thalassa.Client, dbClusterID string) (*dbaas.DbCluster, error) {
	for {
		select {
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				return nil, fmt.Errorf("timeout waiting for db cluster %q to become ready", dbClusterID)
			}
			return nil, ctx.Err()
		default:
		}

		dbCluster, err := client.DBaaS().GetDbCluster(ctx, dbClusterID)
		if err != nil {
			return nil, err
		}
		if dbCluster == nil {
			return nil, fmt.Errorf("db cluster %q not found", dbClusterID)
		}
		if dbCluster.Status == dbaas.DbClusterStatusReady {
			return dbCluster, nil
		}

		time.Sleep(dbClusterReadyPollInterval)
	}
}

func waitForDeletedDbCluster(ctx context.Context, client thalassa.Client, dbClusterID string) error {
	for {
		select {
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				return fmt.Errorf("timeout waiting for db cluster %q to be deleted", dbClusterID)
			}
			return ctx.Err()
		default:
		}

		dbCluster, err := client.DBaaS().GetDbCluster(ctx, dbClusterID)
		if err != nil {
			if tcclient.IsNotFound(err) {
				return nil
			}
			return err
		}
		if dbCluster == nil {
			return nil
		}

		switch dbCluster.Status {
		case dbaas.DbClusterStatusDeleted:
			return nil
		case dbaas.DbClusterStatusDeleting:
			tflog.Debug(ctx, "db cluster deletion in progress", map[string]any{
				"cluster_id": dbClusterID,
				"status":     dbCluster.Status,
			})
		case dbaas.DbClusterStatusReady, dbaas.DbClusterStatusUpdating:
			// The delete API can return before the cluster transitions to deleting.
			if err := client.DBaaS().DeleteDbCluster(ctx, dbClusterID); err != nil && !tcclient.IsNotFound(err) {
				return fmt.Errorf("failed to re-issue delete for db cluster %q: %w", dbClusterID, err)
			}
		case dbaas.DbClusterStatusFailed:
			return fmt.Errorf("db cluster %q failed to delete (status: %s)", dbClusterID, dbCluster.Status)
		default:
			tflog.Debug(ctx, "waiting for db cluster deletion", map[string]any{
				"cluster_id": dbClusterID,
				"status":     dbCluster.Status,
			})
		}

		time.Sleep(dbClusterDeletePollInterval)
	}
}
