package iam

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"

	iam "github.com/thalassa-cloud/client-go/iam"
)

func DataSourceRole() *schema.Resource {
	return &schema.Resource{
		Description: "Get an organisation role",
		ReadContext: dataSourceRoleRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity of the Organisation Role",
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the Role. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.StringLenBetween(1, 255),
				Description:  "Name of the Organisation Role",
			},
			"slug": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Slug of the Organisation Role",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.StringLenBetween(0, 255),
				Description:  "A human readable description about the role",
			},
			"labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Labels for the Organisation Role",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Annotations for the Organisation Role",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation timestamp of the Organisation Role",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last update timestamp of the Organisation Role",
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
				Type:        schema.TypeList,
				Computed:    true,
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
							Computed:    true,
							Description: "List of resources this rule applies to",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"resource_identities": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of specific resource identities this rule applies to",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"permissions": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of permissions (create, read, update, delete, list, *)",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"note": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Human-readable note for the permission rule",
						},
					},
				},
			},
		},
	}
}

func dataSourceRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	slug := d.Get("slug").(string)
	name := d.Get("name").(string)

	if name == "" && slug == "" {
		return diag.FromErr(fmt.Errorf("either 'name' or 'slug' must be provided to look up a role"))
	}

	roles, err := client.IAM().ListOrganisationRoles(ctx, &iam.ListOrganisationRolesRequest{})
	if err != nil {
		return diag.FromErr(err)
	}

	var role *iam.OrganisationRole

	if name != "" && slug != "" {
		// Both name and slug provided: match both
		for _, r := range roles {
			if r.Name == name && r.Slug == slug {
				role = &r
				break
			}
		}
		if role == nil {
			return diag.FromErr(fmt.Errorf("no role found with name '%s' and slug '%s'", name, slug))
		}
	} else if name != "" {
		// Only name provided: match all by name
		var matchingRoles []iam.OrganisationRole
		for _, r := range roles {
			if r.Name == name {
				matchingRoles = append(matchingRoles, r)
			}
		}
		if len(matchingRoles) == 0 {
			return diag.FromErr(fmt.Errorf("no role found with name '%s'", name))
		}
		if len(matchingRoles) > 1 {
			var slugs []string
			for _, r := range matchingRoles {
				slugs = append(slugs, r.Slug)
			}
			return diag.FromErr(fmt.Errorf("multiple roles found with name '%s', please specify one of these slugs: %v", name, slugs))
		}
		role = &matchingRoles[0]
	} else if slug != "" {
		// Only slug provided: match by slug
		for _, r := range roles {
			if r.Slug == slug {
				role = &r
				break
			}
		}
		if role == nil {
			return diag.FromErr(fmt.Errorf("no role found with slug '%s'", slug))
		}
	}

	// Set fields if found
	d.SetId(role.Identity)
	d.Set("id", role.Identity)
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
	ruleList := make([]map[string]any, len(role.Rules))
	for i, rule := range role.Rules {
		ruleMap := map[string]any{
			"identity":            rule.Identity,
			"resources":           toListOfInterfaces(rule.Resources),
			"resource_identities": toListOfInterfaces(rule.ResourceIdentities),
			"permissions":         toListOfInterfaces(convertPermissionsToStrings(rule.Permissions)),
			"note":                rule.Note,
		}
		ruleList[i] = ruleMap
	}
	d.Set("rules", ruleList)

	return diag.Diagnostics{}
}
