package iaas

import (
	"time"

	"github.com/thalassa-cloud/client-go/pkg/base"
)

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
	// EndpointAddress is the address of the endpoint.
	EndpointAddress string `json:"endpointAddress,omitempty"`
	// EndpointHostname is the hostname of the endpoint.
	EndpointHostname string `json:"endpointHostname,omitempty"`
	// EndpointAddressType is the type of the endpoint address.
	EndpointAddressType EndpointAddressType `json:"endpointAddressType,omitempty"`
	EndpointType        EndpointType        `json:"endpointType,omitempty"`
	// MacAddress is the MAC address of the endpoint. Not set for all endpoints.
	MacAddress string `json:"macAddress,omitempty"`
	// Interface is the interface of the endpoint. Not set for all endpoints.
	Interface string `json:"interface,omitempty"`
	// ResourceType is the type of the resource that the endpoint is attached to.
	ResourceType string `json:"resourceType"`
	// ResourceIdentity is the identity of the resource that the endpoint is attached to
	ResourceIdentity       string                    `json:"resourceIdentity"`
	DatabaseCluster        *EndpointResourceMetadata `json:"databaseCluster,omitempty"`
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
	CreatedAt     time.Time   `json:"createdAt,omitempty"`
	UpdatedAt     time.Time   `json:"updatedAt,omitempty"`
	ObjectVersion int         `json:"objectVersion,omitempty"`
	Labels        Labels      `json:"labels,omitempty"`
	Annotations   Annotations `json:"annotations,omitempty"`
}
