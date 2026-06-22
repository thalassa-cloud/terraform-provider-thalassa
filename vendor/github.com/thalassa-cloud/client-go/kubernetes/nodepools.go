package kubernetes

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

const (
	KubernetesNodePoolEndpoint = "nodepools"
)

func getKubernetesNodePoolEndpoint(identity string) string {
	return fmt.Sprintf("%s/%s/%s", KubernetesClusterEndpoint, identity, KubernetesNodePoolEndpoint)
}

type ListKubernetesNodePoolsRequest struct {
	Filters []filters.Filter
}

// ListKubernetesNodePools lists all KubernetesNodePools for a given organisation.
func (c *Client) ListKubernetesNodePools(ctx context.Context, clusterIdentity string, request *ListKubernetesNodePoolsRequest) ([]KubernetesNodePool, error) {
	subnets := []KubernetesNodePool{}
	req := c.R().SetResult(&subnets)
	if request != nil {
		for _, filter := range request.Filters {
			for k, v := range filter.ToParams() {
				req.SetQueryParam(k, v)
			}
		}
	}

	resp, err := c.Do(ctx, req, client.GET, getKubernetesNodePoolEndpoint(clusterIdentity))
	if err != nil {
		return nil, err
	}

	if err := c.Check(resp); err != nil {
		return subnets, err
	}
	return subnets, nil
}

// GetKubernetesNodePool retrieves a specific KubernetesNodePool by its identity.
func (c *Client) GetKubernetesNodePool(ctx context.Context, clusterIdentity string, identity string) (*KubernetesNodePool, error) {
	var subnet *KubernetesNodePool
	req := c.R().SetResult(&subnet)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", getKubernetesNodePoolEndpoint(clusterIdentity), identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return subnet, err
	}
	return subnet, nil
}

// CreateKubernetesNodePool creates a new KubernetesNodePool.
func (c *Client) CreateKubernetesNodePool(ctx context.Context, clusterIdentity string, create CreateKubernetesNodePool) (*KubernetesNodePool, error) {
	var subnet *KubernetesNodePool
	req := c.R().
		SetBody(create).SetResult(&subnet)

	resp, err := c.Do(ctx, req, client.POST, getKubernetesNodePoolEndpoint(clusterIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return subnet, err
	}
	return subnet, nil
}

// UpdateKubernetesNodePool updates an existing KubernetesNodePool.
func (c *Client) UpdateKubernetesNodePool(ctx context.Context, clusterIdentity string, identity string, update UpdateKubernetesNodePool) (*KubernetesNodePool, error) {
	var subnet *KubernetesNodePool
	req := c.R().
		SetBody(update).SetResult(&subnet)

	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s", getKubernetesNodePoolEndpoint(clusterIdentity), identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return subnet, err
	}
	return subnet, nil
}

// DeleteKubernetesNodePool deletes a specific KubernetesNodePool by its identity.
func (c *Client) DeleteKubernetesNodePool(ctx context.Context, clusterIdentity string, identity string) error {
	req := c.R()

	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s", getKubernetesNodePoolEndpoint(clusterIdentity), identity))
	if err != nil {
		return err
	}
	if err := c.Check(resp); err != nil {
		return err
	}
	return nil
}

// WaitUntilKubernetesNodePoolReady waits until the node pool is ready. A node pool is ready when all nodes within the node pool are up-to-date and running.
// It returns the node pool when it is ready or an error if the node pool is not being ready or if the context is cancelled.
// You are responsible for providing a context that can be cancelled, and for handling the error case.
// Example: ctxt, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
// defer cancel()
// nodePool, err := c.WaitUntilKubernetesNodePoolReady(ctxt, "cluster-identity1234", "node-pool-identity1234")
//
//	if err != nil {
//		log.Fatalf("Failed to wait for node pool to be ready: %v", err)
//	}
func (c *Client) WaitUntilKubernetesNodePoolReady(ctx context.Context, clusterIdentity string, identity string) (*KubernetesNodePool, error) {
	nodePool, err := c.GetKubernetesNodePool(ctx, clusterIdentity, identity)
	if err != nil {
		return nil, err
	}
	if nodePool.Status == KubernetesNodePoolStatusReady {
		return nodePool, nil
	}

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(5 * time.Second):
			nodePool, err = c.GetKubernetesNodePool(ctx, clusterIdentity, identity)
			if err != nil {
				return nil, err
			}
			if nodePool.Status == KubernetesNodePoolStatusReady {
				return nodePool, nil
			}
		}
	}
}

// WaitUntilKubernetesNodePoolDeleted waits until the node pool is deleted.
// It returns an error if the node pool is not being deleted or if the context is cancelled.
// You are responsible for providing a context that can be cancelled, and for handling the error case.
// Example: ctxt, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
// defer cancel()
// err := c.WaitUntilKubernetesNodePoolDeleted(ctxt, "cluster-identity1234", "node-pool-identity1234")
//
//	if err != nil {
//		log.Fatalf("Failed to wait for node pool to be deleted: %v", err)
//	}
func (c *Client) WaitUntilKubernetesNodePoolDeleted(ctx context.Context, clusterIdentity string, identity string) error {
	nodePool, err := c.GetKubernetesNodePool(ctx, clusterIdentity, identity)
	if err != nil {
		return err
	}
	if nodePool == nil {
		return nil
	}
	if nodePool.Status == KubernetesNodePoolStatusDeleted {
		return nil
	}
	if nodePool.Status != KubernetesNodePoolStatusDeleting {
		return fmt.Errorf("node pool %s is not being deleted", identity)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(5 * time.Second):
			nodePool, err := c.GetKubernetesNodePool(ctx, clusterIdentity, identity)
			if err != nil {
				if errors.Is(err, client.ErrNotFound) {
					return nil
				}
				return err
			}
			if nodePool.Status == KubernetesNodePoolStatusDeleted {
				return nil
			}
			if nodePool.Status != KubernetesNodePoolStatusDeleting {
				return fmt.Errorf("node pool %s is not being deleted", identity)
			}
		}
	}
}
