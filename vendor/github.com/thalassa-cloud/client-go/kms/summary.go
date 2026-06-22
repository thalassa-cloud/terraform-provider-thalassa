package kms

import (
	"context"

	"github.com/thalassa-cloud/client-go/pkg/client"
)

// GetSummary returns KMS availability and regional configuration for the organisation.
func (c *Client) GetSummary(ctx context.Context) (*KmsSummary, error) {
	var summary KmsSummary
	req := c.R().SetResult(&summary)
	resp, err := c.Do(ctx, req, client.GET, KmsEndpoint+"/summary")
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &summary, err
	}
	return &summary, nil
}
