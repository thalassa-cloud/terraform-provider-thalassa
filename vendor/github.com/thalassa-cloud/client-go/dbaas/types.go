package dbaas

import (
	"time"

	"github.com/thalassa-cloud/client-go/iaas"
	"github.com/thalassa-cloud/client-go/iam"
	"github.com/thalassa-cloud/client-go/objectstorage"
	"github.com/thalassa-cloud/client-go/pkg/base"
)

type DbClusterDatabaseEngine string

const (
	DbClusterDatabaseEnginePostgres DbClusterDatabaseEngine = "postgres"
)

type DbCluster struct {
	Identity      string      `json:"identity"`
	Name          string      `json:"name"`
	Slug          string      `json:"slug"`
	Description   string      `json:"description"`
	CreatedAt     time.Time   `json:"createdAt"`
	UpdatedAt     time.Time   `json:"updatedAt"`
	ObjectVersion int         `json:"objectVersion"`
	Labels        Labels      `json:"labels"`
	Annotations   Annotations `json:"annotations"`

	Organisation *base.Organisation `json:"organisation,omitempty"`
	// Vpc is the VPC the cluster is deployed in
	Vpc    *iaas.Vpc    `json:"vpc,omitempty"`
	Region *iaas.Region `json:"region,omitempty"`
	// Subnet is the subnet the cluster is deployed in
	Subnet *iaas.Subnet `json:"subnet,omitempty"`
	// DatabaseInstanceType is the instance type used to determine the size of the cluster instances
	DatabaseInstanceType *DatabaseInstanceType `json:"database_instance_type,omitempty"`
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
	// DatabaseSize is the size of the database in bytes. This is the total size of all databases on the cluster.
	DatabaseSizeBytes uint64 `json:"databaseSize"`
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

	// Endpoint registration
	InternalEndpointIpv4 *iaas.Endpoint `json:"internalEndpointIpv4,omitempty"`
	InternalEndpointIpv6 *iaas.Endpoint `json:"internalEndpointIpv6,omitempty"`

	// PostgresRoles is a list of PostgreSQL roles associated with the cluster
	PostgresRoles []DbClusterPostgresRole `json:"postgresRoles,omitempty"`
	// PostgresDatabases is a list of PostgreSQL databases associated with the cluster
	PostgresDatabases []DbClusterPostgresDatabase `json:"postgresDatabases,omitempty"`
	// DatabaseInstancesStatus is the status of the database instances in the cluster
	DatabaseInstancesStatus DatabaseInstancesStatus `json:"databaseInstancesStatus"`

	// AutoUpgradePolicy is the auto upgrade policy for the cluster
	AutoUpgradePolicy DbClusterAutoUpgradePolicy `json:"autoUpgradePolicy,omitempty"`
	// MaintenanceDay is the day of the week for the maintenance window. 0 is Sunday, 6 is Saturday.
	MaintenanceDay *uint `json:"maintenanceDay,omitempty"`
	// MaintenanceStartAt is the start time of the maintenance window on the maintenance day in UTC. 0 is 00:00, 23 is 23:00.
	MaintenanceStartAt *uint `json:"maintenanceStartAt,omitempty"`

	// ScheduledMaintenances is the list of scheduled maintenances for the cluster
	ScheduledMaintenances []DbClusterScheduledMaintenance `json:"scheduledMaintenances"`
	// DbObjectStore is the DB object store used for barman backups
	DbObjectStore *DbObjectStore `json:"dbObjectStore,omitempty"`
}

type DatabaseInstancesStatus struct {
	Instances []DatabaseInstanceStatus `json:"instances,omitempty"`
}

