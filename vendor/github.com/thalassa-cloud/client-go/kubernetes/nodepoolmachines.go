package kubernetes

import (
	"context"
	"fmt"

	"github.com/thalassa-cloud/client-go/pkg/client"
)

func getKubernetesNodePoolMachinesEndpoint(clusterIdentity string, nodePoolIdentity string) string {
	return fmt.Sprintf("%s/%s/%s/%s/machines", KubernetesClusterEndpoint, clusterIdentity, KubernetesNodePoolEndpoint, nodePoolIdentity)
}

// ListNodePoolMachines lists all machines for a specific node pool in a cluster
func (c *Client) ListNodePoolMachines(ctx context.Context, clusterIdentity string, nodePoolIdentity string) ([]KubernetesNodePoolMachine, error) {
	machines := []KubernetesNodePoolMachine{}
	req := c.R().SetResult(&machines)

	resp, err := c.Do(ctx, req, client.GET, getKubernetesNodePoolMachinesEndpoint(clusterIdentity, nodePoolIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return machines, err
	}
	return machines, nil
}

// DeleteNodePoolMachine deletes a specific machine from a node pool in a cluster.
// If scaleDown is true, the node pool will be scaled down by one as part of the deletion.
func (c *Client) DeleteNodePoolMachine(ctx context.Context, clusterIdentity string, nodePoolIdentity string, machineIdentity string, scaleDown bool) error {
	req := c.R()
	if scaleDown {
		req.SetQueryParam("scaleDown", "true")
	}

	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s", getKubernetesNodePoolMachinesEndpoint(clusterIdentity, nodePoolIdentity), machineIdentity))
	if err != nil {
		return err
	}
	if err := c.Check(resp); err != nil {
		return err
	}
	return nil
}
