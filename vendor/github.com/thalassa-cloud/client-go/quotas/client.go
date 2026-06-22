package quotas

import (
	"context"
	"fmt"

	"github.com/thalassa-cloud/client-go/pkg/client"
)

const (
	QuotasEndpoint = "/v1/quotas"
)

// Client represents the quotas client
type Client struct {
	client.Client
}

// New creates a new quotas client
func New(c client.Client, opts ...client.Option) (*Client, error) {
	c.WithOptions(opts...)
	return &Client{c}, nil
}

// ListOrganisationQuotas lists all organisation quotas
func (c *Client) ListOrganisationQuotas(ctx context.Context) ([]OrganisationQuota, error) {
	quotas := []OrganisationQuota{}
	req := c.R().SetResult(&quotas)
	resp, err := c.Do(ctx, req, client.GET, QuotasEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return quotas, err
	}
	return quotas, nil
}

// GetOrganisationQuota retrieves a specific organisation quota by name
func (c *Client) GetOrganisationQuota(ctx context.Context, quotaName string) (*OrganisationQuota, error) {
	var quota *OrganisationQuota
	req := c.R().SetResult(&quota)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", QuotasEndpoint, quotaName))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return quota, err
	}
	return quota, nil
}

// RequestQuotaIncrease requests an increase for a specific quota
func (c *Client) RequestQuotaIncrease(ctx context.Context, request RequestQuotaIncreaseRequest) error {
	req := c.R().SetBody(request)
	resp, err := c.Do(ctx, req, client.POST, fmt.Sprintf("%s/increase-requests", QuotasEndpoint))
	if err != nil {
		return err
	}
	if err := c.Check(resp); err != nil {
		return err
	}
	return nil
}
