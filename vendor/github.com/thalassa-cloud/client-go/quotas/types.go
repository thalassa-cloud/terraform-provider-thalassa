package quotas

import (
	"time"

	"github.com/thalassa-cloud/client-go/pkg/base"
)

// OrganisationQuota represents a quota for an organisation
type OrganisationQuota struct {
	// Name is the name of the quota
	Name string `json:"name"`
	// Description is the description of the quota
	Description string `json:"description"`
	// Service is the service that the quota applies to
	Service *string `json:"service,omitempty"`
	// CreatedAt is the timestamp when the quota was created
	CreatedAt time.Time `json:"createdAt"`
	// MaxUsage is the maximum usage allowed for this quota
	MaxUsage int64 `json:"maxUsage"`
	// CurrentUsage is the current usage of this quota
	CurrentUsage int64 `json:"currentUsage"`
	// QuotaType is the type of the quota
	QuotaType string `json:"quotaType"`
	// IncreaseRequests is a list of increase requests for this quota
	IncreaseRequests []IncreaseOrganisationQuotaRequest `json:"increaseRequests,omitempty"`
	// Organisation is the organisation that this quota belongs to
	Organisation *base.Organisation `json:"organisation,omitempty"`
}

// IncreaseOrganisationQuotaRequest represents a request to increase an organisation quota
type IncreaseOrganisationQuotaRequest struct {
	// Identity is the unique identifier for the increase request
	Identity string `json:"identity"`
	// NewMaxUsageRequested is the new maximum usage requested
	NewMaxUsageRequested int64 `json:"newMaxUsageRequested"`
	// RequestedReasonMessage is the reason for the quota increase request
	RequestedReasonMessage string `json:"requestedReasonMessage"`
	// CreatedAt is the timestamp when the request was created
	CreatedAt time.Time `json:"createdAt"`
	// UpdatedAt is the timestamp when the request was last updated
	UpdatedAt time.Time `json:"updatedAt"`
	// Decision is the decision on the request (approved, denied, pending)
	Decision string `json:"decision"`
	// DecisionReason is the reason for the decision
	DecisionReason *string `json:"decisionReason,omitempty"`
	// OrganisationQuota is the quota that this request is for
	OrganisationQuota *OrganisationQuota `json:"organisationQuota,omitempty"`
}

// RequestQuotaIncreaseRequest represents a request to increase a quota
type RequestQuotaIncreaseRequest struct {
	// Name is the name of the quota to increase
	Name string `json:"name"`
	// NewMaxUsage is the new maximum usage requested
	NewMaxUsage int64 `json:"newMaxUsage"`
	// Reason is the reason for the quota increase request
	Reason string `json:"reason"`
}
