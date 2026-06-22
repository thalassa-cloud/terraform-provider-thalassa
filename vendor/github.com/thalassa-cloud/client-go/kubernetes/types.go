package kubernetes

import (
	"time"

	"github.com/thalassa-cloud/client-go/iaas"
	"github.com/thalassa-cloud/client-go/pkg/base"
)

// KubernetesClusterSessionToken represents the authentication and connection details for a Kubernetes cluster session.
type KubernetesClusterSessionToken struct {
	// Username for cluster authentication
	Username string `json:"username"`
	// URL of the Kubernetes API server
	APIServerURL string `json:"apiServerUrl"`
	// CA certificate for API server verification
	CACertificate string `json:"caCertificate"`
	// Unique identifier for the session
	Identity string `json:"identity"`
	// Authentication token
	Token string `json:"token"`
	// Complete kubeconfig file content
	Kubeconfig string `json:"kubeconfig"`
}

// KubernetesVersion represents a supported Kubernetes version configuration.
type KubernetesVersion struct {
	// Unique identifier for the version
	Identity string `json:"identity"`
	// Display name of the version
	Name string `json:"name"`
	// URL-friendly identifier
	Slug string `json:"slug"`
	// Detailed description of the version
	Description string `json:"description"`
	// Creation timestamp
	CreatedAt time.Time `json:"createdAt"`
	// Custom annotations
	Annotations map[string]string `json:"annotations"`
	// Whether this version is available for use
	Enabled bool `json:"enabled"`
	// Supported indicates if this version is currently supported / stable. This flag is used for auto upgrade policies and indicates stable versions.
	// After new versions are released, the old versions are marked as unsupported and cannot be used for auto upgrade policies.
	// You can still use them for manual upgrades should they be enabled, however it is recommend to upgrade to the next available supported version.
	Supported bool `json:"supported"`

	// Core Kubernetes version
	KubernetesVersion string `json:"kubernetesVersion"`
	// Container runtime version
	ContainerdVersion string `json:"containerdVersion"`
	// CNI plugins version
	CNIPluginsVersion string `json:"cniPluginsVersion"`
	// CRI tools version
	CrictlVersion string `json:"crictlVersion"`
	// Container runtime spec version
	RuncVersion string `json:"runcVersion"`
	// Cilium CNI version
	CiliumVersion string `json:"ciliumVersion"`
	// Cloud controller manager version
	CloudControllerManagerVersion string `json:"cloudControllerManagerVersion"`
	// Istio service mesh version
	IstioVersion string `json:"istioVersion"`
}

