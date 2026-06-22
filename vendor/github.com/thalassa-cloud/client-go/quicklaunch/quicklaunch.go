package quicklaunch

import (
	"context"
	"fmt"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

const QuickLaunchEndpoint = "/v1/quick-launch"

// ListQuickLaunchesRequest holds query filters for ListQuickLaunches.
type ListQuickLaunchesRequest struct {
	Filters []filters.Filter
}

// ListQuickLaunches lists quick launch jobs for the organisation (auth context).
func (c *Client) ListQuickLaunches(ctx context.Context, listRequest *ListQuickLaunchesRequest) ([]QuickLaunch, error) {
	var out []QuickLaunch
	req := c.R().SetResult(&out)
	if listRequest != nil {
		for _, filter := range listRequest.Filters {
			for k, v := range filter.ToParams() {
				req = req.SetQueryParam(k, v)
			}
		}
	}
	resp, err := c.Do(ctx, req, client.GET, QuickLaunchEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return out, err
	}
	return out, nil
}

// CreateQuickLaunch starts a quick launch (201 Created; provisioning continues asynchronously).
func (c *Client) CreateQuickLaunch(ctx context.Context, body QuickLaunchRequest) (*QuickLaunch, error) {
	if body.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if body.CloudRegionIdentity == "" {
		return nil, fmt.Errorf("cloudRegionIdentity is required")
	}
	var out QuickLaunch
	req := c.R().SetBody(body).SetResult(&out)
	resp, err := c.Do(ctx, req, client.POST, QuickLaunchEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &out, err
	}
	return &out, nil
}

// GetQuickLaunch returns a quick launch job by identity.
func (c *Client) GetQuickLaunch(ctx context.Context, identity string) (*QuickLaunch, error) {
	if identity == "" {
		return nil, fmt.Errorf("identity is required")
	}
	var out QuickLaunch
	req := c.R().SetResult(&out)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", QuickLaunchEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return &out, err
	}
	return &out, nil
}

// GetQuickLaunchLogs returns the response body from the logs endpoint (typically JSON).
func (c *Client) GetQuickLaunchLogs(ctx context.Context, identity string) (QuickLaunchLogs, error) {
	if identity == "" {
		return nil, fmt.Errorf("identity is required")
	}
	resp, err := c.Do(ctx, c.R(), client.GET, fmt.Sprintf("%s/%s/logs", QuickLaunchEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return nil, err
	}
	return QuickLaunchLogs(resp.Body()), nil
}

// DeleteQuickLaunch deletes a quick launch. When cascade is non-empty, it is sent as the cascade query parameter (Delete or Orphan).
func (c *Client) DeleteQuickLaunch(ctx context.Context, identity string, cascade QuickLaunchCascade) error {
	if identity == "" {
		return fmt.Errorf("identity is required")
	}
	req := c.R()
	if cascade != "" {
		req.SetQueryParam("cascade", string(cascade))
	}
	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s", QuickLaunchEndpoint, identity))
	if err != nil {
		return err
	}
	return c.Check(resp)
}
