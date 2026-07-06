package quicklaunch

import (
	"time"

	"github.com/thalassa-cloud/client-go/pkg/base"
)

// QuickLaunchTemplateType represents the type of quick launch template.
type QuickLaunchTemplateType string

const (
	// QuickLaunchTemplateVPC creates a VPC with subnets and NAT Gateway.
	QuickLaunchTemplateVPC QuickLaunchTemplateType = "vpc"
	// QuickLaunchTemplateKubernetes creates a VPC with subnets, NAT Gateway, and Kubernetes cluster.
	QuickLaunchTemplateKubernetes QuickLaunchTemplateType = "kubernetes"
)

// QuickLaunchRequest is the request body for the quick launch API.
type QuickLaunchRequest struct {
	// Template is the template type to use (e.g., "vpc", "kubernetes").
	// Defaults to "vpc" if not specified on the server when omitted.
	Template QuickLaunchTemplateType `json:"template,omitempty"`
	// Name is the base name prefix for all created resources.
	Name string `json:"name"`
	// Description is an optional description for the resources.
	Description string `json:"description,omitempty"`
	// CloudRegionIdentity is the identity of the cloud region where resources will be created.
	CloudRegionIdentity string `json:"cloudRegionIdentity"`
	// VpcCidr is an optional CIDR block for the VPC (e.g., "10.0.0.0/16").
	// If not provided, a default CIDR will be auto-generated on the server.
	VpcCidr string `json:"vpcCidr,omitempty"`
	// SubnetCidrs is an optional list of subnet CIDRs. If not provided, default subnets will be created.
	SubnetCidrs []string `json:"subnetCidrs,omitempty"`
	// MachineType is an optional machine type identity or slug for Kubernetes node pools.
	// Only used when template is "kubernetes".
	MachineType string `json:"machineType,omitempty"`
	// Labels is a map of key-value pairs used for filtering and grouping objects.
	Labels map[string]string `json:"labels,omitempty"`
	// Annotations is a map of key-value pairs used for storing additional information.
	Annotations map[string]string `json:"annotations,omitempty"`
}

// QuickLaunchResource represents a created resource in the quick launch response.
type QuickLaunchResource struct {
	// Type is the resource type (e.g., "vpc", "subnet", "natgateway", "kubernetes_cluster", "kubernetes_node_pool").
	Type string `json:"type"`
	// Name is the name of the resource.
	Name string `json:"name"`
	// Identity is the unique identity of the resource.
	Identity string `json:"identity"`

	LastSeenAt time.Time `json:"lastSeenAt"`
	LastStatus string    `json:"lastStatus,omitempty"`
}

// QuickLaunchResources is a collection of QuickLaunchResource with helper methods.
type QuickLaunchResources []QuickLaunchResource

// Add adds a resource to the collection if it doesn't already exist (based on Identity).
func (r *QuickLaunchResources) Add(resource QuickLaunchResource) {
	for _, existing := range *r {
		if existing.Identity == resource.Identity {
			return
		}
	}
	*r = append(*r, resource)
}

// ToSlice returns the underlying slice.
func (r QuickLaunchResources) ToSlice() []QuickLaunchResource {
	return []QuickLaunchResource(r)
}

// QuickLaunch is a quick launch job returned by the API (persisted state, async provisioning).
type QuickLaunch struct {
	Identity            string                  `json:"identity"`
	Name                string                  `json:"name"`
	Slug                string                  `json:"slug,omitempty"`
	Description         string                  `json:"description,omitempty"`
	Template            QuickLaunchTemplateType `json:"template,omitempty"`
	Status              string                  `json:"status"`
	StatusMessage       string                  `json:"statusMessage,omitempty"`
	CloudRegionIdentity string                  `json:"cloudRegionIdentity,omitempty"`
	VpcCidr             string                  `json:"vpcCidr,omitempty"`
	SubnetCidrs         []string                `json:"subnetCidrs,omitempty"`
	MachineType         string                  `json:"machineType,omitempty"`
	Resources           QuickLaunchResources    `json:"resources,omitempty"`
	Labels              map[string]string       `json:"labels,omitempty"`
	Annotations         map[string]string       `json:"annotations,omitempty"`
	CreatedAt           time.Time               `json:"createdAt"`
	UpdatedAt           *time.Time              `json:"updatedAt,omitempty"`
	ObjectVersion       int                     `json:"objectVersion,omitempty"`
	Organisation        *base.Organisation      `json:"organisation,omitempty"`
}

// QuickLaunchCascade controls child resource handling on delete (query: cascade).
type QuickLaunchCascade string

const (
	// QuickLaunchCascadeDelete deletes provisioned resources with the quick launch.
	QuickLaunchCascadeDelete QuickLaunchCascade = "Delete"
	// QuickLaunchCascadeOrphan keeps provisioned resources and only removes the quick launch record.
	QuickLaunchCascadeOrphan QuickLaunchCascade = "Orphan"
)

// QuickLaunchLogs is the raw response body from GET /v1/quick-launch/{identity}/logs.
// Unmarshal as JSON into your own type if needed.
type QuickLaunchLogs []byte
