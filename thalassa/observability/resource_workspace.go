package observability

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	tcobs "github.com/thalassa-cloud/client-go/observability"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func ResourceWorkspace() *schema.Resource {
	return &schema.Resource{
		Description:   "Create and manage an observability workspace (metrics and logs) in Thalassa Cloud. Observability workspaces are beta and require sign up to the beta program before they can be used by your organisation.",
		CreateContext: resourceWorkspaceCreate,
		ReadContext:   resourceWorkspaceRead,
		UpdateContext: resourceWorkspaceUpdate,
		DeleteContext: resourceWorkspaceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceWorkspaceImport,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
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
				Required:     true,
				ValidateFunc: validate.StringLenBetween(1, 255),
				Description:  "Name of the observability workspace",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				ValidateFunc: validate.StringLenBetween(0, 1024),
				Description:  "Human-readable description of the workspace",
			},
			"region": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StringLenBetween(1, 255),
				Description:  "Region slug where the workspace is created (e.g. nl-01). Import as region/identity.",
			},
			"retention_days": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.IntBetween(1, 1095),
				Description:  "Retention period in days for metrics and logs (1-1095). Omit to use the service default.",
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Labels for the workspace",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
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
			"wait_for_deleted_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validate.IntAtLeast(0),
				Description:  "Timeout in minutes to wait for the workspace and its underlying data to be fully deleted. Set to 0 (default) to return after the delete is accepted without waiting.",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
	}
}

func resourceWorkspaceImport(_ context.Context, d *schema.ResourceData, _ any) ([]*schema.ResourceData, error) {
	region, identity := parseWorkspaceImportID(d.Id())
	if identity == "" {
		return nil, fmt.Errorf("invalid import id %q; expected \"region/identity\" or \"identity\"", d.Id())
	}
	if region != "" {
		_ = d.Set("region", region)
	}
	d.SetId(identity)
	return []*schema.ResourceData{d}, nil
}

func parseWorkspaceImportID(id string) (region, identity string) {
	id = strings.TrimSpace(id)
	if id == "" {
		return "", ""
	}
	parts := strings.SplitN(id, "/", 2)
	if len(parts) == 1 {
		return "", parts[0]
	}
	return parts[0], parts[1]
}

func resourceWorkspaceCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	createReq := tcobs.CreateObservabilityWorkspaceRequest{
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		Labels:         convert.ConvertToMap(d.Get("labels")),
		Annotations:    convert.ConvertToMap(d.Get("annotations")),
		RegionIdentity: d.Get("region").(string),
	}
	if v, ok := d.GetOk("retention_days"); ok {
		createReq.RetentionDays = convert.Ptr(v.(int))
	}

	workspace, err := client.Observability().CreateObservabilityWorkspace(ctx, createReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating observability workspace: %w", err))
	}
	if workspace == nil {
		return diag.FromErr(fmt.Errorf("error creating observability workspace: empty response"))
	}

	d.SetId(workspace.Identity)

	waitCtx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutCreate))
	defer cancel()
	workspace, err = client.Observability().WaitUntilObservabilityWorkspaceReady(waitCtx, workspace.Identity)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error waiting for observability workspace to become ready: %w", err))
	}
	if workspace != nil {
		if diags := setWorkspaceState(d, workspace); diags != nil {
			return diags
		}
	}

	return resourceWorkspaceRead(ctx, d, m)
}

func resourceWorkspaceRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	workspace, err := client.Observability().GetObservabilityWorkspace(ctx, d.Id())
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error getting observability workspace: %w", err))
	}
	if workspace == nil {
		d.SetId("")
		return nil
	}

	return setWorkspaceState(d, workspace)
}

func setWorkspaceState(d *schema.ResourceData, workspace *tcobs.ObservabilityWorkspace) diag.Diagnostics {
	d.SetId(workspace.Identity)
	_ = d.Set("name", workspace.Name)
	_ = d.Set("description", workspace.Description)
	_ = d.Set("labels", workspace.Labels)
	_ = d.Set("annotations", workspace.Annotations)
	_ = d.Set("status", string(workspace.Status))
	_ = d.Set("status_message", workspace.StatusMessage)
	_ = d.Set("prometheus_enabled", workspace.PrometheusEnabled)
	_ = d.Set("loki_enabled", workspace.LokiEnabled)
	_ = d.Set("remote_write_url", workspace.RemoteWriteURL)
	_ = d.Set("remote_write_otlp_url", workspace.RemoteWriteOTLPURL)
	_ = d.Set("prometheus_query_url", workspace.PrometheusQueryURL)
	_ = d.Set("alerting_url", workspace.AlertingURL)
	_ = d.Set("push_url", workspace.PushURL)
	_ = d.Set("push_otlp_url", workspace.PushOTLPURL)
	_ = d.Set("loki_query_url", workspace.LokiQueryURL)
	_ = d.Set("created_at", workspace.CreatedAt.Format(TimeFormatRFC3339))
	_ = d.Set("updated_at", workspace.UpdatedAt.Format(TimeFormatRFC3339))
	_ = d.Set("object_version", workspace.ObjectVersion)

	if workspace.RetentionDays != nil {
		_ = d.Set("retention_days", *workspace.RetentionDays)
	}

	if workspace.Region != nil {
		switch {
		case workspace.Region.Slug != "":
			_ = d.Set("region", workspace.Region.Slug)
		case workspace.Region.Identity != "":
			_ = d.Set("region", workspace.Region.Identity)
		case workspace.Region.Name != "":
			_ = d.Set("region", workspace.Region.Name)
		}
	}

	return nil
}

func resourceWorkspaceUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	updateReq := tcobs.UpdateObservabilityWorkspaceRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Labels:      convert.ConvertToMap(d.Get("labels")),
		Annotations: convert.ConvertToMap(d.Get("annotations")),
	}
	if v, ok := d.GetOk("retention_days"); ok {
		updateReq.RetentionDays = convert.Ptr(v.(int))
	}

	_, err = client.Observability().UpdateObservabilityWorkspace(ctx, d.Id(), updateReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating observability workspace: %w", err))
	}

	waitCtx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutUpdate))
	defer cancel()
	workspace, err := client.Observability().WaitUntilObservabilityWorkspaceReady(waitCtx, d.Id())
	if err != nil {
		return diag.FromErr(fmt.Errorf("error waiting for observability workspace to become ready: %w", err))
	}
	if workspace != nil {
		if diags := setWorkspaceState(d, workspace); diags != nil {
			return diags
		}
	}

	return resourceWorkspaceRead(ctx, d, m)
}

func resourceWorkspaceDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := client.Observability().DeleteObservabilityWorkspace(ctx, d.Id()); err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error deleting observability workspace: %w", err))
	}

	if timeoutMinutes := d.Get("wait_for_deleted_timeout").(int); timeoutMinutes > 0 {
		waitCtx, cancel := context.WithTimeout(ctx, time.Duration(timeoutMinutes)*time.Minute)
		defer cancel()
		for {
			select {
			case <-waitCtx.Done():
				return diag.FromErr(fmt.Errorf("timeout waiting for observability workspace %s to be deleted", d.Id()))
			case <-time.After(5 * time.Second):
				_, err := client.Observability().GetObservabilityWorkspace(waitCtx, d.Id())
				if err != nil {
					if tcclient.IsNotFound(err) {
						d.SetId("")
						return nil
					}
					return diag.FromErr(fmt.Errorf("error waiting for observability workspace deletion: %w", err))
				}
			}
		}
	}

	d.SetId("")
	return nil
}
