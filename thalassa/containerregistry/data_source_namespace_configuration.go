package containerregistry

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func DataSourceNamespaceConfiguration() *schema.Resource {
	return &schema.Resource{
		Description: "Get configuration for a container registry namespace",
		ReadContext: dataSourceNamespaceConfigurationRead,
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
				Description: "Identity of the container registry namespace",
			},
			"visibility": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Visibility of the namespace",
			},
			"retention_policy": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Retention policy for the namespace",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the retention policy is active",
						},
						"delete_untagged_images": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to delete untagged images",
						},
						"rules": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of retention rules",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"days": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Retain tags based on tag creation age in days",
									},
									"days_since_created": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Retain artifacts based on creation/push age in days",
									},
									"days_since_pulled": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Retain artifacts based on last pull age in days",
									},
									"count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Number of most recent versions/tags to keep",
									},
									"repository_patterns": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Repository name patterns this rule applies to",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"tag_patterns": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Tag patterns this rule applies to",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"scope": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Scope of the retention rule",
									},
								},
							},
						},
					},
				},
			},
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

func dataSourceNamespaceConfigurationRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	namespaceID := d.Get("namespace_id").(string)
	if namespaceID == "" {
		return diag.FromErr(fmt.Errorf("'namespace_id' must be provided to look up a namespace configuration"))
	}

	cfg, err := client.ContainerRegistry().GetNamespaceConfiguration(ctx, namespaceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting namespace configuration: %w", err))
	}
	if cfg == nil {
		return diag.FromErr(fmt.Errorf("no namespace configuration found for namespace '%s'", namespaceID))
	}

	return setNamespaceConfigurationState(d, namespaceID, cfg)
}
