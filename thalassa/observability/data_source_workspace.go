package observability

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	tcobs "github.com/thalassa-cloud/client-go/observability"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func DataSourceWorkspace() *schema.Resource {
	return &schema.Resource{
		Description: "Get an observability workspace by identity or by name (optionally scoped by region)",
		ReadContext: dataSourceWorkspaceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Identity of the observability workspace",
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Organisation of the workspace. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.StringLenBetween(1, 255),
				Description:  "Name of the observability workspace",
			},
			"region": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.StringLenBetween(1, 255),
				Description:  "Region slug of the workspace. Recommended when looking up by name.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Human-readable description of the workspace",
			},
			"retention_days": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Retention period in days for metrics and logs",
			},
			"labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Labels for the workspace",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"annotations": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Annotations for the workspace",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Lifecycle status of the workspace",
			},
			"status_message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Additional status details when available",
			},
			"prometheus_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether Prometheus metrics are enabled for this workspace",
			},
			"loki_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether Loki logs are enabled for this workspace",
			},
			"remote_write_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Prometheus remote write URL",
			},
			"remote_write_otlp_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "OTLP remote write URL for metrics",
			},
			"prometheus_query_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Prometheus query URL",
			},
			"alerting_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Alerting URL",
			},
			"push_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Loki push URL",
			},
			"push_otlp_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "OTLP push URL for logs",
			},
			"loki_query_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Loki query URL",
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
				Description: "Object version of the workspace",
			},
		},
	}
}

func dataSourceWorkspaceRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	identity := d.Get("id").(string)
	name := d.Get("name").(string)
	region := d.Get("region").(string)

	if identity == "" && name == "" {
		return diag.FromErr(fmt.Errorf("either 'id' or 'name' must be provided to look up an observability workspace"))
	}

	if identity != "" {
		workspace, err := client.Observability().GetObservabilityWorkspace(ctx, identity)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error getting observability workspace by id: %w", err))
		}
		if workspace == nil {
			return diag.FromErr(fmt.Errorf("no observability workspace found with id '%s'", identity))
		}
		return setWorkspaceDataSourceState(d, workspace)
	}

	workspaces, err := client.Observability().ListObservabilityWorkspaces(ctx, &tcobs.ListObservabilityWorkspacesRequest{})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error listing observability workspaces: %w", err))
	}

	var matches []tcobs.ObservabilityWorkspace
	for _, ws := range workspaces {
		if ws.Name != name {
			continue
		}
		if region != "" && !workspaceMatchesRegion(ws, region) {
			continue
		}
		matches = append(matches, ws)
	}

	if len(matches) == 0 {
		if region != "" {
			return diag.FromErr(fmt.Errorf("no observability workspace found with name '%s' in region '%s'", name, region))
		}
		return diag.FromErr(fmt.Errorf("no observability workspace found with name '%s'", name))
	}
	if len(matches) > 1 {
		return diag.FromErr(fmt.Errorf("multiple observability workspaces found with name '%s'; specify 'region' or 'id' to disambiguate", name))
	}

	return setWorkspaceDataSourceState(d, &matches[0])
}

func workspaceMatchesRegion(ws tcobs.ObservabilityWorkspace, region string) bool {
	if ws.Region == nil {
		return false
	}
	return ws.Region.Slug == region || ws.Region.Identity == region || ws.Region.Name == region
}

func setWorkspaceDataSourceState(d *schema.ResourceData, workspace *tcobs.ObservabilityWorkspace) diag.Diagnostics {
	_ = d.Set("id", workspace.Identity)
	return setWorkspaceState(d, workspace)
}
