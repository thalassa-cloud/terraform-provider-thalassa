package kubernetes

import (
	"context"
	"fmt"

	"github.com/thalassa-cloud/client-go/pkg/client"
)

const (
	KubernetesVersionEndpoint = "/v1/kubernetes/versions"
)

// ListKubernetesVersions lists all KubernetesVersions for a given organisation.
func (c *Client) ListKubernetesVersions(ctx context.Context) ([]KubernetesVersion, error) {
	items := []KubernetesVersion{}
	req := c.R().SetResult(&items)

	resp, err := c.Do(ctx, req, client.GET, KubernetesVersionEndpoint)
	if err != nil {
		return nil, err
	}

	if err := c.Check(resp); err != nil {
		return items, err
	}
	return items, nil
}

// GetKubernetesVersion retrieves a specific KubernetesVersion by its identity.
func (c *Client) GetKubernetesVersion(ctx context.Context, identity string) (*KubernetesVersion, error) {
	var subnet *KubernetesVersion
	req := c.R().SetResult(&subnet)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", KubernetesVersionEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return subnet, err
	}
	return subnet, nil
}
