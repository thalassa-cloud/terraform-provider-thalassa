package audit

import (
	"time"

	"github.com/thalassa-cloud/client-go/iam"
	"github.com/thalassa-cloud/client-go/pkg/base"
)

type AuditLog struct {
	CreatedAt              time.Time              `json:"createdAt"`
	EventID                string                 `json:"eventID"`
	UserIdentity           *string                `json:"userIdentity,omitempty"`
	User                   *base.AppUser          `json:"user,omitempty"`
	ServiceAccountIdentity *string                `json:"serviceAccountIdentity,omitempty"`
	ServiceAccount         *iam.ServiceAccount    `json:"serviceAccount,omitempty"`
	OrganizationIdentity   *string                `json:"organizationIdentity,omitempty"`
	Organization           *base.Organisation     `json:"organization,omitempty"`
	ImpersonatorIdentity   *string                `json:"impersonatorIdentity,omitempty"`
	Impersonator           *base.AppUser          `json:"impersonator,omitempty"`
	Action                 string                 `json:"action"`
	Description            *string                `json:"description,omitempty"`
	ResourceType           *string                `json:"resourceType,omitempty"`
	ResourceIdentity       *string                `json:"resourceIdentity,omitempty"`
	Context                map[string]interface{} `json:"context,omitempty"`
}

type AuditLogFilter struct {
	SearchText string `json:"searchText,omitempty"`

	ServiceAccount        string   `json:"serviceAccount,omitempty"`
	UserIdentity          string   `json:"userIdentity,omitempty"`
	ImpersonatorIdentity  string   `json:"impersonatorIdentity,omitempty"`
	Actions               []string `json:"actions,omitempty"`
	ResourceTypes         []string `json:"resourceTypes,omitempty"`
	ResourceIdentity      string   `json:"resourceIdentity,omitempty"`
	OrganizationIdentity  string   `json:"organizationIdentity,omitempty"`
	IncludeSystemServices bool     `json:"includeSystemServices,omitempty"`

	ResponseStatus int `json:"responseStatus,omitempty"`
}

type PagedResult[T any] struct {
	Items      []T `json:"items"`
	TotalItems int `json:"totalItems"`

	Page         int `json:"page"`
	ItemsPerPage int `json:"itemsPerPage"`

	TotalPages int `json:"totalPages"`
	TotalCount int `json:"totalCount"`
}