// KubernetesCluster represents a Kubernetes cluster in the Thalassa Cloud Platform.
type KubernetesCluster struct {
	Identity                 string                         `json:"identity"`                  // Unique identifier for the cluster
	Name                     string                         `json:"name"`                      // Display name of the cluster
	Slug                     string                         `json:"slug"`                      // URL-friendly identifier
	Description              string                         `json:"description"`               // Detailed description of the cluster
	Labels                   map[string]string              `json:"labels"`                    // Custom labels
	Annotations              map[string]string              `json:"annotations"`               // Custom annotations
	CreatedAt                time.Time                      `json:"createdAt"`                 // Creation timestamp
	ObjectVersion            int                            `json:"objectVersion"`             // Version for optimistic locking
	Organisation             *base.Organisation             `json:"organisation"`              // Associated organization
	Status                   string                         `json:"status"`                    // Current cluster status
	StatusMessage            string                         `json:"statusMessage"`             // Detailed status message
	LastStatusTransitionedAt time.Time                      `json:"lastStatusTransitioned_at"` // Last status change timestamp
	ClusterType              KubernetesClusterType          `json:"clusterType"`               // Type of cluster deployment
	ClusterVersion           KubernetesVersion              `json:"clusterVersion"`            // Kubernetes version configuration
	APIServerURL             string                         `json:"apiServerURL"`              // Kubernetes API server URL
	APIServerCA              string                         `json:"apiServerCA"`               // API server CA certificate
	Configuration            KubernetesClusterConfiguration `json:"configuration"`             // Cluster configuration

	VPC    *iaas.Vpc    `json:"vpc"`    // Associated VPC (not set for hosted-control-plane)
	Subnet *iaas.Subnet `json:"subnet"` // Associated subnet (not set for hosted-control-plane)
	Region *iaas.Region `json:"region"` // Associated region

	PodSecurityStandardsProfile KubernetesClusterPodSecurityStandards `json:"podSecurityStandardsProfile"` // Pod security standards configuration
	AuditLogProfile             KubernetesClusterAuditLoggingProfile  `json:"auditLogProfile"`             // Audit logging configuration
	DefaultNetworkPolicy        KubernetesDefaultNetworkPolicies      `json:"defaultNetworkPolicy"`        // Default network policy
	DeleteProtection            bool                                  `json:"deleteProtection"`            // Whether deletion protection is enabled

	ApiServerACLs         KubernetesApiServerACLs            `json:"apiServerACL"`                 // ApiServerACLs is the ACLs for the API server
	DisablePublicEndpoint bool                               `json:"disablePublicEndpoint"`        // Whether public endpoint is disabled
	AutoUpgradePolicy     KubernetesClusterAutoUpgradePolicy `json:"autoUpgradePolicy"`            // AutoUpgradePolicy is the auto upgrade policy for the cluster
	MaintenanceDay        *uint                              `json:"maintenanceDay,omitempty"`     // MaintenanceDay is the day of the week when the cluster will be upgraded. Optional.
	MaintenanceStartAt    *uint                              `json:"maintenanceStartAt,omitempty"` // MaintenanceStartAt is the time of day when the cluster will be upgraded. Optional.

	// ScheduledMaintenances is the list of scheduled maintenances for the cluster
	ScheduledMaintenances []KubernetesClusterScheduledMaintenance `json:"scheduledMaintenances,omitempty"`

	// AutoscalerConfig is the configuration for the cluster autoscaler
	// These values can also be configured using annotations on a KubernetesNodePool object
	// cluster-autoscaler.kubernetes.io/<setting-name>
	// For more information, see the Cluster Autoscaler documentation: https://github.com/kubernetes/autoscaler/blob/master/cluster-autoscaler/FAQ.md
	AutoscalerConfig *AutoscalerConfig `json:"autoscalerConfig,omitempty"`

	InternalEndpoint *string `json:"internalEndpoint,omitempty"` // VPC-internal endpoint for the cluster
	AdvertisePort    *int    `json:"advertisePort,omitempty"`    // Advertise port for the cluster within the VPC
	KonnectivityPort *int    `json:"konnectivityPort,omitempty"` // Konnectivity port for the cluster within the VPC

	// SecurityGroups is a list of security groups that are attached to the Control Plane.
	SecurityGroups []iaas.SecurityGroup `json:"securityGroups,omitempty"`
}

