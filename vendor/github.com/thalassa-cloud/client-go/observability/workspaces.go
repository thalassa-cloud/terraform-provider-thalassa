package observability

import (
	"context"
	"fmt"
	"time"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

// ListObservabilityWorkspacesRequest is the request for listing workspaces.
type ListObservabilityWorkspacesRequest struct {
	Filters []filters.Filter
}

// ListObservabilityWorkspaces lists observability workspaces for the organisation/project.
func (c *Client) ListObservabilityWorkspaces(ctx context.Context, listRequest *ListObservabilityWorkspacesRequest) ([]ObservabilityWorkspace, error) {
	workspaces := []ObservabilityWorkspace{}
	req := c.R().SetResult(&workspaces)

	if listRequest != nil {
		for _, filter := range listRequest.Filters {
			for k, v := range filter.ToParams() {
				req = req.SetQueryParam(k, v)
			}
		}
	}

	resp, err := c.Do(ctx, req, client.GET, WorkspaceEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return workspaces, err
	}
	return workspaces, nil
}

// GetObservabilityWorkspace retrieves a workspace by identity.
func (c *Client) GetObservabilityWorkspace(ctx context.Context, identity string) (*ObservabilityWorkspace, error) {
	if identity == "" {
		return nil, fmt.Errorf("workspace identity is required")
	}

	var workspace *ObservabilityWorkspace
	req := c.R().SetResult(&workspace)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", WorkspaceEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return workspace, err
	}
	return workspace, nil
}

// CreateObservabilityWorkspace creates a new observability workspace.
func (c *Client) CreateObservabilityWorkspace(ctx context.Context, create CreateObservabilityWorkspaceRequest) (*ObservabilityWorkspace, error) {
	if create.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if create.RetentionDays != nil {
		if err := validateRetentionDays(*create.RetentionDays); err != nil {
			return nil, err
		}
	}

	var workspace *ObservabilityWorkspace
	req := c.R().SetBody(create).SetResult(&workspace)
	resp, err := c.Do(ctx, req, client.POST, WorkspaceEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return workspace, err
	}
	return workspace, nil
}

// UpdateObservabilityWorkspace updates an existing workspace.
func (c *Client) UpdateObservabilityWorkspace(ctx context.Context, identity string, update UpdateObservabilityWorkspaceRequest) (*ObservabilityWorkspace, error) {
	if identity == "" {
		return nil, fmt.Errorf("workspace identity is required")
	}
	if update.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if update.RetentionDays != nil {
		if err := validateRetentionDays(*update.RetentionDays); err != nil {
			return nil, err
		}
	}

	var workspace *ObservabilityWorkspace
	req := c.R().SetBody(update).SetResult(&workspace)
	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s", WorkspaceEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return workspace, err
	}
	return workspace, nil
}

// DeleteObservabilityWorkspace deletes a workspace by identity.
func (c *Client) DeleteObservabilityWorkspace(ctx context.Context, identity string) error {
	if identity == "" {
		return fmt.Errorf("workspace identity is required")
	}

	req := c.R()
	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s", WorkspaceEndpoint, identity))
	if err != nil {
		return err
	}
	return c.Check(resp)
}

// WaitUntilObservabilityWorkspaceReady polls until the workspace status is ready or the context is cancelled.
func (c *Client) WaitUntilObservabilityWorkspaceReady(ctx context.Context, identity string) (*ObservabilityWorkspace, error) {
	if identity == "" {
		return nil, fmt.Errorf("workspace identity is required")
	}

	workspace, err := c.GetObservabilityWorkspace(ctx, identity)
	if err != nil {
		return nil, err
	}
	if workspace.Status == ObservabilityWorkspaceStatusReady {
		return workspace, nil
	}

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(5 * time.Second):
			workspace, err = c.GetObservabilityWorkspace(ctx, identity)
			if err != nil {
				return nil, err
			}
			if workspace.Status == ObservabilityWorkspaceStatusReady {
				return workspace, nil
			}
			if workspace.Status == ObservabilityWorkspaceStatusFailed {
				return workspace, fmt.Errorf("observability workspace %s failed: %s", identity, workspace.StatusMessage)
			}
		}
	}
}

func validateRetentionDays(days int) error {
	const minRetentionDays = 1
	const maxRetentionDays = 1095 // 3 years
	if days < minRetentionDays {
		return fmt.Errorf("retentionDays must be at least %d, got %d", minRetentionDays, days)
	}
	if days > maxRetentionDays {
		return fmt.Errorf("retentionDays cannot be greater than %d (3 years), got %d", maxRetentionDays, days)
	}
	return nil
}
