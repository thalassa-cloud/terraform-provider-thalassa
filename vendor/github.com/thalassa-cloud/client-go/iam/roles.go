package iam

import (
	"context"
	"fmt"
	"time"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/base"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

const (
	OrganisationRoleEndpoint = "/v1/iam/roles"
)

// ListOrganisationRoles lists all organisation roles for a given organisation.
func (c *Client) ListOrganisationRoles(ctx context.Context, request *ListOrganisationRolesRequest) ([]OrganisationRole, error) {
	roles := []OrganisationRole{}
	req := c.R().SetResult(&roles)
	if request != nil {
		for _, filter := range request.Filters {
			for k, v := range filter.ToParams() {
				req.SetQueryParam(k, v)
			}
		}
	}

	resp, err := c.Do(ctx, req, client.GET, OrganisationRoleEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return roles, err
	}
	return roles, nil
}

// GetOrganisationRole retrieves a specific organisation role by its identity.
func (c *Client) GetOrganisationRole(ctx context.Context, identity string) (*OrganisationRole, error) {
	var role *OrganisationRole
	req := c.R().SetResult(&role)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", OrganisationRoleEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return role, err
	}
	return role, nil
}

// CreateOrganisationRole creates a new organisation role.
func (c *Client) CreateOrganisationRole(ctx context.Context, create CreateOrganisationRoleRequest) (*OrganisationRole, error) {
	var role *OrganisationRole
	req := c.R().SetBody(create).SetResult(&role)
	resp, err := c.Do(ctx, req, client.POST, OrganisationRoleEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return role, err
	}
	return role, nil
}

// DeleteOrganisationRole deletes a specific organisation role by its identity.
func (c *Client) DeleteOrganisationRole(ctx context.Context, identity string) error {
	req := c.R()
	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s", OrganisationRoleEndpoint, identity))
	if err != nil {
		return err
	}
	if err := c.Check(resp); err != nil {
		return err
	}
	return nil
}

// AddRoleRule adds a new permission rule to an organisation role.
func (c *Client) AddRoleRule(ctx context.Context, roleIdentity string, rule OrganisationRolePermissionRule) (*OrganisationRolePermissionRule, error) {
	var result *OrganisationRolePermissionRule
	req := c.R().SetBody(rule).SetResult(&result)
	resp, err := c.Do(ctx, req, client.POST, fmt.Sprintf("%s/%s/rules", OrganisationRoleEndpoint, roleIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return result, err
	}
	return result, nil
}

// DeleteRuleFromRole deletes a permission rule from an organisation role.
func (c *Client) DeleteRuleFromRole(ctx context.Context, roleIdentity string, ruleIdentity string) error {
	req := c.R()
	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s/rules/%s", OrganisationRoleEndpoint, roleIdentity, ruleIdentity))
	if err != nil {
		return err
	}
	if err := c.Check(resp); err != nil {
		return err
	}
	return nil
}

// ListRoleBindings lists all role bindings for a specific organisation role.
func (c *Client) ListRoleBindings(ctx context.Context, roleIdentity string, request *ListRoleBindingsRequest) ([]OrganisationRoleBinding, error) {
	bindings := []OrganisationRoleBinding{}
	req := c.R().SetResult(&bindings)
	if request != nil {
		for _, filter := range request.Filters {
			for k, v := range filter.ToParams() {
				req.SetQueryParam(k, v)
			}
		}
	}

	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s/bindings", OrganisationRoleEndpoint, roleIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return bindings, err
	}
	return bindings, nil
}

