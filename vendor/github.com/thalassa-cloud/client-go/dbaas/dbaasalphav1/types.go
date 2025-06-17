package dbaasalphav1

import (
	"time"

	"github.com/thalassa-cloud/client-go/iaas"
	"github.com/thalassa-cloud/client-go/pkg/base"
)

type DbClusterDatabaseEngine string

const (
	DbClusterDatabaseEnginePostgres DbClusterDatabaseEngine = "postgres"
)

type DbCluster struct {
	Identity      string            `json:"identity"`
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	CreatedAt     time.Time         `json:"createdAt"`
	UpdatedAt     time.Time         `json:"updatedAt"`
	ObjectVersion int               `json:"objectVersion"`
	Labels        map[string]string `json:"labels"`
	Annotations   map[string]string `json:"annotations"`

	Organisation *base.Organisation `json:"organisation,omitempty"`
	// Vpc is the VPC the cluster is deployed in
	Vpc    *iaas.Vpc    `json:"vpc,omitempty"`
	Region *iaas.Region `json:"region,omitempty"`
	// Subnet is the subnet the cluster is deployed in
	Subnet *iaas.Subnet `json:"subnet,omitempty"`
	// DatabaseInstanceType is the instance type used to determine the size of the cluster instances
	DatabaseInstanceType *iaas.MachineType `json:"database_instance_type,omitempty"`
	// Replicas is the number of instances in the cluster
	Replicas int `json:"replicas"`
	// Engine is the database engine of the cluster
	Engine DbClusterDatabaseEngine `json:"engine"`
	// EngineVersion is the version of the database engine
	EngineVersion string `json:"engineVersion"`
	// DatabaseEngineVersion is the version of the database engine
	DatabaseEngineVersion *DbClusterEngineVersion `json:"database_engine_version,omitempty"`
	// // DbParameterGroupId is the ID of the database parameter group
	// Parameters is a map of parameter name to database engine specific parameter value
	Parameters map[string]string `json:"parameters"`
	// AllocatedStorage is the amount of storage allocated to the cluster in GB
	AllocatedStorage uint64 `json:"allocatedStorage"`
	// VolumeTypeClass is the storage type used to determine the size of the cluster storage
	VolumeTypeClass *iaas.VolumeType `json:"volume_type_class,omitempty"`
	// AutoMinorVersionUpgrade is a flag indicating if the cluster should automatically upgrade to the latest minor version
	AutoMinorVersionUpgrade bool `json:"autoMinorVersionUpgrade"`
	// DatabaseName is the name of the database on the cluster. Optional name. If provided, it will be used as the name of the database on the cluster.
	DatabaseName *string `json:"databaseName"`
	// DeleteProtection is a flag indicating if the cluster should be protected from deletion. The database cannot be deleted if this is true.
	DeleteProtection bool `json:"deleteProtection"`
	// SecurityGroups is a list of security groups associated with the cluster
	SecurityGroups []iaas.SecurityGroup `json:"securityGroups,omitempty"`
	// Status is the status of the cluster
	Status DbClusterStatus `json:"status"`
	// EndpointIpv4 is the IPv4 address of the cluster endpoint
	EndpointIpv4 string `json:"endpointIpv4"`
	// EndpointIpv6 is the IPv6 address of the cluster endpoint
	EndpointIpv6 string `json:"endpointIpv6"`
	// Port is the port of the cluster endpoint
	Port int `json:"port"`
}

type DbClusterEngineVersion struct {
	Identity string `json:"identity"`
	// CreatedAt is the date and time the object was created
	CreatedAt time.Time `json:"createdAt"`
	// Engine is the database engine
	Engine DbClusterDatabaseEngine `json:"engine"`
	// EngineVersion is the version of the database engine
	EngineVersion string `json:"engineVersion"`
	// MajorVersion is the major version of the engine
	MajorVersion int `json:"majorVersion"`
	// MinorVersion is the minor version of the engine
	MinorVersion int `json:"minorVersion"`
	// Supported is a flag indicating if the engine version is supported
	Supported bool `json:"supported"`
	// MinMajorVersionUpgradeFrom is the minimum major version required to upgrade from
	MinMajorVersionUpgradeFrom *int `json:"minMajorVersionUpgradeFrom"`
	// MinMinorVersionUpgradeFrom is the minimum minor version required to upgrade from
	MinMinorVersionUpgradeFrom *int `json:"minMinorVersionUpgradeFrom"`
	// MaxMajorVersionUpgradeTo is the maximum major version that can be upgraded to
	MaxMajorVersionUpgradeTo *int `json:"maxMajorVersionUpgradeTo"`
	// MaxMinorVersionUpgradeTo is the maximum minor version that can be upgraded to
	MaxMinorVersionUpgradeTo *int `json:"maxMinorVersionUpgradeTo"`
	// Enabled is a flag indicating if the engine version is enabled
	Enabled bool `json:"enabled"`
	// DefaultParameters is a map of parameter name to database engine specific parameter value
	DefaultParameters map[string]string `json:"defaultParameters"`
}

