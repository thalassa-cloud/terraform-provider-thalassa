package prometheus

import (
	"context"
	"fmt"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

// ListPrometheusTenants lists all Prometheus tenants for the organisation.
func (c *Client) ListPrometheusTenants(ctx context.Context, listRequest *ListPrometheusTenantsRequest) ([]PrometheusTenant, error) {
	tenants := []PrometheusTenant{}
	req := c.R().SetResult(&tenants)

	if listRequest != nil {
		for _, filter := range listRequest.Filters {
			for k, v := range filter.ToParams() {
				req = req.SetQueryParam(k, v)
			}
		}
	}

	resp, err := c.Do(ctx, req, client.GET, PrometheusEndpoint+"/tenants")
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return tenants, err
	}
	return tenants, nil
}

// GetPrometheusTenant retrieves a specific Prometheus tenant by its identity.
func (c *Client) GetPrometheusTenant(ctx context.Context, tenantIdentity string) (*PrometheusTenant, error) {
	if tenantIdentity == "" {
		return nil, fmt.Errorf("tenant identity is required")
	}

	var tenant *PrometheusTenant
	req := c.R().SetResult(&tenant)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/tenants/%s", PrometheusEndpoint, tenantIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return tenant, err
	}
	return tenant, nil
}

// CreatePrometheusTenant creates a new Prometheus tenant.
func (c *Client) CreatePrometheusTenant(ctx context.Context, create CreatePrometheusTenantRequest) (*PrometheusTenant, error) {
	if create.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	if create.Retention != "" {
		// verify the retention period is valid
		if err := validatePrometheusDuration(create.Retention); err != nil {
			return nil, fmt.Errorf("invalid retention period: %w", err)
		}
	}

	var tenant *PrometheusTenant
	req := c.R().SetBody(create).SetResult(&tenant)
	resp, err := c.Do(ctx, req, client.POST, PrometheusEndpoint+"/tenants")
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return tenant, err
	}
	return tenant, nil
}

// UpdatePrometheusTenant updates an existing Prometheus tenant.
func (c *Client) UpdatePrometheusTenant(ctx context.Context, tenantIdentity string, update UpdatePrometheusTenantRequest) (*PrometheusTenant, error) {
	if tenantIdentity == "" {
		return nil, fmt.Errorf("tenant identity is required")
	}

	if update.Retention != "" {
		// verify the retention period is valid
		if err := validatePrometheusDuration(update.Retention); err != nil {
			return nil, fmt.Errorf("invalid retention period: %w", err)
		}
	}

	var tenant *PrometheusTenant
	req := c.R().SetBody(update).SetResult(&tenant)
	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/tenants/%s", PrometheusEndpoint, tenantIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return tenant, err
	}
	return tenant, nil
}

// DeletePrometheusTenant deletes a specific Prometheus tenant by its identity.
func (c *Client) DeletePrometheusTenant(ctx context.Context, tenantIdentity string) error {
	if tenantIdentity == "" {
		return fmt.Errorf("tenant identity is required")
	}

	req := c.R()
	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/tenants/%s", PrometheusEndpoint, tenantIdentity))
	if err != nil {
		return err
	}
	return c.Check(resp)
}

// ListPrometheusTenantsRequest is the request for listing Prometheus tenants.
type ListPrometheusTenantsRequest struct {
	Filters []filters.Filter
}
