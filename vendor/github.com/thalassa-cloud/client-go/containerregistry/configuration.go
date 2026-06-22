package containerregistry

import (
	"context"
	"fmt"

	"github.com/thalassa-cloud/client-go/pkg/client"
)

// GetNamespaceConfiguration retrieves the configuration for a container registry namespace.
func (c *Client) GetNamespaceConfiguration(ctx context.Context, namespaceIdentity string) (*ContainerRegistryNamespaceConfiguration, error) {
	if namespaceIdentity == "" {
		return nil, fmt.Errorf("namespace identity is required")
	}

	var configuration *ContainerRegistryNamespaceConfiguration
	req := c.R().SetResult(&configuration)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s/configuration", ContainerRegistryEndpoint, namespaceIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return configuration, err
	}
	return configuration, nil
}

// CreateNamespaceConfiguration creates a configuration for a container registry namespace.
func (c *Client) CreateNamespaceConfiguration(ctx context.Context, namespaceIdentity string, create CreateNamespaceConfigurationRequest) (*ContainerRegistryNamespaceConfiguration, error) {
	if namespaceIdentity == "" {
		return nil, fmt.Errorf("namespace identity is required")
	}

	var configuration *ContainerRegistryNamespaceConfiguration
	req := c.R().SetBody(create).SetResult(&configuration)
	resp, err := c.Do(ctx, req, client.POST, fmt.Sprintf("%s/%s/configuration", ContainerRegistryEndpoint, namespaceIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return configuration, err
	}
	return configuration, nil
}

// UpdateNamespaceConfiguration updates the configuration for a container registry namespace.
func (c *Client) UpdateNamespaceConfiguration(ctx context.Context, namespaceIdentity string, update UpdateNamespaceConfigurationRequest) (*ContainerRegistryNamespaceConfiguration, error) {
	if namespaceIdentity == "" {
		return nil, fmt.Errorf("namespace identity is required")
	}

	var configuration *ContainerRegistryNamespaceConfiguration
	req := c.R().SetBody(update).SetResult(&configuration)
	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s/configuration", ContainerRegistryEndpoint, namespaceIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return configuration, err
	}
	return configuration, nil
}

// DeleteNamespaceConfiguration deletes the configuration for a container registry namespace.
func (c *Client) DeleteNamespaceConfiguration(ctx context.Context, namespaceIdentity string) error {
	if namespaceIdentity == "" {
		return fmt.Errorf("namespace identity is required")
	}

	req := c.R()
	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s/configuration", ContainerRegistryEndpoint, namespaceIdentity))
	if err != nil {
		return err
	}
	return c.Check(resp)
}