// CreateKubernetesCluster represents the configuration for creating a new Kubernetes cluster.
type CreateKubernetesCluster struct {
	Name                        string                                `json:"name"`                        // Display name for the new cluster
	Description                 string                                `json:"description"`                 // Cluster description
	Labels                      map[string]string                     `json:"labels"`                      // Custom labels
	Annotations                 map[string]string                     `json:"annotations"`                 // Custom annotations
	RegionIdentity              string                                `json:"regionIdentity"`              // Target region identifier
	ClusterType                 KubernetesClusterType                 `json:"clusterType"`                 // Type of cluster deployment
	KubernetesVersionIdentity   string                                `json:"kubernetesVersionIdentity"`   // Kubernetes version identifier
	DeleteProtection            bool                                  `json:"deleteProtection"`            // Whether deletion protection is enabled
	Subnet                      string                                `json:"subnet"`                      // Subnet for cluster deployment
	Networking                  KubernetesClusterNetworking           `json:"networking"`                  // Network configuration
	PodSecurityStandardsProfile KubernetesClusterPodSecurityStandards `json:"podSecurityStandardsProfile"` // Pod security standards
	AuditLogProfile             KubernetesClusterAuditLoggingProfile  `json:"auditLogProfile"`             // Audit logging configuration
	DefaultNetworkPolicy        KubernetesDefaultNetworkPolicies      `json:"defaultNetworkPolicy"`        // Default network policy
	DisablePublicEndpoint       bool                                  `json:"disablePublicEndpoint"`       // Whether public endpoint is disabled
	ApiServerACLs               KubernetesApiServerACLs               `json:"apiServerACL"`                // ApiServerACLs is the ACLs for the API server

	// KubeProxyMode is the mode of the kube proxy. Default is ipvs.
	KubeProxyMode *KubernetesClusterKubeProxyMode `json:"kubeProxyMode,omitempty"`
	// KubeProxyDeployment is the deployment mode of the kube proxy. Default is managed.
	KubeProxyDeployment *KubeProxyDeployment `json:"kubeProxyDeployment,omitempty"`

	AutoUpgradePolicy  KubernetesClusterAutoUpgradePolicy `json:"autoUpgradePolicy"`            // AutoUpgradePolicy is the auto upgrade policy for the cluster
	MaintenanceDay     *uint                              `json:"maintenanceDay,omitempty"`     // MaintenanceDay is the day of the week when the cluster will be upgraded. Optional.
	MaintenanceStartAt *uint                              `json:"maintenanceStartAt,omitempty"` // MaintenanceStartAt is the time of day when the cluster will be upgraded. Optional.

	// AutoscalerConfig is the configuration for the cluster autoscaler
	// These values can also be configured using annotations on a KubernetesNodePool object
	// cluster-autoscaler.kubernetes.io/<setting-name>
	// For more information, see the Cluster Autoscaler documentation: https://github.com/kubernetes/autoscaler/blob/master/cluster-autoscaler/FAQ.md
	AutoscalerConfig *AutoscalerConfig `json:"autoscalerConfig,omitempty"`

	// SecurityGroupAttachments is a list of security group identities that will be attached to the Control Plane VPC-internal endpoint.
	// These do not apply to the public endpoint. If you wish to configure ACLs for the public endpoint, you can use the AllowedCIDRs field.
	SecurityGroupAttachments []string `json:"securityGroupAttachments,omitempty"`
}

type KubernetesApiServerACLs struct {
	// AllowedCIDRs is a list of allowed CIDRs. Either a CIDR or an IP address.
	// These CIDRs will be allowed to access the API server on the public endpoint. These ACLs are not applied to the VPC-internal endpoint.
	// If you wish to configure ACLs for the VPC-internal endpoint, you can use the SecurityGroupAttachments field.
	AllowedCIDRs []string `json:"allowedCIDRs"`
}

type KubernetesClusterAutoUpgradePolicy string

const (
	// KubernetesClusterAutoUpgradePolicyNone does not perform any auto upgrades. User is expected to manually upgrade the cluster.
	KubernetesClusterAutoUpgradePolicyNone KubernetesClusterAutoUpgradePolicy = "none"
	// KubernetesClusterAutoUpgradePolicyLatestVersion is the auto upgrade policy for the cluster.
	// It will upgrade to the latest release of the latest supported minor version.
	// This upgrade strategy is recommended for development clusters.
	KubernetesClusterAutoUpgradePolicyLatestVersion KubernetesClusterAutoUpgradePolicy = "latest-version"
	// KubernetesClusterAutoUpgradePolicyLatestStable is the auto upgrade policy for the cluster.
	// It will upgrade to the latest stable version of the current minor version.
	// Once the current minor version becomes unsupported and no stable version is available, it will perform a minor upgrade.
	// This upgrade strategy is recommended for production clusters.
	KubernetesClusterAutoUpgradePolicyLatestStable KubernetesClusterAutoUpgradePolicy = "latest-stable"
)

