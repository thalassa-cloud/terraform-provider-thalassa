package tfs

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/thalassa-cloud/client-go/pkg/client"
)

const (
	TfsEndpoint = "/v1/tfs"
)

// ListTfsInstances lists all TFS instances for a given organisation.
func (c *Client) ListTfsInstances(ctx context.Context, request *ListTfsInstancesRequest) ([]TfsInstance, error) {
	tfsInstances := []TfsInstance{}
	req := c.R().SetResult(&tfsInstances)
	if request != nil {
		for _, filter := range request.Filters {
			for k, v := range filter.ToParams() {
				req = req.SetQueryParam(k, v)
			}
		}
	}

	resp, err := c.Do(ctx, req, client.GET, TfsEndpoint)
	if err != nil {
		return nil, err
	}

	if err := c.Check(resp); err != nil {
		return tfsInstances, err
	}
	return tfsInstances, nil
}

// GetTfsInstance retrieves a specific TFS instance by its identity.
func (c *Client) GetTfsInstance(ctx context.Context, identity string) (*TfsInstance, error) {
	var tfsInstance *TfsInstance
	req := c.R().SetResult(&tfsInstance)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", TfsEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return tfsInstance, err
	}
	return tfsInstance, nil
}

// CreateTfsInstance creates a new TFS instance.
func (c *Client) CreateTfsInstance(ctx context.Context, create CreateTfsInstanceRequest) (*TfsInstance, error) {
	var tfsInstance *TfsInstance
	req := c.R().
		SetBody(create).SetResult(&tfsInstance)

	resp, err := c.Do(ctx, req, client.POST, TfsEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return tfsInstance, err
	}
	return tfsInstance, nil
}

// UpdateTfsInstance updates an existing TFS instance.
func (c *Client) UpdateTfsInstance(ctx context.Context, identity string, update UpdateTfsInstanceRequest) (*TfsInstance, error) {
	var tfsInstance *TfsInstance
	req := c.R().
		SetBody(update).SetResult(&tfsInstance)

	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s", TfsEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return tfsInstance, err
	}
	return tfsInstance, nil
}

// DeleteTfsInstance deletes a specific TFS instance by its identity.
func (c *Client) DeleteTfsInstance(ctx context.Context, identity string) error {
	req := c.R()

	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s", TfsEndpoint, identity))
	if err != nil {
		return err
	}
	if err := c.Check(resp); err != nil {
		return err
	}
	return nil
}

// WaitUntilTfsInstanceIsAvailable waits until a TFS instance is available.
// The user is expected to provide a timeout context.
func (c *Client) WaitUntilTfsInstanceIsAvailable(ctx context.Context, tfsIdentity string) error {
	return c.WaitUntilTfsInstanceIsStatus(ctx, tfsIdentity, TfsStatusAvailable)
}

// WaitUntilTfsInstanceIsStatus waits until a TFS instance is in a specific status.
// The user is expected to provide a timeout context.
func (c *Client) WaitUntilTfsInstanceIsStatus(ctx context.Context, tfsIdentity string, status TfsStatus) error {
	tfsInstance, err := c.GetTfsInstance(ctx, tfsIdentity)
	if err != nil {
		return err
	}
	if strings.EqualFold(string(tfsInstance.Status), string(status)) {
		return nil
	}
	// wait until the TFS instance is in the desired status
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(DefaultPollIntervalForWaiting):
		}

		tfsInstance, err = c.GetTfsInstance(ctx, tfsIdentity)
		if err != nil {
			return err
		}
		if strings.EqualFold(string(tfsInstance.Status), string(status)) {
			return nil
		}
	}
}

// WaitUntilTfsInstanceIsDeleted waits until a TFS instance is deleted.
// The user is expected to provide a timeout context.
func (c *Client) WaitUntilTfsInstanceIsDeleted(ctx context.Context, tfsIdentity string) error {
	tfsInstance, err := c.GetTfsInstance(ctx, tfsIdentity)
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			return nil
		}
		return err
	}
	if strings.EqualFold(string(tfsInstance.Status), string(TfsStatusDeleted)) {
		return nil
	}
	if !strings.EqualFold(string(tfsInstance.Status), string(TfsStatusDeleting)) {
		return fmt.Errorf("TFS instance %s is not being deleted (status: %s)", tfsIdentity, tfsInstance.Status)
	}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(DefaultPollIntervalForWaiting):
			tfsInstance, err := c.GetTfsInstance(ctx, tfsIdentity)
			if err != nil {
				if errors.Is(err, client.ErrNotFound) {
					return nil
				}
				return err
			}
			if strings.EqualFold(string(tfsInstance.Status), string(TfsStatusDeleted)) {
				return nil
			}
		}
	}
}
