package tfs

import (
	"time"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/iaas"
	"github.com/thalassa-cloud/client-go/pkg/base"
)

// Tfs represents a Thalassa Filesystem Service (TFS) instance
// TFS provides a high-availability, multi-availability zone Network File System (NFS) service
// for shared storage across your infrastructure. TFS supports NFSv4 and NFSv4.1 protocols.
type TfsInstance struct {
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

	Region *iaas.Region `json:"region,omitempty"`

	Organisation *base.Organisation `json:"organisation,omitempty"`

	Vpc *iaas.Vpc `json:"vpc,omitempty"`

	// SubnetId is the subnet where the TFS instance is deployed
	Subnet *iaas.Subnet `json:"subnet,omitempty"`

	// Endpoints is a list of endpoints that are associated with the TFS instance
	Endpoints []iaas.Endpoint `json:"endpoints,omitempty"`

	// SecurityGroups is a list of security groups attached to the TFS instance
	SecurityGroups []iaas.SecurityGroup `json:"securityGroups,omitempty"`

	// SizeGB is the size of the TFS instance in GB
	SizeGB int `json:"size"`

	// DeleteProtection is a flag that prevents the TFS instance from being deleted
	DeleteProtection bool `json:"deleteProtection"`

	// Status is the status of the TFS instance
	Status TfsStatus `json:"status"`

	// LastStatusChangedAt is the time the status of the TFS instance was last changed
	LastStatusChangedAt *time.Time `json:"lastStatusChangedAt,omitempty"`
}

type TfsStatus string

const (
	// TfsStatusCreating is the status of the TFS instance that is being created
	TfsStatusCreating     TfsStatus = "Creating"
	TfsStatusProvisioning TfsStatus = "Provisioning"
	TfsStatusAvailable    TfsStatus = "Available"
	TfsStatusDeleting     TfsStatus = "Deleting"
	TfsStatusDeleted      TfsStatus = "Deleted"
	TfsStatusError        TfsStatus = "Error"
	TfsStatusUnknown      TfsStatus = "Unknown"
)

type CreateTfsInstanceRequest struct {
	// Name is the name of the TFS instance
	Name string `json:"name"`

	// Description is a human-readable description of the object
	Description string `json:"description,omitempty"`
	// Annotations is a map of key-value pairs used for storing additional information
	Annotations Annotations `json:"annotations,omitempty"`

	// Labels is a map of key-value pairs used for filtering and grouping objects
	Labels Labels `json:"labels,omitempty"`

	// CloudRegionIdentity is the identity of the cloud region to create the TFS instance in
	CloudRegionIdentity string `json:"cloudRegionIdentity"`

	// VpcIdentity is the identity of the VPC to create the TFS instance in
	VpcIdentity string `json:"vpcIdentity"`

	// SubnetIdentity is the identity of the subnet to create the TFS instance in
	SubnetIdentity string `json:"subnetIdentity"`

	// SizeGB is the size of the TFS instance in GB
	SizeGB int `json:"size"`

	// SecurityGroupAttachments is a list of security group identities to attach to the TFS instance
	SecurityGroupAttachments []string `json:"securityGroupAttachments,omitempty"`

	// DeleteProtection is a flag that prevents the TFS instance from being deleted
	DeleteProtection bool `json:"deleteProtection"`
}

type UpdateTfsInstanceRequest struct {
	// Name is the name of the TFS instance
	Name string `json:"name"`

	// Description is a human-readable description of the object
	Description string `json:"description,omitempty"`
	// Annotations is a map of key-value pairs used for storing additional information
	Annotations Annotations `json:"annotations,omitempty"`

	// Labels is a map of key-value pairs used for filtering and grouping objects
	Labels Labels `json:"labels,omitempty"`

	// SizeGB is the size of the TFS instance in GB
	SizeGB int `json:"size"`

	// SecurityGroupAttachments is a list of security group identities to attach to the TFS instance
	SecurityGroupAttachments []string `json:"securityGroupAttachments,omitempty"`

	// DeleteProtection is a flag that prevents the TFS instance from being deleted
	DeleteProtection bool `json:"deleteProtection"`
}

type ListTfsInstancesRequest struct {
	// Filters is a list of filters to apply to the list of TFS instances
	Filters []filters.Filter
}
