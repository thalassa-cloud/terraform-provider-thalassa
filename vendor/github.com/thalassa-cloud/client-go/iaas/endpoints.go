package iaas

import (
	"context"
	"fmt"
	"time"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/base"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

const (
	EndpointEndpoint = "/v1/endpoints"
)

// ListEndpointsRequest is the request for listing endpoints.
type ListEndpointsRequest struct {
	Filters []filters.Filter
}

// ListEndpoints lists all endpoints for the current organisation.
func (c *Client) ListEndpoints(ctx context.Context, request *ListEndpointsRequest) ([]Endpoint, error) {
	endpoints := []Endpoint{}
	req := c.R().SetResult(&endpoints)
	if request != nil {
		for _, filter := range request.Filters {
			for k, v := range filter.ToParams() {
				req = req.SetQueryParam(k, v)
			}
		}
	}

	resp, err := c.Do(ctx, req, client.GET, EndpointEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return endpoints, err
	}
	return endpoints, nil
}

// GetEndpoint retrieves an endpoint by identity.
func (c *Client) GetEndpoint(ctx context.Context, identity string) (*Endpoint, error) {
	if identity == "" {
		return nil, fmt.Errorf("endpoint identity is required")
	}

	var endpoint *Endpoint
	req := c.R().SetResult(&endpoint)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", EndpointEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return endpoint, err
	}
	return endpoint, nil
}

type Endpoint struct {
	// Identity is a unique identifier for the endpoint
	Identity string `json:"identity"`
	// Name is a human-readable name of the endpoint
	Name string `json:"name"`
	// Labels is a map of key-value pairs used for filtering and grouping endpoints
	Labels Labels `json:"labels,omitempty"`
	// CreatedAt is the timestamp when the endpoint was created
	CreatedAt time.Time `json:"createdAt"`
	// Vpc is the VPC that the endpoint is attached to.
	Vpc *Vpc `json:"vpc,omitempty"`
	// Subnet is the subnet that the endpoint is attached to.
	Subnet *Subnet `json:"subnet,omitempty"`
	// Organisation is the organisation that the endpoint belongs to.
	Organisation *base.Organisation `json:"organisation,omitempty"`
	// ProjectId is the project identity that the endpoint belongs to.
	ProjectId string `json:"projectId,omitempty"`
	// EndpointAddress is the address of the endpoint.
	EndpointAddress string `json:"endpointAddress,omitempty"`
	// EndpointHostname is the hostname of the endpoint.
	EndpointHostname string `json:"endpointHostname,omitempty"`
	// EndpointAddressType is the type of the endpoint address.
	EndpointAddressType EndpointAddressType `json:"endpointAddressType,omitempty"`
	// EndpointPort is the port of the endpoint.
	EndpointPort int `json:"endpointPort,omitempty"`
	// EndpointProtocol is the protocol of the endpoint.
	EndpointProtocol *EndpointProtocol `json:"endpointProtocol,omitempty"`
	EndpointType     EndpointType      `json:"endpointType,omitempty"`
	// MacAddress is the MAC address of the endpoint. Not set for all endpoints.
	MacAddress string `json:"macAddress,omitempty"`
	// Interface is the interface of the endpoint. Not set for all endpoints.
	Interface string `json:"interface,omitempty"`
	// ResourceType is the type of the resource that the endpoint is attached to.
	ResourceType string `json:"resourceType"`
	// ResourceIdentity is the identity of the resource that the endpoint is attached to
	ResourceIdentity       string                    `json:"resourceIdentity"`
	DatabaseCluster        *EndpointResourceMetadata `json:"databaseCluster,omitempty"`
	KubernetesCluster      *EndpointResourceMetadata `json:"kubernetesCluster,omitempty"`
	VirtualMachineInstance *Machine                  `json:"virtualMachineInstance,omitempty"`
	NatGateway             *VpcNatGateway            `json:"natGateway,omitempty"`
}
type EndpointType string

const (
	EndpointTypeVpc      EndpointType = "vpc"
	EndpointTypePlatform EndpointType = "platform"
)

type EndpointProtocol string

const (
	EndpointProtocolTCP   EndpointProtocol = "tcp"
	EndpointProtocolUDP   EndpointProtocol = "udp"
	EndpointProtocolHTTP  EndpointProtocol = "http"
	EndpointProtocolHTTPS EndpointProtocol = "https"
	EndpointProtocolICMP  EndpointProtocol = "icmp"
)

type EndpointAddressType string

const (
	EndpointAddressTypeIPv4 EndpointAddressType = "IPv4"
	EndpointAddressTypeIPv6 EndpointAddressType = "IPv6"
)

type EndpointResourceMetadata struct {
	Identity      string      `json:"identity"`
	Name          string      `json:"name,omitempty"`
	Slug          string      `json:"slug,omitempty"`
	Description   string      `json:"description,omitempty"`
	CreatedAt     time.Time   `json:"createdAt"`
	UpdatedAt     time.Time   `json:"updatedAt"`
	ObjectVersion int         `json:"objectVersion,omitempty"`
	Labels        Labels      `json:"labels,omitempty"`
	Annotations   Annotations `json:"annotations,omitempty"`
}
