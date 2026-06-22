package containerregistry

import (
	"github.com/thalassa-cloud/client-go/pkg/client"
)

type Client struct {
	client.Client
}

func New(c client.Client, opts ...client.Option) (*Client, error) {
	c.WithOptions(opts...)
	return &Client{c}, nil
}

const (
	ContainerRegistryEndpoint = "/v1/container-registry"
)