// UpdateKubernetesCluster represents the configuration for updating an existing Kubernetes cluster.
type UpdateKubernetesCluster struct {
	Name                      *string           `json:"name,omitempty"`                      // New display name
	Description               *string           `json:"description,omitempty"`               // New description
	Labels                    map[string]string `json:"labels,omitempty"`                    // Updated labels
	Annotations               map[string]string `json:"annotations,omitempty"`               // Updated annotations
	KubernetesVersionIdentity *string           `json:"kubernetesVersionIdentity,omitempty"` // New Kubernetes version identifier
	DeleteProtection          *bool             `json:"deleteProtection,omitempty"`          // Updated deletion protection setting

	DefaultNetworkPolicy        *KubernetesDefaultNetworkPolicies      `json:"defaultNetworkPolicy,omitempty"`        // Updated default network policy
	PodSecurityStandardsProfile *KubernetesClusterPodSecurityStandards `json:"podSecurityStandardsProfile,omitempty"` // Updated pod security standards
	AuditLogProfile             *KubernetesClusterAuditLoggingProfile  `json:"auditLogProfile,omitempty"`             // Updated audit logging configuration
	DisablePublicEndpoint       *bool                                  `json:"disablePublicEndpoint,omitempty"`       // Updated public endpoint setting
	ApiServerACLs               KubernetesApiServerACLs                `json:"apiServerACL"`                          // ApiServerACLs is the ACLs for the API server

	// KubeProxyMode is the mode of the kube proxy. Default is ipvs.
	KubeProxyMode *KubernetesClusterKubeProxyMode `json:"kubeProxyMode,omitempty"`
	// KubeProxyDeployment is the deployment mode of the kube proxy. Default is managed.
	KubeProxyDeployment *KubeProxyDeployment `json:"kubeProxyDeployment,omitempty"`

	AutoUpgradePolicy  KubernetesClusterAutoUpgradePolicy `json:"autoUpgradePolicy"`            // AutoUpgradePolicy is the auto upgrade policy for the cluster
	MaintenanceDay     *uint                              `json:"maintenanceDay,omitempty"`     // MaintenanceDay is the day of the week when the cluster will be upgraded. Optional.
	MaintenanceStartAt *uint                              `json:"maintenanceStartAt,omitempty"` // MaintenanceStartAt is the time of day when the cluster will be upgraded. Optional.

	// AutoscalerConfig is the configuration for the cluster autoscaler
	// These values can also be configured using annotations on a KubernetesNodePool object
	// cluster-autoscaler.kubernetes.io/<setting-name>
	// For more information, see the Cluster Autoscaler documentation: https://github.com/kubernetes/autoscaler/blob/master/cluster-autoscaler/FAQ.md
	AutoscalerConfig *AutoscalerConfig `json:"autoscalerConfig,omitempty"`

	// SecurityGroupAttachments is a list of security group identities that will be attached to the Control Plane VPC-internal endpoint.
	// These do not apply to the public endpoint. If you wish to configure ACLs for the public endpoint, you can use the AllowedCIDRs field.
	SecurityGroupAttachments []string `json:"securityGroupAttachments,omitempty"`
}

// KubernetesClusterNetworking represents the network configuration for a Kubernetes cluster.
type KubernetesClusterNetworking struct {
	CNI string `json:"cni"` // CNI, default is cilium.

	// KubeProxyMode is the mode of the kube proxy. Default is ipvs.
	KubeProxyMode *KubernetesClusterKubeProxyMode `json:"kubeProxyMode,omitempty"`
	// KubeProxyDeployment is the deployment mode of the kube proxy. Default is managed.
	KubeProxyDeployment *KubeProxyDeployment `json:"kubeProxyDeployment,omitempty"`

	// ServiceCidr is the service CIDR for the cluster. Must be a valid CIDR block. Must not overlap with the pod CIDR or the VPC / Subnet CIDRs.
	ServiceCIDR string `json:"serviceCidr"`
	// PodCidr is the pod CIDR for the cluster. Must be a valid CIDR block. Must not overlap with the service CIDR or the VPC / Subnet CIDRs.
	PodCIDR string `json:"podCidr"`
}

// KubernetesClusterType represents the type of Kubernetes cluster deployment.
type KubernetesClusterType string

const (
	Managed            KubernetesClusterType = "managed"              // Fully managed cluster
	HostedControlPlane KubernetesClusterType = "hosted-control-plane" // Cluster with hosted control plane
)

// KubernetesClusterConfiguration represents the configuration of a Kubernetes cluster.
type KubernetesClusterConfiguration struct {
	Networking KubernetesClusterNetworking `json:"networking"` // Network configuration
}

// KubernetesClusterPodSecurityStandards represents the pod security standards profile for a cluster.
type KubernetesClusterPodSecurityStandards string

const (
	KubernetesClusterPodSecurityStandardRestricted KubernetesClusterPodSecurityStandards = "restricted" // Most restrictive security profile
	KubernetesClusterPodSecurityStandardBaseline   KubernetesClusterPodSecurityStandards = "baseline"   // Standard security profile
	KubernetesClusterPodSecurityStandardPrivileged KubernetesClusterPodSecurityStandards = "privileged" // Least restrictive security profile
)

