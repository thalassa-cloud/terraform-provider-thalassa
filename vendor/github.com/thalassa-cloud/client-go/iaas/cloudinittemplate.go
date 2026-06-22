package iaas

import (
	"context"
	"fmt"
	"time"

	"github.com/thalassa-cloud/client-go/pkg/client"
)

// CloudInitTemplateEndpoint defines the API endpoint for cloud-init template operations
const (
	CloudInitTemplateEndpoint = "/v1/cloudinits"
)

// CreateCloudInitTemplate creates a new cloud-init template with the provided configuration.
// Cloud-init templates contain scripts that run when instances are first booted.
func (c *Client) CreateCloudInitTemplate(ctx context.Context, create CreateCloudInitTemplateRequest) (*CloudInitTemplate, error) {
	var cloudInitTemplate CloudInitTemplate
	req := c.R().SetBody(create).SetResult(&cloudInitTemplate)

	resp, err := c.Do(ctx, req, client.POST, CloudInitTemplateEndpoint)
	if err != nil {
		return nil, err
	}

	if err := c.Check(resp); err != nil {
		return nil, err
	}

	return &cloudInitTemplate, nil
}

// ListCloudInitTemplates retrieves all available cloud-init templates.
// Returns a slice of CloudInitTemplate objects that can be used for instance initialization.
func (c *Client) ListCloudInitTemplates(ctx context.Context) ([]CloudInitTemplate, error) {
	var cloudInitTemplates []CloudInitTemplate
	req := c.R().SetResult(&cloudInitTemplates)
	resp, err := c.Do(ctx, req, client.GET, CloudInitTemplateEndpoint)
	if err != nil {
		return nil, err
	}

	if err := c.Check(resp); err != nil {
		return nil, err
	}

	return cloudInitTemplates, nil
}

// GetCloudInitTemplate retrieves a specific cloud-init template by its unique identity.
// The identity is a UUID that uniquely identifies the template in the system.
func (c *Client) GetCloudInitTemplate(ctx context.Context, identity string) (*CloudInitTemplate, error) {
	var cloudInitTemplate CloudInitTemplate

	req := c.R().SetResult(&cloudInitTemplate)

	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", CloudInitTemplateEndpoint, identity))
	if err != nil {
		return nil, err
	}

	if err := c.Check(resp); err != nil {
		return nil, err
	}

	return &cloudInitTemplate, nil
}

// UpdateCloudInitTemplate updates an existing cloud-init template with the provided configuration.
func (c *Client) UpdateCloudInitTemplate(ctx context.Context, identity string, update UpdateCloudInitTemplateRequest) (*CloudInitTemplate, error) {
	var cloudInitTemplate CloudInitTemplate
	req := c.R().SetBody(update).SetResult(&cloudInitTemplate)
	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s", CloudInitTemplateEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return nil, err
	}
	return &cloudInitTemplate, nil
}

// DeleteCloudInitTemplate permanently removes a cloud-init template from the system.
// This operation cannot be undone and will affect any instances using this template.
func (c *Client) DeleteCloudInitTemplate(ctx context.Context, identity string) error {
	req := c.R()

	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s", CloudInitTemplateEndpoint, identity))
	if err != nil {
		return err
	}

	if err := c.Check(resp); err != nil {
		return err
	}

	return nil
}

// CloudInitTemplate represents a cloud-init template that can be used to configure instances.
// Cloud-init templates contain scripts that run during the first boot of a virtual machine.
type CloudInitTemplate struct {
	// Identity is the unique identifier (UUID) for this template
	Identity string `json:"identity"`

	// Name is the human-readable name of the template
	Name string `json:"name"`

	// Labels are key-value pairs for categorizing and organizing templates
	// Example: {"environment": "production", "team": "backend"}
	Labels map[string]string `json:"labels"`

	// Annotations are additional metadata that don't affect template behavior
	// Example: {"description": "Template for web servers", "author": "team-a"}
	Annotations map[string]string `json:"annotations"`

	// CreatedAt is the timestamp when this template was created
	CreatedAt time.Time `json:"created_at"`

	// Slug is a URL-friendly version of the template name
	// Used in API endpoints and for easier reference
	Slug string `json:"slug"`

	// Content contains the actual cloud-init script content
	// This can be a shell script, YAML configuration, or other cloud-init compatible format
	Content string `json:"content"`
}

// CreateCloudInitTemplateRequest defines the parameters needed to create a new cloud-init template.
type CreateCloudInitTemplateRequest struct {
	// Name is the display name for the template (required)
	// Should be descriptive and unique within your organization
	Name string `json:"name"`

	// Content is the cloud-init script content (required)
	// Can include shell scripts, YAML configurations, or other cloud-init directives
	// Example: "#cloud-config\npackages:\n  - nginx\n  - git"
	Content string `json:"content"`

	// Labels are optional key-value pairs for organizing templates
	// Useful for filtering and grouping related templates
	Labels Labels `json:"labels"`

	// Annotations are optional metadata that don't affect template behavior
	// Can include descriptions, usage notes, or other contextual information
	Annotations Annotations `json:"annotations"`
}

type UpdateCloudInitTemplateRequest struct {
	// Name is the display name for the template (required)
	// Should be descriptive and unique within your organization
	Name string `json:"name"`

	// Labels are optional key-value pairs for organizing templates
	// Useful for filtering and grouping related templates
	Labels Labels `json:"labels"`

	// Annotations are optional metadata that don't affect template behavior
	// Can include descriptions, usage notes, or other contextual information
	Annotations Annotations `json:"annotations"`

	// Content is the cloud-init script content (required)
	// Can include shell scripts, YAML configurations, or other cloud-init directives
	// Example: "#cloud-config\npackages:\n  - nginx\n  - git"
	Content string `json:"content"`
}