type DatabaseInstanceStatus struct {
	// Name is the name of the instance.
	// +optional
	Name string `json:"name,omitempty"`
	// IsPrimary is a flag to indicate if the instance is the primary instance.
	// +optional
	IsPrimary bool `json:"isPrimary,omitempty"`
	// IsPrimaryTarget is a flag to indicate if the instance is the target primary instance.
	// +optional
	IsPrimaryTarget bool `json:"isPrimaryTarget,omitempty"`
	// Replicating is a flag to indicate if the instance is replicating.
	// +optional
	Replicating bool `json:"replicating,omitempty"`
	// TimeLineID is the timeline ID of the instance.
	// +optional
	TimeLineID string `json:"timeLineID,omitempty"`
	// Healthy is a flag to indicate if the instance is healthy.
	// +optional
	Healthy bool `json:"healthy,omitempty"`
	// Joining is a flag to indicate if the instance is joining the cluster.
	// +optional
	Joining bool `json:"joining,omitempty"`
	// AvailabilityZone is the availability zone of the instance.
	// +optional
	AvailabilityZone string `json:"availabilityZone,omitempty"`
	// Version is the version of the instance.
	Version string `json:"version,omitempty"`
	// AllocatedStorage is the amount of storage allocated to the instance in GB.
	AllocatedStorage int `json:"allocatedStorage,omitempty"`
	// UsedStorage is the amount of storage used by the instance in GB.
	// May not always be available - feature gated functionality.
	UsedStorage int `json:"usedStorage,omitempty"`
	// Memory is the memory of the instance in MB.
	Memory int `json:"memory,omitempty"`
	// Cpu is the cpu of the instance in cores.
	Cpu int `json:"cpu,omitempty"`
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

type ListDbClusterEngineVersionsResponse struct {
	Engines map[DbClusterDatabaseEngine][]DbClusterEngineVersion `json:"engines"`
}

type CreateDbClusterRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	// Annotations is a map of key-value pairs used for storing additional information
	Annotations map[string]string `json:"annotations,omitempty"`
	// Labels is a map of key-value pairs used for filtering and grouping objects
	Labels map[string]string `json:"labels,omitempty"`
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
	// AutoUpgradePolicy is the auto upgrade policy for the cluster
	AutoUpgradePolicy *DbClusterAutoUpgradePolicy `json:"autoUpgradePolicy,omitempty"`
	// MaintenanceDay is the day of the week for the maintenance window. 0 is Sunday, 6 is Saturday.
	MaintenanceDay *uint `json:"maintenanceDay,omitempty"`
	// MaintenanceStartAt is the start time of the maintenance window on the maintenance day in UTC. 0 is 00:00, 23 is 23:00.
	MaintenanceStartAt *uint `json:"maintenanceStartAt,omitempty"`
	// Replicas is the number of instances in the cluster
	Replicas int `json:"replicas"`
	// PostgresInitDb is the initial database to create on the cluster
	PostgresInitDb *PostgresInitDb `json:"postgresInitDb,omitempty"`

	// RestoreFromBackupIdentity is the identity of the backup to restore from
	RestoreFromBackupIdentity *string `json:"restoreFromBackupIdentity,omitempty"`
	// RestoreRecoveryTarget is the recovery target for Point-In-Time Recovery (PITR)
	// Only used when RestoreFromBackupIdentity is specified
	RestoreRecoveryTarget *RestoreRecoveryTarget `json:"restoreRecoveryTarget,omitempty"`
	// DbObjectStoreIdentity is the identity of the DB object store used for barman backups (optional)
	DbObjectStoreIdentity *string `json:"dbObjectStoreIdentity,omitempty"`
	// ProvisionDbObjectStore is a flag to indicate if the DB object store should be provisioned for the cluster. Defaults to false.
	// if true, the DbObjectStoreIdentity will be ignored.
	ProvisionDbObjectStore bool `json:"provisionDbObjectStore,omitempty"`
	// InitialDbBackupSchedule is the initial PostgreSQL backup schedule to create for the cluster.
	// Only for clusters with engine `postgres` using barman backups.
	InitialDbBackupSchedule *CreateDbBackupScheduleRequest `json:"initialDbBackupSchedule,omitempty"`
}

// RestoreRecoveryTarget specifies the recovery target for Point-In-Time Recovery
type RestoreRecoveryTarget struct {
	// TargetTime is the timestamp to restore to (RFC3339 format)
	// Example: "2023-12-25T10:00:00Z"
	TargetTime *string `json:"targetTime,omitempty"`
	// TargetLSN is the Log Sequence Number to restore to
	// Example: "0/1234567"
	TargetLSN *string `json:"targetLSN,omitempty"`
}

type DbClusterAutoUpgradePolicy string

