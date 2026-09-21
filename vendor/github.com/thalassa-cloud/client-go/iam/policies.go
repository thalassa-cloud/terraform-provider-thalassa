package iam

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/base"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

const (
	IamPolicyEndpoint = "/v1/projects/iam/policies"

	ProjectMemberSourceIamPolicy        = "iam_policy"
	ProjectMemberSourceOrganisationRole = "organisation_role"
)

// ListIamPoliciesRequest configures filters for GET /v1/projects/iam/policies.
type ListIamPoliciesRequest struct {
	Filters []filters.Filter
}

// ListIamPolicyBindingsRequest configures filters for listing bindings on a policy.
type ListIamPolicyBindingsRequest struct {
	Filters []filters.Filter
}

// PolicyConditionals defines conditional allow/deny rules for an IAM policy.
type PolicyConditionals struct {
	Allowed []PolicyCondition `json:"allowed,omitempty"`
	Denied  []PolicyCondition `json:"denied,omitempty"`
}

// PolicyCondition is a single ip or time condition.
type PolicyCondition struct {
	Type string         `json:"type"`
	IPs  []string       `json:"ips,omitempty"`
	Time *TimeCondition `json:"time,omitempty"`
}

// TimeCondition restricts access to a daily hour window in a timezone.
type TimeCondition struct {
	StartHour int    `json:"startHour"`
	EndHour   int    `json:"endHour"`
	Timezone  string `json:"timezone,omitempty"`
}

