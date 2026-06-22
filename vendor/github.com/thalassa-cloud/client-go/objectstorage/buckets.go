package objectstorage

import (
	"context"
	"fmt"

	"github.com/thalassa-cloud/client-go/pkg/client"
)

const (
	BucketEndpoint = "/v1/object-storage/buckets"
)

func (c *Client) ListBuckets(ctx context.Context) ([]ObjectStorageBucket, error) {
	buckets := []ObjectStorageBucket{}
	req := c.R().SetResult(&buckets)

	resp, err := c.Do(ctx, req, client.GET, BucketEndpoint)
	if err != nil {
		return nil, err
	}

	if err := c.Check(resp); err != nil {
		return nil, err
	}

	return buckets, nil
}

func (c *Client) GetBucket(ctx context.Context, bucketName string) (*ObjectStorageBucket, error) {
	bucket := ObjectStorageBucket{}
	req := c.R().SetResult(&bucket)

	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", BucketEndpoint, bucketName))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return nil, err
	}
	return &bucket, nil
}

func (c *Client) CreateBucket(ctx context.Context, create CreateBucketRequest) (*ObjectStorageBucket, error) {
	bucket := ObjectStorageBucket{}
	req := c.R().SetBody(create).SetResult(&bucket)

	resp, err := c.Do(ctx, req, client.POST, BucketEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return nil, err
	}
	return &bucket, nil
}

func (c *Client) UpdateBucket(ctx context.Context, bucketName string, update UpdateBucketRequest) (*ObjectStorageBucket, error) {
	bucket := ObjectStorageBucket{}
	req := c.R().SetBody(update).SetResult(&bucket)

	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s", BucketEndpoint, bucketName))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return nil, err
	}
	return &bucket, nil
}

func (c *Client) DeleteBucket(ctx context.Context, bucketName string) error {
	req := c.R()
	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s", BucketEndpoint, bucketName))
	if err != nil {
		return err
	}
	return c.Check(resp)
}