const (
	// DbClusterAutoUpgradePolicyNone does not perform any auto upgrades. User is expected to manually upgrade the cluster.
	DbClusterAutoUpgradePolicyNone DbClusterAutoUpgradePolicy = "none"
	// DbClusterAutoUpgradePolicyLatestVersion is the auto upgrade policy for the cluster.
	// It will upgrade to the latest release of the latest supported minor version.
	// This upgrade strategy is recommended for development clusters.
	DbClusterAutoUpgradePolicyLatestVersion DbClusterAutoUpgradePolicy = "latest-version"
	// DbClusterAutoUpgradePolicyLatestStable is the auto upgrade policy for the cluster.
	// It will upgrade only when the current version is not supported anymore, and then only to supported versions.
	// This upgrade strategy is recommended for production clusters that want minimal changes.
	DbClusterAutoUpgradePolicyLatestStable DbClusterAutoUpgradePolicy = "latest-stable"
	// DbClusterAutoUpgradePolicyLatestPatch is the auto upgrade policy for the cluster.
	// It will stay in the same minor version, only upgrading to the highest available minor if no new patches
	// are available and the version isn't supported anymore.
	DbClusterAutoUpgradePolicyLatestPatch DbClusterAutoUpgradePolicy = "latest-patch"
	// DbClusterAutoUpgradePolicyLatestMinor is the auto upgrade policy for the cluster.
	// It will attempt to stay in the same major version, only updating to the next major once it's not supported anymore.
	DbClusterAutoUpgradePolicyLatestMinor DbClusterAutoUpgradePolicy = "latest-minor"
	// DbClusterAutoUpgradePolicyLatestMajor is the auto upgrade policy for the cluster.
	// It will update to the highest available and supported major version.
	DbClusterAutoUpgradePolicyLatestMajor DbClusterAutoUpgradePolicy = "latest-major"
)

type PostgresInitDb struct {
	// DataChecksums is a flag to indicate if data checksums should be enabled
	DataChecksums bool `json:"dataChecksums,omitempty"`
	// Maps to the `ENCODING` parameter of `CREATE DATABASE`. This setting
	// cannot be changed. Character set encoding to use in the database.
	Encoding string `json:"encoding,omitempty"`

	// Maps to the `LOCALE` parameter of `CREATE DATABASE`. This setting
	// cannot be changed. Sets the default collation order and character
	// classification in the new database.
	Locale string `json:"locale,omitempty"`

	// Maps to the `LOCALE_PROVIDER` parameter of `CREATE DATABASE`. This
	// setting cannot be changed. This option sets the locale provider for
	// databases created in the new cluster. Available from PostgreSQL 16.
	LocaleProvider string `json:"localeProvider,omitempty"`

	// Maps to the `LC_COLLATE` parameter of `CREATE DATABASE`. This setting cannot be changed.
	LcCollate string `json:"localeCollate,omitempty"`

	// Maps to the `LC_CTYPE` parameter of `CREATE DATABASE`. This setting cannot be changed.
	LcCtype string `json:"localeCType,omitempty"`

	// Maps to the `ICU_LOCALE` parameter of `CREATE DATABASE`. This setting cannot be changed.
	// Specifies the ICU locale when the ICU provider is used.
	// This option requires `localeProvider` to be set to `icu`. Available from PostgreSQL 15.
	IcuLocale string `json:"icuLocale,omitempty"`

	// Maps to the `ICU_RULES` parameter of `CREATE DATABASE`. This setting cannot be changed.
	// Specifies additional collation rules to customize the behavior of the default collation.
	// This option requires `localeProvider` to be set to `icu`. Available from PostgreSQL 16.
	IcuRules string `json:"icuRules,omitempty"`

	// Maps to the `BUILTIN_LOCALE` parameter of `CREATE DATABASE`. This setting cannot be changed.
	// Specifies the locale name when the builtin provider is used. This option requires `localeProvider` to be set to `builtin`.
	// Available from PostgreSQL 17.
	BuiltinLocale string `json:"builtinLocale,omitempty"`
	// Maps to the `COLLATION_VERSION` parameter of `CREATE DATABASE`. This setting cannot be changed.
	// CollationVersion string `json:"collationVersion,omitempty"`
	// The value in megabytes (1 to 1024) to be passed to the `--wal-segsize`
	// option for initdb (default: empty, resulting in PostgreSQL default: 16MB)
	// +optional
	WalSegmentSize int `json:"walSegmentSize,omitempty"`
}

type UpdateDbClusterRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	// Annotations is a map of key-value pairs used for storing additional information
	Annotations Annotations `json:"annotations,omitempty"`
	// Labels is a map of key-value pairs used for filtering and grouping objects
	Labels                   Labels   `json:"labels,omitempty"`
	SecurityGroupAttachments []string `json:"securityGroupAttachments"`
	DeleteProtection         bool     `json:"deleteProtection"`
	// EngineVersion is the version of the database engine
	EngineVersion *string `json:"engineVersion,omitempty"`
	// Parameters is a map of parameter name to database engine specific parameter value
	Parameters map[string]string `json:"parameters"`
	// AllocatedStorage is the amount of storage allocated to the cluster in GB
	AllocatedStorage uint64 `json:"allocatedStorage"`
	// AutoUpgradePolicy is the auto upgrade policy for the cluster
	AutoUpgradePolicy *DbClusterAutoUpgradePolicy `json:"autoUpgradePolicy,omitempty"`
	// MaintenanceDay is the day of the week for the maintenance window. 0 is Sunday, 6 is Saturday.
	MaintenanceDay *uint `json:"maintenanceDay,omitempty"`
	// MaintenanceStartAt is the start time of the maintenance window on the maintenance day in UTC. 0 is 00:00, 23 is 23:00.
	MaintenanceStartAt *uint `json:"maintenanceStartAt,omitempty"`
	// Replicas is the number of instances in the cluster
	Replicas int `json:"replicas"`
	// DatabaseInstanceTypeIdentity is the identity of the database instance type. Optional identity. If provided, it will be used as the database instance type for the cluster.
	DatabaseInstanceTypeIdentity *string `json:"databaseInstanceTypeIdentity,omitempty"`
	// DbObjectStoreIdentity is the identity of the DB object store used for barman backups (optional)
	DbObjectStoreIdentity *string `json:"dbObjectStoreIdentity,omitempty"`
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

type CreatePgDatabaseRequest struct {
	// Name is the name of the database.
	Name string `json:"name"`
	// Maps to the `OWNER` parameter of `CREATE DATABASE`.
	// Maps to the `OWNER TO` command of `ALTER DATABASE`.
	// The role name of the user who owns the database inside PostgreSQL.
	Owner string `json:"owner"`

	// Maps to the `ALLOW_CONNECTIONS` parameter of `CREATE DATABASE` and
	// `ALTER DATABASE`. If false then no one can connect to this database.
	AllowConnections *bool `json:"allowConnections,omitempty"`

	// Maps to the `CONNECTION LIMIT` clause of `CREATE DATABASE` and
	// `ALTER DATABASE`. How many concurrent connections can be made to
	// this database. -1 (the default) means no limit.
	ConnectionLimit *int `json:"connectionLimit,omitempty"`
	// Extensions for the database
	Extensions *PgDatabaseExtensions `json:"extensions,omitempty"`
}

type UpdatePgDatabaseRequest struct {
	// Maps to the `ALLOW_CONNECTIONS` parameter of `CREATE DATABASE` and
	// `ALTER DATABASE`. If false then no one can connect to this database.
	AllowConnections *bool `json:"allowConnections,omitempty"`

	// Maps to the `CONNECTION LIMIT` clause of `CREATE DATABASE` and
	// `ALTER DATABASE`. How many concurrent connections can be made to
	// this database. -1 (the default) means no limit.
	ConnectionLimit *int `json:"connectionLimit,omitempty"`
	// Extensions for the database
	Extensions *PgDatabaseExtensions `json:"extensions,omitempty"`
}

type PgDatabaseExtensions struct {
	// Extensions to install
	Extensions []DatabasePostgresExtension `json:"extensions"`
}

type DatabasePostgresExtension struct {
	// Name is the name of the extension
	Name string `json:"name"`
}

type CreatePgRoleRequest struct {
	Name string `json:"name"`
	// Login is a flag to indicate if the role can login
	Login bool `json:"login"`

	// CreateDb is a flag to indicate if the role can create databases
	CreateDb bool `json:"createDb"`

	// CreateRole is a flag to indicate if the role can create roles
	CreateRole bool `json:"createRole"`

	// ConnectionLimit is the maximum number of concurrent connections for the role. Default is -1, as per PostgreSQL default.
	ConnectionLimit int64 `json:"connectionLimit,omitempty"`

	// ValidUntil is the date and time the role will expire
	ValidUntil *time.Time `json:"validUntil,omitempty"`

	// Password is the password for the role
	Password string `json:"password,omitempty"`
}

type UpdatePgRoleRequest struct {
	// ConnectionLimit is the maximum number of concurrent connections for the role. Default is -1, as per PostgreSQL default.
	ConnectionLimit int64 `json:"connectionLimit,omitempty"`

	// ValidUntil is the date and time the role will expire
	ValidUntil *time.Time `json:"validUntil,omitempty"`

	// Password is the password for the role. If provided, the password will be updated. If not provided, the password will not be updated.
	Password *string `json:"password,omitempty"`
}

