package observability

import (
	"github.com/thalassa-cloud/client-go/pkg/client"
)

// Client provides access to Thalassa Cloud Observability workspaces.
type Client struct {
	client.Client
}

// New creates a new Observability client.
func New(c client.Client, opts ...client.Option) (*Client, error) {
	c.WithOptions(opts...)
	return &Client{c}, nil
}

const (
	// WorkspaceEndpoint is the base path for observability workspaces.
	WorkspaceEndpoint = "/v1/observability/workspaces"
)
