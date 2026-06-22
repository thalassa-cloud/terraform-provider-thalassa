package audit

import (
	"context"
	"strconv"
	"strings"

	"github.com/thalassa-cloud/client-go/pkg/client"
)

const (
	AuditLogEndpoint = "/v1/audit"
)

// ListAuditLogs lists all audit logs for the current organisation.
// The current organisation is determined by the client's organisation identity.
func (c *Client) ListAuditLogs(ctx context.Context, listRequest *ListAuditLogsRequest) (*PagedResult[AuditLog], error) {
	auditLogs := &PagedResult[AuditLog]{}

	queryParams := map[string]string{}

	if listRequest != nil && listRequest.Page > 0 {
		queryParams["page"] = strconv.Itoa(listRequest.Page)
	}
	if listRequest != nil && listRequest.Limit > 0 {
		queryParams["limit"] = strconv.Itoa(listRequest.Limit)
	}

	if listRequest != nil && listRequest.Filter != nil {
		filter := listRequest.Filter

		if filter.SearchText != "" {
			queryParams["searchText"] = filter.SearchText
		}
		if filter.ServiceAccount != "" {
			queryParams["serviceAccount"] = filter.ServiceAccount
		}
		if filter.UserIdentity != "" {
			queryParams["userIdentity"] = filter.UserIdentity
		}
		if filter.ImpersonatorIdentity != "" {
			queryParams["impersonatorIdentity"] = filter.ImpersonatorIdentity
		}
		if len(filter.Actions) > 0 {
			queryParams["action"] = strings.Join(filter.Actions, ",")
		}
		if len(filter.ResourceTypes) > 0 {
			queryParams["resourceType"] = strings.Join(filter.ResourceTypes, ",")
		}
		if filter.ResourceIdentity != "" {
			queryParams["resourceIdentity"] = filter.ResourceIdentity
		}
		if filter.OrganizationIdentity != "" {
			queryParams["organizationIdentity"] = filter.OrganizationIdentity
		}
		if filter.IncludeSystemServices {
			queryParams["includeSystemServices"] = strconv.FormatBool(filter.IncludeSystemServices)
		}
		if filter.ResponseStatus != 0 {
			queryParams["responseStatus"] = strconv.Itoa(filter.ResponseStatus)
		}
	}

	req := c.R().SetQueryParams(queryParams).SetResult(&auditLogs)
	resp, err := c.Do(ctx, req, client.GET, AuditLogEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return auditLogs, err
	}
	return auditLogs, nil
}

type ListAuditLogsRequest struct {
	Page   int `json:"page,omitempty"`
	Limit  int `json:"limit,omitempty"`
	Filter *AuditLogFilter
}
