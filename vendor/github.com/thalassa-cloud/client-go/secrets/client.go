package secrets

import (
	"github.com/thalassa-cloud/client-go/pkg/client"
)

// Client represents the Secrets Manager client.
//
// Create one via thalassa.NewClient(...).Secrets() or secrets.New(baseClient).
// See package documentation and Example functions for usage.
type Client struct {
	client.Client
}

// New creates a new Secrets Manager client.
func New(c client.Client, opts ...client.Option) (*Client, error) {
	c.WithOptions(opts...)
	return &Client{c}, nil
}
