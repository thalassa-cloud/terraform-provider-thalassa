package dns

import (
	"github.com/thalassa-cloud/client-go/pkg/client"
)

// Client represents the DNS client.
//
// Create one via thalassa.NewClient(...).DNS() or dns.New(baseClient).
// See package documentation and Example functions for usage.
type Client struct {
	client.Client
}

// New creates a new DNS client.
func New(c client.Client, opts ...client.Option) (*Client, error) {
	c.WithOptions(opts...)
	return &Client{c}, nil
}