type CreateDbBackupScheduleRequest struct {
	// Name is the name of the backup schedule
	Name string `json:"name"`
	// Description is the description of the backup schedule
	Description *string `json:"description,omitempty"`
	// Annotations is a map of annotations for the backup schedule
	Annotations Annotations `json:"annotations,omitempty"`
	// Labels is a map of labels for the backup schedule
	Labels Labels `json:"labels,omitempty"`
	// Schedule is the schedule of the backup. Cron expression.
	Schedule string `json:"schedule"`
	// RetentionPolicy is the retention policy of the backup
	RetentionPolicy string `json:"retentionPolicy"`
	// // Target is the target of the backup schedule. Primary or prefer-standby.
	// Target DbClusterBackupScheduleTarget `json:"target,omitempty"`
	// Method is the method of the backup schedule
	Method DbClusterBackupScheduleMethod `json:"method"`
}

type UpdateDbBackupScheduleRequest struct {
	// Name is the name of the backup schedule
	Name string `json:"name"`
	// Description is the description of the backup schedule
	Description string `json:"description"`
	// Annotations is a map of annotations for the backup schedule
	Annotations Annotations `json:"annotations,omitempty"`
	// Labels is a map of labels for the backup schedule
	Labels Labels `json:"labels,omitempty"`
	// Schedule is the schedule of the backup. Cron expression.
	Schedule string `json:"schedule"`
	// RetentionPolicy is the retention policy of the backup
	RetentionPolicy string `json:"retentionPolicy"`
	// // Target is the target of the backup schedule
	// Target DbClusterBackupScheduleTarget `json:"target,omitempty"`
}

type DbClusterBackupScheduleMethod string

const (
	DbClusterBackupScheduleMethodSnapshot DbClusterBackupScheduleMethod = "snapshot"
	DbClusterBackupScheduleMethodBarman   DbClusterBackupScheduleMethod = "barman"
)

type DbClusterBackupScheduleTarget string

const (
	DbClusterBackupScheduleTargetPrimary       DbClusterBackupScheduleTarget = "primary"
	DbClusterBackupScheduleTargetPreferStandby DbClusterBackupScheduleTarget = "prefer-standby"
)

type DbClusterBackupSchedule struct {
	// Identity is a unique identifier for the backup schedule
	Identity string `json:"identity"`

	// Name is a human-readable name of the backup schedule
	Name string `json:"name"`

	// Description is the description of the backup schedule
	Description *string `json:"description,omitempty"`

	// Annotations is a map of annotations for the backup schedule
	Annotations Annotations `json:"annotations,omitempty"`

	// Labels is a map of labels for the backup schedule
	Labels Labels `json:"labels,omitempty"`

	// Status is the status of the role
	Status ObjectStatus `json:"status"`

	// StatusMessage is the message of the role status
	StatusMessage string `json:"statusMessage,omitempty"`

	// CreatedAt is the date and time the object was created
	CreatedAt time.Time `json:"createdAt"`

	// DbCluster is the cluster the backup schedule belongs to
	DbCluster *DbCluster `json:"dbCluster,omitempty"`

	// Organisation is the organisation the backup schedule belongs to
	Organisation *base.Organisation `json:"organisation,omitempty"`

	// Method is the method of the backup schedule
	Method DbClusterBackupScheduleMethod `json:"method"`

	// Schedule is the schedule of the backup. Cron expression.
	// see https://pkg.go.dev/github.com/robfig/cron#hdr-CRON_Expression_Format
	Schedule string `json:"schedule"`

	// RetentionPolicy is the retention policy of the backup
	RetentionPolicy string `json:"retentionPolicy"`

	// NextBackupAt is the date and time the next backup will be taken
	NextBackupAt *time.Time `json:"nextBackupAt,omitempty"`

	// LastBackupAt is the date and time the last backup was taken
	LastBackupAt *time.Time `json:"lastBackupAt,omitempty"`

	// BackupCount is the number of backups
	BackupCount int64 `json:"backupCount"`

	// Suspended is a flag to indicate if the backup schedule is suspended
	Suspended bool `json:"suspended"`

	// Target is the target of the backup schedule
	Target DbClusterBackupScheduleTarget `json:"target"`

	// DeleteScheduledAt is the date and time the backup schedule will be deleted
	DeleteScheduledAt *time.Time `json:"deleteScheduledAt,omitempty"`
}

