package containerregistry

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	tcregistry "github.com/thalassa-cloud/client-go/containerregistry"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func retentionPolicySchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Retention policy for the namespace",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"enabled": {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Whether the retention policy is active",
				},
				"delete_untagged_images": {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Whether to delete untagged images",
				},
				"rules": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "List of retention rules",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"days": {
								Type:        schema.TypeInt,
								Optional:    true,
								Description: "Retain tags based on tag creation age in days",
							},
							"days_since_created": {
								Type:        schema.TypeInt,
								Optional:    true,
								Description: "Retain artifacts based on creation/push age in days",
							},
							"days_since_pulled": {
								Type:        schema.TypeInt,
								Optional:    true,
								Description: "Retain artifacts based on last pull age in days",
							},
							"count": {
								Type:        schema.TypeInt,
								Optional:    true,
								Description: "Number of most recent versions/tags to keep",
							},
							"repository_patterns": {
								Type:        schema.TypeList,
								Optional:    true,
								Description: "Repository name patterns this rule applies to",
								Elem:        &schema.Schema{Type: schema.TypeString},
							},
							"tag_patterns": {
								Type:        schema.TypeList,
								Optional:    true,
								Description: "Tag patterns this rule applies to",
								Elem:        &schema.Schema{Type: schema.TypeString},
							},
							"scope": {
								Type:         schema.TypeString,
								Optional:     true,
								Default:      string(tcregistry.RetentionPolicyScopeTags),
								ValidateFunc: validate.StringInSlice([]string{string(tcregistry.RetentionPolicyScopeTags)}, false),
								Description:  "Scope of the retention rule (currently only tags)",
							},
						},
					},
				},
			},
		},
	}
}

func ResourceNamespaceConfiguration() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage configuration for a container registry namespace in Thalassa Cloud",
		CreateContext: resourceNamespaceConfigurationCreate,
		ReadContext:   resourceNamespaceConfigurationRead,
		UpdateContext: resourceNamespaceConfigurationUpdate,
		DeleteContext: resourceNamespaceConfigurationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity of the namespace this configuration belongs to",
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Organisation of the namespace. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"namespace_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Identity of the container registry namespace",
			},
			"visibility": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.StringInSlice([]string{string(tcregistry.NamespaceVisibilityPrivate)}, false),
				Description:  "Visibility of the namespace (currently only private)",
			},
			"retention_policy": retentionPolicySchema(),
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation timestamp (RFC3339)",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last update timestamp (RFC3339)",
			},
			"object_version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Object version of the configuration",
			},
		},
	}
}

func expandRetentionPolicy(v any) *tcregistry.RetentionPolicy {
	list, ok := v.([]any)
	if !ok || len(list) == 0 || list[0] == nil {
		return nil
	}

	m, ok := list[0].(map[string]any)
	if !ok || m == nil {
		return nil
	}

	policy := &tcregistry.RetentionPolicy{
		Rules: []tcregistry.RetentionPolicyRule{},
	}
	if enabled, ok := m["enabled"].(bool); ok {
		policy.Enabled = enabled
	}
	if deleteUntagged, ok := m["delete_untagged_images"].(bool); ok {
		policy.DeleteUntaggedImages = deleteUntagged
	}

	rulesRaw, ok := m["rules"].([]any)
	if !ok {
		return policy
	}

	for _, ruleRaw := range rulesRaw {
		rm, ok := ruleRaw.(map[string]any)
		if !ok || rm == nil {
			continue
		}
		rule := tcregistry.RetentionPolicyRule{
			RepositoryPatterns: stringSliceFromAny(rm["repository_patterns"]),
			TagPatterns:        stringSliceFromAny(rm["tag_patterns"]),
			Scope:              tcregistry.RetentionPolicyScopeTags,
		}
		if scope, ok := rm["scope"].(string); ok && scope != "" {
			rule.Scope = tcregistry.RetentionPolicyScope(scope)
		}
		if days, ok := rm["days"].(int); ok && days != 0 {
			rule.Days = &days
		}
		if days, ok := rm["days_since_created"].(int); ok && days != 0 {
			rule.DaysSinceCreated = &days
		}
		if days, ok := rm["days_since_pulled"].(int); ok && days != 0 {
			rule.DaysSincePulled = &days
		}
		if count, ok := rm["count"].(int); ok && count != 0 {
			rule.Count = &count
		}
		policy.Rules = append(policy.Rules, rule)
	}

	return policy
}

func stringSliceFromAny(v any) []string {
	if v == nil {
		return []string{}
	}
	switch typed := v.(type) {
	case []string:
		return typed
	case []any:
		return convert.ConvertToStringSlice(typed)
	default:
		return []string{}
	}
}

