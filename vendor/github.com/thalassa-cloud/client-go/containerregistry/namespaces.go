package containerregistry

import (
	"context"
	"fmt"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

// ListContainerRegistryNamespaces lists all container registry namespaces for the organisation.
func (c *Client) ListContainerRegistryNamespaces(ctx context.Context, listRequest *ListContainerRegistryNamespacesRequest) ([]ContainerRegistryNamespace, error) {
	namespaces := []ContainerRegistryNamespace{}
	req := c.R().SetResult(&namespaces)

	if listRequest != nil {
		for _, filter := range listRequest.Filters {
			for k, v := range filter.ToParams() {
				req = req.SetQueryParam(k, v)
			}
		}
	}

	resp, err := c.Do(ctx, req, client.GET, ContainerRegistryEndpoint+"")
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return namespaces, err
	}
	return namespaces, nil
}

// GetContainerRegistryNamespace retrieves a specific container registry namespace by its identity.
func (c *Client) GetContainerRegistryNamespace(ctx context.Context, namespaceIdentity string) (*ContainerRegistryNamespace, error) {
	if namespaceIdentity == "" {
		return nil, fmt.Errorf("namespace identity is required")
	}

	var namespace *ContainerRegistryNamespace
	req := c.R().SetResult(&namespace)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", ContainerRegistryEndpoint, namespaceIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return namespace, err
	}
	return namespace, nil
}

// CreateContainerRegistryNamespace creates a new container registry namespace.
func (c *Client) CreateContainerRegistryNamespace(ctx context.Context, create CreateContainerRegistryNamespaceRequest) (*ContainerRegistryNamespace, error) {
	if create.Namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}
	if create.Region == "" {
		return nil, fmt.Errorf("region is required")
	}

	var namespace *ContainerRegistryNamespace
	req := c.R().SetBody(create).SetResult(&namespace)
	resp, err := c.Do(ctx, req, client.POST, ContainerRegistryEndpoint+"")
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return namespace, err
	}
	return namespace, nil
}

// UpdateContainerRegistryNamespace updates an existing container registry namespace.
func (c *Client) UpdateContainerRegistryNamespace(ctx context.Context, namespaceIdentity string, update UpdateContainerRegistryNamespaceRequest) (*ContainerRegistryNamespace, error) {
	if namespaceIdentity == "" {
		return nil, fmt.Errorf("namespace identity is required")
	}

	var namespace *ContainerRegistryNamespace
	req := c.R().SetBody(update).SetResult(&namespace)
	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s", ContainerRegistryEndpoint, namespaceIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return namespace, err
	}
	return namespace, nil
}

// DeleteContainerRegistryNamespace deletes a specific container registry namespace by its identity.
func (c *Client) DeleteContainerRegistryNamespace(ctx context.Context, namespaceIdentity string) error {
	if namespaceIdentity == "" {
		return fmt.Errorf("namespace identity is required")
	}

	req := c.R()
	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s", ContainerRegistryEndpoint, namespaceIdentity))
	if err != nil {
		return err
	}
	return c.Check(resp)
}

// ListContainerRegistryNamespacesRequest is the request for listing container registry namespaces.
type ListContainerRegistryNamespacesRequest struct {
	Filters []filters.Filter
}
