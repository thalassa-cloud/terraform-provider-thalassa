package iam

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"

	iam "github.com/thalassa-cloud/client-go/iam"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/client-go/thalassa"
)

func ResourceRole() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage an organisation role in Thalassa Cloud",
		CreateContext: resourceRoleCreate,
		ReadContext:   resourceRoleRead,
		UpdateContext: resourceRoleUpdate,
		DeleteContext: resourceRoleDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.StringLenBetween(1, 255),
				Description:  "Name of the organisation role",
				ForceNew:     true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.StringLenBetween(0, 255),
				Description:  "Description of the organisation role",
				ForceNew:     true,
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Labels for the organisation role",
				ForceNew:    true,
			},
			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Annotations for the organisation role",
				ForceNew:    true,
			},
			"slug": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Slug of the organisation role",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation timestamp of the organisation role",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last update timestamp of the organisation role",
			},
			"role_is_read_only": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the role is read-only and cannot be modified.",
			},
			"system": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the role is a system role",
			},
			"rules": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Permission rules for the organisation role",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"identity": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identity of the permission rule",
						},
						"resources": {
							Type:        schema.TypeList,
							Required:    true,
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
							Description: "Human-readable note for the permission rule",
						},
					},
				},
			},
		},
	}
}

func resourceRoleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	createReq := iam.CreateOrganisationRoleRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Labels:      convert.ConvertToMap(d.Get("labels")),
		Annotations: convert.ConvertToMap(d.Get("annotations")),
	}

	role, err := client.IAM().CreateOrganisationRole(ctx, createReq)
	if err != nil {
		return diag.FromErr(err)
	}
	if role != nil {
		d.SetId(role.Identity)

		// Handle rules after role creation
		if d.HasChange("rules") {
			err = updateRoleRules(ctx, client, role.Identity, d)
			if err != nil {
				return diag.FromErr(fmt.Errorf("error updating role rules: %s", err))
			}
		}

		return resourceRoleRead(ctx, d, m)
	}
	return diag.FromErr(fmt.Errorf("failed to create organisation role"))
}

func resourceRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	identity := d.Get("id").(string)
	role, err := client.IAM().GetOrganisationRole(ctx, identity)
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

	d.SetId(role.Identity)
	d.Set("name", role.Name)
	d.Set("slug", role.Slug)
	d.Set("description", role.Description)
	d.Set("labels", role.Labels)
	d.Set("annotations", role.Annotations)
	d.Set("created_at", role.CreatedAt.Format(TimeFormatRFC3339))
	d.Set("updated_at", role.UpdatedAt.Format(TimeFormatRFC3339))
	d.Set("role_is_read_only", role.IsReadOnly)
	d.Set("system", role.System)

	// Set rules
	rulesSet := schema.NewSet(schema.HashResource(&schema.Resource{
		Schema: map[string]*schema.Schema{
			"identity":            {Type: schema.TypeString},
			"resources":           {Type: schema.TypeList, Elem: &schema.Schema{Type: schema.TypeString}},
			"resource_identities": {Type: schema.TypeList, Elem: &schema.Schema{Type: schema.TypeString}},
			"permissions":         {Type: schema.TypeList, Elem: &schema.Schema{Type: schema.TypeString}},
			"note":                {Type: schema.TypeString},
		},
	}), []interface{}{})

	for _, rule := range role.Rules {
		ruleMap := map[string]any{
			"identity":            rule.Identity,
			"resources":           toListOfInterfaces(rule.Resources),
			"resource_identities": toListOfInterfaces(rule.ResourceIdentities),
			"permissions":         toListOfInterfaces(convertPermissionsToStrings(rule.Permissions)),
			"note":                rule.Note,
		}
		rulesSet.Add(ruleMap)
	}
	d.Set("rules", rulesSet)

	return nil
}

func resourceRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}
	identity := d.Get("id").(string)
	err = client.IAM().DeleteOrganisationRole(ctx, identity)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func resourceRoleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	identity := d.Get("id").(string)

	// Handle rules changes
	if d.HasChange("rules") {
		err = updateRoleRules(ctx, client, identity, d)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error updating role rules: %s", err))
		}
	}

	return resourceRoleRead(ctx, d, m)
}

func updateRoleRules(ctx context.Context, client thalassa.Client, roleIdentity string, d *schema.ResourceData) error {

	oldRules, newRules := d.GetChange("rules")
	oldSet := oldRules.(*schema.Set)
	newSet := newRules.(*schema.Set)

	// Find rules to remove
	toRemove := oldSet.Difference(newSet)
	for _, rule := range toRemove.List() {
		ruleMap := rule.(map[string]interface{})
		ruleIdentity := ruleMap["identity"].(string)

		if ruleIdentity != "" {
			err := client.IAM().DeleteRuleFromRole(ctx, roleIdentity, ruleIdentity)
			if err != nil {
				return fmt.Errorf("error removing role rule: %s", err)
			}
		}
	}

	// Find rules to add
	toAdd := newSet.Difference(oldSet)
	for _, rule := range toAdd.List() {
		ruleMap := rule.(map[string]interface{})

		// Convert resources list
		resources := make([]string, 0)
		for _, r := range ruleMap["resources"].([]interface{}) {
			resources = append(resources, r.(string))
		}

		// Convert resource identities list
		resourceIdentities := make([]string, 0)
		if ri, ok := ruleMap["resource_identities"].([]interface{}); ok {
			for _, r := range ri {
				resourceIdentities = append(resourceIdentities, r.(string))
			}
		}

		// Convert permissions list
		permissions := make([]string, 0)
		for _, p := range ruleMap["permissions"].([]interface{}) {
			permissions = append(permissions, p.(string))
		}

		ruleReq := iam.OrganisationRolePermissionRule{
			Resources:          resources,
			ResourceIdentities: resourceIdentities,
			Permissions:        convertStringsToPermissions(permissions),
			Note:               ruleMap["note"].(string),
		}

		_, err := client.IAM().AddRoleRule(ctx, roleIdentity, ruleReq)
		if err != nil {
			return fmt.Errorf("error adding role rule: %s", err)
		}
	}

	return nil
}
