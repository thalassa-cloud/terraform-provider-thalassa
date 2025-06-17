package iaas

import (
	"time"

	"github.com/thalassa-cloud/client-go/pkg/base"
)

type Region struct {
	Identity      string      `json:"identity"`
	Name          string      `json:"name"`
	Slug          string      `json:"slug"`
	Description   string      `json:"description"`
	CreatedAt     time.Time   `json:"createdAt"`
	UpdatedAt     time.Time   `json:"updatedAt"`
	ObjectVersion int         `json:"objectVersion"`
	Labels        Labels      `json:"labels"`
	Annotations   Annotations `json:"annotations"`
	Zones         []Zone      `json:"zones"`
}

type Zone struct {
	Identity            string      `json:"identity"`
	Name                string      `json:"name"`
	Slug                string      `json:"slug"`
	Description         string      `json:"description"`
	CreatedAt           time.Time   `json:"createdAt"`
	UpdatedAt           time.Time   `json:"updatedAt"`
	ObjectVersion       int         `json:"objectVersion"`
	Labels              Labels      `json:"labels"`
	Annotations         Annotations `json:"annotations"`
	CloudRegionIdentity string      `json:"cloudRegionIdentity"`
	CloudRegion         *Region     `json:"CloudRegion"`
}

type VpcFirewallRule struct {
}

type Vpc struct {
	Identity      string    `json:"identity"`
	Name          string    `json:"name"`
	Slug          string    `json:"slug"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	ObjectVersion int       `json:"objectVersion"`
	Status        string    `json:"status"`

	Labels      Labels      `json:"labels"`
	Annotations Annotations `json:"annotations"`
	CIDRs       []string    `json:"cidrs"`

	Organisation  *base.Organisation `json:"organisation"`
	CloudRegion   *Region            `json:"cloudRegion"`
	Subnets       []Subnet           `json:"subnets"`
	FirewallRules []VpcFirewallRule  `json:"firewallRules"`
}

// Subnet
type Subnet struct {
	Identity      string    `json:"identity"`
	Name          string    `json:"name"`
	Slug          string    `json:"slug"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	ObjectVersion int       `json:"objectVersion"`

	Labels      Labels      `json:"labels"`
	Annotations Annotations `json:"annotations"`

	Type        SubnetType   `json:"type"`
	VpcIdentity string       `json:"vpcIdentity"`
	Vpc         *Vpc         `json:"vpc"`
	Cidr        string       `json:"cidr"`
	Status      SubnetStatus `json:"status"`

	RouteTable *RouteTable `json:"routeTable,omitempty"`

	V4usingIPs     int `json:"v4usingIPs"`
	V4availableIPs int `json:"v4availableIPs"`
	V6usingIPs     int `json:"v6usingIPs"`
	V6availableIPs int `json:"v6availableIPs"`
}

type VpcNatGateway struct {
	Identity      string    `json:"identity"`
	Name          string    `json:"name"`
	Slug          string    `json:"slug"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	ObjectVersion int       `json:"objectVersion"`
	Status        string    `json:"status"`

	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`

	Organisation   *base.Organisation `json:"organisation"`
	VpcIdentity    string             `json:"vpcIdentity"`
	Vpc            *Vpc               `json:"vpc"`
	SubnetIdentity string             `json:"subnetIdentity"`
	Subnet         *Subnet            `json:"subnet"`
	EndpointIP     string             `json:"endpointIP"`

	V4IP string `json:"v4IP"`
	V6IP string `json:"v6IP"`
}

type VpcLoadbalancer struct {
	Identity      string      `json:"identity"`
	Name          string      `json:"name"`
	Slug          string      `json:"slug"`
	Description   string      `json:"description"`
	CreatedAt     time.Time   `json:"createdAt"`
	UpdatedAt     time.Time   `json:"updatedAt"`
	ObjectVersion int         `json:"objectVersion"`
	Labels        Labels      `json:"labels"`
	Annotations   Annotations `json:"annotations"`
	Status        string      `json:"status"`

	Organisation   *base.Organisation `json:"organisation"`
	VpcIdentity    string             `json:"vpcIdentity"`
	Vpc            *Vpc               `json:"vpc"`
	SubnetIdentity string             `json:"subnetIdentity"`
	Subnet         *Subnet            `json:"subnet"`

	ExternalIpAddresses []string `json:"externalIpAddresses"`
	InternalIpAddresses []string `json:"internalIpAddresses"`
	Hostname            string   `json:"hostname"`

	LoadbalancerListeners []VpcLoadbalancerListener `json:"loadbalancerListeners"`
}

