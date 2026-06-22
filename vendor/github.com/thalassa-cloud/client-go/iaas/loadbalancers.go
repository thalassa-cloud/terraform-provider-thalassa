package iaas

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

const (
	LoadbalancerEndpoint = "/v1/loadbalancers"
)

// ListLoadbalancers lists all loadbalancers for a given organisation.
func (c *Client) ListLoadbalancers(ctx context.Context, listRequest *ListLoadbalancersRequest) ([]VpcLoadbalancer, error) {
	loadbalancers := []VpcLoadbalancer{}
	req := c.R().SetResult(&loadbalancers)

	if listRequest != nil {
		for _, filter := range listRequest.Filters {
			for k, v := range filter.ToParams() {
				req = req.SetQueryParam(k, v)
			}
		}
	}

	resp, err := c.Do(ctx, req, client.GET, LoadbalancerEndpoint)
	if err != nil {
		return nil, err
	}

	if err := c.Check(resp); err != nil {
		return loadbalancers, err
	}
	return loadbalancers, nil
}

// GetLoadbalancer retrieves a specific loadbalancer by its identity.
func (c *Client) GetLoadbalancer(ctx context.Context, loadbalancerIdentity string) (*VpcLoadbalancer, error) {
	if loadbalancerIdentity == "" {
		return nil, fmt.Errorf("identity is required")
	}

	var loadbalancer *VpcLoadbalancer
	req := c.R().SetResult(&loadbalancer)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", LoadbalancerEndpoint, loadbalancerIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return loadbalancer, err
	}
	return loadbalancer, nil
}

// CreateLoadbalancer creates a new loadbalancer.
func (c *Client) CreateLoadbalancer(ctx context.Context, create CreateLoadbalancer) (*VpcLoadbalancer, error) {
	if create.Subnet == "" {
		return nil, fmt.Errorf("subnet is required")
	}
	if create.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	var loadbalancer *VpcLoadbalancer
	req := c.R().
		SetBody(create).SetResult(&loadbalancer)

	resp, err := c.Do(ctx, req, client.POST, LoadbalancerEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return loadbalancer, err
	}
	return loadbalancer, nil
}

// UpdateLoadbalancer updates an existing loadbalancer.
func (c *Client) UpdateLoadbalancer(ctx context.Context, loadbalancerIdentity string, update UpdateLoadbalancer) (*VpcLoadbalancer, error) {
	if loadbalancerIdentity == "" {
		return nil, fmt.Errorf("identity of the loadbalancer to update is required")
	}

	var loadbalancer *VpcLoadbalancer
	req := c.R().
		SetBody(update).SetResult(&loadbalancer)

	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s", LoadbalancerEndpoint, loadbalancerIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return loadbalancer, err
	}
	return loadbalancer, nil
}

// DeleteLoadbalancer deletes a specific loadbalancer by its identity.
func (c *Client) DeleteLoadbalancer(ctx context.Context, loadbalancerIdentity string) error {
	if loadbalancerIdentity == "" {
		return fmt.Errorf("identity is required")
	}

	req := c.R()

	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s", LoadbalancerEndpoint, loadbalancerIdentity))
	if err != nil {
		return err
	}
	if err := c.Check(resp); err != nil {
		return err
	}
	return nil
}

// WaitUntilLoadbalancerIsReady waits until a loadbalancer is ready.
// The user is expected to provide a timeout context.
func (c *Client) WaitUntilLoadbalancerIsReady(ctx context.Context, loadbalancerIdentity string) error {
	return c.WaitUntilLoadbalancerIsStatus(ctx, loadbalancerIdentity, "ready")
}

// WaitUntilLoadbalancerIsStatus waits until a loadbalancer is in a specific status.
// The user is expected to provide a timeout context.
func (c *Client) WaitUntilLoadbalancerIsStatus(ctx context.Context, loadbalancerIdentity string, status string) error {
	loadbalancer, err := c.GetLoadbalancer(ctx, loadbalancerIdentity)
	if err != nil {
		return err
	}
	if strings.EqualFold(loadbalancer.Status, status) {
		return nil
	}
	// wait until the loadbalancer is in the desired status
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(DefaultPollIntervalForWaiting):
		}

		loadbalancer, err = c.GetLoadbalancer(ctx, loadbalancerIdentity)
		if err != nil {
			return err
		}
		if strings.EqualFold(loadbalancer.Status, status) {
			return nil
		}
	}
}

