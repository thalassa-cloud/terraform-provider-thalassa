package objectstorage

import (
	"context"
	"fmt"

	"github.com/thalassa-cloud/client-go/pkg/client"
)

func bucketLifecyclePath(bucketName string) string {
	return fmt.Sprintf("%s/%s/lifecycle", BucketEndpoint, bucketName)
}

// GetBucketLifecycle returns the lifecycle configuration for a bucket.
func (c *Client) GetBucketLifecycle(ctx context.Context, bucketName string) (*BucketLifecycle, error) {
	if bucketName == "" {
		return nil, fmt.Errorf("bucket name is required")
	}

	var lifecycle BucketLifecycle
	req := c.R().SetResult(&lifecycle)
	resp, err := c.Do(ctx, req, client.GET, bucketLifecyclePath(bucketName))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &lifecycle, err
	}
	return &lifecycle, nil
}

// SetBucketLifecycle replaces all lifecycle rules on a bucket.
func (c *Client) SetBucketLifecycle(ctx context.Context, bucketName string, set SetBucketLifecycleRequest) (*BucketLifecycle, error) {
	if bucketName == "" {
		return nil, fmt.Errorf("bucket name is required")
	}

	var lifecycle BucketLifecycle
	req := c.R().SetBody(set).SetResult(&lifecycle)
	resp, err := c.Do(ctx, req, client.PUT, bucketLifecyclePath(bucketName))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &lifecycle, err
	}
	return &lifecycle, nil
}

// DeleteBucketLifecycle removes all lifecycle rules from a bucket.
func (c *Client) DeleteBucketLifecycle(ctx context.Context, bucketName string) error {
	if bucketName == "" {
		return fmt.Errorf("bucket name is required")
	}

	resp, err := c.Do(ctx, c.R(), client.DELETE, bucketLifecyclePath(bucketName))
	if err != nil {
		return err
	}
	return c.Check(resp)
}
