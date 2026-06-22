package containerregistry

import (
	"context"
	"fmt"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

// ListContainerRegistryRepositories lists all repositories for a specific container registry namespace.
func (c *Client) ListContainerRegistryRepositories(ctx context.Context, namespaceIdentity string, listRequest *ListContainerRegistryRepositoriesRequest) ([]ContainerRegistryRepository, error) {
	if namespaceIdentity == "" {
		return nil, fmt.Errorf("namespace identity is required")
	}

	repositories := []ContainerRegistryRepository{}
	req := c.R().SetResult(&repositories)

	if listRequest != nil {
		for _, filter := range listRequest.Filters {
			for k, v := range filter.ToParams() {
				req = req.SetQueryParam(k, v)
			}
		}
	}

	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s/repositories", ContainerRegistryEndpoint, namespaceIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return repositories, err
	}
	return repositories, nil
}

// GetContainerRegistryRepository retrieves a specific container registry repository by its identity.
func (c *Client) GetContainerRegistryRepository(ctx context.Context, namespaceIdentity string, repositoryIdentity string) (*ContainerRegistryRepository, error) {
	if namespaceIdentity == "" {
		return nil, fmt.Errorf("namespace identity is required")
	}
	if repositoryIdentity == "" {
		return nil, fmt.Errorf("repository identity is required")
	}

	var repository *ContainerRegistryRepository
	req := c.R().SetResult(&repository)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s/repositories/%s", ContainerRegistryEndpoint, namespaceIdentity, repositoryIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return repository, err
	}
	return repository, nil
}

// DeleteContainerRegistryRepositoryWithAllArtifacts deletes a container registry repository and all its artifacts.
func (c *Client) DeleteContainerRegistryRepositoryWithAllArtifacts(ctx context.Context, namespaceIdentity string, repositoryIdentity string) error {
	if namespaceIdentity == "" {
		return fmt.Errorf("namespace identity is required")
	}
	if repositoryIdentity == "" {
		return fmt.Errorf("repository identity is required")
	}

	req := c.R()
	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s/repositories/%s", ContainerRegistryEndpoint, namespaceIdentity, repositoryIdentity))
	if err != nil {
		return err
	}
	return c.Check(resp)
}

// DeleteContainerRegistryRepositoryArtifact deletes artifacts from a container registry repository.
func (c *Client) DeleteContainerRegistryRepositoryArtifact(ctx context.Context, namespaceIdentity string, repositoryIdentity string) error {
	if namespaceIdentity == "" {
		return fmt.Errorf("namespace identity is required")
	}
	if repositoryIdentity == "" {
		return fmt.Errorf("repository identity is required")
	}

	req := c.R()
	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s/repositories/%s/delete-request", ContainerRegistryEndpoint, namespaceIdentity, repositoryIdentity))
	if err != nil {
		return err
	}
	return c.Check(resp)
}

// ListContainerRegistryRepositoriesRequest is the request for listing container registry repositories.
type ListContainerRegistryRepositoriesRequest struct {
	Filters []filters.Filter
}