type CreateDbClusterBackupRequest struct {
	// Name is the name of the backup
	Name string `json:"name"`
	// Description is the description of the backup
	Description *string `json:"description,omitempty"`
	// Labels is a map of labels for the backup
	Labels Labels `json:"labels"`
	// Annotations is a map of annotations for the backup
	Annotations Annotations `json:"annotations"`
	// RetentionPolicy is the retention policy of the backup
	RetentionPolicy *string `json:"retentionPolicy,omitempty"`
}

type DbClusterBackup struct {
	// Identity is a unique identifier for the backup
	Identity string `json:"identity"`
	// DbCluster is the cluster the backup belongs to
	DbCluster *DbCluster `json:"dbCluster,omitempty"`
	// Organisation is the organisation the backup belongs to
	Organisation *base.Organisation `json:"organisation,omitempty"`
	// Labels is a map of labels for the backup
	Labels Labels `json:"labels,omitempty"`
	// Annotations is a map of annotations for the backup
	Annotations Annotations `json:"annotations,omitempty"`
	// Region is the region the backup belongs to
	Region *iaas.Region `json:"region,omitempty"`
	// BackupSchedule is the backup schedule the backup belongs to
	// +optional
	BackupSchedule *DbClusterBackupSchedule `json:"backupSchedule,omitempty"`
	// BackupTrigger is the trigger of the backup
	BackupTrigger DbClusterBackupTrigger `json:"backupTrigger"`
	// EngineType is the type of the database engine, used for back-up restore purposes
	EngineType DbClusterDatabaseEngine `json:"engineType"`
	// DbObjectStore is the object store the backup belongs to. Only set for type barman backups.
	DbObjectStore *DbObjectStore `json:"dbObjectStore,omitempty"`

	// EngineVersion is the version of the database engine, used for back-up restore purposes
	EngineVersion string `json:"engineVersion"`
	// DeleteProtection is a flag to indicate if the backup is protected from deletion. The backup cannot be deleted if this is true.
	DeleteProtection bool `json:"deleteProtection"`
	// BackupType is the type of the backup
	BackupType string `json:"backupType"`

	// Online is a flag to indicate if the backup is an online backup or offline/cold backup
	Online bool `json:"online"`

	// BeginLSN is the starting LSN of the backup
	BeginLSN string `json:"beginLSN"`

	// EndLSN is the ending LSN of the backup
	EndLSN string `json:"endLSN"`

	// BeginWAL is the starting WAL of the backup
	BeginWAL string `json:"beginWAL"`

	// EndWAL is the ending WAL of the backup
	EndWAL string `json:"endWAL"`

	// StartedAt is the date and time the backup started
	StartedAt *time.Time `json:"startedAt,omitempty"`

	// StoppedAt is the date and time the backup stopped
	StoppedAt *time.Time `json:"stoppedAt,omitempty"`

	// Status is the status of the backup
	Status ObjectStatus `json:"status"`

	// StatusMessage is the message of the backup status
	StatusMessage string `json:"statusMessage,omitempty"`

	// CreatedAt is the date and time the object was created
	CreatedAt time.Time `json:"createdAt"`

	// DeleteScheduledAt is the date and time the backup will be deleted
	DeleteScheduledAt *time.Time `json:"deleteScheduledAt,omitempty"`
}

type DbClusterBackupTrigger string

const (
	DbClusterBackupTriggerManual   DbClusterBackupTrigger = "manual"
	DbClusterBackupTriggerSchedule DbClusterBackupTrigger = "schedule"
	DbClusterBackupTriggerSystem   DbClusterBackupTrigger = "system"
)

type DbClusterPostgresGrant struct {
	// Identity is a unique identifier for the grant
	Identity string `json:"identity"`
	// Name is the name of the grant
	Name string `json:"name"`
	// Status is the status of the grant
	Status ObjectStatus `json:"status"`
	// StatusMessage is the message of the grant status
	StatusMessage string `json:"statusMessage,omitempty"`
	// CreatedAt is the date and time the object was created
	CreatedAt time.Time `json:"createdAt"`
	// DbCluster is the cluster the grant belongs to
	DbCluster *DbCluster `json:"dbCluster,omitempty"`
	// Database is the database the grant belongs to
	Database *DbClusterPostgresDatabase `json:"database,omitempty"`
	// Role is the role the grant belongs to
	Role *DbClusterPostgresRole `json:"role,omitempty"`
	// Read is a flag to indicate if the role can read from the database
	Read bool `json:"read"`
	// Write is a flag to indicate if the role can write to the database
	Write bool `json:"write"`
}

