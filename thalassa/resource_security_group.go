package thalassa

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	iaas "github.com/thalassa-cloud/client-go/iaas"
)

func resourceSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Description:   "A security group is a collection of rules that control the traffic to and from a virtual machine instance or other cloud resource within a VPC.",
		CreateContext: resourceSecurityGroupCreate,
		ReadContext:   resourceSecurityGroupRead,
		UpdateContext: resourceSecurityGroupUpdate,
		DeleteContext: resourceSecurityGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the security group. Must be between 1 and 16 characters and contain only ASCII characters.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the security group",
			},
			"vpc_identity": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Identity of the VPC that the security group belongs to",
			},
			"allow_same_group_traffic": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Flag that indicates if the security group allows traffic between instances in the same security group",
			},
			"ingress_rules": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of ingress rules for the security group",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the rule",
						},
						"ip_version": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.StringInSlice([]string{"ipv4", "ipv6"}, false),
							Description:  "IP version of the rule (ipv4 or ipv6)",
						},
						"protocol": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.StringInSlice([]string{"all", "tcp", "udp", "icmp"}, false),
							Description:  "Protocol of the rule (all, tcp, udp, icmp)",
						},
						"priority": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validate.IntBetween(1, 199),
							Description:  "Priority of the rule. Must be greater than 0 and less than 200.",
						},
						"remote_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.StringInSlice([]string{"address", "securityGroup"}, false),
							Description:  "Type of the remote address (address or securityGroup)",
						},
						"remote_address": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "IP address or CIDR block that the rule applies to",
						},
						"remote_security_group_identity": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Identity of the security group that the rule applies to",
						},
						"port_range_min": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validate.IntBetween(1, 65535),
							Description:  "Minimum port of the rule. Must be greater than 0 and less than 65535.",
						},
						"port_range_max": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validate.IntBetween(1, 65535),
							Description:  "Maximum port of the rule. Must be greater than 0 and less than 65535.",
						},
						"policy": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.StringInSlice([]string{"allow", "drop"}, false),
							Description:  "Policy of the rule (allow or drop)",
						},
					},
				},
			},
			"egress_rules": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of egress rules for the security group",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the rule",
						},
						"ip_version": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.StringInSlice([]string{"ipv4", "ipv6"}, false),
							Description:  "IP version of the rule (ipv4 or ipv6)",
						},
						"protocol": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.StringInSlice([]string{"all", "tcp", "udp", "icmp"}, false),
							Description:  "Protocol of the rule (all, tcp, udp, icmp)",
						},
						"priority": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validate.IntBetween(1, 199),
							Description:  "Priority of the rule. Must be greater than 0 and less than 200.",
						},
						"remote_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.StringInSlice([]string{"address", "securityGroup"}, false),
							Description:  "Type of the remote address (address or securityGroup)",
						},
						"remote_address": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "IP address or CIDR block that the rule applies to",
						},
						"remote_security_group_identity": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Identity of the security group that the rule applies to",
						},
						"port_range_min": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validate.IntBetween(1, 65535),
							Description:  "Minimum port of the rule. Must be greater than 0 and less than 65535.",
						},
						"port_range_max": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validate.IntBetween(1, 65535),
							Description:  "Maximum port of the rule. Must be greater than 0 and less than 65535.",
						},
						"policy": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.StringInSlice([]string{"allow", "drop"}, false),
							Description:  "Policy of the rule (allow or drop)",
						},
					},
				},
			},
			"identity": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity of the security group",
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

func resourceSecurityGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getClient(getProvider(meta), d)
	if err != nil {
		return diag.FromErr(err)
	}

	createReq := iaas.CreateSecurityGroupRequest{
		Name:                  d.Get("name").(string),
		Description:           d.Get("description").(string),
		VpcIdentity:           d.Get("vpc_identity").(string),
		AllowSameGroupTraffic: d.Get("allow_same_group_traffic").(bool),
	}

	if v, ok := d.GetOk("ingress_rules"); ok {
		createReq.IngressRules = expandSecurityGroupRules(v.([]interface{}))
	}

	if v, ok := d.GetOk("egress_rules"); ok {
		createReq.EgressRules = expandSecurityGroupRules(v.([]interface{}))
	}

	securityGroup, err := client.IaaS().CreateSecurityGroup(ctx, createReq)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(securityGroup.Identity)
	return resourceSecurityGroupRead(ctx, d, meta)
}

func resourceSecurityGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getClient(getProvider(meta), d)
	if err != nil {
		return diag.FromErr(err)
	}

	securityGroup, err := client.IaaS().GetSecurityGroup(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

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
	if err := d.Set("identity", securityGroup.Identity); err != nil {
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

func resourceSecurityGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getClient(getProvider(meta), d)
	if err != nil {
		return diag.FromErr(err)
	}

	updateReq := iaas.UpdateSecurityGroupRequest{
		Name:                  d.Get("name").(string),
		Description:           d.Get("description").(string),
		AllowSameGroupTraffic: d.Get("allow_same_group_traffic").(bool),
		ObjectVersion:         d.Get("object_version").(int),
	}

	if v, ok := d.GetOk("ingress_rules"); ok {
		updateReq.IngressRules = expandSecurityGroupRules(v.([]interface{}))
	}

	if v, ok := d.GetOk("egress_rules"); ok {
		updateReq.EgressRules = expandSecurityGroupRules(v.([]interface{}))
	}

	_, err = client.IaaS().UpdateSecurityGroup(ctx, d.Id(), updateReq)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceSecurityGroupRead(ctx, d, meta)
}

func resourceSecurityGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getClient(getProvider(meta), d)
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.IaaS().DeleteSecurityGroup(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func expandSecurityGroupRules(rules []interface{}) []iaas.SecurityGroupRule {
	expandedRules := make([]iaas.SecurityGroupRule, len(rules))
	for i, rule := range rules {
		ruleMap := rule.(map[string]interface{})
		expandedRule := iaas.SecurityGroupRule{
			Name:         ruleMap["name"].(string),
			IPVersion:    iaas.SecurityGroupIPVersion(ruleMap["ip_version"].(string)),
			Protocol:     iaas.SecurityGroupRuleProtocol(ruleMap["protocol"].(string)),
			Priority:     int32(ruleMap["priority"].(int)),
			RemoteType:   iaas.SecurityGroupRuleRemoteType(ruleMap["remote_type"].(string)),
			PortRangeMin: int32(ruleMap["port_range_min"].(int)),
			PortRangeMax: int32(ruleMap["port_range_max"].(int)),
			Policy:       iaas.SecurityGroupRulePolicy(ruleMap["policy"].(string)),
		}

		if v, ok := ruleMap["remote_address"].(string); ok && v != "" {
			expandedRule.RemoteAddress = &v
		}

		if v, ok := ruleMap["remote_security_group_identity"].(string); ok && v != "" {
			expandedRule.RemoteSecurityGroupIdentity = &v
		}

		expandedRules[i] = expandedRule
	}
	return expandedRules
}

func flattenSecurityGroupRules(rules []iaas.SecurityGroupRule) []map[string]interface{} {
	flattenedRules := make([]map[string]interface{}, len(rules))
	for i, rule := range rules {
		flattenedRule := map[string]interface{}{
			"name":           rule.Name,
			"ip_version":     rule.IPVersion,
			"protocol":       rule.Protocol,
			"priority":       rule.Priority,
			"remote_type":    rule.RemoteType,
			"port_range_min": rule.PortRangeMin,
			"port_range_max": rule.PortRangeMax,
			"policy":         rule.Policy,
		}

		if rule.RemoteAddress != nil {
			flattenedRule["remote_address"] = *rule.RemoteAddress
		}

		if rule.RemoteSecurityGroupIdentity != nil {
			flattenedRule["remote_security_group_identity"] = *rule.RemoteSecurityGroupIdentity
		}

		flattenedRules[i] = flattenedRule
	}
	return flattenedRules
}