// KubernetesClusterAuditLoggingProfile represents the audit logging configuration for a cluster.
type KubernetesClusterAuditLoggingProfile string

const (
	KubernetesClusterAuditLoggingProfileNone     KubernetesClusterAuditLoggingProfile = "none"     // No audit logging
	KubernetesClusterAuditLoggingProfileBasic    KubernetesClusterAuditLoggingProfile = "basic"    // Basic audit logging
	KubernetesClusterAuditLoggingProfileAdvanced KubernetesClusterAuditLoggingProfile = "advanced" // Advanced audit logging
)

type KubeProxyDeployment string

const (
	KubeProxyDeploymentCustom   KubeProxyDeployment = "custom"
	KubeProxyDeploymentManaged  KubeProxyDeployment = "managed"
	KubeProxyDeploymentDisabled KubeProxyDeployment = "disabled"
)

type KubernetesClusterKubeProxyMode string

const (
	KubernetesClusterKubeProxyModeIPVS     KubernetesClusterKubeProxyMode = "ipvs"
	KubernetesClusterKubeProxyModeIptables KubernetesClusterKubeProxyMode = "iptables"
)

// KubernetesDefaultNetworkPolicies represents the default network policy for a cluster.
type KubernetesDefaultNetworkPolicies string

const (
	KubernetesDefaultNetworkPolicyNone     KubernetesDefaultNetworkPolicies = ""          // No default policy
	KubernetesDefaultNetworkPolicyAllowAll KubernetesDefaultNetworkPolicies = "allow-all" // Allow all traffic
	KubernetesDefaultNetworkPolicyDenyAll  KubernetesDefaultNetworkPolicies = "deny-all"  // Deny all traffic
)

// KubernetesNodePool represents a group of nodes in a Kubernetes cluster with identical configuration.
type KubernetesNodePool struct {
	Identity         string            `json:"identity"`         // Unique identifier for the node pool
	Name             string            `json:"name"`             // Display name of the node pool
	Slug             string            `json:"slug"`             // URL-friendly identifier
	Description      string            `json:"description"`      // Detailed description
	CreatedAt        time.Time         `json:"createdAt"`        // Creation timestamp
	UpdatedAt        *time.Time        `json:"updatedAt"`        // Last update timestamp
	ObjectVersion    int               `json:"objectVersion"`    // Version for optimistic locking
	Labels           map[string]string `json:"labels"`           // Custom labels
	Annotations      map[string]string `json:"annotations"`      // Custom annotations
	AvailabilityZone string            `json:"availabilityZone"` // Availability zone for the node pool

	Status KubernetesNodePoolStatus `json:"status"` // Current status of the node pool

	Vpc    *iaas.Vpc    `json:"vpc"`    // Associated VPC
	Subnet *iaas.Subnet `json:"subnet"` // Associated subnet

	EnableAutoscaling     bool `json:"enableAutoscaling"`     // Whether autoscaling is enabled
	EnableAutoHealing     bool `json:"enableAutoHealing"`     // Whether auto-healing is enabled
	Replicas              int  `json:"replicas"`              // Current number of nodes
	MinReplicas           int  `json:"minReplicas"`           // Minimum number of nodes for autoscaling
	MaxReplicas           int  `json:"maxReplicas"`           // Maximum number of nodes for autoscaling
	ManageNodeAllocatable bool `json:"manageNodeAllocatable"` // ManageNodeAllocatable is a flag to manage the node allocatable resources.

	MachineType       iaas.MachineType                  `json:"machineType"`       // Type of machine for nodes
	NodeSettings      KubernetesNodeSettings            `json:"nodeSettings"`      // Node-specific settings
	KubernetesVersion *KubernetesVersion                `json:"kubernetesVersion"` // Kubernetes version for node pool
	UpgradeStrategy   KubernetesNodePoolUpgradeStrategy `json:"upgradeStrategy"`   // Upgrade strategy for node pool

	// SecurityGroups is a list of security groups that are attached to the nodes / vmi in the node pool.
	SecurityGroups []iaas.SecurityGroup `json:"securityGroups,omitempty"`
}