type CreatePgGrantRequest struct {
	// Name is the name of the grant
	Name string `json:"name"`
	// RoleName is the name of the role
	RoleName string `json:"roleName"`
	// DatabaseName is the name of the database
	DatabaseName string `json:"databaseName"`
	// Read is a flag to indicate if the role can read from the database
	Read bool `json:"read"`
	// Write is a flag to indicate if the role can write to the database
	Write bool `json:"write"`
}

type UpdatePgGrantRequest struct {
	// Read is a flag to indicate if the role can read from the database
	Read *bool `json:"read"`
	// Write is a flag to indicate if the role can write to the database
	Write *bool `json:"write"`
}

type DbClusterPostgresRole struct {
	// Identity is a unique identifier for the role
	Identity string `json:"identity"`
	// Name is a human-readable name of the role
	Name string `json:"name"`
	// Status is the status of the role
	Status ObjectStatus `json:"status"`
	// StatusMessage is the message of the role status
	StatusMessage string `json:"statusMessage,omitempty"`
	// CreatedAt is the date and time the object was created
	CreatedAt time.Time `json:"createdAt"`
	// DbCluster is the cluster the role belongs to
	DbCluster *DbCluster `json:"dbCluster,omitempty"`
	// Login is a flag to indicate if the role can login
	Login bool `json:"login"`
	// CreateDb is a flag to indicate if the role can create databases
	CreateDb bool `json:"createDb"`
	// CreateRole is a flag to indicate if the role can create roles
	CreateRole bool `json:"createRole"`
	// ConnectionLimit is the maximum number of concurrent connections for the role. Default is -1, as per PostgreSQL default.
	ConnectionLimit int64 `json:"connectionLimit,omitempty"`
	// ValidUntil is the date and time the role will expire
	ValidUntil *time.Time `json:"validUntil,omitempty"`
	// DeleteScheduledAt is the date and time the role will be deleted
	DeleteScheduledAt *time.Time `json:"deleteScheduledAt,omitempty"`
}

type DbClusterPostgresDatabase struct {
	// Identity is a unique identifier for the database
	Identity string `json:"identity"`
	// Name is a human-readable name of the database
	Name string `json:"name"`

	// Owner is the name of the owner role
	Owner string `json:"owner"`

	// ConnectionLimit is the maximum number of concurrent connections for the database. Default is -1, as per PostgreSQL default.
	ConnectionLimit int64 `json:"connectionLimit,omitempty"`

	// DeleteScheduledAt is the date and time the database will be deleted
	DeleteScheduledAt *time.Time `json:"deleteScheduledAt,omitempty"`

	Status string `json:"status"`
}

type DbClusterScheduledMaintenance struct {
	// Identity is a unique identifier for the scheduled maintenance
	Identity string `json:"identity"`
	// CreatedAt is the date and time the scheduled maintenance was created
	CreatedAt time.Time `json:"createdAt"`
	// ScheduledAt is the date and time the scheduled maintenance was scheduled
	ScheduledAt time.Time `json:"scheduledAt"`
	// StartedAt is the date and time the scheduled maintenance started
	StartedAt *time.Time `json:"startedAt,omitempty"`
	// CompletedAt is the date and time the scheduled maintenance completed
	CompletedAt *time.Time `json:"completedAt,omitempty"`
	// CanceledAt is the date and time the scheduled maintenance was canceled
	CanceledAt *time.Time `json:"canceledAt,omitempty"`
	// FailedAt is the date and time the scheduled maintenance failed
	FailedAt *time.Time `json:"failedAt,omitempty"`
	// DbCluster is the cluster the scheduled maintenance belongs to
	DbCluster *DbCluster `json:"dbCluster"`
	// Status is the status of the scheduled maintenance
	Status DbClusterScheduledMaintenanceStatus `json:"status"`
	// StatusReason is the reason for the status of the scheduled maintenance
	StatusReason string `json:"statusReason,omitempty"`
	// CurrentVersion is the current version of the engine
	CurrentVersion *DbClusterEngineVersion `json:"currentVersion,omitempty"`
	// TargetVersion is the target version of the engine
	TargetVersion *DbClusterEngineVersion `json:"targetVersion,omitempty"`
}

