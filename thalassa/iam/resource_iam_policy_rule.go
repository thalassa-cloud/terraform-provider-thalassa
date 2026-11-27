package iam

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	iam "github.com/thalassa-cloud/client-go/iam"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func ResourceIamPolicyRule() *schema.Resource {
	return &schema.Resource{
		Description: "Manage a permission rule for a project IAM policy in Thalassa Cloud. " +
			"Requires projects to be enabled for the organisation (currently private beta).",
		CreateContext: resourceIamPolicyRuleCreate,
		ReadContext:   resourceIamPolicyRuleRead,
		DeleteContext: resourceIamPolicyRuleDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity of the permission rule",
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Organisation of the IAM policy rule. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Identity of the IAM policy this rule belongs to",
			},
			"resources": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "List of resources this rule applies to",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"resource_identities": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "List of specific resource identities this rule applies to",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"permissions": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "List of permissions (create, read, update, delete, list, *)",
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validate.StringInSlice([]string{
						"create", "read", "update", "delete", "list", "*",
					}, false),
				},
			},
			"note": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Human-readable note for the permission rule",
			},
		},
	}
}

func resourceIamPolicyRuleCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	policyIdentity := d.Get("policy_id").(string)

	resources := make([]string, 0)
	if rList, ok := d.Get("resources").([]any); ok {
		for _, r := range rList {
			resources = append(resources, r.(string))
		}
	}

	resourceIdentities := make([]string, 0)
	if ri, ok := d.Get("resource_identities").([]any); ok {
		for _, r := range ri {
			resourceIdentities = append(resourceIdentities, r.(string))
		}
	}

	permissions := make([]string, 0)
	if pList, ok := d.Get("permissions").([]any); ok {
		for _, p := range pList {
			permissions = append(permissions, p.(string))
		}
	}

	note := ""
	if n, ok := d.Get("note").(string); ok {
		note = n
	}

	ruleReq := iam.AddIamPolicyRuleRequest{
		Resources:          resources,
		ResourceIdentities: resourceIdentities,
		Permissions:        convertStringsToPermissions(permissions),
		Note:               note,
	}

	rule, err := client.IAM().AddIamPolicyRule(ctx, policyIdentity, ruleReq)
	if err != nil {
		return diag.FromErr(err)
	}
	if rule == nil {
		return diag.FromErr(fmt.Errorf("failed to create IAM policy rule"))
	}

	d.SetId(rule.Identity)
	return resourceIamPolicyRuleRead(ctx, d, m)
}

func resourceIamPolicyRuleRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	policyIdentity := d.Get("policy_id").(string)
	ruleIdentity := d.Id()

	policy, err := client.IAM().GetIamPolicy(ctx, policyIdentity)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	if policy == nil {
		d.SetId("")
		return nil
	}

	var rule *iam.IamPolicyPermissionRule
	for i := range policy.Rules {
		if policy.Rules[i].Identity == ruleIdentity {
			rule = &policy.Rules[i]
			break
		}
	}
	if rule == nil {
		d.SetId("")
		return nil
	}

	d.SetId(rule.Identity)
	_ = d.Set("resources", toListOfInterfaces(rule.Resources))
	_ = d.Set("resource_identities", toListOfInterfaces(rule.ResourceIdentities))
	_ = d.Set("permissions", toListOfInterfaces(convertPermissionsToStrings(rule.Permissions)))
	_ = d.Set("note", rule.Note)

	return nil
}

func resourceIamPolicyRuleDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	policyIdentity := d.Get("policy_id").(string)
	ruleIdentity := d.Id()

	err = client.IAM().DeleteIamPolicyRule(ctx, policyIdentity, ruleIdentity)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
