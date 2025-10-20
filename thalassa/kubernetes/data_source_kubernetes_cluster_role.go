package kubernetes

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	kubernetes "github.com/thalassa-cloud/client-go/kubernetes"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func dataSourceKubernetesClusterRole() *schema.Resource {
	return &schema.Resource{
		Description: "Get a Kubernetes cluster role by name or slug",
		ReadContext: dataSourceKubernetesClusterRoleRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the Kubernetes cluster role",
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the Kubernetes Cluster Role. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the Kubernetes cluster role to look up",
			},
			"slug": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The slug of the Kubernetes cluster role to look up",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A human-readable description of the Kubernetes cluster role",
			},
			"labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Labels for the Kubernetes cluster role",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"annotations": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Annotations for the Kubernetes cluster role",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"system": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether this is a system role",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp when the Kubernetes cluster role was created",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp when the Kubernetes cluster role was last updated",
			},
			"rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of permission rules for this Kubernetes cluster role",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier of the permission rule",
						},
						"resources": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of resources that the rule applies to",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"verbs": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of verbs that the rule applies to",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"api_groups": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of API groups that the rule applies to",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"resource_names": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of resource names that the rule applies to",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"non_resource_urls": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of non-resource URLs that the rule applies to",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"note": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A human-readable note for the permission rule",
						},
					},
				},
			},
		},
	}
}

func dataSourceKubernetesClusterRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	roles, err := client.Kubernetes().ListKubernetesClusterRoles(ctx, &kubernetes.ListKubernetesClusterRolesRequest{})
	if err != nil {
		return diag.FromErr(err)
	}

	var matchingRole *kubernetes.KubernetesClusterRole
	name := d.Get("name").(string)
	slug := d.Get("slug").(string)

	// Find matching role by name or slug
	for _, role := range roles {
		if (name != "" && role.Name == name) || (slug != "" && role.Slug == slug) {
			if matchingRole != nil {
				return diag.Errorf("multiple roles found with the same name/slug")
			}
			matchingRole = &role
		}
	}

	if matchingRole == nil {
		return diag.Errorf("no Kubernetes cluster role found with the specified name or slug")
	}

	// Set the resource data
	d.SetId(matchingRole.Identity)
	d.Set("name", matchingRole.Name)
	d.Set("slug", matchingRole.Slug)
	d.Set("description", matchingRole.Description)
	d.Set("labels", matchingRole.Labels)
	d.Set("annotations", matchingRole.Annotations)
	d.Set("system", matchingRole.System)
	d.Set("created_at", matchingRole.CreatedAt.Format(time.RFC3339))
	if matchingRole.UpdatedAt != nil {
		d.Set("updated_at", matchingRole.UpdatedAt.Format(time.RFC3339))
	}

	// Set rules
	if len(matchingRole.Rules) > 0 {
		rules := make([]map[string]interface{}, len(matchingRole.Rules))
		for i, rule := range matchingRole.Rules {
			rules[i] = map[string]interface{}{
				"id":                rule.Identity,
				"resources":         rule.Resources,
				"verbs":             convertVerbsToStringSlice(rule.Verbs),
				"api_groups":        rule.ApiGroups,
				"resource_names":    rule.ResourceNames,
				"non_resource_urls": rule.NonResourceURLs,
				"note":              rule.Note,
			}
		}
		d.Set("rules", rules)
	}

	return nil
}
