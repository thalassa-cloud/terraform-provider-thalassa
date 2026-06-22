package secrets

import (
	"context"
	"fmt"
	"strconv"

	"github.com/thalassa-cloud/client-go/pkg/client"
)

// PutSecretValue stores a new secret version.
func (c *Client) PutSecretValue(ctx context.Context, region, path string, put PutSecretValueRequest) (*PutSecretValueResponse, error) {
	normalized, err := NormalizePath(path)
	if err != nil {
		return nil, err
	}
	url, err := SecretResourceURL(region, normalized, "/versions")
	if err != nil {
		return nil, err
	}
	put.Path = normalized
	var result PutSecretValueResponse
	req := c.R().SetBody(put).SetResult(&result)
	resp, err := c.Do(ctx, req, client.POST, url)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &result, err
	}
	return &result, nil
}

// PutSecretString stores a new string secret version, encoding plaintext as base64.
func (c *Client) PutSecretString(ctx context.Context, region, path string, plaintext []byte) (*PutSecretValueResponse, error) {
	return c.PutSecretValue(ctx, region, path, PutSecretValueRequest{
		SecretString: EncodeBytes(plaintext),
	})
}

// GetSecretValue retrieves a secret value.
// Callers must not log or persist decrypted secret material from the response.
func (c *Client) GetSecretValue(ctx context.Context, region, path string, version *int) (*GetSecretValueResponse, error) {
	normalized, err := NormalizePath(path)
	if err != nil {
		return nil, err
	}
	url, err := SecretResourceURL(region, normalized, "/value")
	if err != nil {
		return nil, err
	}
	body := GetSecretValueRequest{Path: normalized}
	if version != nil {
		body.Version = version
	}
	var result GetSecretValueResponse
	req := c.R().SetBody(body).SetResult(&result)
	resp, err := c.Do(ctx, req, client.POST, url)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &result, err
	}
	return &result, nil
}

// GetSecretString retrieves and decodes a string secret value.
func (c *Client) GetSecretString(ctx context.Context, region, path string, version *int) ([]byte, int, error) {
	result, err := c.GetSecretValue(ctx, region, path, version)
	if err != nil {
		return nil, 0, err
	}
	if result.SecretString == "" {
		return nil, result.Version, fmt.Errorf("secret at %s has no secretString value", path)
	}
	plaintext, err := DecodeBytes("secretString", result.SecretString)
	if err != nil {
		return nil, result.Version, err
	}
	return plaintext, result.Version, nil
}

// DestroySecretVersion permanently destroys a specific secret version.
func (c *Client) DestroySecretVersion(ctx context.Context, region, path string, version int) error {
	url, err := SecretResourceURL(region, path, "/versions")
	if err != nil {
		return err
	}
	req := c.R().SetQueryParam("version", strconv.Itoa(version))
	resp, err := c.Do(ctx, req, client.DELETE, url)
	if err != nil {
		return err
	}
	return c.Check(resp)
}