// IamPolicyProject is the nested project reference on IAM policy responses.
type IamPolicyProject struct {
	Identity string `json:"identity"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
}

// IamPolicy is a project-scoped (or organisation-root) IAM policy.
type IamPolicy struct {
	Identity      string            `json:"identity"`
	Name          string            `json:"name"`
	Slug          string            `json:"slug"`
	Description   string            `json:"description"`
	Labels        map[string]string `json:"labels"`
	Annotations   map[string]string `json:"annotations"`
	CreatedAt     time.Time         `json:"createdAt"`
	UpdatedAt     *time.Time        `json:"updatedAt,omitempty"`
	ObjectVersion int64             `json:"objectVersion"`

	IsReadOnly          bool `json:"isReadOnly,omitempty"`
	ReplicateToChildren bool `json:"replicateToChildren,omitempty"`
	System              bool `json:"system,omitempty"`

	SourceIamPolicy *IamPolicy          `json:"sourceIamPolicy,omitempty"`
	Organisation    *base.Organisation  `json:"organisation,omitempty"`
	Project         *IamPolicyProject   `json:"project,omitempty"`
	Conditionals    *PolicyConditionals `json:"conditionals,omitempty"`

	Rules    []IamPolicyPermissionRule `json:"rules,omitempty"`
	Bindings []IamPolicyBinding        `json:"bindings,omitempty"`
}

// IamPolicyPermissionRule grants permissions on resources within a policy.
type IamPolicyPermissionRule struct {
	Identity           string           `json:"identity"`
	IamPolicy          *IamPolicy       `json:"iamPolicy,omitempty"`
	Resources          []string         `json:"resources"`
	ResourceIdentities []string         `json:"resourceIdentities,omitempty"`
	Permissions        []PermissionType `json:"permissions"`
	Note               string           `json:"note,omitempty"`
}

// IamPolicyBinding binds a principal to an IAM policy.
type IamPolicyBinding struct {
	Identity      string            `json:"identity"`
	Name          string            `json:"name"`
	Slug          string            `json:"slug"`
	Description   string            `json:"description"`
	Labels        map[string]string `json:"labels"`
	Annotations   map[string]string `json:"annotations"`
	CreatedAt     time.Time         `json:"createdAt"`
	UpdatedAt     *time.Time        `json:"updatedAt,omitempty"`
	ObjectVersion int64             `json:"objectVersion"`

	IamPolicy              *IamPolicy        `json:"iamPolicy,omitempty"`
	ProjectId              string            `json:"projectId,omitempty"`
	Project                *IamPolicyProject `json:"project,omitempty"`
	User                   *base.AppUser     `json:"user,omitempty"`
	ServiceAccount         *ServiceAccount   `json:"serviceAccount,omitempty"`
	SourceIamPolicyBinding *IamPolicyBinding `json:"sourceIamPolicyBinding,omitempty"`
	ExpiresAt              *time.Time        `json:"expiresAt,omitempty"`
}

// CreateIamPolicyRequest is the body for POST /v1/projects/iam/policies.
type CreateIamPolicyRequest struct {
	Name                string              `json:"name"`
	Description         string              `json:"description"`
	Annotations         map[string]string   `json:"annotations"`
	Labels              map[string]string   `json:"labels"`
	ReplicateToChildren bool                `json:"replicateToChildren,omitempty"`
	Conditionals        *PolicyConditionals `json:"conditionals,omitempty"`
}

// UpdateIamPolicyRequest is the body for PUT /v1/projects/iam/policies/{identity}.
type UpdateIamPolicyRequest struct {
	Description  string              `json:"description"`
	Annotations  map[string]string   `json:"annotations"`
	Labels       map[string]string   `json:"labels"`
	Conditionals *PolicyConditionals `json:"conditionals,omitempty"`
}

// AddIamPolicyRuleRequest is the body for POST /v1/projects/iam/policies/{identity}/rules.
type AddIamPolicyRuleRequest struct {
	Resources          []string         `json:"resources"`
	ResourceIdentities []string         `json:"resourceIdentities,omitempty"`
	Permissions        []PermissionType `json:"permissions"`
	Note               string           `json:"note,omitempty"`
}

// CreateIamPolicyBindingRequest is the body for POST .../bindings.
type CreateIamPolicyBindingRequest struct {
	Name                   string            `json:"name"`
	Description            string            `json:"description"`
	Annotations            map[string]string `json:"annotations"`
	Labels                 map[string]string `json:"labels"`
	UserIdentity           *string           `json:"userIdentity"`
	ServiceAccountIdentity *string           `json:"serviceAccountIdentity"`
	ExpiresAt              *time.Time        `json:"expiresAt,omitempty"`
	Duration               string            `json:"duration,omitempty"`
}

// UpdateIamPolicyBindingRequest is the body for PUT .../bindings/{bindingIdentity}.
type UpdateIamPolicyBindingRequest struct {
	Description string            `json:"description"`
	Annotations map[string]string `json:"annotations"`
	Labels      map[string]string `json:"labels"`
}

// ProjectMember is a user with one or more active access bindings in the IAM scope.
type ProjectMember struct {
	User     *base.AppUser               `json:"user"`
	Bindings []ProjectMemberPolicyAccess `json:"bindings"`
}

// ProjectMemberPolicyAccess summarises a single policy or organisation-role binding.
type ProjectMemberPolicyAccess struct {
	Source          string     `json:"source"`
	BindingIdentity string     `json:"bindingIdentity"`
	BindingName     string     `json:"bindingName,omitempty"`
	ExpiresAt       *time.Time `json:"expiresAt,omitempty"`
	PolicyIdentity  string     `json:"policyIdentity"`
	PolicyName      string     `json:"policyName,omitempty"`
	PolicySlug      string     `json:"policySlug,omitempty"`
}

// ListResourceTypes lists API resource types usable in IAM policy rules.
func (c *Client) ListResourceTypes(ctx context.Context) ([]string, error) {
	resourceTypes := []string{}
	req := c.R().SetResult(&resourceTypes)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/resources", IamPolicyEndpoint))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return nil, err
	}
	return resourceTypes, nil
}

// ListProjectMembers lists users with active IAM (and at root, organisation-role) bindings.
func (c *Client) ListProjectMembers(ctx context.Context) ([]ProjectMember, error) {
	members := []ProjectMember{}
	req := c.R().SetResult(&members)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/members", IamPolicyEndpoint))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return nil, err
	}
	return members, nil
}

// ListIamPolicies lists IAM policies in the current organisation/project scope.
func (c *Client) ListIamPolicies(ctx context.Context, request *ListIamPoliciesRequest) ([]IamPolicy, error) {
	policies := []IamPolicy{}
	req := c.R().SetResult(&policies)
	if request != nil {
		for _, filter := range request.Filters {
			for k, v := range filter.ToParams() {
				req.SetQueryParam(k, v)
			}
		}
	}

	resp, err := c.Do(ctx, req, client.GET, IamPolicyEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return nil, err
	}
	return policies, nil
}

// GetIamPolicy retrieves an IAM policy by identity.
func (c *Client) GetIamPolicy(ctx context.Context, identity string) (*IamPolicy, error) {
	if strings.TrimSpace(identity) == "" {
		return nil, fmt.Errorf("identity is required")
	}

	policy := IamPolicy{}
	req := c.R().SetResult(&policy)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", IamPolicyEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return nil, err
	}
	return &policy, nil
}

// CreateIamPolicy creates an IAM policy in the current scope.
func (c *Client) CreateIamPolicy(ctx context.Context, create CreateIamPolicyRequest) (*IamPolicy, error) {
	if strings.TrimSpace(create.Name) == "" {
		return nil, fmt.Errorf("name is required")
	}

	policy := IamPolicy{}
	req := c.R().SetBody(create).SetResult(&policy)
	resp, err := c.Do(ctx, req, client.POST, IamPolicyEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return nil, err
	}
	return &policy, nil
}

// UpdateIamPolicy updates an IAM policy by identity.
func (c *Client) UpdateIamPolicy(ctx context.Context, identity string, update UpdateIamPolicyRequest) (*IamPolicy, error) {
	if strings.TrimSpace(identity) == "" {
		return nil, fmt.Errorf("identity is required")
	}

	policy := IamPolicy{}
	req := c.R().SetBody(update).SetResult(&policy)
	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s", IamPolicyEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return nil, err
	}
	return &policy, nil
}

// DeleteIamPolicy deletes an IAM policy by identity.
func (c *Client) DeleteIamPolicy(ctx context.Context, identity string) error {
	if strings.TrimSpace(identity) == "" {
		return fmt.Errorf("identity is required")
	}

	resp, err := c.Do(ctx, c.R(), client.DELETE, fmt.Sprintf("%s/%s", IamPolicyEndpoint, identity))
	if err != nil {
		return err
	}
	return c.Check(resp)
}

// AddIamPolicyRule adds a permission rule to an IAM policy.
func (c *Client) AddIamPolicyRule(ctx context.Context, policyIdentity string, rule AddIamPolicyRuleRequest) (*IamPolicyPermissionRule, error) {
	if strings.TrimSpace(policyIdentity) == "" {
		return nil, fmt.Errorf("policy identity is required")
	}
	if len(rule.Resources) == 0 {
		return nil, fmt.Errorf("resources is required")
	}
	if len(rule.Permissions) == 0 {
		return nil, fmt.Errorf("permissions is required")
	}

	result := IamPolicyPermissionRule{}
	req := c.R().SetBody(rule).SetResult(&result)
	resp, err := c.Do(ctx, req, client.POST, fmt.Sprintf("%s/%s/rules", IamPolicyEndpoint, policyIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteIamPolicyRule deletes a permission rule from an IAM policy.
func (c *Client) DeleteIamPolicyRule(ctx context.Context, policyIdentity, ruleIdentity string) error {
	if strings.TrimSpace(policyIdentity) == "" {
		return fmt.Errorf("policy identity is required")
	}
	if strings.TrimSpace(ruleIdentity) == "" {
		return fmt.Errorf("rule identity is required")
	}

	resp, err := c.Do(ctx, c.R(), client.DELETE, fmt.Sprintf("%s/%s/rules/%s", IamPolicyEndpoint, policyIdentity, ruleIdentity))
	if err != nil {
		return err
	}
	return c.Check(resp)
}

// ListIamPolicyBindings lists bindings for an IAM policy.
func (c *Client) ListIamPolicyBindings(ctx context.Context, policyIdentity string, request *ListIamPolicyBindingsRequest) ([]IamPolicyBinding, error) {
	if strings.TrimSpace(policyIdentity) == "" {
		return nil, fmt.Errorf("policy identity is required")
	}

	bindings := []IamPolicyBinding{}
	req := c.R().SetResult(&bindings)
	if request != nil {
		for _, filter := range request.Filters {
			for k, v := range filter.ToParams() {
				req.SetQueryParam(k, v)
			}
		}
	}

	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s/bindings", IamPolicyEndpoint, policyIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return nil, err
	}
	return bindings, nil
}

// CreateIamPolicyBinding creates a binding on an IAM policy.
func (c *Client) CreateIamPolicyBinding(ctx context.Context, policyIdentity string, create CreateIamPolicyBindingRequest) (*IamPolicyBinding, error) {
	if strings.TrimSpace(policyIdentity) == "" {
		return nil, fmt.Errorf("policy identity is required")
	}
	if strings.TrimSpace(create.Name) == "" {
		return nil, fmt.Errorf("name is required")
	}
	hasUser := create.UserIdentity != nil && strings.TrimSpace(*create.UserIdentity) != ""
	hasSA := create.ServiceAccountIdentity != nil && strings.TrimSpace(*create.ServiceAccountIdentity) != ""
	if hasUser == hasSA {
		return nil, fmt.Errorf("must bind to exactly one of userIdentity or serviceAccountIdentity")
	}

	binding := IamPolicyBinding{}
	req := c.R().SetBody(create).SetResult(&binding)
	resp, err := c.Do(ctx, req, client.POST, fmt.Sprintf("%s/%s/bindings", IamPolicyEndpoint, policyIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return nil, err
	}
	return &binding, nil
}

// UpdateIamPolicyBinding updates a binding on an IAM policy.
func (c *Client) UpdateIamPolicyBinding(ctx context.Context, policyIdentity, bindingIdentity string, update UpdateIamPolicyBindingRequest) (*IamPolicyBinding, error) {
	if strings.TrimSpace(policyIdentity) == "" {
		return nil, fmt.Errorf("policy identity is required")
	}
	if strings.TrimSpace(bindingIdentity) == "" {
		return nil, fmt.Errorf("binding identity is required")
	}

	binding := IamPolicyBinding{}
	req := c.R().SetBody(update).SetResult(&binding)
	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s/bindings/%s", IamPolicyEndpoint, policyIdentity, bindingIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return nil, err
	}
	return &binding, nil
}

// DeleteIamPolicyBinding deletes a binding from an IAM policy.
func (c *Client) DeleteIamPolicyBinding(ctx context.Context, policyIdentity, bindingIdentity string) error {
	if strings.TrimSpace(policyIdentity) == "" {
		return fmt.Errorf("policy identity is required")
	}
	if strings.TrimSpace(bindingIdentity) == "" {
		return fmt.Errorf("binding identity is required")
	}

	resp, err := c.Do(ctx, c.R(), client.DELETE, fmt.Sprintf("%s/%s/bindings/%s", IamPolicyEndpoint, policyIdentity, bindingIdentity))
	if err != nil {
		return err
	}
	return c.Check(resp)
}
