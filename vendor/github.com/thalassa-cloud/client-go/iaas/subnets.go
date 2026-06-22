package iaas

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

const (
	SubnetEndpoint = "/v1/subnets"
)

type ListSubnetsRequest struct {
	Filters []filters.Filter
}

// ListSubnets lists all Subnets for a given organisation.
func (c *Client) ListSubnets(ctx context.Context, listRequest *ListSubnetsRequest) ([]Subnet, error) {
	subnets := []Subnet{}
	req := c.R().SetResult(&subnets)

	if listRequest != nil {
		for _, filter := range listRequest.Filters {
			for k, v := range filter.ToParams() {
				req = req.SetQueryParam(k, v)
			}
		}
	}

	resp, err := c.Do(ctx, req, client.GET, SubnetEndpoint)
	if err != nil {
		return nil, err
	}

	if err := c.Check(resp); err != nil {
		return subnets, err
	}
	return subnets, nil
}

// GetSubnet retrieves a specific Subnet by its identity.
// It returns an error if the subnet is not found.
// Example: subnet, err := c.GetSubnet(ctx, "subnet-identity1234")
//
//	if err != nil {
//		log.Fatalf("Failed to get subnet: %v", err)
//	}
func (c *Client) GetSubnet(ctx context.Context, identity string) (*Subnet, error) {
	var subnet *Subnet
	req := c.R().SetResult(&subnet)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", SubnetEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return subnet, err
	}
	return subnet, nil
}

// WaitUntilSubnetDeleted waits until the subnet is deleted.
// It returns an error if the subnet fails to delete.
// You are responsible for providing a context that can be cancelled, and for handling the error case.
// Example: ctxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// err := c.WaitUntilSubnetDeleted(ctxt, "subnet-123")
//
//	if err != nil {
//		log.Fatalf("Failed to wait for subnet to be deleted: %v", err)
//	}
//	defer cancel()
func (c *Client) WaitUntilSubnetDeleted(ctx context.Context, identity string) error {
	subnet, err := c.GetSubnet(ctx, identity)
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			return nil
		}
		return err
	}
	if subnet.Status == SubnetStatusDeleted {
		return nil
	}
	if subnet.Status != SubnetStatusDeleting {
		return fmt.Errorf("subnet %s is not being deleted (status: %s)", identity, subnet.Status)
	}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(DefaultPollIntervalForWaiting):
			subnet, err := c.GetSubnet(ctx, identity)
			if err != nil {
				if errors.Is(err, client.ErrNotFound) {
					return nil
				}
				return err
			}
			switch subnet.Status {
			case SubnetStatusFailed:
				return fmt.Errorf("subnet %s failed to delete (status: %s)", identity, subnet.Status)
			case SubnetStatusDeleted:
				return nil
			}
		}
	}
}

// WaitUntilSubnetReady waits until the subnet is ready.
// It returns the subnet when it is ready or an error if the subnet fails to become ready. This could happen if the subnet is being deleted, or entered a failed state.
// You are responsible for providing a context that can be cancelled, and for handling the error case.
// Example: ctxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// subnet, err := c.WaitUntilSubnetReady(ctxt, "subnet-123")
//
//	if err != nil {
//		log.Fatalf("Failed to wait for subnet to become ready: %v", err)
//	}
//	defer cancel()
func (c *Client) WaitUntilSubnetReady(ctx context.Context, identity string) (*Subnet, error) {
	subnet, err := c.GetSubnet(ctx, identity)
	if err != nil {
		return nil, err
	}
	if subnet.Status == SubnetStatusReady {
		return subnet, nil
	}
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(DefaultPollIntervalForWaiting):
			subnet, err := c.GetSubnet(ctx, identity)
			if err != nil {
				return nil, err
			}
			switch subnet.Status {
			case SubnetStatusReady:
				return subnet, nil
			case SubnetStatusFailed:
				return nil, fmt.Errorf("subnet %s failed to become ready", identity)
			case SubnetStatusDeleting:
				return nil, fmt.Errorf("subnet %s is being deleted", identity)
			case SubnetStatusDeleted:
				return nil, fmt.Errorf("subnet %s is deleted", identity)
			}
		}
	}
}

// CreateSubnet creates a new Subnet.
func (c *Client) CreateSubnet(ctx context.Context, create CreateSubnet) (*Subnet, error) {
	var subnet *Subnet
	req := c.R().
		SetBody(create).SetResult(&subnet)

	resp, err := c.Do(ctx, req, client.POST, SubnetEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return subnet, err
	}
	return subnet, nil
}

// UpdateSubnet updates an existing Subnet.
func (c *Client) UpdateSubnet(ctx context.Context, identity string, update UpdateSubnet) (*Subnet, error) {
	var subnet *Subnet
	req := c.R().
		SetBody(update).SetResult(&subnet)

	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s", SubnetEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return subnet, err
	}
	return subnet, nil
}

// DeleteSubnet deletes a specific Subnet by its identity.
func (c *Client) DeleteSubnet(ctx context.Context, identity string) error {
	req := c.R()

	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s", SubnetEndpoint, identity))
	if err != nil {
		return err
	}
	if err := c.Check(resp); err != nil {
		return err
	}
	return nil
}

type SubnetStatus string

const (
	SubnetStatusCreating SubnetStatus = "creating"
	SubnetStatusUpdating SubnetStatus = "updating"
	SubnetStatusReady    SubnetStatus = "ready"
	SubnetStatusActive   SubnetStatus = "active"
	SubnetStatusDeleting SubnetStatus = "deleting"
	SubnetStatusDeleted  SubnetStatus = "deleted"
	SubnetStatusFailed   SubnetStatus = "failed"
)

type SubnetType string

const (
	SubnetTypeIPv4 SubnetType = "IPv4"
	SubnetTypeIPv6 SubnetType = "IPv6"
	SubnetTypeDual SubnetType = "Dual"
)
