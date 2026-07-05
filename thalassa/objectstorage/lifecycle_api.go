package objectstorage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/thalassa-cloud/client-go/objectstorage"
	"github.com/thalassa-cloud/client-go/thalassa"
)

func fetchBucketLifecycle(ctx context.Context, client thalassa.Client, bucketName string) (*objectstorage.BucketLifecycle, error) {
	path := fmt.Sprintf("/v1/object-storage/buckets/%s/lifecycle", bucketName)
	resp, err := client.GetClient().RawRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	if err := client.GetClient().Check(resp); err != nil {
		return nil, err
	}

	return decodeBucketLifecycleResponse(resp.Body())
}

func decodeBucketLifecycleResponse(body []byte) (*objectstorage.BucketLifecycle, error) {
	body = bytes.TrimSpace(body)
	if len(body) == 0 {
		return &objectstorage.BucketLifecycle{}, nil
	}

	var payload struct {
		Lifecycle *objectstorage.BucketLifecycle      `json:"lifecycle"`
		Rules     []objectstorage.BucketLifecycleRule `json:"rules"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("decoding bucket lifecycle: %w", err)
	}
	if len(payload.Rules) > 0 {
		return &objectstorage.BucketLifecycle{Rules: payload.Rules}, nil
	}
	if payload.Lifecycle != nil {
		return payload.Lifecycle, nil
	}

	var direct objectstorage.BucketLifecycle
	if err := json.Unmarshal(body, &direct); err != nil {
		return nil, fmt.Errorf("decoding bucket lifecycle: %w", err)
	}

	return &direct, nil
}
