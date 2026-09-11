package projects

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	tcprojects "github.com/thalassa-cloud/client-go/projects"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func DataSourceProject() *schema.Resource {
	return &schema.Resource{
		Description: "Get a project by identity, name, or slug",
		ReadContext: dataSourceProjectRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
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
				Optional:     true,
				ValidateFunc: validate.StringLenBetween(1, 255),
				Description:  "Name of the Project",
			},
			"slug": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Slug of the Project",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Human-readable description of the Project",
			},
			"labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Labels for the Project",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"annotations": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Annotations for the Project",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"parent_project_identity": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity of the parent project, if nested",
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

func dataSourceProjectRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	identity := d.Get("id").(string)
	name := d.Get("name").(string)
	slug := d.Get("slug").(string)

	if identity == "" && name == "" && slug == "" {
		return diag.FromErr(fmt.Errorf("one of 'id', 'name', or 'slug' must be provided to look up a project"))
	}

	if identity != "" {
		project, err := client.Projects().GetProject(ctx, identity)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error getting project by id: %w", err))
		}
		if project == nil {
			return diag.FromErr(fmt.Errorf("no project found with id '%s'", identity))
		}
		return setProjectDataSourceState(d, project)
	}

	projects, err := client.Projects().ListProjects(ctx, &tcprojects.ListProjectsRequest{})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error listing projects: %w", err))
	}

	var matches []tcprojects.Project
	for _, p := range projects {
		if name != "" && slug != "" {
			if p.Name == name && p.Slug == slug {
				matches = append(matches, p)
			}
			continue
		}
		if name != "" && p.Name == name {
			matches = append(matches, p)
			continue
		}
		if slug != "" && p.Slug == slug {
			matches = append(matches, p)
		}
	}

	if len(matches) == 0 {
		switch {
		case name != "" && slug != "":
			return diag.FromErr(fmt.Errorf("no project found with name '%s' and slug '%s'", name, slug))
		case name != "":
			return diag.FromErr(fmt.Errorf("no project found with name '%s'", name))
		default:
			return diag.FromErr(fmt.Errorf("no project found with slug '%s'", slug))
		}
	}
	if len(matches) > 1 {
		var slugs []string
		for _, p := range matches {
			slugs = append(slugs, p.Slug)
		}
		return diag.FromErr(fmt.Errorf("multiple projects found matching the lookup criteria, please specify one of these slugs: %v", slugs))
	}

	return setProjectDataSourceState(d, &matches[0])
}

func setProjectDataSourceState(d *schema.ResourceData, project *tcprojects.Project) diag.Diagnostics {
	d.SetId(project.Identity)
	_ = d.Set("id", project.Identity)
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
