package iaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	iaas "github.com/thalassa-cloud/client-go/iaas"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func DataSourceVpcFirewallRule() *schema.Resource {
	return &schema.Resource{
		Description: "Get a VPC firewall rule",
		ReadContext: dataSourceVpcFirewallRuleRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity of the VPC firewall rule",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Identity of the VPC that the firewall rule belongs to",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.StringLenBetween(1, 16),
				Description:  "Name of the VPC firewall rule",
			},
			"identity": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Identity of the VPC firewall rule",
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
	}
}

func dataSourceVpcFirewallRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := provider.GetProvider(m)
	vpcIdentity := d.Get("vpc_id").(string)
	name := d.Get("name").(string)
	identity := d.Get("identity").(string)

	// If identity is provided, get the specific rule
	if identity != "" {
		firewallRule, err := provider.Client.IaaS().GetVpcFirewallRule(ctx, vpcIdentity, identity)
		if err != nil {
			return diag.FromErr(err)
		}
		return setVpcFirewallRuleData(d, firewallRule)
	}

	// If name is provided, list rules and find by name
	if name != "" {
		firewallRules, err := provider.Client.IaaS().ListVpcFirewallRule(ctx, vpcIdentity, &iaas.ListVpcFirewallRulesRequest{})
		if err != nil {
			return diag.FromErr(err)
		}

		var firewallRule *iaas.VpcFirewallRule
		for _, rule := range firewallRules {
			if rule.Name == name {
				firewallRule = &rule
				break
			}
		}

		if firewallRule == nil {
			return diag.FromErr(fmt.Errorf("VPC firewall rule with name '%s' not found", name))
		}

		return setVpcFirewallRuleData(d, firewallRule)
	}

	return diag.FromErr(fmt.Errorf("either 'identity' or 'name' must be provided"))
}
