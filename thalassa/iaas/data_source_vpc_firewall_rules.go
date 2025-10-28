package iaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	iaas "github.com/thalassa-cloud/client-go/iaas"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func DataSourceVpcFirewallRules() *schema.Resource {
	return &schema.Resource{
		Description: "Get all VPC firewall rules for a VPC",
		ReadContext: dataSourceVpcFirewallRulesRead,
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Identity of the VPC to get firewall rules for",
			},
			"firewall_rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of VPC firewall rules",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identity of the VPC firewall rule",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the VPC firewall rule",
						},
						"protocols": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Protocols that the firewall rule applies to",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tcp": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether the rule applies to TCP protocol",
									},
									"udp": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether the rule applies to UDP protocol",
									},
									"icmp": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether the rule applies to ICMP protocol",
									},
									"any": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether the rule applies to any protocol",
									},
								},
							},
						},
						"source": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Source CIDR of the firewall rule",
						},
						"source_ports": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Source ports of the firewall rule",
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"destination": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Destination CIDR of the firewall rule",
						},
						"destination_ports": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Destination ports of the firewall rule",
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"action": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Action of the firewall rule (allow or drop)",
						},
						"priority": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Priority of the firewall rule",
						},
						"source_subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identity of the source subnet",
						},
						"destination_subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identity of the destination subnet",
						},
						"interface_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identity of the interface",
						},
						"direction": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Direction of the firewall rule (inbound or outbound)",
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "State of the firewall rule",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Time when the firewall rule was created",
						},
					},
				},
			},
		},
	}
}

func dataSourceVpcFirewallRulesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := provider.GetProvider(m)
	vpcIdentity := d.Get("vpc_id").(string)

	firewallRules, err := provider.Client.IaaS().ListVpcFirewallRule(ctx, vpcIdentity, &iaas.ListVpcFirewallRulesRequest{})
	if err != nil {
		return diag.FromErr(err)
	}

	// Convert firewall rules to the expected format
	rules := make([]map[string]interface{}, len(firewallRules))
	for i, rule := range firewallRules {
		ruleMap := map[string]interface{}{
			"id":         rule.Identity,
			"name":       rule.Name,
			"action":     rule.Action,
			"priority":   rule.Priority,
			"direction":  rule.Direction,
			"state":      rule.State,
			"created_at": rule.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}

		// Set protocols
		protocols := []map[string]interface{}{
			{
				"tcp":  rule.Protocols.TCP,
				"udp":  rule.Protocols.UDP,
				"icmp": rule.Protocols.ICMP,
				"any":  rule.Protocols.Any,
			},
		}
		ruleMap["protocols"] = protocols

		// Set optional fields
		if rule.Source != nil {
			ruleMap["source"] = *rule.Source
		}
		if rule.Destination != nil {
			ruleMap["destination"] = *rule.Destination
		}
		if rule.SourceSubnet != nil {
			ruleMap["source_subnet_id"] = rule.SourceSubnet.Identity
		}
		if rule.DestinationSubnet != nil {
			ruleMap["destination_subnet_id"] = rule.DestinationSubnet.Identity
		}
		if rule.Interface != nil {
			ruleMap["interface_id"] = rule.Interface.Identity
		}

		// Set source ports
		if rule.SourcePorts != nil {
			sourcePorts := make([]int, len(rule.SourcePorts))
			for j, port := range rule.SourcePorts {
				sourcePorts[j] = int(port)
			}
			ruleMap["source_ports"] = sourcePorts
		}

		// Set destination ports
		if rule.DestinationPorts != nil {
			destinationPorts := make([]int, len(rule.DestinationPorts))
			for j, port := range rule.DestinationPorts {
				destinationPorts[j] = int(port)
			}
			ruleMap["destination_ports"] = destinationPorts
		}

		rules[i] = ruleMap
	}

	d.SetId(fmt.Sprintf("vpc-%s-firewall-rules", vpcIdentity))
	d.Set("firewall_rules", rules)

	return nil
}