// WaitUntilLoadbalancerIsDeleted waits until a loadbalancer is deleted.
// The user is expected to provide a timeout context.
func (c *Client) WaitUntilLoadbalancerIsDeleted(ctx context.Context, loadbalancerIdentity string) error {
	loadbalancer, err := c.GetLoadbalancer(ctx, loadbalancerIdentity)
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			return nil
		}
		return err
	}
	if strings.EqualFold(loadbalancer.Status, "deleted") {
		return nil
	}
	if !strings.EqualFold(loadbalancer.Status, "deleting") {
		return fmt.Errorf("loadbalancer %s is not being deleted (status: %s)", loadbalancerIdentity, loadbalancer.Status)
	}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(DefaultPollIntervalForWaiting):
			loadbalancer, err := c.GetLoadbalancer(ctx, loadbalancerIdentity)
			if err != nil {
				if errors.Is(err, client.ErrNotFound) {
					return nil
				}
				return err
			}
			if strings.EqualFold(loadbalancer.Status, "deleted") {
				return nil
			}
		}
	}
}

type ListLoadbalancersRequest struct {
	Filters []filters.Filter
}

type CreateLoadbalancer struct {
	// Name is the name of the loadbalancer.
	Name string `json:"name"`
	// Description is the description of the loadbalancer.
	Description string `json:"description"`
	// Labels are arbitrary key-value pairs that can be used to store additional information about the loadbalancer, and are used for matching resources.
	Labels Labels `json:"labels,omitempty"`
	// Annotations are arbitrary key-value pairs that can be used to store additional information about the loadbalancer, and are used for matching resources.
	Annotations Annotations `json:"annotations,omitempty"`

	// Subnet is the subnet in which the loadbalancer will be deployed.
	Subnet string `json:"subnet"`

	// InternalLoadbalancer is a flag that indicates whether the loadbalancer should be an internal loadbalancer.
	// If true, the loadbalancer will be an internal loadbalancer and will not be accessible from the public internet.
	// It will not be assigned a public IP address, and instead will be assigned a (private) IP address from the subnet.
	InternalLoadbalancer bool `json:"internalLoadbalancer"`

	// DeleteProtection is a flag that indicates whether the loadbalancer should be protected from deletion.
	// Meaning delete protection will require to be disabled explicitly before the loadbalancer can be deleted.
	DeleteProtection bool `json:"deleteProtection"`

	// Listeners is a list of listeners that will be created on the loadbalancer during creation.
	Listeners []CreateListener `json:"listeners"`

	// SecurityGroupAttachments is a list of security group identities that will be attached to the loadbalancer.
	SecurityGroupAttachments []string `json:"securityGroupAttachments,omitempty"`
	// ReservedIpID, if set, attaches this reserved IP after create (must be available, same region as subnet).
	ReservedIpID *string `json:"reservedIpId,omitempty"`
}

type UpdateLoadbalancer struct {
	// Name is the name of the loadbalancer.
	Name string `json:"name"`
	// Description is the description of the loadbalancer.
	Description string `json:"description"`
	// Labels are arbitrary key-value pairs that can be used to store additional information about the loadbalancer, and are used for matching resources.
	Labels Labels `json:"labels,omitempty"`
	// Annotations are arbitrary key-value pairs that can be used to store additional information about the loadbalancer, and are used for matching resources.
	Annotations Annotations `json:"annotations,omitempty"`
	// DeleteProtection is a flag that indicates whether the loadbalancer should be protected from deletion.
	DeleteProtection bool `json:"deleteProtection"`
	// Subnet is the subnet in which the loadbalancer will be deployed. Optional. If provided, the loadbalancer will be moved to the new subnet.
	// Only supported for load balancers that are not internal.
	Subnet *string `json:"subnet,omitempty"`
	// SecurityGroupAttachments is a list of security group identities that will be attached to the loadbalancer.
	SecurityGroupAttachments []string `json:"securityGroupAttachments,omitempty"`
	// ReservedIpID: nil = leave attachment unchanged, empty string = detach, non-empty = attach/replace.
	ReservedIpID *string `json:"reservedIpId,omitempty"`
}
