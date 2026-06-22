package kms

import (
	"context"

	"github.com/thalassa-cloud/client-go/pkg/client"
)

// ListKeys lists KMS keys in a region.
func (c *Client) ListKeys(ctx context.Context, region string, req *ListKeysRequest) ([]KmsKey, error) {
	keys := []KmsKey{}
	r := c.R().SetResult(&keys)
	if req != nil {
		for _, filter := range req.Filters {
			for k, v := range filter.ToParams() {
				r = r.SetQueryParam(k, v)
			}
		}
	}
	resp, err := c.Do(ctx, r, client.GET, regionPath(region, "keys"))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return keys, err
	}
	return keys, nil
}

// GetWrappingKey returns the regional wrapping public key for BYOK import.
func (c *Client) GetWrappingKey(ctx context.Context, region string) (*WrappingKeyResponse, error) {
	var wrappingKey WrappingKeyResponse
	req := c.R().SetResult(&wrappingKey)
	resp, err := c.Do(ctx, req, client.GET, regionPath(region, "wrapping-key"))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &wrappingKey, err
	}
	return &wrappingKey, nil
}

// CreateKey creates a new KMS key in a region.
func (c *Client) CreateKey(ctx context.Context, region string, create CreateKmsKeyRequest) (*KmsKey, error) {
	var key KmsKey
	req := c.R().SetBody(create).SetResult(&key)
	resp, err := c.Do(ctx, req, client.POST, regionPath(region, "keys"))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &key, err
	}
	return &key, nil
}

// GetKey retrieves a KMS key by identity.
func (c *Client) GetKey(ctx context.Context, region, identity string) (*KmsKey, error) {
	var key KmsKey
	req := c.R().SetResult(&key)
	resp, err := c.Do(ctx, req, client.GET, regionPath(region, "keys", identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &key, err
	}
	return &key, nil
}

// DeleteKey schedules a KMS key for deletion.
func (c *Client) DeleteKey(ctx context.Context, region, identity string) error {
	resp, err := c.Do(ctx, c.R(), client.DELETE, regionPath(region, "keys", identity))
	if err != nil {
		return err
	}
	return c.Check(resp)
}

// CancelDeletion cancels a pending KMS key deletion.
func (c *Client) CancelDeletion(ctx context.Context, region, identity string) error {
	resp, err := c.Do(ctx, c.R(), client.DELETE, regionPath(region, "keys", identity, "cancel-deletion"))
	if err != nil {
		return err
	}
	return c.Check(resp)
}

// UpdateRotation updates automatic key rotation settings.
func (c *Client) UpdateRotation(ctx context.Context, region, identity string, update UpdateRotationRequest) (*KmsKey, error) {
	var key KmsKey
	req := c.R().SetBody(update).SetResult(&key)
	resp, err := c.Do(ctx, req, client.PATCH, regionPath(region, "keys", identity, "rotation"))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &key, err
	}
	return &key, nil
}

// RotateKey rotates a KMS key on demand.
func (c *Client) RotateKey(ctx context.Context, region, identity string) (*KmsKey, error) {
	var key KmsKey
	req := c.R().SetResult(&key)
	resp, err := c.Do(ctx, req, client.POST, regionPath(region, "keys", identity, "rotate"))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &key, err
	}
	return &key, nil
}

// DisableKey disables a KMS key.
func (c *Client) DisableKey(ctx context.Context, region, identity string) (*KmsKey, error) {
	var key KmsKey
	req := c.R().SetResult(&key)
	resp, err := c.Do(ctx, req, client.POST, regionPath(region, "keys", identity, "disable"))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &key, err
	}
	return &key, nil
}

// EnableKey enables a KMS key.
func (c *Client) EnableKey(ctx context.Context, region, identity string) (*KmsKey, error) {
	var key KmsKey
	req := c.R().SetResult(&key)
	resp, err := c.Do(ctx, req, client.POST, regionPath(region, "keys", identity, "enable"))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &key, err
	}
	return &key, nil
}
