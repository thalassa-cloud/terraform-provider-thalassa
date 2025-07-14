package iam

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"

	iam "github.com/thalassa-cloud/client-go/iam"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

func ResourceRole() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage an organisation role in Thalassa Cloud",
		CreateContext: resourceRoleCreate,
		ReadContext:   resourceRoleRead,
		DeleteContext: resourceRoleDelete,
		// UpdateContext: resourceRoleUpdate,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.StringLenBetween(1, 255),
				Description:  "Name of the organisation role",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.StringLenBetween(0, 255),
				Description:  "Description of the organisation role",
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Labels for the organisation role",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Annotations for the organisation role",
			},
			"slug": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Slug of the organisation role",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation timestamp of the organisation role",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last update timestamp of the organisation role",
			},
			"role_is_read_only": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the role is read-only and cannot be modified.",
			},
			"system": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the role is a system role",
			},
		},
	}
}

func resourceRoleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	createReq := iam.CreateOrganisationRoleRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Labels:      convert.ConvertToMap(d.Get("labels")),
		Annotations: convert.ConvertToMap(d.Get("annotations")),
	}

	role, err := client.IAM().CreateOrganisationRole(ctx, createReq)
	if err != nil {
		return diag.FromErr(err)
	}
	if role != nil {
		d.SetId(role.Identity)
		return resourceRoleRead(ctx, d, m)
	}
	return diag.FromErr(fmt.Errorf("failed to create organisation role"))
}

func resourceRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	identity := d.Get("id").(string)
	role, err := client.IAM().GetOrganisationRole(ctx, identity)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	if role == nil {
		d.SetId("")
		return nil
	}

	d.SetId(role.Identity)
	d.Set("name", role.Name)
	d.Set("slug", role.Slug)
	d.Set("description", role.Description)
	d.Set("labels", role.Labels)
	d.Set("annotations", role.Annotations)
	d.Set("created_at", role.CreatedAt.Format(time.RFC3339))
	d.Set("updated_at", role.UpdatedAt.Format(time.RFC3339))
	d.Set("role_is_read_only", role.IsReadOnly)
	d.Set("system", role.System)
	return nil
}

func resourceRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}
	identity := d.Get("id").(string)
	err = client.IAM().DeleteOrganisationRole(ctx, identity)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

// func resourceRoleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
// 	client, err := provider.GetClient(provider.GetProvider(m), d)
// 	if err != nil {
// 		return diag.FromErr(err)
// 	}
// 	identity := d.Get("id").(string)
// 	updateReq := iam.UpdateOrganisationRoleRequest{
// 		Name:        d.Get("name").(string),
// 		Description: d.Get("description").(string),
// 		Labels:      convert.ConvertToMap(d.Get("labels")),
// 		Annotations: convert.ConvertToMap(d.Get("annotations")),
// 	}
// 	role, err := client.IAM().UpdateOrganisationRole(ctx, identity, updateReq)
// 	if err != nil {
// 		return diag.FromErr(err)
// 	}
// 	if role != nil {
// 		d.SetId(role.Identity)
// 		return resourceRoleRead(ctx, d, m)
// 	}
// 	return diag.FromErr(fmt.Errorf("failed to update organisation role"))

// }