// CreateKubernetesNodePool represents the configuration for creating a new node pool.
type CreateKubernetesNodePool struct {
	Name        string            `json:"name"`        // Display name for the node pool
	Description string            `json:"description"` // Detailed description
	Labels      map[string]string `json:"labels"`      // Custom labels
	Annotations map[string]string `json:"annotations"` // Custom annotations

	MachineType      string  `json:"machineType"`      // Type of machine for nodes
	Replicas         int     `json:"replicas"`         // Initial number of nodes
	MinReplicas      int     `json:"minReplicas"`      // Minimum nodes for autoscaling
	MaxReplicas      int     `json:"maxReplicas"`      // Maximum nodes for autoscaling
	SubnetIdentity   *string `json:"subnetIdentity"`   // Subnet for node pool deployment
	AvailabilityZone string  `json:"availabilityZone"` // Availability zone for the node pool

	KubernetesVersionIdentity *string                            `json:"kubernetesVersionIdentity"` // Kubernetes version for node pool
	UpgradeStrategy           *KubernetesNodePoolUpgradeStrategy `json:"upgradeStrategy"`           // Upgrade strategy for node pool
	EnableAutoscaling         bool                               `json:"enableAutoscaling"`         // Whether to enable autoscaling
	EnableAutoHealing         bool                               `json:"enableAutoHealing"`         // Whether auto-healing is enabled
	NodeSettings              KubernetesNodeSettings             `json:"nodeSettings"`              // Node-specific settings
	ManageNodeAllocatable     bool                               `json:"manageNodeAllocatable"`     // ManageNodeAllocatable is a flag to manage the node allocatable resources.

	// SecurityGroupAttachments is a list of security group identities that will be attached to the nodes / vmi in the node pool.
	SecurityGroupAttachments []string `json:"securityGroupAttachments,omitempty"`
}

// UpdateKubernetesNodePool represents the configuration for updating an existing node pool.
type UpdateKubernetesNodePool struct {
	Description string            `json:"description"` // New description
	Labels      map[string]string `json:"labels"`      // Custom labels
	Annotations map[string]string `json:"annotations"` // Custom annotations

	MachineType               string  `json:"machineType"`               // New machine type
	Replicas                  *int    `json:"replicas"`                  // New number of nodes
	MinReplicas               *int    `json:"minReplicas"`               // New minimum nodes for autoscaling
	MaxReplicas               *int    `json:"maxReplicas"`               // New maximum nodes for autoscaling
	KubernetesVersionIdentity *string `json:"kubernetesVersionIdentity"` // Kubernetes version for node pool
	AvailabilityZone          string  `json:"availabilityZone"`          // Availability zone for the node pool

	UpgradeStrategy       *KubernetesNodePoolUpgradeStrategy `json:"upgradeStrategy"`       // Upgrade strategy for node pool
	EnableAutoHealing     *bool                              `json:"enableAutoHealing"`     // Whether auto-healing is enabled
	EnableAutoscaling     *bool                              `json:"enableAutoscaling"`     // Updated autoscaling setting
	ManageNodeAllocatable bool                               `json:"manageNodeAllocatable"` // ManageNodeAllocatable is a flag to manage the node allocatable resources.

	NodeSettings *KubernetesNodeSettings `json:"nodeSettings"` // Updated node settings

	// SecurityGroupAttachments is a list of security group identities that will be attached to the nodes / vmi in the node pool.
	SecurityGroupAttachments []string `json:"securityGroupAttachments,omitempty"`
}

// KubernetesNodeSettings represents the configuration settings for nodes in a node pool.
type KubernetesNodeSettings struct {
	Annotations map[string]string `json:"nodeAnnotations"` // Kubernetes node annotations
	Labels      map[string]string `json:"nodeLabels"`      // Kubernetes node labels
	Taints      []NodeTaint       `json:"nodeTaints"`      // Node taints for pod scheduling
}

// NodeTaint represents a taint that can be applied to nodes to control pod scheduling.
type NodeTaint struct {
	Key      string `json:"key"`      // Taint key
	Value    string `json:"value"`    // Taint value
	Operator string `json:"operator"` // Taint operator (Equal, Exists)
	Effect   string `json:"effect"`   // Taint effect (NoSchedule, PreferNoSchedule, NoExecute)
}

type KubernetesNodePoolUpgradeStrategy string

