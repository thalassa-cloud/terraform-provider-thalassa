package kubernetes

import (
	"context"
	"fmt"
	"time"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

const (
	KubernetesClusterEndpoint = "/v1/kubernetes/clusters"
)

type ListKubernetesClustersRequest struct {
	Filters []filters.Filter
}

// ListKubernetesClusters lists all KubernetesClusters for a given organisation.
func (c *Client) ListKubernetesClusters(ctx context.Context, request *ListKubernetesClustersRequest) ([]KubernetesCluster, error) {
	subnets := []KubernetesCluster{}
	req := c.R().SetResult(&subnets)
	if request != nil {
		for _, filter := range request.Filters {
			for k, v := range filter.ToParams() {
				req.SetQueryParam(k, v)
			}
		}
	}

	resp, err := c.Do(ctx, req, client.GET, KubernetesClusterEndpoint)
	if err != nil {
		return nil, err
	}

	if err := c.Check(resp); err != nil {
		return subnets, err
	}
	return subnets, nil
}

// GetKubernetesCluster retrieves a specific KubernetesCluster by its identity.
func (c *Client) GetKubernetesCluster(ctx context.Context, identity string) (*KubernetesCluster, error) {
	var subnet *KubernetesCluster
	req := c.R().SetResult(&subnet)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", KubernetesClusterEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return subnet, err
	}
	return subnet, nil
}

// CreateKubernetesCluster creates a new KubernetesCluster.
func (c *Client) CreateKubernetesCluster(ctx context.Context, create CreateKubernetesCluster) (*KubernetesCluster, error) {
	var subnet *KubernetesCluster
	req := c.R().
		SetBody(create).SetResult(&subnet)

	resp, err := c.Do(ctx, req, client.POST, KubernetesClusterEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return subnet, err
	}
	return subnet, nil
}

// WaitUntilKubernetesClusterReady waits until the KubernetesCluster is ready.
// It returns the KubernetesCluster when it is ready or an error if the KubernetesCluster is not being ready or if the context is cancelled.
// You are responsible for providing a context that can be cancelled, and for handling the error case.
// Example: ctxt, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
// defer cancel()
// kubernetesCluster, err := c.WaitUntilKubernetesClusterReady(ctxt, "kubernetes-cluster-identity1234")
//
//	if err != nil {
//		log.Fatalf("Failed to wait for KubernetesCluster to be ready: %v", err)
//	}
func (c *Client) WaitUntilKubernetesClusterReady(ctx context.Context, identity string) (*KubernetesCluster, error) {
	kubernetesCluster, err := c.GetKubernetesCluster(ctx, identity)
	if err != nil {
		return nil, err
	}
	if kubernetesCluster.Status == "ready" {
		return kubernetesCluster, nil
	}

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(5 * time.Second):
			kubernetesCluster, err = c.GetKubernetesCluster(ctx, identity)
			if err != nil {
				return nil, err
			}
			if kubernetesCluster.Status == "ready" {
				return kubernetesCluster, nil
			}
		}
	}
}

// UpdateKubernetesCluster updates an existing KubernetesCluster.
func (c *Client) UpdateKubernetesCluster(ctx context.Context, identity string, update UpdateKubernetesCluster) (*KubernetesCluster, error) {
	var subnet *KubernetesCluster
	req := c.R().
		SetBody(update).SetResult(&subnet)

	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s", KubernetesClusterEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return subnet, err
	}
	return subnet, nil
}

// DeleteKubernetesCluster deletes a specific KubernetesCluster by its identity.
func (c *Client) DeleteKubernetesCluster(ctx context.Context, identity string) error {
	req := c.R()

	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s", KubernetesClusterEndpoint, identity))
	if err != nil {
		return err
	}
	if err := c.Check(resp); err != nil {
		return err
	}
	return nil
}