type CreateDbClusterRequest struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	// Subnet is the subnet identity of the cloud subnet
	SubnetIdentity           string   `json:"subnetIdentity"`
	SecurityGroupAttachments []string `json:"securityGroupAttachments"`
	DeleteProtection         bool     `json:"deleteProtection"`
	// Engine is the database engine
	Engine DbClusterDatabaseEngine `json:"engine"`
	// EngineVersion is the version of the database engine
	EngineVersion string `json:"engineVersion"`
	// Parameters is a map of parameter name to database engine specific parameter value
	Parameters map[string]string `json:"parameters"`
	// AllocatedStorage is the amount of storage allocated to the cluster in GB
	AllocatedStorage uint64 `json:"allocatedStorage"`
	// VolumeTypeClassIdentity is the identity of the storage type
	VolumeTypeClassIdentity string `json:"volumeTypeClassIdentity"`
	// DatabaseInstanceTypeIdentity is the identity of the database instance type
	DatabaseInstanceTypeIdentity string `json:"databaseInstanceTypeIdentity"`
	// AutoMinorVersionUpgrade is a flag indicating if the cluster should automatically upgrade to the latest minor version
	AutoMinorVersionUpgrade bool `json:"autoMinorVersionUpgrade"`
	// DatabaseName is the name of the database on the cluster. Optional name. If provided, it will be used as the name of the database on the cluster.
	DatabaseName *string `json:"databaseName"`
	// Replicas is the number of instances in the cluster
	Replicas int `json:"replicas"`
}

type UpdateDbClusterRequest struct {
	Name                     string            `json:"name"`
	Description              string            `json:"description"`
	Labels                   map[string]string `json:"labels"`
	Annotations              map[string]string `json:"annotations"`
	SecurityGroupAttachments []string          `json:"securityGroupAttachments"`
	DeleteProtection         bool              `json:"deleteProtection"`
	// EngineVersion is the version of the database engine
	EngineVersion string `json:"engineVersion"`
	// Parameters is a map of parameter name to database engine specific parameter value
	Parameters map[string]string `json:"parameters"`
	// AllocatedStorage is the amount of storage allocated to the cluster in GB
	AllocatedStorage uint64 `json:"allocatedStorage"`
	// AutoMinorVersionUpgrade is a flag indicating if the cluster should automatically upgrade to the latest minor version
	AutoMinorVersionUpgrade bool `json:"autoMinorVersionUpgrade"`
	// DatabaseName is the name of the database on the cluster. Optional name. If provided, it will be used as the name of the database on the cluster.
	DatabaseName *string `json:"databaseName"`
	// Replicas is the number of instances in the cluster
	Replicas int `json:"replicas"`
	// DatabaseInstanceTypeIdentity is the identity of the database instance type. Optional identity. If provided, it will be used as the database instance type for the cluster.
	DatabaseInstanceTypeIdentity *string `json:"databaseInstanceTypeIdentity"`
}

type DbClusterStatus string

const (
	DbClusterStatusPending               DbClusterStatus = "pending"
	DbClusterStatusCreating              DbClusterStatus = "creating"
	DbClusterStatusReady                 DbClusterStatus = "ready"
	DbClusterStatusUpdating              DbClusterStatus = "updating"
	DbClusterStatusUpgradingMajorVersion DbClusterStatus = "upgrading-major-version"
	DbClusterStatusUpgradingMinorVersion DbClusterStatus = "upgrading-minor-version"
	DbClusterStatusFailed                DbClusterStatus = "failed"
	DbClusterStatusDeleting              DbClusterStatus = "deleting"
	DbClusterStatusDeleted               DbClusterStatus = "deleted"
	DbClusterStatusUnknown               DbClusterStatus = "unknown"
)
