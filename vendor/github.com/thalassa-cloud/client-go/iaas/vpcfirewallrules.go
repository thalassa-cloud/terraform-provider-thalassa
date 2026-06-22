package iaas

import (
	"context"
	"fmt"
	"time"

	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/pkg/base"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

type VpcFirewallRuleProtocols struct {
	// TCP is a flag that indicates that the rule applies to the TCP protocol
	TCP bool `json:"tcp"`
	// UDP is a flag that indicates that the rule applies to the UDP protocol
	UDP bool `json:"udp"`
	// ICMP is a flag that indicates that the rule applies to the ICMP protocol
	ICMP bool `json:"icmp"`
	// Any is a flag that indicates that the rule applies to any protocol. Matches all protocols if set to true
	Any bool `json:"any"`
}

// VpcFirewallRuleDirection is the direction of the firewall rule
type VpcFirewallRuleDirection string

const (
	VpcFirewallRuleDirectionInbound  VpcFirewallRuleDirection = "inbound"
	VpcFirewallRuleDirectionOutbound VpcFirewallRuleDirection = "outbound"
)

type VpcFirewallObject struct {
	SubnetIdentity *string `json:"subnetIdentity"`
	VpcIdentity    *string `json:"vpcIdentity"`
}

// FirewallRuleState is the state of the firewall rule
type FirewallRuleState string

const (
	FirewallRuleStateActive   FirewallRuleState = "active"
	FirewallRuleStateInactive FirewallRuleState = "inactive"
	FirewallRuleStateDeleted  FirewallRuleState = "deleted"
)

type CreateVpcFirewallRuleRequest struct {
	// Name of the firewall rule. Must be between 1 and 16 characters and contain only ASCII characters.
	Name string `json:"name"`
	// VpcIdentity is the identity of the VPC that the firewall rule belongs to
	VpcIdentity string `json:"vpcIdentity"`

	// Protocols determines the protocols that match the firewall rule
	Protocols VpcFirewallRuleProtocols `json:"protocols"`
	// Source CIDR of the firewall rule
	Source *string `json:"source"`

	// Source ports of the firewall rule. Must be between 0 and 65535
	SourcePorts []int32 `json:"sourcePorts"`
	// Destination CIDR of the firewall rule
	Destination *string `json:"destination"`
	// Destination ports of the firewall rule. Must be between 0 and 65535
	DestinationPorts []int32 `json:"destinationPorts"`
	// Action of the firewall rule. One of allow, drop
	Action FirewallRuleAction `json:"action"`
	// Priority of the firewall rule. Must be between 1 and 1000
	Priority *int32 `json:"priority"`
	// SourceSubnetIdentity is the identity of the source subnet
	SourceSubnetIdentity *string `json:"sourceSubnetIdentity,omitempty"`
	// DestinationSubnetIdentity is the identity of the destination subnet
	DestinationSubnetIdentity *string `json:"destinationSubnetIdentity,omitempty"`
	// InterfaceIdentity is the identity of the interface. Leaving empty will apply to all interfaces
	InterfaceIdentity *string `json:"interfaceIdentity,omitempty"`
	// Direction is the direction of the firewall rule
	Direction VpcFirewallRuleDirection `json:"direction"`
	// State is the state of the firewall rule
	State FirewallRuleState `json:"state"`
}

type FirewallRuleAction string

const (
	FirewallRuleActionAllow FirewallRuleAction = "allow"
	FirewallRuleActionDrop  FirewallRuleAction = "drop"
)

type UpdateVpcFirewallRuleRequest struct {
	Identity string `json:"identity,omitempty"`

	// Name of the firewall rule
	Name string `json:"name"`
	// Protocols determines the protocols that match the firewall rule
	Protocols VpcFirewallRuleProtocols `json:"protocols"`
	// Source CIDR of the firewall rule. Must be a valid CIDR block.
	Source *string `json:"source,omitempty"`
	// Source ports of the firewall rule. Must be between 0 and 65535
	SourcePorts []int32 `json:"sourcePorts,omitempty"`
	// Destination CIDR of the firewall rule. Must be a valid CIDR block.
	Destination *string `json:"destination,omitempty"`
	// Destination ports of the firewall rule. Must be between 0 and 65535
	DestinationPorts []int32 `json:"destinationPorts,omitempty"`
	// Action of the firewall rule. One of allow, drop
	Action FirewallRuleAction `json:"action"`
	// Priority of the firewall rule. Must be between 1 and 1000
	Priority int32 `json:"priority"`
	// InterfaceIdentity is the identity of the interface. Leaving empty will apply to all interfaces
	InterfaceIdentity *string `json:"interfaceIdentity,omitempty"`
	// SourceSubnetIdentity is the identity of the source subnet
	SourceSubnetIdentity *string `json:"sourceSubnetIdentity,omitempty"`
	// DestinationSubnetIdentity is the identity of the destination subnet
	DestinationSubnetIdentity *string `json:"destinationSubnetIdentity,omitempty"`
	// Direction is the direction of the firewall rule
	Direction VpcFirewallRuleDirection `json:"direction"`
	// State is the state of the firewall rule
	State FirewallRuleState `json:"state"`
}

// BulkUpdateVpcFirewallRuleRequest is a request to update multiple firewall rules
type BulkUpdateVpcFirewallRuleRequest struct {
	// FirewallRules is a list of firewall rules to update
	FirewallRules []UpdateVpcFirewallRuleRequest `json:"firewallRules"`
}

type VpcFirewallRule struct {
	// Identity is a unique identifier for the object
	Identity string `json:"identity"`
	// Name is a human-readable name of the firewall rule
	Name string `json:"name"`
	// Organisation is the organisation that the firewall belongs to
	Organisation *base.Organisation `json:"organisation,omitempty"`
	// Vpc is the VPC that the firewall belongs to
	Vpc *Vpc `json:"vpc,omitempty"`
	// CreatedAt is the time the firewall rule was created
	CreatedAt time.Time `json:"createdAt"`
	// Protocols is the protocols that the firewall rule applies to
	Protocols VpcFirewallRuleProtocols `json:"protocols"`
	// Interface is the interface that the firewall rule applies to
	Interface *Subnet `json:"interface,omitempty"`
	// Direction is the direction of the firewall rule
	Direction VpcFirewallRuleDirection `json:"direction"`
	// CIDR block that the rule applies to
	Source      *string `json:"source"`
	Destination *string `json:"destination"`
	// SourceSubnet is the subnet that the firewall rule applies to
	SourceSubnet *Subnet `json:"sourceSubnet,omitempty"`
	// DestinationSubnet is the subnet that the firewall rule applies to
	DestinationSubnet *Subnet `json:"destinationSubnet,omitempty"`
	// Source ports of the firewall rule
	SourcePorts []int32 `json:"sourcePorts"`
	// Destination ports of the firewall rule
	DestinationPorts []int32 `json:"destinationPorts"`
	// Action is the action of the firewall rule
	Action FirewallRuleAction `json:"action"`
	// Priority is the priority of the firewall rule
	Priority int32 `json:"priority"`
	// State is the state of the firewall rule
	State FirewallRuleState `json:"state"`
}

const (
	VpcFirewallRuleEndpoint = "/v1/vpc-firewall-rules"
)

type ListVpcFirewallRulesRequest struct {
	Filters []filters.Filter
}

// ListVpcFirewallRule lists all VPC firewall rules for a given VPC identity.
func (c *Client) ListVpcFirewallRule(ctx context.Context, identity string, request *ListVpcFirewallRulesRequest) ([]VpcFirewallRule, error) {
	firewallRules := []VpcFirewallRule{}
	req := c.R().SetResult(&firewallRules)
	if request != nil {
		for _, filter := range request.Filters {
			for k, v := range filter.ToParams() {
				req = req.SetQueryParam(k, v)
			}
		}
	}

	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s", VpcFirewallRuleEndpoint, identity))
	if err != nil {
		return nil, err
	}

	if err := c.Check(resp); err != nil {
		return firewallRules, err
	}
	return firewallRules, nil
}

// CreateVpcFirewallRule creates a new VPC firewall rule.
func (c *Client) CreateVpcFirewallRule(ctx context.Context, identity string, create CreateVpcFirewallRuleRequest) (*VpcFirewallRule, error) {
	var firewallRule *VpcFirewallRule
	req := c.R().
		SetBody(create).SetResult(&firewallRule)

	resp, err := c.Do(ctx, req, client.POST, fmt.Sprintf("%s/%s", VpcFirewallRuleEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return firewallRule, err
	}
	return firewallRule, nil
}

// BulkUpdateVpcFirewallRule updates multiple VPC firewall rules.
func (c *Client) BulkUpdateVpcFirewallRule(ctx context.Context, identity string, update BulkUpdateVpcFirewallRuleRequest) ([]VpcFirewallRule, error) {
	firewallRules := []VpcFirewallRule{}
	req := c.R().
		SetBody(update).SetResult(&firewallRules)

	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s", VpcFirewallRuleEndpoint, identity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return firewallRules, err
	}
	return firewallRules, nil
}

// GetVpcFirewallRule retrieves a specific VPC firewall rule by its identity.
func (c *Client) GetVpcFirewallRule(ctx context.Context, identity string, firewallRuleIdentity string) (*VpcFirewallRule, error) {
	var firewallRule *VpcFirewallRule
	req := c.R().SetResult(&firewallRule)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s/%s", VpcFirewallRuleEndpoint, identity, firewallRuleIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return firewallRule, err
	}
	return firewallRule, nil
}

// UpdateVpcFirewallRule updates an existing VPC firewall rule.
func (c *Client) UpdateVpcFirewallRule(ctx context.Context, identity string, firewallRuleIdentity string, update UpdateVpcFirewallRuleRequest) (*VpcFirewallRule, error) {
	var firewallRule *VpcFirewallRule
	req := c.R().
		SetBody(update).SetResult(&firewallRule)

	resp, err := c.Do(ctx, req, client.PUT, fmt.Sprintf("%s/%s/%s", VpcFirewallRuleEndpoint, identity, firewallRuleIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return firewallRule, err
	}
	return firewallRule, nil
}

// DeleteVpcFirewallRule deletes a specific VPC firewall rule by its identity.
func (c *Client) DeleteVpcFirewallRule(ctx context.Context, identity string, firewallRuleIdentity string) error {
	req := c.R()

	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s/%s", VpcFirewallRuleEndpoint, identity, firewallRuleIdentity))
	if err != nil {
		return err
	}
	if err := c.Check(resp); err != nil {
		return err
	}
	return nil
}