func flattenRetentionPolicy(policy *tcregistry.RetentionPolicy) []any {
	if policy == nil {
		return []any{}
	}

	rules := make([]any, 0, len(policy.Rules))
	for _, rule := range policy.Rules {
		rm := map[string]any{
			"repository_patterns": rule.RepositoryPatterns,
			"tag_patterns":        rule.TagPatterns,
			"scope":               string(rule.Scope),
		}
		if rm["scope"] == "" {
			rm["scope"] = string(tcregistry.RetentionPolicyScopeTags)
		}
		if rule.RepositoryPatterns == nil {
			rm["repository_patterns"] = []string{}
		}
		if rule.TagPatterns == nil {
			rm["tag_patterns"] = []string{}
		}
		if rule.Days != nil {
			rm["days"] = *rule.Days
		}
		if rule.DaysSinceCreated != nil {
			rm["days_since_created"] = *rule.DaysSinceCreated
		}
		if rule.DaysSincePulled != nil {
			rm["days_since_pulled"] = *rule.DaysSincePulled
		}
		if rule.Count != nil {
			rm["count"] = *rule.Count
		}
		rules = append(rules, rm)
	}

	return []any{
		map[string]any{
			"enabled":                policy.Enabled,
			"delete_untagged_images": policy.DeleteUntaggedImages,
			"rules":                  rules,
		},
	}
}

func resourceNamespaceConfigurationCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	namespaceID := d.Get("namespace_id").(string)
	createReq := tcregistry.CreateNamespaceConfigurationRequest{
		Visibility:      tcregistry.NamespaceVisibility(d.Get("visibility").(string)),
		RetentionPolicy: expandRetentionPolicy(d.Get("retention_policy")),
	}

	cfg, err := client.ContainerRegistry().CreateNamespaceConfiguration(ctx, namespaceID, createReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating namespace configuration: %w", err))
	}
	if cfg == nil {
		return diag.FromErr(fmt.Errorf("error creating namespace configuration: empty response"))
	}

	d.SetId(namespaceID)
	return resourceNamespaceConfigurationRead(ctx, d, m)
}

func resourceNamespaceConfigurationRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	namespaceID := d.Id()
	if namespaceID == "" {
		namespaceID = d.Get("namespace_id").(string)
	}

	cfg, err := client.ContainerRegistry().GetNamespaceConfiguration(ctx, namespaceID)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error getting namespace configuration: %w", err))
	}
	if cfg == nil {
		d.SetId("")
		return nil
	}

	return setNamespaceConfigurationState(d, namespaceID, cfg)
}

func setNamespaceConfigurationState(d *schema.ResourceData, namespaceID string, cfg *tcregistry.ContainerRegistryNamespaceConfiguration) diag.Diagnostics {
	d.SetId(namespaceID)
	_ = d.Set("namespace_id", namespaceID)
	_ = d.Set("visibility", string(cfg.Visibility))
	_ = d.Set("retention_policy", flattenRetentionPolicy(cfg.RetentionPolicy))
	_ = d.Set("created_at", cfg.CreatedAt.Format(TimeFormatRFC3339))
	_ = d.Set("updated_at", cfg.UpdatedAt.Format(TimeFormatRFC3339))
	_ = d.Set("object_version", cfg.ObjectVersion)
	if cfg.Namespace != nil && cfg.Namespace.Identity != "" {
		_ = d.Set("namespace_id", cfg.Namespace.Identity)
		d.SetId(cfg.Namespace.Identity)
	}
	return nil
}

func resourceNamespaceConfigurationUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	namespaceID := d.Id()
	updateReq := tcregistry.UpdateNamespaceConfigurationRequest{
		Visibility:      tcregistry.NamespaceVisibility(d.Get("visibility").(string)),
		RetentionPolicy: expandRetentionPolicy(d.Get("retention_policy")),
	}

	cfg, err := client.ContainerRegistry().UpdateNamespaceConfiguration(ctx, namespaceID, updateReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating namespace configuration: %w", err))
	}
	if cfg != nil {
		if diags := setNamespaceConfigurationState(d, namespaceID, cfg); diags != nil {
			return diags
		}
	}

	return resourceNamespaceConfigurationRead(ctx, d, m)
}

func resourceNamespaceConfigurationDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := client.ContainerRegistry().DeleteNamespaceConfiguration(ctx, d.Id()); err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error deleting namespace configuration: %w", err))
	}

	d.SetId("")
	return nil
}
