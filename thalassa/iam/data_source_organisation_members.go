package iam

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"

	iam "github.com/thalassa-cloud/client-go/iam"
)

func DataSourceOrganisationMembers() *schema.Resource {
	return &schema.Resource{
		Description: "Get all members of an organisation, optionally filtered by email address",
		ReadContext: dataSourceOrganisationMembersRead,
		Schema: map[string]*schema.Schema{
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"email_filter": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Filter members by email address. If provided, only members with matching email addresses will be returned.",
			},
			"members": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of organisation members",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"identity": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identity of the organisation member",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation timestamp of the organisation member",
						},
						"member_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the organisation member (OWNER or MEMBER)",
						},
						"user": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "User information for the organisation member",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"subject": {
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
									"created_at": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Creation timestamp of the user",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceOrganisationMembersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	members, err := client.IAM().ListOrganisationMembers(ctx, &iam.ListMembersRequest{})
	if err != nil {
		return diag.FromErr(err)
	}

	// Apply email filter if provided
	emailFilter := d.Get("email_filter").(string)
	var filteredMembers []iam.OrganisationMember
	if emailFilter != "" {
		for _, member := range members {
			if member.User != nil && member.User.Email == emailFilter {
				filteredMembers = append(filteredMembers, member)
			}
		}
		members = filteredMembers
	}

	// Convert members to Terraform schema format
	memberList := make([]map[string]interface{}, len(members))
	for i, member := range members {
		memberMap := map[string]interface{}{
			"identity":    member.Identity,
			"created_at":  member.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			"member_type": string(member.MemberType),
		}

		// Add user information if available
		if member.User != nil {
			userMap := map[string]interface{}{
				"subject":    member.User.Subject,
				"name":       member.User.Name,
				"email":      member.User.Email,
				"created_at": member.User.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			}
			memberMap["user"] = []map[string]interface{}{userMap}
		}

		memberList[i] = memberMap
	}

	// Set the data
	resourceID := fmt.Sprintf("organisation-members-%d", len(members))
	if emailFilter != "" {
		resourceID = fmt.Sprintf("organisation-members-%s-%d", emailFilter, len(members))
	}
	d.SetId(resourceID)
	d.Set("members", memberList)

	return diag.Diagnostics{}
}