// CreateRoleBinding creates a new role binding for an organisation role.
func (c *Client) CreateRoleBinding(ctx context.Context, roleIdentity string, create CreateRoleBinding) (*OrganisationRoleBinding, error) {
	var binding *OrganisationRoleBinding
	req := c.R().SetBody(create).SetResult(&binding)
	resp, err := c.Do(ctx, req, client.POST, fmt.Sprintf("%s/%s/bindings", OrganisationRoleEndpoint, roleIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return binding, err
	}
	return binding, nil
}

// DeleteRoleBinding deletes a specific role binding from an organisation role.
func (c *Client) DeleteRoleBinding(ctx context.Context, roleIdentity string, bindingIdentity string) error {
	req := c.R()
	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s/bindings/%s", OrganisationRoleEndpoint, roleIdentity, bindingIdentity))
	if err != nil {
		return err
	}
	if err := c.Check(resp); err != nil {
		return err
	}
	return nil
}

type ListOrganisationRolesRequest struct {
	Filters []filters.Filter
}

type ListRoleBindingsRequest struct {
	Filters []filters.Filter
}

type CreateOrganisationRoleRequest struct {
	// Name of the organisationRole
	Name string `json:"name"`

	// Description of the organisationRole
	Description string `json:"description"`

	// Annotations for the organisationRole
	Annotations map[string]string `json:"annotations"`

	// Labels for the organisationRole
	Labels map[string]string `json:"labels"`
}

type CreateRoleBinding struct {
	// Name of the organisationRole
	Name string `json:"name"`

	// Description of the organisationRole
	Description string `json:"description"`

	// Annotations for the organisationRole
	Annotations map[string]string `json:"annotations"`

	// Labels for the organisationRole
	Labels map[string]string `json:"labels"`

	// UserIdentity is the identity of the user to bind
	UserIdentity *string `json:"userIdentity"`

	// TeamIdentity is the identity of the team to bind
	TeamIdentity *string `json:"teamIdentity"`

	// ServiceAccountIdentity is the identity of the service account to bind
	ServiceAccountIdentity *string `json:"serviceAccountIdentity"`

	// Scopes is the scopes to bind the role binding to
	Scopes []string `json:"scopes"`
}

type OrganisationRole struct {
	Identity      string            `json:"identity"`
	Name          string            `json:"name"`
	Slug          string            `json:"slug"`
	Description   string            `json:"description"`
	CreatedAt     time.Time         `json:"createdAt"`
	UpdatedAt     time.Time         `json:"updatedAt"`
	ObjectVersion int               `json:"objectVersion"`
	Labels        map[string]string `json:"labels"`
	Annotations   map[string]string `json:"annotations"`
	IsReadOnly    bool              `json:"isReadOnly,omitempty"`
	// Organisation
	Organisation *base.Organisation               `json:"organisation,omitempty"`
	Rules        []OrganisationRolePermissionRule `json:"rules,omitempty"`
	Bindings     []OrganisationRoleBinding        `json:"bindings,omitempty"`
	// System is a flag that indicates if the role is a system role. Cannot be modified by the user. Default is false. Can only be set by the system.
	System bool `json:"system,omitempty"`
}

type OrganisationRolePermissionRule struct {
	// Identity is a unique identifier for the object
	Identity         string            `json:"identity"`
	OrganisationRole *OrganisationRole `json:"organisationRole,omitempty"`

	// Permission
	Resources          []string         `json:"resources"`
	ResourceIdentities []string         `json:"resourceIdentities"`
	Permissions        []PermissionType `json:"permissions"`
	// Note is a human-readable note for the permission rule
	Note string `json:"note,omitempty"`
}

type PermissionType string

const (
	PermissionTypeCreate   PermissionType = "create"
	PermissionTypeRead     PermissionType = "read"
	PermissionTypeUpdate   PermissionType = "update"
	PermissionTypeDelete   PermissionType = "delete"
	PermissionTypeList     PermissionType = "list"
	PermissionTypeWildcard PermissionType = "*"
)

var (
	PermissionTypes = []PermissionType{
		PermissionTypeCreate,
		PermissionTypeRead,
		PermissionTypeUpdate,
		PermissionTypeDelete,
		PermissionTypeList,
		PermissionTypeWildcard,
	}
)

type OrganisationRoleBinding struct {
	Identity         string            `json:"identity"`
	Name             string            `json:"name"`
	Slug             string            `json:"slug"`
	Description      string            `json:"description"`
	CreatedAt        time.Time         `json:"createdAt"`
	UpdatedAt        time.Time         `json:"updatedAt"`
	ObjectVersion    int               `json:"objectVersion"`
	Labels           map[string]string `json:"labels"`
	Annotations      map[string]string `json:"annotations"`
	OrganisationRole *OrganisationRole `json:"organisationRole,omitempty"`
	AppUser          *base.AppUser     `json:"user,omitempty"`
	OrganisationTeam *Team             `json:"team,omitempty"`
	ServiceAccount   *ServiceAccount   `json:"serviceAccount,omitempty"`
}