type Volume struct {
	Identity      string    `json:"identity"`
	Name          string    `json:"name"`
	Slug          string    `json:"slug"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	ObjectVersion int       `json:"objectVersion"`
	Status        string    `json:"status"`

	Labels      Labels      `json:"labels"`
	Annotations Annotations `json:"annotations"`

	// SourceMachineImage is the machine image that was used to create the volume. Only set if the volume was created from a machine image.
	SourceMachineImage *MachineImage      `json:"sourceMachineImage"`
	VolumeType         *VolumeType        `json:"volumeType"`
	Attachments        []VolumeAttachment `json:"attachments"`
	Organisation       *base.Organisation `json:"organisation"`
	Region             *Region            `json:"cloudRegion"`
	// AvailabilityZones is a list of availability zones that the volume can be attached to.
	AvailabilityZones []Zone `json:"availabilityZones,omitempty"`
	Size              int    `json:"size"`
	DeleteProtection  bool   `json:"deleteProtection"`
}

type VolumeType struct {
	Identity    string `json:"identity"`
	Name        string `json:"name"`
	Description string `json:"description"`
	StorageType string `json:"storageType"`
	AllowResize bool   `json:"allowResize"`
}

type VolumeAttachment struct {
	Identity               string     `json:"identity"`
	CreatedAt              time.Time  `json:"createdAt"`
	Description            string     `json:"description"`
	Serial                 string     `json:"serial"`
	AttachedToIdentity     string     `json:"attachedToIdentity"`
	AttachedToResourceType string     `json:"attachedToResourceType"`
	DetachmentRequestedAt  *time.Time `json:"detachmentRequestedAt,omitempty"`
	CanDetach              bool       `json:"canDetach"`

	// Only set if attachedToResourceType == "cloud_virtual_machine"
	VirtualMachine *Machine `json:"virtualMachine"`

	PersistentVolume *Volume `json:"persistentVolume"`
}

type VpcGatewayEndpoint struct {
	Identity      string      `json:"identity"`
	Name          string      `json:"name"`
	Slug          string      `json:"slug"`
	Description   *string     `json:"description,omitempty"`
	Labels        Labels      `json:"labels,omitempty"`
	Annotations   Annotations `json:"annotations,omitempty"`
	CreatedAt     time.Time   `json:"createdAt"`
	UpdatedAt     *time.Time  `json:"updatedAt,omitempty"`
	DeletedAt     *time.Time  `json:"deletedAt,omitempty"`
	ObjectVersion int         `json:"objectVersion"`

	EndpointAddress  string             `json:"endpointAddress"`
	EndpointHostname string             `json:"endpointHostname"`
	Vpc              *Vpc               `json:"vpc,omitempty"`
	Organisation     *base.Organisation `json:"organisation,omitempty"`
	CloudRegion      *Region            `json:"cloudRegion,omitempty"`
	Subnet           *Subnet            `json:"subnet,omitempty"`
	Status           string             `json:"status"`
}

type RouteTable struct {
	Identity      string      `json:"identity"`
	Name          string      `json:"name"`
	Slug          string      `json:"slug"`
	Description   *string     `json:"description,omitempty"`
	Labels        Labels      `json:"labels,omitempty"`
	Annotations   Annotations `json:"annotations,omitempty"`
	CreatedAt     time.Time   `json:"createdAt"`
	UpdatedAt     *time.Time  `json:"updatedAt,omitempty"`
	DeletedAt     *time.Time  `json:"deletedAt,omitempty"`
	ObjectVersion int         `json:"objectVersion"`

	Organisation      *base.Organisation `json:"organisation,omitempty"`
	Vpc               *Vpc               `json:"vpc"`
	Routes            []RouteEntry       `json:"routes,omitempty"`
	IsDefault         bool               `json:"isDefault"`
	AssociatedSubnets []Subnet           `json:"associatedSubnets"`
}

type RouteEntry struct {
	Identity                 string              `json:"identity"`
	Note                     *string             `json:"note,omitempty"`
	RouteTable               *RouteTable         `json:"routeTable,omitempty"`
	DestinationCidrBlock     string              `json:"destinationCidrBlock"`
	TargetGatewayIdentity    *string             `json:"targetGatewayIdentity,omitempty"`
	TargetGateway            *VpcGatewayEndpoint `json:"targetGateway,omitempty"`
	TargetNatGatewayIdentity *string             `json:"targetNatGatewayIdentity,omitempty"`
	TargetNatGateway         *VpcNatGateway      `json:"targetNatGateway,omitempty"`
	GatewayAddress           *string             `json:"gatewayAddress,omitempty"`
	TargetGatewayEndpoint    *VpcGatewayEndpoint `json:"targetGatewayEndpoint,omitempty"`
	Type                     string              `json:"type"`
}

type ResourceStatus struct {
	Status             string    `json:"status"`
	StatusMessage      string    `json:"statusMessage"`
	LastTransitionTime time.Time `json:"lastTransitionTime"`
}

type MachineTypeCategory struct {
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	MachineTypes []MachineType `json:"machineTypes"`
}

type MachineType struct {
	Identity    string `json:"identity"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Vcpus       int    `json:"vcpus"`
	RamMb       int    `json:"ramMb"`
	DiskGb      int    `json:"diskGb"`
	SwapMb      int    `json:"swapMb"`
}

