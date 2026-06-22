package containerregistry

import (
	"context"
	"fmt"

	"github.com/thalassa-cloud/client-go/pkg/client"
)

// RunRetentionPolicy runs the retention policy for a container registry namespace.
func (c *Client) RunRetentionPolicy(ctx context.Context, namespaceIdentity string) error {
	if namespaceIdentity == "" {
		return fmt.Errorf("namespace identity is required")
	}

	req := c.R()
	resp, err := c.Do(ctx, req, client.POST, fmt.Sprintf("%s/%s/retention-policy/run", ContainerRegistryEndpoint, namespaceIdentity))
	if err != nil {
		return err
	}
	return c.Check(resp)
}
