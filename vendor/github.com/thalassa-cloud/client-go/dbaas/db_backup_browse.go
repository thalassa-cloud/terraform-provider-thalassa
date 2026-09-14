package dbaas

import (
	"context"
	"fmt"

	"github.com/thalassa-cloud/client-go/objectstorage"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

// BrowseDbClusterBackupStorage lists objects and prefixes under a cluster's backup path
// in its linked DB object store. Prefix is relative to the cluster backup root.
func (c *Client) BrowseDbClusterBackupStorage(ctx context.Context, clusterIdentity string, browseRequest *objectstorage.BrowseObjectsRequest) (*objectstorage.BrowseResponse, error) {
	if clusterIdentity == "" {
		return nil, fmt.Errorf("cluster identity is required")
	}

	var result *objectstorage.BrowseResponse
	req := c.R().SetResult(&result)
	applyBrowseObjectsQuery(req, browseRequest)

	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s/backups/browse", DbClusterEndpoint, clusterIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return result, err
	}
	return result, nil
}

// GetDbClusterBackupObjectMetadata returns metadata for an object under a cluster's backup path.
// Key is relative to the cluster backup root.
func (c *Client) GetDbClusterBackupObjectMetadata(ctx context.Context, clusterIdentity, key string) (*objectstorage.ObjectMetadata, error) {
	if clusterIdentity == "" {
		return nil, fmt.Errorf("cluster identity is required")
	}
	if key == "" {
		return nil, fmt.Errorf("object key is required")
	}

	var result *objectstorage.ObjectMetadata
	req := c.R().SetQueryParam("key", key).SetResult(&result)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s/backups/browse/objects", DbClusterEndpoint, clusterIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return result, err
	}
	return result, nil
}
