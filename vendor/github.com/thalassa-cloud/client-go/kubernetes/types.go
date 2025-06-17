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
}

// UpdateKubernetesCluster represents the configuration for updating an existing Kubernetes cluster.
type UpdateKubernetesCluster struct {
	Name                        *string                                `json:"name,omitempty"`                        // New display name
	Description                 *string                                `json:"description,omitempty"`                 // New description
	Labels                      map[string]string                      `json:"labels,omitempty"`                      // Updated labels
	Annotations                 map[string]string                      `json:"annotations,omitempty"`                 // Updated annotations
	KubernetesVersionIdentity   *string                                `json:"kubernetesVersionIdentity,omitempty"`   // New Kubernetes version identifier
	DeleteProtection            *bool                                  `json:"deleteProtection,omitempty"`            // Updated deletion protection setting
	Subnet                      *string                                `json:"subnet,omitempty"`                      // New subnet
	Networking                  *KubernetesClusterNetworking           `json:"networking,omitempty"`                  // Updated network configuration
	PodSecurityStandardsProfile *KubernetesClusterPodSecurityStandards `json:"podSecurityStandardsProfile,omitempty"` // Updated pod security standards
	AuditLogProfile             *KubernetesClusterAuditLoggingProfile  `json:"auditLogProfile,omitempty"`             // Updated audit logging configuration
	DefaultNetworkPolicy        *KubernetesDefaultNetworkPolicies      `json:"defaultNetworkPolicy,omitempty"`        // Updated default network policy
}

// KubernetesClusterNetworking represents the network configuration for a Kubernetes cluster.
type KubernetesClusterNetworking struct {
	CNI         string `json:"cni"`         // Container Network Interface type
	ServiceCIDR string `json:"serviceCIDR"` // CIDR range for Kubernetes services
	PodCIDR     string `json:"podCIDR"`     // CIDR range for Kubernetes pods
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

	EnableAutoscaling bool `json:"enableAutoscaling"` // Whether autoscaling is enabled
	EnableAutoHealing bool `json:"enableAutoHealing"` // Whether auto-healing is enabled
	Replicas          int  `json:"replicas"`          // Current number of nodes
	MinReplicas       int  `json:"minReplicas"`       // Minimum number of nodes for autoscaling
	MaxReplicas       int  `json:"maxReplicas"`       // Maximum number of nodes for autoscaling

	MachineType       iaas.MachineType                  `json:"machineType"`       // Type of machine for nodes
	NodeSettings      KubernetesNodeSettings            `json:"nodeSettings"`      // Node-specific settings
	KubernetesVersion *KubernetesVersion                `json:"kubernetesVersion"` // Kubernetes version for node pool
	UpgradeStrategy   KubernetesNodePoolUpgradeStrategy `json:"upgradeStrategy"`   // Upgrade strategy for node pool
}

// CreateKubernetesNodePool represents the configuration for creating a new node pool.
type CreateKubernetesNodePool struct {
	Name        string `json:"name"`        // Display name for the node pool
	Description string `json:"description"` // Detailed description

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
}

// UpdateKubernetesNodePool represents the configuration for updating an existing node pool.
type UpdateKubernetesNodePool struct {
	Description string `json:"description"` // New description

	MachineType               string  `json:"machineType"`               // New machine type
	Replicas                  *int    `json:"replicas"`                  // New number of nodes
	MinReplicas               *int    `json:"minReplicas"`               // New minimum nodes for autoscaling
	MaxReplicas               *int    `json:"maxReplicas"`               // New maximum nodes for autoscaling
	KubernetesVersionIdentity *string `json:"kubernetesVersionIdentity"` // Kubernetes version for node pool
	AvailabilityZone          string  `json:"availabilityZone"`          // Availability zone for the node pool

	UpgradeStrategy   *KubernetesNodePoolUpgradeStrategy `json:"upgradeStrategy"`   // Upgrade strategy for node pool
	EnableAutoHealing *bool                              `json:"enableAutoHealing"` // Whether auto-healing is enabled
	EnableAutoscaling *bool                              `json:"enableAutoscaling"` // Updated autoscaling setting

	NodeSettings *KubernetesNodeSettings `json:"nodeSettings"` // Updated node settings
}

// KubernetesNodeSettings represents the configuration settings for nodes in a node pool.
type KubernetesNodeSettings struct {
	Annotations map[string]string `json:"annotations"` // Kubernetes node annotations
	Labels      map[string]string `json:"labels"`      // Kubernetes node labels
	Taints      []NodeTaint       `json:"taints"`      // Node taints for pod scheduling
}

// NodeTaint represents a taint that can be applied to nodes to control pod scheduling.
type NodeTaint struct {
	Key    string `json:"key"`    // Taint key
	Value  string `json:"value"`  // Taint value
	Effect string `json:"effect"` // Taint effect (NoSchedule, PreferNoSchedule, NoExecute)
}

type KubernetesNodePoolUpgradeStrategy string

const (
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
