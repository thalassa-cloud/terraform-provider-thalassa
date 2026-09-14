package kubernetes

import (
	"context"
	"fmt"

	"github.com/thalassa-cloud/client-go/pkg/client"
)

// ListKubeconfigSessions lists active kubeconfig sessions for a cluster.
func (c *Client) ListKubeconfigSessions(ctx context.Context, clusterIdentity string) ([]KubernetesClusterSession, error) {
	if clusterIdentity == "" {
		return nil, fmt.Errorf("cluster identity is required")
	}

	sessions := []KubernetesClusterSession{}
	req := c.R().SetResult(&sessions)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s/kubeconfig-sessions", KubernetesClusterEndpoint, clusterIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return sessions, err
	}
	return sessions, nil
}

// DeleteKubeconfigSession revokes a kubeconfig session for a cluster.
func (c *Client) DeleteKubeconfigSession(ctx context.Context, clusterIdentity, sessionIdentity string) error {
	if clusterIdentity == "" {
		return fmt.Errorf("cluster identity is required")
	}
	if sessionIdentity == "" {
		return fmt.Errorf("session identity is required")
	}

	req := c.R()
	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s/kubeconfig-sessions/%s", KubernetesClusterEndpoint, clusterIdentity, sessionIdentity))
	if err != nil {
		return err
	}
	return c.Check(resp)
}
