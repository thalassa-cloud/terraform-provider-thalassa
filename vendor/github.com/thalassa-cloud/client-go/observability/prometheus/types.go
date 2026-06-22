package prometheus

import (
	"time"

	"github.com/thalassa-cloud/client-go/iaas"
	"github.com/thalassa-cloud/client-go/pkg/base"
)

type CreatePrometheusTenantRequest struct {
	// Name of the object
	Name string `json:"name"`
	// Description of the object
	Description string `json:"description"`
	// Annotations is a map of key-value pairs used for storing additional information
	Annotations map[string]string `json:"annotations,omitempty"`
	// Labels is a map of key-value pairs used for filtering and grouping objects
	Labels         map[string]string `json:"labels,omitempty"`
	RegionIdentity string            `json:"regionIdentity,omitempty"`
	// Retention is the retention period for metrics data in Prometheus duration format
	// Examples: "1h", "24h", "7d", "30d", "90d", "1y"
	Retention string `json:"retention,omitempty"`
}

type UpdatePrometheusTenantRequest struct {
	Name string `json:"name"`
	// Description of the object
	Description string `json:"description"`
	// Annotations is a map of key-value pairs used for storing additional information
	Annotations map[string]string `json:"annotations,omitempty"`
	// Labels is a map of key-value pairs used for filtering and grouping objects
	Labels         map[string]string `json:"labels,omitempty"`
	RegionIdentity string            `json:"regionIdentity,omitempty"`
	// Retention is the retention period for metrics data in Prometheus duration format
	// Examples: "1h", "24h", "7d", "30d", "90d", "1y"
	Retention string `json:"retention,omitempty"`
}

// PrometheusTenant represents a managed Prometheus tenant in a multitenant Cortex cluster
type PrometheusTenant struct {
	Identity      string            `json:"identity"`
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	Annotations   map[string]string `json:"annotations,omitempty"`
	Labels        map[string]string `json:"labels,omitempty"`
	CreatedAt     time.Time         `json:"createdAt"`
	UpdatedAt     time.Time         `json:"updatedAt"`
	ObjectVersion int               `json:"objectVersion"`

	Organisation *base.Organisation `json:"organisation,omitempty"`

	Region *iaas.Region `json:"region,omitempty"`

	// Status is the current status of the tenant
	Status PrometheusTenantStatus `json:"status"`

	// StatusMessage provides additional information about the current status
	StatusMessage string `json:"statusMessage,omitempty"`

	// RemoteWriteURL is the URL endpoint for remote write operations
	RemoteWriteURL string `json:"remoteWriteUrl,omitempty"`

	// RemoteWriteOTLPURL is the URL endpoint for remote write operations using OTLP
	// This is used for sending metrics to the Prometheus service using the OTLP protocol
	RemoteWriteOTLPURL string `json:"remoteWriteOtlpUrl,omitempty"`

	// QueryURL is the URL endpoint for querying metrics
	QueryURL string `json:"queryUrl,omitempty"`

	// AlertingURL is the URL endpoint for alerting operations
	AlertingURL string `json:"alertingUrl,omitempty"`

	// Usage statistics
	// TotalMetricsIngested is the total number of metric samples ingested
	TotalMetricsIngested int64 `json:"totalMetricsIngested,omitempty"`

	// TotalStorageBytes is the total storage used in bytes
	TotalStorageBytes int64 `json:"totalStorageBytes,omitempty"`

	// ActiveSeriesCount is the current number of active time series
	ActiveSeriesCount int64 `json:"activeSeriesCount,omitempty"`

	// TotalQueriesExecuted is the total number of queries executed
	TotalQueriesExecuted int64 `json:"totalQueriesExecuted,omitempty"`

	// TotalRemoteWriteRequests is the total number of remote write requests
	TotalRemoteWriteRequests int64 `json:"totalRemoteWriteRequests,omitempty"`

	// LastUsageUpdateAt is the timestamp when usage stats were last updated
	LastUsageUpdateAt *time.Time `json:"lastUsageUpdateAt,omitempty"`

	// DeleteScheduledAt is the date and time when the tenant will be permanently deleted
	// This provides a grace period before actual data destruction (minimum 1 day)
	DeleteScheduledAt *time.Time `json:"deleteScheduledAt,omitempty"`

	// Retention is the retention period for metrics data in Prometheus duration format
	// Examples: "1h", "24h", "7d", "30d", "90d", "1y"
	// If not set, a default retention period will be used
	Retention string `json:"retention,omitempty" validate:"omitempty,min=2,max=50"`
}

type PrometheusTenantStatus string

const (
	PrometheusTenantStatusCreating PrometheusTenantStatus = "creating"
	PrometheusTenantStatusReady    PrometheusTenantStatus = "ready"
	PrometheusTenantStatusUpdating PrometheusTenantStatus = "updating"
	PrometheusTenantStatusFailed   PrometheusTenantStatus = "failed"
	PrometheusTenantStatusDeleting PrometheusTenantStatus = "deleting"
	PrometheusTenantStatusDeleted  PrometheusTenantStatus = "deleted"
)
