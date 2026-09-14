package dbaas

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/thalassa-cloud/client-go/objectstorage"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

// BrowseDbObjectStoreStorage lists objects and prefixes in a DBaaS backup object store.
// Prefix is relative to the object store destination root.
func (c *Client) BrowseDbObjectStoreStorage(ctx context.Context, identity string, browseRequest *objectstorage.BrowseObjectsRequest) (*objectstorage.BrowseResponse, error) {
	if identity == "" {
		return nil, fmt.Errorf("identity is required")
	}

	var result *objectstorage.BrowseResponse
	req := c.R().SetResult(&result)
	applyBrowseObjectsQuery(req, browseRequest)

	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s/browse", DbObjectStoreEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return result, err
	}
	return result, nil
}

// GetDbObjectStoreObjectMetadata returns metadata for an object in a DBaaS backup object store.
// Key is relative to the object store destination root.
func (c *Client) GetDbObjectStoreObjectMetadata(ctx context.Context, identity, key string) (*objectstorage.ObjectMetadata, error) {
	if identity == "" {
		return nil, fmt.Errorf("identity is required")
	}
	if key == "" {
		return nil, fmt.Errorf("object key is required")
	}

	var result *objectstorage.ObjectMetadata
	req := c.R().SetQueryParam("key", key).SetResult(&result)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s/browse/objects", DbObjectStoreEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return result, err
	}
	return result, nil
}

// GetDbObjectStoreRecoveryOverview returns the aggregated recovery overview for a DB object store.
func (c *Client) GetDbObjectStoreRecoveryOverview(ctx context.Context, identity string, overviewRequest *GetDbObjectStoreRecoveryOverviewRequest) (*DbObjectStoreRecoveryOverview, error) {
	if identity == "" {
		return nil, fmt.Errorf("identity is required")
	}

	var result *DbObjectStoreRecoveryOverview
	req := c.R().SetResult(&result)
	if overviewRequest != nil {
		if overviewRequest.IncludeUnavailable != nil {
			req = req.SetQueryParam("includeUnavailable", strconv.FormatBool(*overviewRequest.IncludeUnavailable))
		}
		if overviewRequest.ClusterID != "" {
			req = req.SetQueryParam("clusterId", overviewRequest.ClusterID)
		}
		if overviewRequest.BackupID != "" {
			req = req.SetQueryParam("backupId", overviewRequest.BackupID)
		}
	}

	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s/recovery", DbObjectStoreEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return result, err
	}
	return result, nil
}

// GetDbObjectStoreRecoveryWalHierarchy returns the WAL hierarchy for a cluster in a DB object store.
func (c *Client) GetDbObjectStoreRecoveryWalHierarchy(ctx context.Context, identity, clusterID string, hierarchyRequest *GetDbObjectStoreRecoveryWalHierarchyRequest) (*DbObjectStoreRecoveryWalHierarchy, error) {
	if identity == "" {
		return nil, fmt.Errorf("identity is required")
	}
	if clusterID == "" {
		return nil, fmt.Errorf("cluster identity is required")
	}

	var result *DbObjectStoreRecoveryWalHierarchy
	req := c.R().SetResult(&result)
	if hierarchyRequest != nil {
		if hierarchyRequest.Timeline != "" {
			req = req.SetQueryParam("timeline", hierarchyRequest.Timeline)
		}
		if len(hierarchyRequest.Kinds) > 0 {
			req = req.SetQueryParam("kinds", strings.Join(hierarchyRequest.Kinds, ","))
		}
	}

	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s/recovery/clusters/%s/wal", DbObjectStoreEndpoint, identity, clusterID))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return result, err
	}
	return result, nil
}

// ListDbObjectStoreRecoveryWalSegments lists WAL segment objects (timeline+log, or view=flat).
func (c *Client) ListDbObjectStoreRecoveryWalSegments(ctx context.Context, identity, clusterID string, listRequest *ListDbObjectStoreRecoveryWalSegmentsRequest) (*DbObjectStoreRecoveryWalSegments, error) {
	if identity == "" {
		return nil, fmt.Errorf("identity is required")
	}
	if clusterID == "" {
		return nil, fmt.Errorf("cluster identity is required")
	}
	if listRequest == nil {
		return nil, fmt.Errorf("list request is required")
	}
	if listRequest.View != "flat" && (listRequest.Timeline == "" || listRequest.Log == "") {
		return nil, fmt.Errorf("timeline and log are required unless view is flat")
	}

	var result *DbObjectStoreRecoveryWalSegments
	req := c.R().SetResult(&result)
	if listRequest.View != "" {
		req = req.SetQueryParam("view", listRequest.View)
	}
	if listRequest.Timeline != "" {
		req = req.SetQueryParam("timeline", listRequest.Timeline)
	}
	if listRequest.Log != "" {
		req = req.SetQueryParam("log", listRequest.Log)
	}
	if listRequest.ContinuationToken != "" {
		req = req.SetQueryParam("continuationToken", listRequest.ContinuationToken)
	}
	if listRequest.Limit > 0 {
		req = req.SetQueryParam("limit", strconv.Itoa(listRequest.Limit))
	}
	if len(listRequest.Kinds) > 0 {
		req = req.SetQueryParam("kinds", strings.Join(listRequest.Kinds, ","))
	}

	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s/recovery/clusters/%s/wal", DbObjectStoreEndpoint, identity, clusterID))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return result, err
	}
	return result, nil
}

// SearchDbObjectStoreRecoveryWal searches WAL archive objects by hex substring or WAL name.
func (c *Client) SearchDbObjectStoreRecoveryWal(ctx context.Context, identity, clusterID string, searchRequest *SearchDbObjectStoreRecoveryWalRequest) (*DbObjectStoreRecoveryWalSearch, error) {
	if identity == "" {
		return nil, fmt.Errorf("identity is required")
	}
	if clusterID == "" {
		return nil, fmt.Errorf("cluster identity is required")
	}
	if searchRequest == nil || searchRequest.Query == "" {
		return nil, fmt.Errorf("search query is required")
	}

	var result *DbObjectStoreRecoveryWalSearch
	req := c.R().SetQueryParam("q", searchRequest.Query).SetResult(&result)
	if searchRequest.Limit > 0 {
		req = req.SetQueryParam("limit", strconv.Itoa(searchRequest.Limit))
	}

	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s/recovery/clusters/%s/wal/search", DbObjectStoreEndpoint, identity, clusterID))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return result, err
	}
	return result, nil
}

func applyBrowseObjectsQuery(req *resty.Request, browseRequest *objectstorage.BrowseObjectsRequest) {
	if browseRequest == nil {
		return
	}
	if browseRequest.Prefix != "" {
		req.SetQueryParam("prefix", browseRequest.Prefix)
	}
	if browseRequest.Delimiter != "" {
		req.SetQueryParam("delimiter", browseRequest.Delimiter)
	}
	if browseRequest.MaxKeys > 0 {
		req.SetQueryParam("maxKeys", strconv.FormatInt(int64(browseRequest.MaxKeys), 10))
	}
	if browseRequest.ContinuationToken != "" {
		req.SetQueryParam("continuationToken", browseRequest.ContinuationToken)
	}
}
