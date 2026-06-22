package secrets

import (
	"context"

	"github.com/thalassa-cloud/client-go/pkg/client"
)

// UpdateAccessPolicy updates the access policy for a secret.
func (c *Client) UpdateAccessPolicy(ctx context.Context, region, path string, update UpdateAccessPolicyRequest) (*Secret, error) {
	url, err := SecretResourceURL(region, path, "/policy")
	if err != nil {
		return nil, err
	}
	var secret Secret
	req := c.R().SetBody(update).SetResult(&secret)
	resp, err := c.Do(ctx, req, client.PUT, url)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &secret, err
	}
	return &secret, nil
}