type MachineImage struct {
	Identity     string            `json:"identity"`
	Name         string            `json:"name"`
	Slug         string            `json:"slug"`
	Labels       map[string]string `json:"labels"`
	Description  string            `json:"description"`
	Architecture string            `json:"architecture"`
}

type CreateVpc struct {
	Name                string      `json:"name"`
	Description         string      `json:"description"`
	Labels              Labels      `json:"labels"`
	Annotations         Annotations `json:"annotations"`
	CloudRegionIdentity string      `json:"cloudRegionIdentity"`
	VpcCidrs            []string    `json:"vpcCidrs"`
}

type UpdateVpc struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Labels      Labels      `json:"labels"`
	Annotations Annotations `json:"annotations"`
	VpcCidrs    []string    `json:"vpcCidrs"`
}

type CreateVolume struct {
	Name                string      `json:"name"`
	Description         string      `json:"description"`
	Labels              Labels      `json:"labels"`
	Annotations         Annotations `json:"annotations"`
	Type                string      `json:"type"`
	Size                int         `json:"size"`
	CloudRegionIdentity string      `json:"cloudRegionIdentity"`
	VolumeTypeIdentity  string      `json:"volumeTypeIdentity"`
	DeleteProtection    bool        `json:"deleteProtection"`
}

type UpdateVolume struct {
	Name             string      `json:"name"`
	Description      string      `json:"description"`
	Labels           Labels      `json:"labels"`
	Annotations      Annotations `json:"annotations"`
	Size             int         `json:"size"`
	DeleteProtection bool        `json:"deleteProtection"`
}

type AttachVolumeRequest struct {
	Description      string `json:"description"`
	ResourceType     string `json:"resourceType"`
	ResourceIdentity string `json:"resourceIdentity"`
}

type DetachVolumeRequest struct {
	ResourceType     string `json:"resourceType"`
	ResourceIdentity string `json:"resourceIdentity"`
}

type CreateSubnet struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Labels      Labels      `json:"labels,omitempty"`
	Annotations Annotations `json:"annotations,omitempty"`
	VpcIdentity string      `json:"vpcIdentity"`
	// Cidr is the CIDR block for the subnet. i.e. 10.0.0.0/24.
	// Supports both IPv4 and IPv6 and Dual-Stack (IPv4 and IPv6). Provide DualStack CIDR by using comma separated values. For example: "10.0.0.0/24,fc00::/64"
	// CIDR subnets must be unique within the VPC and must fall within the VPC CIDR blocks associated with the VPC.
	// The CIDR may not be modified once the subnet is created.
	Cidr string `json:"cidr"`
	// AssociatedRouteTableIdentity is the identity of the route table that will be associated with the subnet.
	AssociatedRouteTableIdentity *string `json:"associatedRouteTableIdentity,omitempty"`
}

type UpdateSubnet struct {
	Name                         string      `json:"name"`
	Description                  string      `json:"description"`
	Labels                       Labels      `json:"labels,omitempty"`
	Annotations                  Annotations `json:"annotations,omitempty"`
	AssociatedRouteTableIdentity *string     `json:"associatedRouteTableIdentity,omitempty"`
}

type UpdateRouteTableRoutes struct {
	Routes []UpdateRouteTableRoute `json:"routes"`
}

type CreateRouteTableRoute struct {
	DestinationCidrBlock     string `json:"destinationCidrBlock"`
	TargetGatewayIdentity    string `json:"targetGatewayIdentity,omitempty"`
	TargetNatGatewayIdentity string `json:"targetNatGatewayIdentity,omitempty"`
	GatewayAddress           string `json:"gatewayAddress,omitempty"`
}

