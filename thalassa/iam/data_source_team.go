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

func DataSourceTeam() *schema.Resource {
	return &schema.Resource{
		Description: "Get a team",
		ReadContext: dataSourceTeamRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity of the Team",
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the Team. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.StringLenBetween(1, 255),
				Description:  "Name of the Team",
			},
			"slug": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Slug of the Team",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.StringLenBetween(0, 255),
				Description:  "A human readable description about the team",
			},
			"labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Labels for the Team",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Annotations for the Team",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation timestamp of the Team",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last update timestamp of the Team",
			},
			"members": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of team members",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"identity": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "team membershhip identity",
						},
						"role": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Role of the team member",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation timestamp of the team member",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last update timestamp of the team member",
						},
						"user_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subject identifier of the user",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the user",
						},
						"email": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Email address of the user",
						},
					},
				},
			},
		},
	}
}

func dataSourceTeamRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	slug := d.Get("slug").(string)
	name := d.Get("name").(string)

	if name == "" && slug == "" {
		return diag.FromErr(fmt.Errorf("either 'name' or 'slug' must be provided to look up a team"))
	}

	teams, err := client.IAM().ListTeams(ctx, &iam.ListTeamsRequest{})
	if err != nil {
		return diag.FromErr(err)
	}

	var team *iam.Team

	if name != "" && slug != "" {
		// Both name and slug provided: match both
		for _, t := range teams {
			if t.Name == name && t.Slug == slug {
				team = &t
				break
			}
		}
		if team == nil {
			return diag.FromErr(fmt.Errorf("no team found with name '%s' and slug '%s'", name, slug))
		}
	} else if name != "" {
		// Only name provided: match all by name
		var matchingTeams []iam.Team
		for _, t := range teams {
			if t.Name == name {
				matchingTeams = append(matchingTeams, t)
			}
		}
		if len(matchingTeams) == 0 {
			return diag.FromErr(fmt.Errorf("no team found with name '%s'", name))
		}
		if len(matchingTeams) > 1 {
			var slugs []string
			for _, t := range matchingTeams {
				slugs = append(slugs, t.Slug)
			}
			return diag.FromErr(fmt.Errorf("multiple teams found with name '%s', please specify one of these slugs: %v", name, slugs))
		}
		team = &matchingTeams[0]
	} else if slug != "" {
		// Only slug provided: match by slug
		for _, t := range teams {
			if t.Slug == slug {
				team = &t
				break
			}
		}
		if team == nil {
			return diag.FromErr(fmt.Errorf("no team found with slug '%s'", slug))
		}
	}

	// Set fields if found
	d.SetId(team.Identity)
	d.Set("id", team.Identity)
	d.Set("name", team.Name)
	d.Set("slug", team.Slug)
	d.Set("description", team.Description)
	d.Set("labels", team.Labels)
	d.Set("annotations", team.Annotations)
	d.Set("created_at", team.CreatedAt.Format(TimeFormatRFC3339))
	if team.UpdatedAt != nil {
		d.Set("updated_at", team.UpdatedAt.Format(TimeFormatRFC3339))
	}

	// Set members data
	memberList := make([]map[string]interface{}, len(team.Members))
	for i, member := range team.Members {
		memberMap := map[string]interface{}{
			"identity":   member.Identity,
			"role":       member.Role,
			"created_at": member.CreatedAt.Format(TimeFormatRFC3339),
			"updated_at": member.UpdatedAt.Format(TimeFormatRFC3339),
			"user_id":    member.User.Subject,
			"name":       member.User.Name,
			"email":      member.User.Email,
		}
		memberList[i] = memberMap
	}
	d.Set("members", memberList)

	return diag.Diagnostics{}
}
