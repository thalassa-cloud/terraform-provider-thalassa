package iaas

import (
	"context"
	"fmt"
	"time"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/base"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

// VpcPeeringConnectionStatus represents the status of a VPC peering connection
type VpcPeeringConnectionStatus string

const (
	// VpcPeeringConnectionStatusPending is the status when the peering request is pending acceptance
	VpcPeeringConnectionStatusPending VpcPeeringConnectionStatus = "pending"
	// VpcPeeringConnectionStatusAccepted is the status when the peering request has been accepted
	VpcPeeringConnectionStatusAccepted VpcPeeringConnectionStatus = "accepted"
	// VpcPeeringConnectionStatusRejected is the status when the peering request has been rejected
	VpcPeeringConnectionStatusRejected VpcPeeringConnectionStatus = "rejected"
	// VpcPeeringConnectionStatusActive is the status when the peering connection is active
	VpcPeeringConnectionStatusActive VpcPeeringConnectionStatus = "active"
	// VpcPeeringConnectionStatusDeleting is the status when the peering connection is being deleted
	VpcPeeringConnectionStatusDeleting VpcPeeringConnectionStatus = "deleting"
	// VpcPeeringConnectionStatusDeleted is the status when the peering connection has been deleted
	VpcPeeringConnectionStatusDeleted VpcPeeringConnectionStatus = "deleted"
	// VpcPeeringConnectionStatusFailed is the status when the peering connection failed
	VpcPeeringConnectionStatusFailed VpcPeeringConnectionStatus = "failed"
	// VpcPeeringConnectionStatusExpired is the status when the peering request has expired
	VpcPeeringConnectionStatusExpired VpcPeeringConnectionStatus = "expired"
)

// VpcPeeringOrganisation represents an organisation in VPC peering context
type VpcPeeringOrganisation struct {
	// Identity is the identity of the organisation
	Identity string `json:"identity,omitempty"`
	// Name is the name of the organisation
	Name string `json:"name,omitempty"`
}

// VpcPeeringVpc represents a VPC in VPC peering context
type VpcPeeringVpc struct {
	// Identity is the identity of the VPC
	Identity string `json:"identity,omitempty"`
	// Name is the name of the VPC
	Name string `json:"name,omitempty"`
}

// VpcPeeringConnection represents a VPC peering connection between two VPCs
type VpcPeeringConnection struct {
	// Identity is a unique identifier for the object
	Identity string `json:"identity"`
	// Name is a human-readable name of the VPC peering connection
	Name string `json:"name"`
	// Slug is a URL-friendly version of the name
	Slug string `json:"slug"`
	// Description provides additional context about the VPC peering connection
	Description string `json:"description"`
	// CreatedAt is the time the VPC peering connection was created
	CreatedAt time.Time `json:"createdAt"`
	// UpdatedAt is the time the VPC peering connection was last updated
	UpdatedAt time.Time `json:"updatedAt"`
	// ObjectVersion is used for optimistic locking
	ObjectVersion int `json:"objectVersion"`
	// Labels are key-value pairs for categorizing and organizing connections
	Labels Labels `json:"labels"`
	// Annotations are additional metadata that don't affect connection behavior
	Annotations Annotations `json:"annotations"`

	// RequesterVpc is the VPC that initiated the peering request
	RequesterVpc *VpcPeeringVpc `json:"requesterVpc,omitempty"`
	// AccepterVpc is the VPC that will accept or deny the peering request
	AccepterVpc *VpcPeeringVpc `json:"accepterVpc,omitempty"`
	// RequesterOrganisation is the organisation that owns the requester VPC
	RequesterOrganisation *VpcPeeringOrganisation `json:"requesterOrganisation,omitempty"`
	// AccepterOrganisation is the organisation that owns the accepter VPC
	AccepterOrganisation *VpcPeeringOrganisation `json:"accepterOrganisation,omitempty"`

	// Status represents the current status of the peering connection
	Status VpcPeeringConnectionStatus `json:"status"`

	// StatusMessage provides additional information about the current status
	StatusMessage *string `json:"statusMessage,omitempty"`

	// ExpiresAt is the time when the peering request expires if not accepted
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`

	// RequesterNextHopIP is the next hop IP address for the requester VPC
	RequesterNextHopIP *string `json:"requesterNextHopIP,omitempty"`
	// AccepterNextHopIP is the next hop IP address for the accepter VPC
	AccepterNextHopIP *string `json:"accepterNextHopIP,omitempty"`

	// PeeringCidr is the CIDR block that is used for the peering connection
	PeeringCidr string `json:"peeringCidr,omitempty"`

	// Organisation is the organisation that owns the VPC peering connection
	Organisation *base.Organisation `json:"organisation,omitempty"`
}

// CreateVpcPeeringConnectionRequest represents a request to create a VPC peering connection
type CreateVpcPeeringConnectionRequest struct {
	// Name is the name of the VPC peering connection. Must be at least 1 character and at most 63 characters. ASCII only.
	Name string `json:"name"`
	// Description is the description of the VPC peering connection. Must be at most 500 characters. ASCII only.
	Description string `json:"description"`
	// Labels is a map of key-value pairs used for filtering and grouping objects
	Labels Labels `json:"labels"`
	// Annotations is a map of key-value pairs used for storing additional information
	Annotations Annotations `json:"annotations"`
	// RequesterVpcIdentity is the identity of the VPC that will initiate the peering request
	RequesterVpcIdentity string `json:"requesterVpcIdentity"`
	// AccepterVpcIdentity is the identity of the VPC that will accept or deny the peering request
	AccepterVpcIdentity string `json:"accepterVpcIdentity"`
	// AccepterOrganisationIdentity is the identity of the organisation that owns the accepter VPC
	AccepterOrganisationIdentity string `json:"accepterOrganisationIdentity"`

	// AutoAccept is a flag to indicate if the peering connection should be automatically accepted
	// This is only allowed if the requester and accepter VPCs are in the same region and the requester VPC is owned by the same organisation as the accepter VPC
	AutoAccept bool `json:"autoAccept"`
}

// UpdateVpcPeeringConnectionRequest represents a request to update a VPC peering connection
type UpdateVpcPeeringConnectionRequest struct {
	// Name is the name of the VPC peering connection. Must be at least 1 character and at most 63 characters. ASCII only.
	Name string `json:"name"`
	// Description is the description of the VPC peering connection. Must be at most 500 characters. ASCII only.
	Description string `json:"description"`
	// Labels is a map of key-value pairs used for filtering and grouping objects
	Labels Labels `json:"labels"`
	// Annotations is a map of key-value pairs used for storing additional information
	Annotations Annotations `json:"annotations"`
}

// AcceptVpcPeeringConnectionRequest represents a request to accept a VPC peering connection
type AcceptVpcPeeringConnectionRequest struct {
}

// RejectVpcPeeringConnectionRequest represents a request to reject a VPC peering connection
type RejectVpcPeeringConnectionRequest struct {
	// Reason is the reason for rejecting the peering connection
	Reason string `json:"reason,omitempty"`
}

const (
	VpcPeeringEndpoint = "/v1/vpc-peering-connections"
)

type ListVpcPeeringConnectionsRequest struct {
	Filters []filters.Filter
}

// ListVpcPeeringConnections lists all VPC peering connections for the current organisation.
func (c *Client) ListVpcPeeringConnections(ctx context.Context, request *ListVpcPeeringConnectionsRequest) ([]VpcPeeringConnection, error) {
	peeringConnections := []VpcPeeringConnection{}
	req := c.R().SetResult(&peeringConnections)
	if request != nil {
		for _, filter := range request.Filters {
			for k, v := range filter.ToParams() {
				req = req.SetQueryParam(k, v)
			}
		}
	}

	resp, err := c.Do(ctx, req, client.GET, VpcPeeringEndpoint)
	if err != nil {
		return nil, err
	}

	if err := c.Check(resp); err != nil {
		return peeringConnections, err
	}
	return peeringConnections, nil
}

// CreateVpcPeeringConnection creates a new VPC peering connection.
func (c *Client) CreateVpcPeeringConnection(ctx context.Context, create CreateVpcPeeringConnectionRequest) (*VpcPeeringConnection, error) {
	var peeringConnection *VpcPeeringConnection
	req := c.R().
		SetBody(create).SetResult(&peeringConnection)

	resp, err := c.Do(ctx, req, client.POST, VpcPeeringEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return peeringConnection, err
	}
	return peeringConnection, nil
}

// GetVpcPeeringConnection retrieves a specific VPC peering connection by its identity.
func (c *Client) GetVpcPeeringConnection(ctx context.Context, identity string) (*VpcPeeringConnection, error) {
	var peeringConnection *VpcPeeringConnection
	req := c.R().SetResult(&peeringConnection)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", VpcPeeringEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return peeringConnection, err
	}
	return peeringConnection, nil
}

// UpdateVpcPeeringConnection updates an existing VPC peering connection.
func (c *Client) UpdateVpcPeeringConnection(ctx context.Context, identity string, update UpdateVpcPeeringConnectionRequest) (*VpcPeeringConnection, error) {
	var peeringConnection *VpcPeeringConnection
	req := c.R().
		SetBody(update).SetResult(&peeringConnection)

	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s", VpcPeeringEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return peeringConnection, err
	}
	return peeringConnection, nil
}

// DeleteVpcPeeringConnection deletes a specific VPC peering connection by its identity.
func (c *Client) DeleteVpcPeeringConnection(ctx context.Context, identity string) error {
	req := c.R()

	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s", VpcPeeringEndpoint, identity))
	if err != nil {
		return err
	}
	if err := c.Check(resp); err != nil {
		return err
	}
	return nil
}

// AcceptVpcPeeringConnection accepts a VPC peering connection.
func (c *Client) AcceptVpcPeeringConnection(ctx context.Context, identity string, accept AcceptVpcPeeringConnectionRequest) (*VpcPeeringConnection, error) {
	var peeringConnection *VpcPeeringConnection
	req := c.R().
		SetBody(accept).SetResult(&peeringConnection)

	resp, err := c.Do(ctx, req, client.POST, fmt.Sprintf("%s/%s/accept", VpcPeeringEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return peeringConnection, err
	}
	return peeringConnection, nil
}

// RejectVpcPeeringConnection rejects a VPC peering connection.
func (c *Client) RejectVpcPeeringConnection(ctx context.Context, identity string, reject RejectVpcPeeringConnectionRequest) (*VpcPeeringConnection, error) {
	var peeringConnection *VpcPeeringConnection
	req := c.R().
		SetBody(reject).SetResult(&peeringConnection)

	resp, err := c.Do(ctx, req, client.POST, fmt.Sprintf("%s/%s/reject", VpcPeeringEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return peeringConnection, err
	}
	return peeringConnection, nil
}
