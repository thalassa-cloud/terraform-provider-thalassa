package projects

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
	tcprojects "github.com/thalassa-cloud/client-go/projects"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func ResourceProject() *schema.Resource {
	return &schema.Resource{
		Description:   "Create and manage a project in the Thalassa Cloud platform",
		CreateContext: resourceProjectCreate,
		ReadContext:   resourceProjectRead,
		UpdateContext: resourceProjectUpdate,
		DeleteContext: resourceProjectDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity of the Project",
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Organisation of the Project. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.StringLenBetween(1, 255),
				Description:  "Name of the Project",
			},
			"slug": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Slug of the Project",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				ValidateFunc: validate.StringLenBetween(0, 255),
				Description:  "Human-readable description of the Project",
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Labels for the Project",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Annotations for the Project",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"parent_project_identity": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Identity of the parent project, if this project is nested",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation timestamp of the Project (RFC3339)",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last update timestamp of the Project (RFC3339)",
			},
			"object_version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Object version of the Project",
			},
		},
	}
}

func resourceProjectCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	createReq := tcprojects.CreateProjectRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Labels:      convert.ConvertToMap(d.Get("labels")),
		Annotations: convert.ConvertToMap(d.Get("annotations")),
	}
	if v, ok := d.GetOk("parent_project_identity"); ok {
		parent := v.(string)
		createReq.ParentProjectIdentity = &parent
	}

	project, err := client.Projects().CreateProject(ctx, createReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating project: %w", err))
	}
	if project == nil {
		return diag.FromErr(fmt.Errorf("error creating project: empty response"))
	}

	d.SetId(project.Identity)
	return resourceProjectRead(ctx, d, m)
}

func resourceProjectRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	identity := d.Id()
	project, err := client.Projects().GetProject(ctx, identity)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error getting project: %w", err))
	}
	if project == nil {
		d.SetId("")
		return nil
	}

	return setProjectResourceState(d, project)
}

func setProjectResourceState(d *schema.ResourceData, project *tcprojects.Project) diag.Diagnostics {
	d.SetId(project.Identity)
	_ = d.Set("name", project.Name)
	_ = d.Set("slug", project.Slug)
	_ = d.Set("description", project.Description)
	_ = d.Set("labels", project.Labels)
	_ = d.Set("annotations", project.Annotations)
	_ = d.Set("created_at", project.CreatedAt.Format(TimeFormatRFC3339))
	_ = d.Set("object_version", project.ObjectVersion)
	if project.UpdatedAt != nil {
		_ = d.Set("updated_at", project.UpdatedAt.Format(TimeFormatRFC3339))
	}
	if project.ParentProject != nil {
		_ = d.Set("parent_project_identity", project.ParentProject.Identity)
	} else {
		_ = d.Set("parent_project_identity", "")
	}
	return nil
}

func resourceProjectUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	updateReq := tcprojects.UpdateProjectRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Labels:      convert.ConvertToMap(d.Get("labels")),
		Annotations: convert.ConvertToMap(d.Get("annotations")),
	}
	if v, ok := d.GetOk("parent_project_identity"); ok {
		parent := v.(string)
		updateReq.ParentProjectIdentity = &parent
	}

	identity := d.Id()
	project, err := client.Projects().UpdateProject(ctx, identity, updateReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating project: %w", err))
	}
	if project != nil {
		if diags := setProjectResourceState(d, project); diags != nil {
			return diags
		}
	}

	return resourceProjectRead(ctx, d, m)
}

func resourceProjectDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	identity := d.Id()
	if err := client.Projects().DeleteProject(ctx, identity); err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error deleting project: %w", err))
	}

	d.SetId("")
	return nil
}
