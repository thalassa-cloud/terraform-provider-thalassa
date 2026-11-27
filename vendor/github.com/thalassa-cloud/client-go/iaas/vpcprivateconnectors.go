package iaas

import (
	"context"
	"fmt"
	"time"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

const (
	VpcPrivateConnectorEndpoint = "/v1/vpc-private-connectors"
)

// VpcPrivateConnectorType is the type of private connector.
type VpcPrivateConnectorType string

const (
	VpcPrivateConnectorTypeCloudServices VpcPrivateConnectorType = "cloud_services"
)

// VpcPrivateConnectorStatus is the status of a VPC private connector.
type VpcPrivateConnectorStatus string

const (
	VpcPrivateConnectorStatusPending  VpcPrivateConnectorStatus = "pending"
	VpcPrivateConnectorStatusActive   VpcPrivateConnectorStatus = "active"
	VpcPrivateConnectorStatusFailed   VpcPrivateConnectorStatus = "failed"
	VpcPrivateConnectorStatusDeleting VpcPrivateConnectorStatus = "deleting"
)

// VpcPrivateConnector is a private connector for cloud services within a VPC.
type VpcPrivateConnector struct {
	Identity               string                    `json:"identity"`
	Annotations            Annotations               `json:"annotations,omitempty"`
	Labels                 Labels                    `json:"labels,omitempty"`
	CreatedAt              time.Time                 `json:"createdAt"`
	UpdatedAt              *time.Time                `json:"updatedAt,omitempty"`
	DeletedAt              *time.Time                `json:"deletedAt,omitempty"`
	ProjectId              string                    `json:"projectId,omitempty"`
	Vpc                    *Vpc                      `json:"vpc,omitempty"`
	Type                   VpcPrivateConnectorType   `json:"type"`
	AttachedSecurityGroups []string                  `json:"attachedSecurityGroups,omitempty"`
	EndpointIPv4           *string                   `json:"endpointIPv4,omitempty"`
	EndpointIPv6           *string                   `json:"endpointIPv6,omitempty"`
	Status                 VpcPrivateConnectorStatus `json:"status"`
	StatusMessage          *string                   `json:"statusMessage,omitempty"`
}

// CreateVpcPrivateConnectorRequest is the request body for creating a VPC private connector.
type CreateVpcPrivateConnectorRequest struct {
	VpcIdentity            string                  `json:"vpcIdentity"`
	Type                   VpcPrivateConnectorType `json:"type"`
	AttachedSecurityGroups []string                `json:"attachedSecurityGroups,omitempty"`
	Labels                 Labels                  `json:"labels,omitempty"`
	Annotations            Annotations             `json:"annotations,omitempty"`
}

// UpdateVpcPrivateConnectorRequest is the request body for updating a VPC private connector.
type UpdateVpcPrivateConnectorRequest struct {
	AttachedSecurityGroups []string    `json:"attachedSecurityGroups,omitempty"`
	Labels                 Labels      `json:"labels,omitempty"`
	Annotations            Annotations `json:"annotations,omitempty"`
}

// ListVpcPrivateConnectorsRequest is the request for listing VPC private connectors.
type ListVpcPrivateConnectorsRequest struct {
	Filters []filters.Filter
}

// ListVpcPrivateConnectors lists all VPC private connectors for the current organisation.
func (c *Client) ListVpcPrivateConnectors(ctx context.Context, request *ListVpcPrivateConnectorsRequest) ([]VpcPrivateConnector, error) {
	connectors := []VpcPrivateConnector{}
	req := c.R().SetResult(&connectors)
	if request != nil {
		for _, filter := range request.Filters {
			for k, v := range filter.ToParams() {
				req = req.SetQueryParam(k, v)
			}
		}
	}

	resp, err := c.Do(ctx, req, client.GET, VpcPrivateConnectorEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return connectors, err
	}
	return connectors, nil
}

// GetVpcPrivateConnector retrieves a VPC private connector by identity.
func (c *Client) GetVpcPrivateConnector(ctx context.Context, identity string) (*VpcPrivateConnector, error) {
	if identity == "" {
		return nil, fmt.Errorf("vpc private connector identity is required")
	}

	var connector *VpcPrivateConnector
	req := c.R().SetResult(&connector)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", VpcPrivateConnectorEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return connector, err
	}
	return connector, nil
}

// CreateVpcPrivateConnector creates a new VPC private connector.
func (c *Client) CreateVpcPrivateConnector(ctx context.Context, create CreateVpcPrivateConnectorRequest) (*VpcPrivateConnector, error) {
	if create.VpcIdentity == "" {
		return nil, fmt.Errorf("vpc identity is required")
	}
	if create.Type == "" {
		return nil, fmt.Errorf("type is required")
	}

	var connector *VpcPrivateConnector
	req := c.R().SetBody(create).SetResult(&connector)
	resp, err := c.Do(ctx, req, client.POST, VpcPrivateConnectorEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return connector, err
	}
	return connector, nil
}

// UpdateVpcPrivateConnector updates an existing VPC private connector.
func (c *Client) UpdateVpcPrivateConnector(ctx context.Context, identity string, update UpdateVpcPrivateConnectorRequest) (*VpcPrivateConnector, error) {
	if identity == "" {
		return nil, fmt.Errorf("vpc private connector identity is required")
	}

	var connector *VpcPrivateConnector
	req := c.R().SetBody(update).SetResult(&connector)
	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s", VpcPrivateConnectorEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return connector, err
	}
	return connector, nil
}

// DeleteVpcPrivateConnector deletes a VPC private connector by identity.
func (c *Client) DeleteVpcPrivateConnector(ctx context.Context, identity string) error {
	if identity == "" {
		return fmt.Errorf("vpc private connector identity is required")
	}

	req := c.R()
	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s", VpcPrivateConnectorEndpoint, identity))
	if err != nil {
		return err
	}
	return c.Check(resp)
}
