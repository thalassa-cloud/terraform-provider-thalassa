package containerregistry

import (
	"time"

	"github.com/thalassa-cloud/client-go/iaas"
	"github.com/thalassa-cloud/client-go/pkg/base"
)

type Labels map[string]string
type Annotations map[string]string

type CreateContainerRegistryNamespaceRequest struct {
	Region      string `json:"region"`
	Namespace   string `json:"namespace"`
	Description string `json:"description"`

	Labels      Labels      `json:"labels"`
	Annotations Annotations `json:"annotations"`
}

type UpdateContainerRegistryNamespaceRequest struct {
	Description string `json:"description"`

	Labels      Labels      `json:"labels"`
	Annotations Annotations `json:"annotations"`
}

type CreateNamespaceConfigurationRequest struct {
	Visibility      NamespaceVisibility `json:"visibility"`
	RetentionPolicy *RetentionPolicy    `json:"retention_policy,omitempty"`
}

type UpdateNamespaceConfigurationRequest struct {
	Visibility      NamespaceVisibility `json:"visibility"`
	RetentionPolicy *RetentionPolicy    `json:"retention_policy,omitempty"`
}

type ContainerRegistryNamespace struct {
	Identity      string      `json:"identity"`
	Namespace     string      `json:"namespace"`
	CreatedAt     time.Time   `json:"createdAt"`
	UpdatedAt     time.Time   `json:"updatedAt"`
	ObjectVersion int64       `json:"objectVersion"`
	Description   string      `json:"description"`
	Annotations   Annotations `json:"annotations,omitempty"`
	Labels        Labels      `json:"labels,omitempty"`

	Region *iaas.Region `json:"region,omitempty"`

	Organisation   *base.Organisation                       `json:"organisation,omitempty"`
	TotalSizeBytes int64                                    `json:"total_size_bytes"`
	Repositories   []ContainerRegistryRepository            `json:"repositories,omitempty"`
	Configuration  *ContainerRegistryNamespaceConfiguration `json:"configuration,omitempty"`
}

type ContainerRegistryType string

type NamespaceVisibility string

const (
	NamespaceVisibilityPrivate NamespaceVisibility = "private"
)

type ContainerRegistryNamespaceConfiguration struct {
	CreatedAt     time.Time                   `json:"created_at"`
	UpdatedAt     time.Time                   `json:"updated_at"`
	ObjectVersion int64                       `json:"object_version"`
	Namespace     *ContainerRegistryNamespace `json:"namespace,omitempty"`
	Visibility    NamespaceVisibility         `json:"visibility"`
	// RetentionPolicy stores retention policy configuration
	RetentionPolicy *RetentionPolicy `json:"retention_policy,omitempty"`
}

// RetentionPolicyScope defines what resources are covered by the retention policy
type RetentionPolicyScope string

const (
	RetentionPolicyScopeTags RetentionPolicyScope = "tags"
)

// RetentionPolicyRule defines a single retention rule within a retention policy
type RetentionPolicyRule struct {
	// Days specifies the number of days to retain resources based on tag creation time. Tags older than this will be deleted.
	// If not set, age-based retention by tag creation is disabled for this rule.
	Days *int `json:"days,omitempty"`

	// DaysSinceCreated specifies the number of days to retain resources based on artifact creation/push time.
	// Artifacts created/pushed within the last X days will be kept.
	// If not set, age-based retention by artifact creation is disabled for this rule.
	DaysSinceCreated *int `json:"days_since_created,omitempty"`

	// DaysSincePulled specifies the number of days to retain resources based on artifact last pull time.
	// Artifacts pulled within the last X days will be kept.
	// If not set, age-based retention by artifact pull is disabled for this rule.
	DaysSincePulled *int `json:"days_since_pulled,omitempty"`

	// Count specifies the number of versions/tags to keep. Only the most recent N items will be retained.
	// If not set, count-based retention is disabled for this rule.
	Count *int `json:"count,omitempty"`

	// RepositoryPatterns is a list of repository name patterns to match. Only repositories matching these patterns will be subject to this rule.
	// Supports wildcards (e.g., "myapp/*", "*-service", "frontend").
	// If empty, all repositories are considered for this rule.
	RepositoryPatterns []string `json:"repository_patterns,omitempty"`

	// TagPatterns is a list of tag patterns to match. Only tags matching these patterns will be subject to this rule.
	// Supports wildcards (e.g., "v*", "*-dev", "latest").
	// If empty, all tags are considered for this rule.
	TagPatterns []string `json:"tag_patterns,omitempty"`

	// Scope defines what resources are covered by this retention rule
	// Can be one of: tags
	// If not set, defaults to "tags"
	Scope RetentionPolicyScope `json:"scope,omitempty"`
}

// RetentionPolicy defines the retention policy configuration similar to Harbor's retention policies
// It supports multiple rules that are evaluated together
type RetentionPolicy struct {
	// Enabled indicates whether the retention policy is active
	Enabled bool `json:"enabled"`
	// DeleteUntaggedImages indicates whether to delete untagged images
	DeleteUntaggedImages bool `json:"delete_untagged_images,omitempty"`
	// Rules is a list of retention rules. Rules are evaluated in order, and resources matching any rule
	// will be subject to retention according to that rule's criteria.
	Rules []RetentionPolicyRule `json:"rules"`
}

type ContainerRegistryRepository struct {
	Identity      string                      `json:"identity"`
	CreatedAt     time.Time                   `json:"createdAt"`
	UpdatedAt     time.Time                   `json:"updatedAt"`
	ObjectVersion int64                       `json:"objectVersion"`
	Namespace     *ContainerRegistryNamespace `json:"namespace,omitempty"`
	Image         string                      `json:"image"`
	FullName      string                      `json:"full_name"`
	Description   string                      `json:"description"`

	LastPulledAt   *time.Time `json:"last_pulled_at"`
	LastPushedAt   *time.Time `json:"last_pushed_at"`
	TotalSizeBytes int64      `json:"total_size_bytes"`
	TagCount       int64      `json:"tag_count"`
	ArtifactCount  int64      `json:"artifact_count"`

	Tags      []ContainerRegistryImageTag `json:"tags,omitempty"`
	Artifacts []ContainerRegistryArtifact `json:"artifacts,omitempty"`
}

type ContainerRegistryImageTag struct {
	CreatedAt     time.Time                    `json:"created_at"`
	UpdatedAt     time.Time                    `json:"updated_at"`
	ObjectVersion int64                        `json:"object_version"`
	Repository    *ContainerRegistryRepository `json:"repository,omitempty"`
	Tag           string                       `json:"tag"`
	Sha256        string                       `json:"sha256"`
	SizeMb        float64                      `json:"size_mb"`
}

type ContainerRegistryArtifact struct {
	CreatedAt    time.Time                    `json:"created_at"`
	LastPulledAt *time.Time                   `json:"last_pulled_at"`
	MediaType    string                       `json:"media_type"`
	Size         float64                      `json:"size"`
	Digest       string                       `json:"digest"`
	Repository   *ContainerRegistryRepository `json:"repository,omitempty"`
	// Not yet supported
	// SBOMs                []ContainerRegistrySBOM                `json:"sboms,omitempty"`
	// VulnerabilityReports []ContainerRegistryVulnerabilityReport `json:"vulnerability_reports,omitempty"`
}
