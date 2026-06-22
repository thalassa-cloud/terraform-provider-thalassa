package secrets

import (
	"context"

	"github.com/thalassa-cloud/client-go/pkg/client"
)

// BrowseSecrets lists child prefixes and secrets at a path.
func (c *Client) BrowseSecrets(ctx context.Context, region, path string) (*BrowseSecretsResponse, error) {
	normalized, err := NormalizePath(path)
	if err != nil {
		return nil, err
	}
	var result BrowseSecretsResponse
	req := c.R().SetQueryParam("path", normalized).SetResult(&result)
	resp, err := c.Do(ctx, req, client.GET, secretsCollectionURL(region))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &result, err
	}
	return &result, nil
}

// ListSecrets recursively lists secrets under a path prefix.
func (c *Client) ListSecrets(ctx context.Context, region, pathPrefix string) ([]Secret, error) {
	normalized, err := NormalizePath(pathPrefix)
	if err != nil {
		return nil, err
	}
	secrets := []Secret{}
	req := c.R().SetQueryParam("pathPrefix", normalized).SetResult(&secrets)
	resp, err := c.Do(ctx, req, client.GET, secretsCollectionURL(region))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return secrets, err
	}
	return secrets, nil
}

// CreateSecret creates a new secret.
func (c *Client) CreateSecret(ctx context.Context, region string, create CreateSecretRequest) (*Secret, error) {
	if _, err := NormalizePath(create.Path); err != nil {
		return nil, err
	}
	var secret Secret
	req := c.R().SetBody(create).SetResult(&secret)
	resp, err := c.Do(ctx, req, client.POST, secretsCollectionURL(region))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &secret, err
	}
	return &secret, nil
}

// GetSecret retrieves secret metadata. Set includeVersions to include version history.
func (c *Client) GetSecret(ctx context.Context, region, path string, includeVersions bool) (*Secret, error) {
	url, err := SecretResourceURL(region, path, "")
	if err != nil {
		return nil, err
	}
	var secret Secret
	req := c.R().SetResult(&secret)
	if includeVersions {
		req = req.SetQueryParam("includeVersions", "true")
	}
	resp, err := c.Do(ctx, req, client.GET, url)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &secret, err
	}
	return &secret, nil
}

// DeleteSecret deletes a secret and all of its versions.
func (c *Client) DeleteSecret(ctx context.Context, region, path string) error {
	url, err := SecretResourceURL(region, path, "")
	if err != nil {
		return err
	}
	resp, err := c.Do(ctx, c.R(), client.DELETE, url)
	if err != nil {
		return err
	}
	return c.Check(resp)
}
