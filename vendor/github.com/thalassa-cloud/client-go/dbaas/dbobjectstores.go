package dbaas

import (
	"context"
	"fmt"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

const (
	DbObjectStoreEndpoint = "/v1/dbaas/object-stores"
)

// ListDbObjectStoresRequest carries query filters for listing DB object stores.
type ListDbObjectStoresRequest struct {
	Filters []filters.Filter
}

// ListDbObjectStores lists all DB object stores for the organisation.
func (c *Client) ListDbObjectStores(ctx context.Context, listRequest *ListDbObjectStoresRequest) ([]DbObjectStore, error) {
	stores := []DbObjectStore{}
	req := c.R().SetResult(&stores)

	if listRequest != nil {
		for _, filter := range listRequest.Filters {
			for k, v := range filter.ToParams() {
				req = req.SetQueryParam(k, v)
			}
		}
	}

	resp, err := c.Do(ctx, req, client.GET, DbObjectStoreEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return stores, err
	}
	return stores, nil
}

// GetDbObjectStore returns a DB object store by identity.
func (c *Client) GetDbObjectStore(ctx context.Context, identity string) (*DbObjectStore, error) {
	if identity == "" {
		return nil, fmt.Errorf("identity is required")
	}

	var store *DbObjectStore
	req := c.R().SetResult(&store)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", DbObjectStoreEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return store, err
	}
	return store, nil
}

// CreateDbObjectStore creates a new DB object store.
func (c *Client) CreateDbObjectStore(ctx context.Context, create CreateDbObjectStoreRequest) (*DbObjectStore, error) {
	if create.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if create.Region == "" {
		return nil, fmt.Errorf("region is required")
	}

	var store *DbObjectStore
	req := c.R().SetBody(create).SetResult(&store)
	resp, err := c.Do(ctx, req, client.POST, DbObjectStoreEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return store, err
	}
	return store, nil
}

// UpdateDbObjectStore updates an existing DB object store.
func (c *Client) UpdateDbObjectStore(ctx context.Context, identity string, update UpdateDbObjectStoreRequest) (*DbObjectStore, error) {
	if identity == "" {
		return nil, fmt.Errorf("identity is required")
	}
	if update.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	var store *DbObjectStore
	req := c.R().SetBody(update).SetResult(&store)
	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s", DbObjectStoreEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return store, err
	}
	return store, nil
}

// DeleteDbObjectStore deletes a DB object store by identity.
func (c *Client) DeleteDbObjectStore(ctx context.Context, identity string) error {
	if identity == "" {
		return fmt.Errorf("identity is required")
	}

	req := c.R()
	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s", DbObjectStoreEndpoint, identity))
	if err != nil {
		return err
	}
	return c.Check(resp)
}
