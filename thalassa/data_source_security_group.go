package thalassa

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thalassa-cloud/client-go/filters"
	iaas "github.com/thalassa-cloud/client-go/iaas"
)

func dataSourceSecurityGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecurityGroupRead,
		Description: "A security group is a collection of rules that control the traffic to and from a virtual machine instance or other cloud resource within a VPC.",
		Schema: map[string]*schema.Schema{
			"identity": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Identity of the security group",
				ExactlyOneOf: []string{"identity", "name"},
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Name of the security group",
				ExactlyOneOf: []string{"identity", "name"},
			},
			"vpc_identity": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Identity of the VPC that the security group belongs to. Required when searching by name.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the security group",
			},
			"allow_same_group_traffic": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag that indicates if the security group allows traffic between instances in the same security group",
			},
			"ingress_rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of ingress rules for the security group",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the rule",
						},
						"ip_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP version of the rule (ipv4 or ipv6)",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Protocol of the rule (all, tcp, udp, icmp)",
						},
						"priority": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Priority of the rule",
						},
						"remote_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the remote address (address or securityGroup)",
						},
						"remote_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP address or CIDR block that the rule applies to",
						},
						"remote_security_group_identity": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identity of the security group that the rule applies to",
						},
						"port_range_min": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Minimum port of the rule",
						},
						"port_range_max": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum port of the rule",
						},
						"policy": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy of the rule (allow or drop)",
						},
					},
				},
			},
			"egress_rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of egress rules for the security group",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the rule",
						},
						"ip_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP version of the rule (ipv4 or ipv6)",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Protocol of the rule (all, tcp, udp, icmp)",
						},
						"priority": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Priority of the rule",
						},
						"remote_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the remote address (address or securityGroup)",
						},
						"remote_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP address or CIDR block that the rule applies to",
						},
						"remote_security_group_identity": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identity of the security group that the rule applies to",
						},
						"port_range_min": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Minimum port of the rule",
						},
						"port_range_max": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum port of the rule",
						},
						"policy": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy of the rule (allow or drop)",
						},
					},
				},
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation timestamp of the security group",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last update timestamp of the security group",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the security group",
			},
		},
	}
}

func dataSourceSecurityGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getClient(getProvider(meta), d)
	if err != nil {
		return diag.FromErr(err)
	}

	var securityGroup *iaas.SecurityGroup

	if v, ok := d.GetOk("identity"); ok {
		// Look up by identity
		securityGroup, err = client.IaaS().GetSecurityGroup(ctx, v.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	} else if v, ok := d.GetOk("name"); ok {
		// Look up by name
		vpcIdentity, ok := d.GetOk("vpc_identity")
		if !ok {
			return diag.Errorf("vpc_identity is required when searching by name")
		}

		securityGroups, err := client.IaaS().ListSecurityGroups(ctx, &iaas.ListSecurityGroupsRequest{
			Filters: []filters.Filter{
				&filters.FilterKeyValue{
					Key:   filters.FilterVpcIdentity,
					Value: vpcIdentity.(string),
				},
			},
		})
		if err != nil {
			return diag.FromErr(err)
		}

		// Find the security group with the matching name
		for _, sg := range securityGroups {
			if sg.Vpc != nil && sg.Vpc.Identity != vpcIdentity.(string) {
				continue
			}
			if sg.Name == v.(string) {
				securityGroup = &sg
				break
			}
		}

		if securityGroup == nil {
			return diag.Errorf("security group with name %s not found in VPC %s", v.(string), vpcIdentity.(string))
		}
	}

	// Set the ID and other attributes
	d.SetId(securityGroup.Identity)
	if err := d.Set("name", securityGroup.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", securityGroup.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("vpc_identity", securityGroup.Vpc.Identity); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("allow_same_group_traffic", securityGroup.AllowSameGroupTraffic); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ingress_rules", flattenSecurityGroupRules(securityGroup.IngressRules)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("egress_rules", flattenSecurityGroupRules(securityGroup.EgressRules)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("created_at", securityGroup.CreatedAt.Format(time.RFC3339)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("updated_at", securityGroup.UpdatedAt.Format(time.RFC3339)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("status", securityGroup.Status); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
