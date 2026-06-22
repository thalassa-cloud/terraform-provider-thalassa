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
	VpcEndpoint = "/v1/vpcs"
)

type ListVpcsRequest struct {
	Filters []filters.Filter
}

// ListVpcs lists all VPCs for a given organisation.
func (c *Client) ListVpcs(ctx context.Context, request *ListVpcsRequest) ([]Vpc, error) {
	vpcs := []Vpc{}
	req := c.R().SetResult(&vpcs)
	if request != nil {
		for _, filter := range request.Filters {
			for k, v := range filter.ToParams() {
				req = req.SetQueryParam(k, v)
			}
		}
	}

	resp, err := c.Do(ctx, req, client.GET, VpcEndpoint)
	if err != nil {
		return nil, err
	}

	if err := c.Check(resp); err != nil {
		return vpcs, err
	}
	return vpcs, nil
}

// GetVpc retrieves a specific VPC by its identity.
func (c *Client) GetVpc(ctx context.Context, identity string) (*Vpc, error) {
	var vpc *Vpc
	req := c.R().SetResult(&vpc)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", VpcEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return vpc, err
	}
	return vpc, nil
}

// CreateVpc creates a new VPC.
func (c *Client) CreateVpc(ctx context.Context, create CreateVpc) (*Vpc, error) {
	var vpc *Vpc
	req := c.R().
		SetBody(create).SetResult(&vpc)

	resp, err := c.Do(ctx, req, client.POST, VpcEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return vpc, err
	}
	return vpc, nil
}

// UpdateVpc updates an existing VPC.
func (c *Client) UpdateVpc(ctx context.Context, identity string, update UpdateVpc) (*Vpc, error) {
	var vpc *Vpc
	req := c.R().
		SetBody(update).SetResult(&vpc)

	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s", VpcEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return vpc, err
	}
	return vpc, nil
}

// DeleteVpc deletes a specific VPC by its identity.
func (c *Client) DeleteVpc(ctx context.Context, identity string) error {
	req := c.R()

	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s", VpcEndpoint, identity))
	if err != nil {
		return err
	}
	if err := c.Check(resp); err != nil {
		return err
	}
	return nil
}

// WaitUntilVpcIsReady waits until a VPC is ready.
// The user is expected to provide a timeout context.
func (c *Client) WaitUntilVpcIsReady(ctx context.Context, vpcIdentity string) error {
	return c.WaitUntilVpcIsStatus(ctx, vpcIdentity, "ready")
}

// WaitUntilVpcIsStatus waits until a VPC is in a specific status.
// The user is expected to provide a timeout context.
func (c *Client) WaitUntilVpcIsStatus(ctx context.Context, vpcIdentity string, status string) error {
	vpc, err := c.GetVpc(ctx, vpcIdentity)
	if err != nil {
		return err
	}
	if strings.EqualFold(vpc.Status, status) {
		return nil
	}
	// wait until the VPC is in the desired status
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(DefaultPollIntervalForWaiting):
		}

		vpc, err = c.GetVpc(ctx, vpcIdentity)
		if err != nil {
			return err
		}
		if strings.EqualFold(vpc.Status, status) {
			return nil
		}
	}
}

// WaitUntilVpcIsDeleted waits until a VPC is deleted.
// The user is expected to provide a timeout context.
func (c *Client) WaitUntilVpcIsDeleted(ctx context.Context, vpcIdentity string) error {
	vpc, err := c.GetVpc(ctx, vpcIdentity)
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			return nil
		}
		return err
	}
	if strings.EqualFold(vpc.Status, "deleted") {
		return nil
	}
	if !strings.EqualFold(vpc.Status, "deleting") {
		return fmt.Errorf("VPC %s is not being deleted (status: %s)", vpcIdentity, vpc.Status)
	}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(DefaultPollIntervalForWaiting):
			vpc, err := c.GetVpc(ctx, vpcIdentity)
			if err != nil {
				if errors.Is(err, client.ErrNotFound) {
					return nil
				}
				return err
			}
			if strings.EqualFold(vpc.Status, "deleted") {
				return nil
			}
		}
	}
}