type DbClusterScheduledMaintenanceStatus string

const (
	DbClusterScheduledMaintenanceStatusScheduled  DbClusterScheduledMaintenanceStatus = "scheduled"
	DbClusterScheduledMaintenanceStatusInProgress DbClusterScheduledMaintenanceStatus = "inProgress"
	DbClusterScheduledMaintenanceStatusCompleted  DbClusterScheduledMaintenanceStatus = "completed"
	DbClusterScheduledMaintenanceStatusFailed     DbClusterScheduledMaintenanceStatus = "failed"
	DbClusterScheduledMaintenanceStatusCancelled  DbClusterScheduledMaintenanceStatus = "cancelled"
	DbClusterScheduledMaintenanceStatusSkipped    DbClusterScheduledMaintenanceStatus = "skipped"
)

// DbObjectStore represents an object storage location for database backups.
// When created, it provisions an object storage bucket with correct policies
// and a service account with object storage access credentials.
type DbObjectStore struct {
	// Identity is a unique identifier for the object store
	Identity string `json:"identity"`
	Name     string `json:"name"`
	// Description is the description of the object store
	Description string `json:"description"`
	// CreatedAt is the date and time the object store was created
	CreatedAt time.Time `json:"createdAt"`
	// UpdatedAt is the date and time the object store was last updated
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	// ObjectVersion is the version of the object store
	ObjectVersion int `json:"objectVersion"`
	// Labels is a map of labels for the object store
	Labels Labels `json:"labels,omitempty"`
	// Annotations is a map of annotations for the object store
	Annotations Annotations `json:"annotations,omitempty"`
	// Organisation is the organisation the object store belongs to
	Organisation *base.Organisation `json:"organisation,omitempty"`
	// Region is the region of the object store
	Region *iaas.Region `json:"region,omitempty"`
	// Status is the status of the object store
	Status ObjectStatus `json:"status"`
	// StatusMessage is the message of the object store status
	StatusMessage string `json:"statusMessage,omitempty"`
	// DeleteProtection is a flag to indicate if the object store is protected from deletion
	// This is used to prevent the object store from being deleted accidentally
	// It is set to true when the object store is created and remains true until the object store is deleted
	DeleteProtection bool `json:"deleteProtection"`
	// ObjectStorageBucket is the underlying object storage bucket
	ObjectStorageBucket *objectstorage.ObjectStorageBucket `json:"objectStorageBucket,omitempty"`
	// ServiceAccount is the service account with object storage access credentials
	ServiceAccount *iam.ServiceAccount `json:"serviceAccount,omitempty"`
	// RetentionPolicy is the retention policy for backups in the format "<number>d" where d is days.
	// For example, "30d" means backups will be retained for 30 days.
	// This is used with barman-cloud-backup-delete command: --retention-policy "RECOVERY WINDOW OF <number> days"
	RetentionPolicy string `json:"retentionPolicy,omitempty"`
}

// CreateDbObjectStoreRequest is the request body for creating a DB object store.
type CreateDbObjectStoreRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	// Annotations is a map of key-value pairs used for storing additional information
	Annotations Annotations `json:"annotations,omitempty"`
	// Labels is a map of key-value pairs used for filtering and grouping objects
	Labels Labels `json:"labels,omitempty"`
	// Region is the identity or slug of the cloud region where the object store will be created
	Region string `json:"region"`
	// RetentionPolicy is the retention policy for backups in the format "<number>d" where d is days.
	// For example, "30d" means backups will be retained for 30 days.
	RetentionPolicy string `json:"retentionPolicy,omitempty"`
	// DeleteProtection is a flag to indicate if the object store is protected from deletion
	DeleteProtection bool `json:"deleteProtection"`
}

// UpdateDbObjectStoreRequest is the request body for updating a DB object store.
type UpdateDbObjectStoreRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	// Annotations is a map of key-value pairs used for storing additional information
	Annotations Annotations `json:"annotations,omitempty"`
	// Labels is a map of key-value pairs used for filtering and grouping objects
	Labels Labels `json:"labels,omitempty"`
	// RetentionPolicy is the retention policy for backups in the format "<number>d" where d is days.
	RetentionPolicy string `json:"retentionPolicy,omitempty"`
	// DeleteProtection is a flag to indicate if the object store is protected from deletion
	DeleteProtection bool `json:"deleteProtection"`
}