const (
	KubernetesNodePoolUpgradeStrategyManual    KubernetesNodePoolUpgradeStrategy = "manual"
	KubernetesNodePoolUpgradeStrategyAuto      KubernetesNodePoolUpgradeStrategy = "auto"
	KubernetesNodePoolUpgradeStrategyMinorOnly KubernetesNodePoolUpgradeStrategy = "minor-only"

	// Backward compatibility
	KubernetesNodePoolUpgradeStrategyAlways   KubernetesNodePoolUpgradeStrategy = "always"
	KubernetesNodePoolUpgradeStrategyOnDelete KubernetesNodePoolUpgradeStrategy = "onDelete"
	KubernetesNodePoolUpgradeStrategyInplace  KubernetesNodePoolUpgradeStrategy = "inPlace"
	KubernetesNodePoolUpgradeStrategyNever    KubernetesNodePoolUpgradeStrategy = "never"
)

type KubernetesNodePoolStatus string

const (
	KubernetesNodePoolStatusProvisioning KubernetesNodePoolStatus = "provisioning"
	KubernetesNodePoolStatusReady        KubernetesNodePoolStatus = "ready"
	KubernetesNodePoolStatusFailed       KubernetesNodePoolStatus = "failed"
	KubernetesNodePoolStatusUpdating     KubernetesNodePoolStatus = "updating"
	KubernetesNodePoolStatusDeleting     KubernetesNodePoolStatus = "deleting"
	KubernetesNodePoolStatusDeleted      KubernetesNodePoolStatus = "deleted"
)

type KubernetesNodePoolMachine struct {
	// CreatedAt is the timestamp when the object was created
	CreatedAt time.Time `json:"createdAt"`
	// UpdatedAt is the timestamp when the object was last updated
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	// Identity is a unique identifier for the Kubernetes Node Pool Machine
	Identity string `json:"identity"`
	// MachineName is the name of the machine
	MachineName string `json:"machineName"`

	Vpc    *iaas.Vpc    `json:"vpc,omitempty"`
	Subnet *iaas.Subnet `json:"subnet,omitempty"`

	// Conditions is a list of conditions for the Kubernetes Node Pool Machine
	Conditions []Condition `json:"conditions,omitempty"`
	// SystemInfo is the system information for the Kubernetes Node Pool Machine
	SystemInfo NodeSystemInfo `json:"systemInfo"`
}

type Condition struct {
	Type               string    `json:"type"`
	Status             string    `json:"status"`
	LastTransitionTime time.Time `json:"lastTransitionTime"`
	Reason             string    `json:"reason"`
	Message            string    `json:"message"`
}

type NodeSystemInfo struct {
	// Architecture is the architecture of the node
	Architecture string `json:"architecture,omitempty"`
	// BootID is the boot ID of the node
	BootID string `json:"bootID,omitempty"`
	// KernelVersion is the kernel version of the node
	KernelVersion string `json:"kernelVersion,omitempty"`
	// OsImage is the operating system image of the node
	OsImage string `json:"osImage,omitempty"`
	// ContainerRuntimeVersion is the container runtime version of the node
	ContainerRuntimeVersion string `json:"containerRuntimeVersion,omitempty"`
	// KubeletVersion is the kubelet version of the node
	KubeletVersion string `json:"kubeletVersion,omitempty"`
	// Addresses is a list of addresses for the node
	Addresses []NodeAddress `json:"addresses,omitempty"`
	// Conditions is a list of conditions for the node
	Conditions []NodeCondition `json:"conditions,omitempty"`
}

type NodeCondition struct {
	// Type of node condition.
	Type string `json:"type"`
	// Status of the condition, one of True, False, Unknown.
	Status string `json:"status"`
	// Last time we got an update on a given condition.
	LastHeartbeatTime  time.Time `json:"lastHeartbeatTime,omitempty"`
	LastTransitionTime time.Time `json:"lastTransitionTime,omitempty"`
	Reason             string    `json:"reason,omitempty"`
	Message            string    `json:"message,omitempty"`
}

type NodeAddress struct {
	// Node address type, one of Hostname, ExternalIP or InternalIP.
	Type string `json:"type"`
	// The node address.
	Address string `json:"address"`
}

