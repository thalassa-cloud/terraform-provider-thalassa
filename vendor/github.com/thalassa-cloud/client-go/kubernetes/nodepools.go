package kubernetes

import (
	"context"
	"fmt"

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