type UpdateRouteTableRoute struct {
	DestinationCidrBlock     string `json:"destinationCidrBlock"`
	TargetGatewayIdentity    string `json:"targetGatewayIdentity,omitempty"`
	TargetNatGatewayIdentity string `json:"targetNatGatewayIdentity,omitempty"`
	GatewayAddress           string `json:"gatewayAddress,omitempty"`
}

type CreateRouteTable struct {
	Name        string      `json:"name"`
	Description *string     `json:"description,omitempty"`
	Labels      Labels      `json:"labels,omitempty"`
	Annotations Annotations `json:"annotations,omitempty"`
	VpcIdentity string      `json:"vpcIdentity"`
}

type UpdateRouteTable struct {
	Name        *string     `json:"name,omitempty"`
	Description *string     `json:"description,omitempty"`
	Labels      Labels      `json:"labels,omitempty"`
	Annotations Annotations `json:"annotations,omitempty"`
}

type CreateVpcNatGateway struct {
	Name           string      `json:"name"`
	Description    string      `json:"description"`
	Labels         Labels      `json:"labels"`
	Annotations    Annotations `json:"annotations"`
	SubnetIdentity string      `json:"subnetIdentity"`
}

type UpdateVpcNatGateway struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Labels      Labels      `json:"labels"`
	Annotations Annotations `json:"annotations"`
}

type CreateMachine struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Labels      Labels      `json:"labels"`
	Annotations Annotations `json:"annotations"`

	// Subnet is the subnet in which the machine will be deployed.
	Subnet string `json:"subnet"`

	// CloudInit is the cloud-init configuration for the machine.
	// If non empty, will be used to populate the cloud-init configuration for the machine.
	CloudInit string `json:"cloudInit"`
	// CloudInitRef is the reference to the cloud-init configuration for the machine.
	// If non empty, will be used to populate the cloud-init configuration for the machine. If cloudInit is also provided, cloudInit will take precedence.
	CloudInitRef string `json:"cloudInitRef"`

	// DeleteProtection is a flag that indicates whether the machine should be protected from deletion.
	// Meaning delete protection will require to be disabled explicitly before the machine can be deleted.
	DeleteProtection bool `json:"deleteProtection"`

	// State is the initial state of the machine. If not provided, the machine will be created in the "running" state. Must be one of running or stopped.
	State *MachineState `json:"state"`

	MachineImage string              `json:"machineImage"`
	MachineType  string              `json:"machineType"`
	RootVolume   CreateMachineVolume `json:"rootVolume"`
	// AvailabilityZone is the availability zone in which the machine will be deployed. This is the slug of the availability zone. Must match the region of the VPC and subnet.
	// If not provided, thee machine will be scheduled within a random zone within the region of the VPC.
	AvailabilityZone *string `json:"availabilityZone"`
	// SecurityGroupAttachments is a list of security group identities that will be attached to the virtual machine instance.
	SecurityGroupAttachments []string `json:"securityGroupAttachments,omitempty"`
}

type CreateMachineVolume struct {
	ExistingVolumeRef  *string `json:"existingVolumeRef,omitempty"`
	VolumeTypeIdentity string  `json:"volumeTypeIdentity"`
	Size               int     `json:"size"`
	Name               *string `json:"name,omitempty"`
	Description        *string `json:"description,omitempty"`
}

type UpdateMachine struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Labels      Labels      `json:"labels"`
	Annotations Annotations `json:"annotations"`
	Subnet      *string     `json:"subnet,omitempty"`

	// State is the new state of the machine. Must be one of running or stopped.
	State *MachineState `json:"state"`

	// AvailabilityZone is the availability zone in which the machine will be deployed.
	// You can use this to move the machine to a different zone.
	// NOTE: Only possible if the cloud region supports migrating machines between zones.
	// Depending on the region, VMs may be live migrated automatically. If not supported, the machine will be stopped and restarted in the new zone.
	AvailabilityZone *string `json:"availabilityZone,omitempty"`

	MachineType      *string `json:"machineType,omitempty"`
	DeleteProtection *bool   `json:"deleteProtection,omitempty"`

	// SecurityGroupAttachments is a list of security group identities that will be attached to the virtual machine instance.
	SecurityGroupAttachments []string `json:"securityGroupAttachments,omitempty"`
}
