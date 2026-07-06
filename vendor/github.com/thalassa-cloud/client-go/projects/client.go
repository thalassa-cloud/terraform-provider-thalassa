package projects

import (
	"github.com/thalassa-cloud/client-go/pkg/client"
)

// Client represents the Projects API client.
type Client struct {
	client.Client
}

// New creates a new Projects client.
func New(c client.Client, opts ...client.Option) (*Client, error) {
	c.WithOptions(opts...)
	return &Client{c}, nil
}
