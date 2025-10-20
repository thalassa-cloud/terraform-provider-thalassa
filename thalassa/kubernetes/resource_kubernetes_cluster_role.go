package kubernetes

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	kubernetes "github.com/thalassa-cloud/client-go/kubernetes"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func resourceKubernetesClusterRole() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages a Kubernetes cluster role for IAM access control. This resource allows you to create, update, and delete Kubernetes cluster roles with associated permission rules.",
		CreateContext: resourceKubernetesClusterRoleCreate,
		ReadContext:   resourceKubernetesClusterRoleRead,
		UpdateContext: resourceKubernetesClusterRoleUpdate,
		DeleteContext: resourceKubernetesClusterRoleDelete,
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
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.StringLenBetween(1, 255),
				Description:  "The name of the Kubernetes cluster role",
			},
			"slug": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The slug of the Kubernetes cluster role",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.StringLenBetween(0, 1000),
				Description:  "A human-readable description of the Kubernetes cluster role",
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Labels for the Kubernetes cluster role",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
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
				Optional:    true,
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
							Required:    true,
							Description: "List of resources that the rule applies to",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"verbs": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "List of verbs that the rule applies to",
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validate.StringInSlice([]string{"*", "get", "list", "watch", "create", "update", "delete", "patch"}, false),
							},
						},
						"api_groups": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of API groups that the rule applies to",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"resource_names": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of resource names that the rule applies to",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"non_resource_urls": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of non-resource URLs that the rule applies to",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"note": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A human-readable note for the permission rule",
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceKubernetesClusterRoleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	createRequest := kubernetes.CreateKubernetesClusterRoleRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Labels:      convert.ConvertToMap(d.Get("labels")),
		Annotations: convert.ConvertToMap(d.Get("annotations")),
	}

	role, err := client.Kubernetes().CreateKubernetesClusterRole(ctx, createRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	if role != nil {
		d.SetId(role.Identity)
		d.Set("slug", role.Slug)
		d.Set("system", role.System)
		d.Set("created_at", role.CreatedAt.Format(time.RFC3339))
		if role.UpdatedAt != nil {
			d.Set("updated_at", role.UpdatedAt.Format(time.RFC3339))
		}
	}

	// Add rules if provided
	if rules, ok := d.GetOk("rules"); ok {
		rulesList := rules.([]interface{})
		for _, ruleInterface := range rulesList {
			ruleMap := ruleInterface.(map[string]interface{})

			rule := kubernetes.AddKubernetesClusterRolePermissionRule{
				Resources:       convert.ConvertToStringSlice(ruleMap["resources"]),
				Verbs:           convertToStringSliceToVerbs(ruleMap["verbs"]),
				ApiGroups:       convert.ConvertToStringSlice(ruleMap["api_groups"]),
				ResourceNames:   convert.ConvertToStringSlice(ruleMap["resource_names"]),
				NonResourceURLs: convert.ConvertToStringSlice(ruleMap["non_resource_urls"]),
				Note:            ruleMap["note"].(string),
			}

			_, err := client.Kubernetes().AddClusterRoleRule(ctx, role.Identity, rule)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to add rule: %w", err))
			}
		}
	}

	return resourceKubernetesClusterRoleRead(ctx, d, m)
}

func resourceKubernetesClusterRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	identity := d.Get("id").(string)
	role, err := client.Kubernetes().GetKubernetesClusterRole(ctx, identity)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error getting kubernetes cluster role: %s", err))
	}

	if role == nil {
		return diag.FromErr(fmt.Errorf("kubernetes cluster role was not found"))
	}

	d.SetId(role.Identity)
	d.Set("name", role.Name)
	d.Set("slug", role.Slug)
	d.Set("description", role.Description)
	d.Set("labels", role.Labels)
	d.Set("annotations", role.Annotations)
	d.Set("system", role.System)
	d.Set("created_at", role.CreatedAt.Format(time.RFC3339))
	if role.UpdatedAt != nil {
		d.Set("updated_at", role.UpdatedAt.Format(time.RFC3339))
	}

	// Set rules
	if len(role.Rules) > 0 {
		rules := make([]map[string]interface{}, len(role.Rules))
		for i, rule := range role.Rules {
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

func resourceKubernetesClusterRoleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Note: The API doesn't seem to have an update method for cluster roles
	// This would need to be implemented if the API supports it
	return diag.Errorf("updating Kubernetes cluster roles is not currently supported")
}

func resourceKubernetesClusterRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	identity := d.Get("id").(string)
	err = client.Kubernetes().DeleteClusterRole(ctx, identity)
	if err != nil {
		if !tcclient.IsNotFound(err) {
			return diag.FromErr(err)
		}
	}

	d.SetId("")
	return nil
}

// Helper functions for converting between string slices and verb types
func convertToStringSliceToVerbs(verbsInterface interface{}) []kubernetes.KubernetesClusterRolePermissionVerb {
	verbsList := verbsInterface.([]interface{})
	verbs := make([]kubernetes.KubernetesClusterRolePermissionVerb, len(verbsList))
	for i, verb := range verbsList {
		verbs[i] = kubernetes.KubernetesClusterRolePermissionVerb(verb.(string))
	}
	return verbs
}

func convertVerbsToStringSlice(verbs []kubernetes.KubernetesClusterRolePermissionVerb) []string {
	result := make([]string, len(verbs))
	for i, verb := range verbs {
		result[i] = string(verb)
	}
	return result
}
