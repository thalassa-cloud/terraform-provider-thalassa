package iam

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/client-go/thalassa"
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
			"members": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of team members",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_identity": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Identity of the user to add to the team",
						},
						"email": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Email address of the user to add to the team. If provided, user_identity will be resolved automatically.",
						},
						"role": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "Role of the team member. Optional. Default: MEMBER.",
							ValidateFunc: validate.StringInSlice([]string{"OWNER", "ADMIN", "MEMBER"}, false),
						},
					},
				},
				Set: func(v interface{}) int {
					m := v.(map[string]interface{})
					userIdentity := m["user_identity"].(string)
					email := m["email"].(string)
					role := m["role"].(string)

					// Set default role if not provided
					if role == "" {
						role = "MEMBER"
					}

					// Use email as primary identifier if available, otherwise use user_identity
					identifier := email
					if identifier == "" {
						identifier = userIdentity
					}

					return schema.HashString(fmt.Sprintf("%s-%s", identifier, role))
				},
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

		err = updateTeamMembers(ctx, client, team.Identity, d)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error updating team members for newly created team: %s", err))
		}
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
	d.Set("created_at", team.CreatedAt.Format(TimeFormatRFC3339))
	if team.UpdatedAt != nil {
		d.Set("updated_at", team.UpdatedAt.Format(TimeFormatRFC3339))
	}

	// Set members data
	memberSet := schema.NewSet(func(v interface{}) int {
		m := v.(map[string]interface{})
		userIdentity := m["user_identity"].(string)
		email := m["email"].(string)
		role := m["role"].(string)

		// Set default role if not provided
		if role == "" {
			role = "MEMBER"
		}

		// Use email as primary identifier if available, otherwise use user_identity
		identifier := email
		if identifier == "" {
			identifier = userIdentity
		}

		return schema.HashString(fmt.Sprintf("%s-%s", identifier, role))
	}, []interface{}{})

	for _, member := range team.Members {
		role := member.Role
		if role == "" {
			role = "MEMBER" // Set default role if empty
		}

		memberMap := map[string]interface{}{
			"user_identity": member.User.Subject, // Using Subject as the user identity
			"email":         member.User.Email,   // Include email for reference
			"role":          role,
		}
		memberSet.Add(memberMap)
	}
	d.Set("members", memberSet)

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
		return diag.FromErr(fmt.Errorf("error updating team: %s", err))
	}
	if team != nil {
		d.Set("name", team.Name)
		d.Set("description", team.Description)
		d.Set("slug", team.Slug)
		d.Set("labels", team.Labels)
		d.Set("annotations", team.Annotations)
		if team.UpdatedAt != nil {
			d.Set("updated_at", team.UpdatedAt.Format(TimeFormatRFC3339))
		}
	}

	// Handle member changes
	if d.HasChange("members") {
		err := updateTeamMembers(ctx, client, identity, d)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error updating team members for updated team: %s", err))
		}
	}
	return resourceTeamRead(ctx, d, m)
}

func resolveUserIdentity(ctx context.Context, client thalassa.Client, email string) (string, error) {
	members, err := client.IAM().ListOrganisationMembers(ctx, &iam.ListMembersRequest{})
	if err != nil {
		return "", fmt.Errorf("error listing organisation members: %s", err)
	}

	for _, member := range members {
		if member.User != nil && member.User.Email == email {
			return member.User.Subject, nil
		}
	}

	return "", fmt.Errorf("no user found with email address: %s", email)
}

func updateTeamMembers(ctx context.Context, client thalassa.Client, teamID string, d *schema.ResourceData) error {
	oldMembers, newMembers := d.GetChange("members")
	oldSet := oldMembers.(*schema.Set)
	newSet := newMembers.(*schema.Set)

	// Find the team member identity by user identity
	// We need to get the current team to find the member identity
	team, err := client.IAM().GetTeam(ctx, teamID, &iam.GetTeamRequest{})
	if err != nil {
		return fmt.Errorf("error getting team to find member identity: %s", err)
	}

	currentTeamMembers := team.Members

	// Find members to remove
	toRemove := oldSet.Difference(newSet)
	for _, member := range toRemove.List() {
		memberMap := member.(map[string]interface{})
		userIdentity := memberMap["user_identity"].(string)
		email := memberMap["email"].(string)

		// If we have an email but no user_identity, resolve it
		if email != "" && userIdentity == "" {
			resolvedIdentity, err := resolveUserIdentity(ctx, client, email)
			if err != nil {
				// If we can't resolve the email, try to find by email directly
				for _, m := range currentTeamMembers {
					if m.User.Email == email {
						userIdentity = m.User.Subject
						break
					}
				}
			} else {
				userIdentity = resolvedIdentity
			}
		}

		var memberIdentity string
		for _, m := range currentTeamMembers {
			if m.User.Subject == userIdentity || (email != "" && m.User.Email == email) {
				memberIdentity = m.Identity
				break
			}
		}

		if memberIdentity != "" {
			err = client.IAM().RemoveTeamMember(ctx, teamID, memberIdentity)
			if err != nil {
				return fmt.Errorf("error removing team member: %s", err)
			}
		}
	}

	// Find members to add
	toAdd := newSet.Difference(oldSet)
	for _, member := range toAdd.List() {
		memberMap := member.(map[string]interface{})
		userIdentity := memberMap["user_identity"].(string)
		email := memberMap["email"].(string)
		role := memberMap["role"].(string)

		// Validate that either user_identity or email is provided
		if userIdentity == "" && email == "" {
			return fmt.Errorf("either user_identity or email must be provided for team member")
		}

		// Set default role if not provided
		if role == "" {
			role = "MEMBER"
		}

		// Resolve email to user identity if needed
		if userIdentity == "" {
			resolvedIdentity, err := resolveUserIdentity(ctx, client, email)
			if err != nil {
				return fmt.Errorf("error resolving email to user identity: %s", err)
			}
			userIdentity = resolvedIdentity
		}

		addRequest := iam.AddTeamMemberRequest{
			UserIdentity: userIdentity,
			Role:         role,
		}

		err := client.IAM().AddTeamMember(ctx, teamID, addRequest)
		if err != nil {
			return fmt.Errorf("error adding team member: %s", err)
		}
	}

	return nil
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
