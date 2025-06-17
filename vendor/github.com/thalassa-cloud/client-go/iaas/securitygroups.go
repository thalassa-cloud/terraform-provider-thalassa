package iaas

import (
	"context"
	"fmt"
	"time"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/base"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

const (
	SecurityGroupEndpoint = "/v1/security-groups"
)

type CreateSecurityGroupRequest struct {
	// Name is the name of the security group. Must be between 1 and 16 characters and contain only ASCII characters.
	Name string `json:"name"`
	// Description is the description of the security group
	Description string `json:"description"`
	// Labels are the labels of the security group
	Labels Labels `json:"labels"`
	// Annotations are the annotations of the security group
	Annotations Annotations `json:"annotations"`
	// VpcIdentity is the identity of the VPC that the security group belongs to
	VpcIdentity string `json:"vpcIdentity"`
	// AllowSameGroupTraffic is a flag that indicates if the security group allows traffic between instances in the same security group
	AllowSameGroupTraffic bool `json:"allowSameGroupTraffic"`
	// IngressRules are the ingress rules of the security group
	IngressRules []SecurityGroupRule `json:"ingressRules"`
	// EgressRules are the egress rules of the security group
	EgressRules []SecurityGroupRule `json:"egressRules"`
}

type UpdateSecurityGroupRequest struct {
	// Name is the name of the security group
	Name string `json:"name"`
	// Description is the description of the security group
	Description string `json:"description"`
	// Labels are the labels of the security group
	Labels Labels `json:"labels"`
	// Annotations are the annotations of the security group
	Annotations Annotations `json:"annotations"`
	// ObjectVersion is the version of the security group
	ObjectVersion int `json:"objectVersion"`
	// AllowSameGroupTraffic is a flag that indicates if the security group allows traffic between instances in the same security group
	AllowSameGroupTraffic bool `json:"allowSameGroupTraffic"`
	// IngressRules are the ingress rules of the security group
	IngressRules []SecurityGroupRule `json:"ingressRules"`
	// EgressRules are the egress rules of the security group
	EgressRules []SecurityGroupRule `json:"egressRules"`
}

type SecurityGroupStatus string

const (
	SecurityGroupStatusProvisioning SecurityGroupStatus = "provisioning"
	SecurityGroupStatusActive       SecurityGroupStatus = "active"
	SecurityGroupStatusReady        SecurityGroupStatus = "ready"
	SecurityGroupStatusDeleting     SecurityGroupStatus = "deleting"
	SecurityGroupStatusError        SecurityGroupStatus = "error"
)

type SecurityGroup struct {
	Identity      string    `json:"identity"`
	Name          string    `json:"name"`
	Slug          string    `json:"slug"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	ObjectVersion int       `json:"objectVersion"`

	Labels      Labels      `json:"labels"`
	Annotations Annotations `json:"annotations"`
	// Organisation is the organisation that the security group belongs to
	Organisation *base.Organisation `json:"organisation,omitempty"`
	// Vpc is the VPC that the security group belongs to
	Vpc *Vpc `json:"vpc,omitempty"`

	// Status is the status of the security group
	Status SecurityGroupStatus `json:"status,omitempty"`
	// AllowSameGroupTraffic is a flag that indicates if the security group allows traffic between instances in the same security group
	AllowSameGroupTraffic bool `json:"allowSameGroupTraffic"`
	// IngressRules are the ingress rules of the security group
	IngressRules []SecurityGroupRule `json:"ingressRules"`
	// EgressRules are the egress rules of the security group
	EgressRules []SecurityGroupRule `json:"egressRules"`
}

type SecurityGroupRule struct {
	// Name is the name of the rule
	Name string `json:"name"`
	// IPVersion is the IP version of the rule
	IPVersion SecurityGroupIPVersion `json:"ipVersion"`
	// Protocol is the protocol of the rule
	Protocol SecurityGroupRuleProtocol `json:"protocol"`
	// Priority is the priority of the rule. Must be greater than 0 and less than 200.
	Priority int32 `json:"priority"`
	// RemoteType is the type of the remote address
	RemoteType SecurityGroupRuleRemoteType `json:"remoteType"`
	// RemoteAddress is the IP address or CIDR block that the rule applies to
	RemoteAddress *string `json:"remoteAddress"`
	// RemoteSecurityGroupIdentity is the identity of the security group that the rule applies to
	RemoteSecurityGroupIdentity *string `json:"remoteSecurityGroupIdentity"`
	// PortRangeMin is the minimum port of the rule. Must be greater than 0 and less than 65535.
	PortRangeMin int32 `json:"portRangeMin"`
	// PortRangeMax is the maximum port of the rule. Must be greater than 0 and less than 65535.
	PortRangeMax int32 `json:"portRangeMax"`
	// Policy is the policy of the rule
	Policy SecurityGroupRulePolicy `json:"policy"`
}

type SecurityGroupRuleProtocol string

const (
	SecurityGroupRuleProtocolAll  SecurityGroupRuleProtocol = "all"
	SecurityGroupRuleProtocolTCP  SecurityGroupRuleProtocol = "tcp"
	SecurityGroupRuleProtocolUDP  SecurityGroupRuleProtocol = "udp"
	SecurityGroupRuleProtocolICMP SecurityGroupRuleProtocol = "icmp"
)

type SecurityGroupRulePolicy string

const (
	SecurityGroupRulePolicyAllow SecurityGroupRulePolicy = "allow"
	SecurityGroupRulePolicyDrop  SecurityGroupRulePolicy = "drop"
)

type SecurityGroupRuleRemoteType string

const (
	SecurityGroupRuleRemoteTypeAddress       SecurityGroupRuleRemoteType = "address"
	SecurityGroupRuleRemoteTypeSecurityGroup SecurityGroupRuleRemoteType = "securityGroup"
)

type SecurityGroupIPVersion string

const (
	SecurityGroupIPVersionIPv4 SecurityGroupIPVersion = "ipv4"
	SecurityGroupIPVersionIPv6 SecurityGroupIPVersion = "ipv6"
)

type ListSecurityGroupsRequest struct {
	Filters []filters.Filter
}

// ListSecurityGroups lists all security groups for a given organisation.
func (c *Client) ListSecurityGroups(ctx context.Context, listRequest *ListSecurityGroupsRequest) ([]SecurityGroup, error) {
	securityGroups := []SecurityGroup{}
	req := c.R().SetResult(&securityGroups)

	if listRequest != nil {
		for _, filter := range listRequest.Filters {
			for k, v := range filter.ToParams() {
				req = req.SetQueryParam(k, v)
			}
		}
	}

	resp, err := c.Do(ctx, req, client.GET, SecurityGroupEndpoint)
	if err != nil {
		return nil, err
	}

	if err := c.Check(resp); err != nil {
		return securityGroups, err
	}
	return securityGroups, nil
}

// GetSecurityGroup retrieves a specific security group by its identity.
func (c *Client) GetSecurityGroup(ctx context.Context, identity string) (*SecurityGroup, error) {
	var securityGroup *SecurityGroup
	req := c.R().SetResult(&securityGroup)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", SecurityGroupEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return securityGroup, err
	}
	return securityGroup, nil
}

// CreateSecurityGroup creates a new security group.
func (c *Client) CreateSecurityGroup(ctx context.Context, create CreateSecurityGroupRequest) (*SecurityGroup, error) {
	var securityGroup *SecurityGroup
	req := c.R().
		SetBody(create).SetResult(&securityGroup)

	resp, err := c.Do(ctx, req, client.POST, SecurityGroupEndpoint)
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return securityGroup, err
	}
	return securityGroup, nil
}

// UpdateSecurityGroup updates an existing security group.
func (c *Client) UpdateSecurityGroup(ctx context.Context, identity string, update UpdateSecurityGroupRequest) (*SecurityGroup, error) {
	var securityGroup *SecurityGroup
	req := c.R().
		SetBody(update).SetResult(&securityGroup)

	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s", SecurityGroupEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return securityGroup, err
	}
	return securityGroup, nil
}

// DeleteSecurityGroup deletes a specific security group by its identity.
func (c *Client) DeleteSecurityGroup(ctx context.Context, identity string) error {
	req := c.R()

	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s", SecurityGroupEndpoint, identity))
	if err != nil {
		return err
	}
	if err := c.Check(resp); err != nil {
		return err
	}
	return nil
}
