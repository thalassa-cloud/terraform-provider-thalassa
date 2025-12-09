package iaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"

	iaas "github.com/thalassa-cloud/client-go/iaas"
)

func resourceSecurityGroupEgressRule() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages egress rules for a security group using the batch API. This resource replaces all egress rules for the security group. This is an optional alternative to managing rules through the thalassa_security_group resource. Warning: Do not use both this resource and egress_rule in thalassa_security_group for the same security group, as this will cause conflicts.",
		CreateContext: resourceSecurityGroupEgressRuleCreate,
		ReadContext:   resourceSecurityGroupEgressRuleRead,
		UpdateContext: resourceSecurityGroupEgressRuleUpdate,
		DeleteContext: resourceSecurityGroupEgressRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"security_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the Security Group to manage egress rules for",
			},
			"rule": {
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
							Description: "ID of the Security Group that the rule applies to",
						},
						"port_range_min": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      1,
							ValidateFunc: validate.IntBetween(1, 65535),
							Description:  "Minimum port of the rule. Must be greater than 0 and less than 65535.",
						},
						"port_range_max": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      65535,
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
		},
	}
}

func resourceSecurityGroupEgressRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(meta), d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting client: %w", err))
	}

	securityGroupID := d.Get("security_group_id").(string)

	// Verify security group exists and check for potential conflicts
	securityGroup, err := client.IaaS().GetSecurityGroup(ctx, securityGroupID)
	if err != nil {
		if tcclient.IsNotFound(err) {
			return diag.FromErr(fmt.Errorf("security group %q was not found", securityGroupID))
		}
		return diag.FromErr(fmt.Errorf("error getting security group: %w", err))
	}

	// Warn if the security group already has egress rules, as this might indicate
	// that rules are being managed by the security group resource
	if len(securityGroup.EgressRules) > 0 {
		tflog.Warn(ctx, "Security group already has egress rules. Ensure you're not managing rules through both thalassa_security_group and thalassa_security_group_egress_rule resources for the same security group, as this will cause conflicts.", map[string]interface{}{
			"security_group_id": securityGroupID,
			"existing_rules":    len(securityGroup.EgressRules),
		})
	}

	var rules []iaas.SecurityGroupRule
	if v, ok := d.GetOk("rule"); ok {
		rules = expandSecurityGroupRules(v.([]interface{}))
	}

	batchReq := iaas.BatchUpdateSecurityGroupRulesRequest{
		Rules: rules,
	}

	_, err = client.IaaS().BatchUpdateSecurityGroupEgressRules(ctx, securityGroupID, batchReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating security group egress rules: %w", err))
	}

	// Use security_group_id as the resource ID since we manage all egress rules for a security group
	d.SetId(securityGroupID)

	return resourceSecurityGroupEgressRuleRead(ctx, d, meta)
}

func resourceSecurityGroupEgressRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(meta), d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting client: %w", err))
	}

	securityGroupID := d.Id()

	securityGroup, err := client.IaaS().GetSecurityGroup(ctx, securityGroupID)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error getting security group: %w", err))
	}

	if err := d.Set("security_group_id", securityGroupID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("rule", flattenSecurityGroupRules(securityGroup.EgressRules)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSecurityGroupEgressRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(meta), d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting client: %w", err))
	}

	securityGroupID := d.Get("security_group_id").(string)

	var rules []iaas.SecurityGroupRule
	if v, ok := d.GetOk("rule"); ok {
		rules = expandSecurityGroupRules(v.([]interface{}))
	}

	batchReq := iaas.BatchUpdateSecurityGroupRulesRequest{
		Rules: rules,
	}

	_, err = client.IaaS().BatchUpdateSecurityGroupEgressRules(ctx, securityGroupID, batchReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating security group egress rules: %w", err))
	}

	return resourceSecurityGroupEgressRuleRead(ctx, d, meta)
}

func resourceSecurityGroupEgressRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(meta), d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting client: %w", err))
	}

	securityGroupID := d.Get("security_group_id").(string)

	// Delete all egress rules by setting an empty list
	batchReq := iaas.BatchUpdateSecurityGroupRulesRequest{
		Rules: []iaas.SecurityGroupRule{},
	}

	_, err = client.IaaS().BatchUpdateSecurityGroupEgressRules(ctx, securityGroupID, batchReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting security group egress rules: %w", err))
	}

	d.SetId("")
	return nil
}
