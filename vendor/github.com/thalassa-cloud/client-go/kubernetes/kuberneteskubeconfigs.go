package kubernetes

import (
	"context"
	"fmt"
	"strings"

	"github.com/thalassa-cloud/client-go/pkg/client"
)

const (
	// KubernetesClusterKubeConfigEndpoint is the API endpoint for retrieving kubeconfig for a Kubernetes cluster.
	KubernetesClusterKubeConfigEndpoint = "/v1/kubernetes/clusters/%s/kubeconfig"
)

// GetKubernetesClusterKubeconfig retrieves the kubeconfig and authentication details for a Kubernetes cluster.
// This method returns a KubernetesClusterSessionToken containing all necessary information to connect to the cluster,
// including the API server URL, authentication token, CA certificate, and complete kubeconfig file content.
//
// Parameters:
//   - ctx: Context for the request, supporting cancellation and timeouts
//   - identity: The unique identifier of the Kubernetes cluster
//
// Returns:
//   - *KubernetesClusterSessionToken: Contains all cluster connection details
//   - error: Returns an error if the request fails, the cluster is not found, or the identity is invalid
//
// Example usage:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//	defer cancel()
//
//	session, err := client.GetKubernetesClusterKubeconfig(ctx, "cluster-123")
//	if err != nil {
//	    log.Fatalf("Failed to get kubeconfig: %v", err)
//	}
func (c *Client) GetKubernetesClusterKubeconfig(ctx context.Context, identity string) (*KubernetesClusterSessionToken, error) {
	// Validate input parameters
	if ctx == nil {
		return nil, fmt.Errorf("context cannot be nil")
	}

	if strings.TrimSpace(identity) == "" {
		return nil, fmt.Errorf("cluster identity cannot be empty")
	}

	var clusterSession *KubernetesClusterSessionToken
	req := c.R().SetResult(&clusterSession)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf(KubernetesClusterKubeConfigEndpoint, identity))
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve kubeconfig for cluster %s: %w", identity, err)
	}
	if err := c.Check(resp); err != nil {
		return nil, fmt.Errorf("failed to get kubeconfig for cluster %s: %w", identity, err)
	}
	return clusterSession, nil
}
