package kubernetes

import (
	"context"
	"fmt"

	"github.com/thalassa-cloud/client-go/pkg/client"
)

const (
	KubernetesClusterKubeConfigEndpoint = "/v1/kubernetes/clusters/%s/kubeconfig"
)

func (c *Client) GetKubernetesClusterKubeconfig(ctx context.Context, identity string) (*KubernetesClusterSessionToken, error) {
	var clusterSession *KubernetesClusterSessionToken
	req := c.R().SetResult(&clusterSession)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf(KubernetesClusterKubeConfigEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return clusterSession, err
	}
	return clusterSession, nil
}
