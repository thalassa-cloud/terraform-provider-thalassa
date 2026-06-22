package kubernetes

import (
	"context"
	"fmt"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

const (
	KubernetesClusterRoleEndpoint = "/v1/kubernetes/iam/roles"
)

type ListKubernetesClusterRolesRequest struct {
	Filters []filters.Filter
}

// ListKubernetesClusterRoles lists all KubernetesClusterRoles for a given organisation.
func (c *Client) ListKubernetesClusterRoles(ctx context.Context, request *ListKubernetesClusterRolesRequest) ([]KubernetesClusterRole, error) {
	roles := []KubernetesClusterRole{}
	req := c.R().SetResult(&roles)
	if request != nil {
		for _, filter := range request.Filters {
			for k, v := range filter.ToParams() {
				req.SetQueryParam(k, v)
			}
		}
	}

	resp, err := c.Do(ctx, req, client.GET, KubernetesClusterRoleEndpoint)
	if err != nil {
		return nil, err
	}

	if err := c.Check(resp); err != nil {
		return roles, err
	}
	return roles, nil
}

// CreateKubernetesClusterRole creates a new KubernetesClusterRole.
func (c *Client) CreateKubernetesClusterRole(ctx context.Context, create CreateKubernetesClusterRoleRequest) (*KubernetesClusterRole, error) {
	var role *KubernetesClusterRole
	req := c.R().
		SetBody(create).SetResult(&role)

	resp, err := c.Do(ctx, req, client.POST, KubernetesClusterRoleEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return role, err
	}
	return role, nil
}

// GetKubernetesClusterRole retrieves a specific KubernetesClusterRole by its identity.
func (c *Client) GetKubernetesClusterRole(ctx context.Context, identity string) (*KubernetesClusterRole, error) {
	var role *KubernetesClusterRole
	req := c.R().SetResult(&role)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", KubernetesClusterRoleEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return role, err
	}
	return role, nil
}

// DeleteClusterRole deletes a specific KubernetesClusterRole by its identity.
func (c *Client) DeleteClusterRole(ctx context.Context, identity string) error {
	req := c.R()

	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s", KubernetesClusterRoleEndpoint, identity))
	if err != nil {
		return err
	}
	if err := c.Check(resp); err != nil {
		return err
	}
	return nil
}

// AddClusterRoleRule adds a permission rule to a KubernetesClusterRole.
func (c *Client) AddClusterRoleRule(ctx context.Context, identity string, rule AddKubernetesClusterRolePermissionRule) (*KubernetesClusterRolePermissionRule, error) {
	var permissionRule *KubernetesClusterRolePermissionRule
	req := c.R().
		SetBody(rule).SetResult(&permissionRule)

	resp, err := c.Do(ctx, req, client.POST, fmt.Sprintf("%s/%s/rules", KubernetesClusterRoleEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return permissionRule, err
	}
	return permissionRule, nil
}

// DeleteClusterRoleRule deletes a permission rule from a KubernetesClusterRole.
func (c *Client) DeleteClusterRoleRule(ctx context.Context, identity string, ruleIdentity string) error {
	req := c.R()

	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s/rules/%s", KubernetesClusterRoleEndpoint, identity, ruleIdentity))
	if err != nil {
		return err
	}
	if err := c.Check(resp); err != nil {
		return err
	}
	return nil
}

// ListClusterRoleBindings lists all bindings for a specific KubernetesClusterRole.
func (c *Client) ListClusterRoleBindings(ctx context.Context, identity string) ([]KubernetesClusterRoleBinding, error) {
	bindings := []KubernetesClusterRoleBinding{}
	req := c.R().SetResult(&bindings)

	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s/bindings", KubernetesClusterRoleEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return bindings, err
	}
	return bindings, nil
}

// CreateClusterRoleBinding creates a new binding for a KubernetesClusterRole.
func (c *Client) CreateClusterRoleBinding(ctx context.Context, identity string, create CreateKubernetesClusterRoleBinding) (*KubernetesClusterRoleBinding, error) {
	var binding *KubernetesClusterRoleBinding
	req := c.R().
		SetBody(create).SetResult(&binding)

	resp, err := c.Do(ctx, req, client.POST, fmt.Sprintf("%s/%s/bindings", KubernetesClusterRoleEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return binding, err
	}
	return binding, nil
}

// DeleteClusterRoleBinding deletes a specific binding from a KubernetesClusterRole.
func (c *Client) DeleteClusterRoleBinding(ctx context.Context, identity string, roleBindingIdentity string) error {
	req := c.R()

	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s/bindings/%s", KubernetesClusterRoleEndpoint, identity, roleBindingIdentity))
	if err != nil {
		return err
	}
	if err := c.Check(resp); err != nil {
		return err
	}
	return nil
}
