package observability

import (
	"time"

	"github.com/thalassa-cloud/client-go/iaas"
	"github.com/thalassa-cloud/client-go/pkg/base"
)

// ObservabilityWorkspaceStatus is the lifecycle status of a workspace.
type ObservabilityWorkspaceStatus string

const (
	ObservabilityWorkspaceStatusProvisioning ObservabilityWorkspaceStatus = "provisioning"
	ObservabilityWorkspaceStatusReady        ObservabilityWorkspaceStatus = "ready"
	ObservabilityWorkspaceStatusUpdating     ObservabilityWorkspaceStatus = "updating"
	ObservabilityWorkspaceStatusFailed       ObservabilityWorkspaceStatus = "failed"
	ObservabilityWorkspaceStatusDeleting     ObservabilityWorkspaceStatus = "deleting"
	ObservabilityWorkspaceStatusDeleted      ObservabilityWorkspaceStatus = "deleted"
)

// CreateObservabilityWorkspaceRequest creates a unified metrics/logs workspace.
type CreateObservabilityWorkspaceRequest struct {
	Name           string            `json:"name"`
	Description    string            `json:"description,omitempty"`
	Annotations    map[string]string `json:"annotations,omitempty"`
	Labels         map[string]string `json:"labels,omitempty"`
	RegionIdentity string            `json:"regionIdentity,omitempty"`
	// RetentionDays applies to all enabled backends. Nil uses the service default.
	RetentionDays *int `json:"retentionDays,omitempty"`
}

// UpdateObservabilityWorkspaceRequest updates an existing workspace.
type UpdateObservabilityWorkspaceRequest struct {
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	// RetentionDays applies to all enabled backends. Nil uses the service default.
	RetentionDays *int `json:"retentionDays,omitempty"`
}

// ObservabilityWorkspace is a unified Prometheus (Cortex) + Loki workspace.
type ObservabilityWorkspace struct {
	Identity      string            `json:"identity"`
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	Annotations   map[string]string `json:"annotations,omitempty"`
	Labels        map[string]string `json:"labels,omitempty"`
	CreatedAt     time.Time         `json:"createdAt"`
	UpdatedAt     time.Time         `json:"updatedAt"`
	ObjectVersion int               `json:"objectVersion"`

	Organisation *base.Organisation `json:"organisation,omitempty"`
	Region       *iaas.Region       `json:"region,omitempty"`

	Status        ObservabilityWorkspaceStatus `json:"status"`
	StatusMessage string                       `json:"statusMessage,omitempty"`

	PrometheusEnabled bool `json:"prometheusEnabled"`
	LokiEnabled       bool `json:"lokiEnabled"`

	RemoteWriteURL     string `json:"remoteWriteUrl,omitempty"`
	RemoteWriteOTLPURL string `json:"remoteWriteOtlpUrl,omitempty"`
	PrometheusQueryURL string `json:"prometheusQueryUrl,omitempty"`
	AlertingURL        string `json:"alertingUrl,omitempty"`

	PushURL      string `json:"pushUrl,omitempty"`
	PushOTLPURL  string `json:"pushOtlpUrl,omitempty"`
	LokiQueryURL string `json:"lokiQueryUrl,omitempty"`

	RetentionDays *int `json:"retentionDays,omitempty"`

	TotalMetricsIngested     int64 `json:"totalMetricsIngested,omitempty"`
	TotalLogLinesIngested    int64 `json:"totalLogLinesIngested,omitempty"`
	TotalStorageBytes        int64 `json:"totalStorageBytes,omitempty"`
	ActiveSeriesCount        int64 `json:"activeSeriesCount,omitempty"`
	TotalQueriesExecuted     int64 `json:"totalQueriesExecuted,omitempty"`
	TotalRemoteWriteRequests int64 `json:"totalRemoteWriteRequests,omitempty"`
	TotalPushRequests        int64 `json:"totalPushRequests,omitempty"`

	LastUsageUpdateAt *time.Time `json:"lastUsageUpdateAt,omitempty"`
	DeleteScheduledAt *time.Time `json:"deleteScheduledAt,omitempty"`
}
