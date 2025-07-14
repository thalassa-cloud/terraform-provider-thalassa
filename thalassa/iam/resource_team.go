package iam

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"

	iam "github.com/thalassa-cloud/client-go/iam"
)

func ResourceTeam() *schema.Resource {
	return &schema.Resource{
		Description:   "Create a team in the Thalassa Cloud platform",
		CreateContext: resourceTeamCreate,
		ReadContext:   resourceTeamRead,
		UpdateContext: resourceTeamUpdate,
		DeleteContext: resourceTeamDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the Team. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.StringLenBetween(1, 255),
				Description:  "Name of the Team",
			},
			"slug": {
				Type:        schema.TypeString,
				Computed:    true,
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
				Optional:    true,
				Description: "Labels for the Team",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
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
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTeamCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	createTeam := iam.CreateTeam{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Labels:      convert.ConvertToMap(d.Get("labels")),
		Annotations: convert.ConvertToMap(d.Get("annotations")),
	}

	team, err := client.IAM().CreateTeam(ctx, createTeam)
	if err != nil {
		return diag.FromErr(err)
	}
	if team != nil {
		d.SetId(team.Identity)
		d.Set("slug", team.Slug)
		return nil
	}
	return resourceTeamRead(ctx, d, m)
}

func resourceTeamRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	identity := d.Get("id").(string)
	team, err := client.IAM().GetTeam(ctx, identity, &iam.GetTeamRequest{})
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error getting team: %s", err))
	}
	if team == nil {
		d.SetId("")
		return nil
	}

	d.SetId(team.Identity)
	d.Set("name", team.Name)
	d.Set("slug", team.Slug)
	d.Set("description", team.Description)
	d.Set("labels", team.Labels)
	d.Set("annotations", team.Annotations)
	d.Set("created_at", team.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
	if team.UpdatedAt != nil {
		d.Set("updated_at", team.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))
	}
	return nil
}

func resourceTeamUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	updateTeam := iam.UpdateTeam{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Labels:      convert.ConvertToMap(d.Get("labels")),
		Annotations: convert.ConvertToMap(d.Get("annotations")),
	}

	identity := d.Get("id").(string)

	team, err := client.IAM().UpdateTeam(ctx, identity, updateTeam)
	if err != nil {
		return diag.FromErr(err)
	}
	if team != nil {
		d.Set("name", team.Name)
		d.Set("description", team.Description)
		d.Set("slug", team.Slug)
		d.Set("labels", team.Labels)
		d.Set("annotations", team.Annotations)
		if team.UpdatedAt != nil {
			d.Set("updated_at", team.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))
		}
		return nil
	}

	return resourceTeamRead(ctx, d, m)
}

func resourceTeamDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	identity := d.Get("id").(string)
	err = client.IAM().DeleteTeam(ctx, identity)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error deleting team: %s", err))
	}

	d.SetId("")
	return nil
}
