package projects

import (
	"context"
	"fmt"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

const (
	ProjectEndpoint = "/v1/projects"
)

// ListProjectsRequest configures filters for GET /v1/projects.
type ListProjectsRequest struct {
	Filters []filters.Filter
}

// ListProjects lists projects visible to the caller within the organisation.
func (c *Client) ListProjects(ctx context.Context, request *ListProjectsRequest) ([]Project, error) {
	projects := []Project{}
	req := c.R().SetResult(&projects)
	if request != nil {
		for _, filter := range request.Filters {
			for k, v := range filter.ToParams() {
				req.SetQueryParam(k, v)
			}
		}
	}

	resp, err := c.Do(client.WithoutProject(ctx), req, client.GET, ProjectEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return nil, err
	}
	return projects, nil
}

// GetProject retrieves a project by ID or slug.
func (c *Client) GetProject(ctx context.Context, identity string) (*Project, error) {
	if identity == "" {
		return nil, fmt.Errorf("identity is required")
	}

	project := Project{}
	req := c.R().SetResult(&project)
	resp, err := c.Do(client.WithoutProject(ctx), req, client.GET, fmt.Sprintf("%s/%s", ProjectEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return nil, err
	}
	return &project, nil
}

// CreateProject creates a new project in the organisation.
func (c *Client) CreateProject(ctx context.Context, create CreateProjectRequest) (*Project, error) {
	if create.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	project := Project{}
	req := c.R().SetBody(create).SetResult(&project)
	resp, err := c.Do(client.WithoutProject(ctx), req, client.POST, ProjectEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return nil, err
	}
	return &project, nil
}

// UpdateProject replaces a project's fields. identity accepts identity or slug.
func (c *Client) UpdateProject(ctx context.Context, identity string, update UpdateProjectRequest) (*Project, error) {
	if identity == "" {
		return nil, fmt.Errorf("identity is required")
	}
	if update.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	project := Project{}
	req := c.R().SetBody(update).SetResult(&project)
	resp, err := c.Do(client.WithoutProject(ctx), req, client.PUT, fmt.Sprintf("%s/%s", ProjectEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return nil, err
	}
	return &project, nil
}

// DeleteProject deletes a project by identity or slug.
func (c *Client) DeleteProject(ctx context.Context, identity string) error {
	if identity == "" {
		return fmt.Errorf("identity is required")
	}

	resp, err := c.Do(client.WithoutProject(ctx), c.R(), client.DELETE, fmt.Sprintf("%s/%s", ProjectEndpoint, identity))
	if err != nil {
		return err
	}
	return c.Check(resp)
}
