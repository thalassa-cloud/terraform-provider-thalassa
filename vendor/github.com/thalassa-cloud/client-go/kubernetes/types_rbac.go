package kubernetes

import (
	"time"

	"github.com/thalassa-cloud/client-go/iam"
	"github.com/thalassa-cloud/client-go/pkg/base"
)

type KubernetesClusterRole struct {
	// Identity is a unique identifier for the Kubernetes Cluster Role
	Identity string `json:"identity"`
	// Name is a human-readable name of the Kubernetes Cluster Role
	Name string `json:"name"`
	// Slug is a human-readable unique identifier for the Kubernetes Cluster Role
	Slug string `json:"slug"`
	// Description is a human-readable description of the Kubernetes Cluster Role
	Description string `json:"description,omitempty"`
	// Annotations is a map of key-value pairs used for storing additional information
	Annotations map[string]string `json:"annotations,omitempty"`

	// Labels is a map of key-value pairs used for filtering and grouping Kubernetes Cluster Roles
	Labels map[string]string `json:"labels,omitempty"`
	// CreatedAt is the timestamp when the Kubernetes Cluster Role was created
	CreatedAt time.Time `json:"createdAt"`
	// UpdatedAt is the timestamp when the Kubernetes Cluster Role was last updated
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`

	// ObjectVersion represents the version of the Kubernetes Cluster Role.
	ObjectVersion int64 `json:"objectVersion"`
	// System is a flag to indicate if the role is a system role
	System bool `json:"system"`
	// Organisation
	Organisation *base.Organisation `json:"organisation,omitempty"`

	// Rules is a list of permission rules for this Kubernetes Cluster Role
	Rules []KubernetesClusterRolePermissionRule `json:"rules,omitempty"`

	// Bindings is a list of bindings for this Kubernetes Cluster Role
	Bindings []KubernetesClusterRoleBinding `json:"bindings,omitempty"`
}

type KubernetesClusterRolePermissionRule struct {
	// Identity is a unique identifier for the Kubernetes Cluster Role Permission Rule
	Identity string `json:"identity"`

	// KubernetesClusterRole is the Kubernetes Cluster Role that the rule is for
	KubernetesClusterRole *KubernetesClusterRole `json:"kubernetesClusterRole,omitempty"`

	// Resources is a list of resources that the rule applies to
	Resources []string `json:"resources"`

	// Verbs is a list of verbs that the rule applies to
	Verbs []KubernetesClusterRolePermissionVerb `json:"verbs"`

	// ApiGroups is a list of API groups that the rule applies to
	ApiGroups []string `json:"apiGroups"`

	// ResourceNames is a list of resource names that the rule applies to
	ResourceNames []string `json:"resourceNames"`
	// NonResourceURLs is a list of non-resource URLs that the rule applies to
	NonResourceURLs []string `json:"nonResourceURLs"`
	// Note is a human-readable note for the permission rule
	Note string `json:"note,omitempty"`
}

type KubernetesClusterRoleBinding struct {
	// Identity is a unique identifier for the Kubernetes Cluster Role Binding
	Identity string `json:"identity"`
	// Name is a human-readable name of the Kubernetes Cluster Role Binding
	Name string `json:"name"`
	// Slug is a human-readable unique identifier for the Kubernetes Cluster Role Binding
	Slug string `json:"slug"`
	// Description is a human-readable description of the Kubernetes Cluster Role Binding
	Description string `json:"description,omitempty"`
	// Annotations is a map of key-value pairs used for storing additional information
	Annotations map[string]string `json:"annotations,omitempty"`

	// Labels is a map of key-value pairs used for filtering and grouping Kubernetes Cluster Role Bindings
	Labels map[string]string `json:"labels,omitempty"`
	// CreatedAt is the timestamp when the Kubernetes Cluster Role Binding was created
	CreatedAt time.Time `json:"createdAt"`
	// UpdatedAt is the timestamp when the Kubernetes Cluster Role was last updated
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	// ObjectVersion represents the version of the Kubernetes Cluster Role Binding.
	ObjectVersion int64 `json:"objectVersion"`

	// KubernetesClusterRole is the Kubernetes Cluster Role that the binding is for
	KubernetesClusterRole *KubernetesClusterRole `json:"kubernetesClusterRole,omitempty"`
	// User is the user that the binding is for
	User *base.AppUser `json:"user,omitempty"`
	// OrganisationTeam is the team that the binding is for
	OrganisationTeam *iam.Team `json:"team,omitempty"`
	// ServiceAccount is the service account that the binding is for
	ServiceAccount *iam.ServiceAccount `json:"serviceAccount,omitempty"`
	// ExpiresAt is the time at which the binding expires. Optional.
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
	// Note is a human-readable note for the binding. Optional.
	Note *string `json:"note,omitempty"`
}

type CreateKubernetesClusterRoleRequest struct {
	// Name of the Kubernetes Cluster Role
	Name string `json:"name"`
	// Description of the Kubernetes Cluster Role
	Description string `json:"description"`
	// Annotations is a map of key-value pairs used for storing additional information
	Annotations map[string]string `json:"annotations,omitempty"`
	// Labels is a map of key-value pairs used for filtering and grouping Kubernetes Cluster Roles
	Labels map[string]string `json:"labels,omitempty"`
}

type CreateKubernetesClusterRoleBinding struct {
	// Name of the object
	Name string `json:"name"`
	// Description of the object
	Description string `json:"description"`
	// Annotations is a map of key-value pairs used for storing additional information
	Annotations map[string]string `json:"annotations,omitempty"`
	// Labels is a map of key-value pairs used for filtering and grouping objects
	Labels map[string]string `json:"labels,omitempty"`

	// UserIdentity is the identity of the user to bind. Must be provided if TeamIdentity is not provided.
	UserIdentity *string `json:"userIdentity"`

	// TeamIdentity is the identity of the team to bind. Must be provided if UserIdentity is not provided.
	TeamIdentity *string `json:"teamIdentity"`

	// ServiceAccountIdentity is the identity of the service account to bind. Must be provided if UserIdentity and TeamIdentity are not provided.
	ServiceAccountIdentity *string `json:"serviceAccountIdentity"`
}

type AddKubernetesClusterRolePermissionRule struct {
	Resources       []string                              `json:"resources"`
	Verbs           []KubernetesClusterRolePermissionVerb `json:"verbs"`
	ApiGroups       []string                              `json:"apiGroups,omitempty"`
	ResourceNames   []string                              `json:"resourceNames,omitempty"`
	NonResourceURLs []string                              `json:"nonResourceURLs,omitempty"`
	Note            string                                `json:"note,omitempty"`
}

// KubernetesClusterRolePermissionVerb is a verb that can be used to describe a permission rule
type KubernetesClusterRolePermissionVerb string

const (
	// KubernetesClusterRolePermissionVerbWildcard is a wildcard verb that can be used to describe a permission rule
	KubernetesClusterRolePermissionVerbWildcard KubernetesClusterRolePermissionVerb = "*"
	KubernetesClusterRolePermissionVerbGet      KubernetesClusterRolePermissionVerb = "get"
	KubernetesClusterRolePermissionVerbList     KubernetesClusterRolePermissionVerb = "list"
	KubernetesClusterRolePermissionVerbWatch    KubernetesClusterRolePermissionVerb = "watch"
	KubernetesClusterRolePermissionVerbCreate   KubernetesClusterRolePermissionVerb = "create"
	KubernetesClusterRolePermissionVerbUpdate   KubernetesClusterRolePermissionVerb = "update"
	KubernetesClusterRolePermissionVerbDelete   KubernetesClusterRolePermissionVerb = "delete"
	KubernetesClusterRolePermissionVerbPatch    KubernetesClusterRolePermissionVerb = "patch"
)

var (
	KubernetesClusterRolePermissionVerbs = []KubernetesClusterRolePermissionVerb{
		KubernetesClusterRolePermissionVerbGet,
		KubernetesClusterRolePermissionVerbList,
		KubernetesClusterRolePermissionVerbWatch,
		KubernetesClusterRolePermissionVerbCreate,
		KubernetesClusterRolePermissionVerbUpdate,
		KubernetesClusterRolePermissionVerbDelete,
		KubernetesClusterRolePermissionVerbPatch,
	}
)