type KubernetesClusterScheduledMaintenance struct {
	Identity    string     `json:"identity"`
	CreatedAt   time.Time  `json:"createdAt"`
	ScheduledAt time.Time  `json:"scheduledAt"`
	StartedAt   *time.Time `json:"startedAt,omitempty"`
	CompletedAt *time.Time `json:"completedAt,omitempty"`
	CanceledAt  *time.Time `json:"canceledAt,omitempty"`
	FailedAt    *time.Time `json:"failedAt,omitempty"`

	Status       KubernetesClusterScheduledMaintenanceStatus `json:"status"`
	StatusReason string                                      `json:"statusReason,omitempty"`

	CurrentVersion *KubernetesVersion `json:"currentVersion,omitempty"`
	TargetVersion  *KubernetesVersion `json:"targetVersion,omitempty"`
}

type KubernetesClusterScheduledMaintenanceStatus string

const (
	KubernetesClusterScheduledMaintenanceStatusScheduled  KubernetesClusterScheduledMaintenanceStatus = "scheduled"
	KubernetesClusterScheduledMaintenanceStatusInProgress KubernetesClusterScheduledMaintenanceStatus = "inProgress"
	KubernetesClusterScheduledMaintenanceStatusCompleted  KubernetesClusterScheduledMaintenanceStatus = "completed"
	KubernetesClusterScheduledMaintenanceStatusFailed     KubernetesClusterScheduledMaintenanceStatus = "failed"
	KubernetesClusterScheduledMaintenanceStatusCancelled  KubernetesClusterScheduledMaintenanceStatus = "cancelled"
	KubernetesClusterScheduledMaintenanceStatusSkipped    KubernetesClusterScheduledMaintenanceStatus = "skipped"
)

// AutoscalerConfig is the configuration for the cluster autoscaler
// These values can also be configured using annotations on a KubernetesNodePool object
// cluster-autoscaler.kubernetes.io/<setting-name>
// For more information, see the Cluster Autoscaler documentation: https://github.com/kubernetes/autoscaler/blob/master/cluster-autoscaler/FAQ.md
type AutoscalerConfig struct {
	// ScaleDownDisabled is a flag to disable the scale down of node pools by the cluster autoscaler
	ScaleDownDisabled bool `json:"scaleDownDisabled"`
	// ScaleDownDelayAfterAdd is the delay after adding a node to the node pool by the cluster autoscaler
	ScaleDownDelayAfterAdd string `json:"scaleDownDelayAfterAdd"`
	// Estimator is the estimator to use for the cluster autoscaler. Available values: binpacking
	Estimator string `json:"estimator"`

	// Expander is the expander to use for the cluster autoscaler
	Expander string `json:"expander"`
	// IgnoreDaemonsetsUtilization is a flag to ignore the utilization of daemonsets by the cluster autoscaler
	IgnoreDaemonsetsUtilization bool `json:"ignoreDaemonsetsUtilization"`
	// BalanceSimilarNodeGroups is a flag to balance the utilization of similar node groups by the cluster autoscaler
	BalanceSimilarNodeGroups bool `json:"balanceSimilarNodeGroups"`
	// ExpendablePodsPriorityCutoff is the priority cutoff for the expendable pods by the cluster autoscaler
	ExpendablePodsPriorityCutoff int `json:"expendablePodsPriorityCutoff"`
	// ScaleDownUnneededTime is the time after which a node can be scaled down by the cluster autoscaler
	ScaleDownUnneededTime string `json:"scaleDownUnneededTime"`
	// ScaleDownUtilizationThreshold is the utilization threshold for the cluster autoscaler
	// The autoscaler might scale down non-empty nodes with utilization below a threshold. To prevent this behavior, set the utilization threshold to 0
	ScaleDownUtilizationThreshold float64 `json:"scaleDownUtilizationThreshold"`
	// MaxGracefulTerminationSec is the maximum graceful termination time for the cluster autoscaler.
	// If the pod is not stopped within these 10 min then the node is terminated anyway. Earlier versions of CA gave 1 minute or didn't respect graceful termination at all.
	MaxGracefulTerminationSec int `json:"maxGracefulTerminationSec"`

	// EnableProactiveScaleUp is a flag to enable the proactive scale up of the cluster autoscaler.
	// Whether to enable/disable proactive scale-ups, defaults to false
	EnableProactiveScaleUp bool `json:"enableProactiveScaleUp"`
}
