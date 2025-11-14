package iam

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"

	iam "github.com/thalassa-cloud/client-go/iam"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

func ResourceRoleRule() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a permission rule for an organisation role in Thalassa Cloud",
		CreateContext: resourceRoleRuleCreate,
		ReadContext:   resourceRoleRuleRead,
		UpdateContext: resourceRoleRuleUpdate,
		DeleteContext: resourceRoleRuleDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity of the permission rule",
			},
			"role_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Identity of the role this rule belongs to",
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

func resourceRoleRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	roleIdentity := d.Get("role_id").(string)

	// Convert resources list
	resources := make([]string, 0)
	if rList, ok := d.Get("resources").([]interface{}); ok {
		for _, r := range rList {
			resources = append(resources, r.(string))
		}
	}

	// Convert resource identities list
	resourceIdentities := make([]string, 0)
	if ri, ok := d.Get("resource_identities").([]interface{}); ok {
		for _, r := range ri {
			resourceIdentities = append(resourceIdentities, r.(string))
		}
	}

	// Convert permissions list
	permissions := make([]string, 0)
	if pList, ok := d.Get("permissions").([]interface{}); ok {
		for _, p := range pList {
			permissions = append(permissions, p.(string))
		}
	}

	note := ""
	if n, ok := d.Get("note").(string); ok {
		note = n
	}

	ruleReq := iam.OrganisationRolePermissionRule{
		Resources:          resources,
		ResourceIdentities: resourceIdentities,
		Permissions:        convertStringsToPermissions(permissions),
		Note:               note,
	}

	rule, err := client.IAM().AddRoleRule(ctx, roleIdentity, ruleReq)
	if err != nil {
		return diag.FromErr(err)
	}

	if rule != nil {
		d.SetId(rule.Identity)
		return resourceRoleRuleRead(ctx, d, m)
	}

	return diag.FromErr(fmt.Errorf("failed to create role rule"))
}

func resourceRoleRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	roleIdentity := d.Get("role_id").(string)
	ruleIdentity := d.Get("id").(string)

	// Get the role to access its rules
	role, err := client.IAM().GetOrganisationRole(ctx, roleIdentity)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if role == nil {
		d.SetId("")
		return nil
	}

	// Find the specific rule
	var rule *iam.OrganisationRolePermissionRule
	for _, r := range role.Rules {
		if r.Identity == ruleIdentity {
			rule = &r
			break
		}
	}

	if rule == nil {
		d.SetId("")
		return nil
	}

	d.SetId(rule.Identity)
	d.Set("resources", toListOfInterfaces(rule.Resources))
	d.Set("resource_identities", toListOfInterfaces(rule.ResourceIdentities))
	d.Set("permissions", toListOfInterfaces(convertPermissionsToStrings(rule.Permissions)))
	d.Set("note", rule.Note)

	return nil
}

func resourceRoleRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Role rules don't support updates via the API, so we need to delete and recreate
	// This is similar to how role bindings work
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	roleIdentity := d.Get("role_id").(string)
	ruleIdentity := d.Get("id").(string)

	// Delete the old rule
	err = client.IAM().DeleteRuleFromRole(ctx, roleIdentity, ruleIdentity)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting role rule for update: %s", err))
	}

	// Clear the ID so create will work
	d.SetId("")

	// Create the new rule with updated values
	return resourceRoleRuleCreate(ctx, d, m)
}

func resourceRoleRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	roleIdentity := d.Get("role_id").(string)
	ruleIdentity := d.Get("id").(string)

	err = client.IAM().DeleteRuleFromRole(ctx, roleIdentity, ruleIdentity)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
