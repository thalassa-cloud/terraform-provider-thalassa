package iam

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"

	iam "github.com/thalassa-cloud/client-go/iam"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

func ResourceRoleBinding() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a role binding in Thalassa Cloud",
		CreateContext: resourceRoleBindingCreate,
		ReadContext:   resourceRoleBindingRead,
		UpdateContext: resourceRoleBindingUpdate,
		DeleteContext: resourceRoleBindingDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Identity of the role to bind",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the role binding",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the role binding",
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Labels for the role binding",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Annotations for the role binding",
			},
			"user_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Identity of the user to bind to this role",
			},
			"team_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Identity of the team to bind to this role",
			},
			"service_account_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Identity of the service account to bind to this role",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation timestamp of the role binding",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last update timestamp of the role binding",
			},
		},
	}
}

func resourceRoleBindingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	roleIdentity := d.Get("role_id").(string)

	createReq := iam.CreateRoleBinding{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Labels:      convert.ConvertToMap(d.Get("labels")),
		Annotations: convert.ConvertToMap(d.Get("annotations")),
	}

	// Set the appropriate identity field
	if userIdentity, ok := d.Get("user_id").(string); ok && userIdentity != "" {
		createReq.UserIdentity = &userIdentity
	}
	if teamIdentity, ok := d.Get("team_id").(string); ok && teamIdentity != "" {
		createReq.TeamIdentity = &teamIdentity
	}
	if serviceAccountIdentity, ok := d.Get("service_account_id").(string); ok && serviceAccountIdentity != "" {
		createReq.ServiceAccountIdentity = &serviceAccountIdentity
	}

	binding, err := client.IAM().CreateRoleBinding(ctx, roleIdentity, createReq)
	if err != nil {
		return diag.FromErr(err)
	}

	if binding != nil {
		d.SetId(binding.Identity)
		return resourceRoleBindingRead(ctx, d, m)
	}

	return diag.FromErr(fmt.Errorf("failed to create role binding"))
}

func resourceRoleBindingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	roleIdentity := d.Get("role_id").(string)
	bindingIdentity := d.Get("id").(string)

	// Get role bindings for the role
	bindings, err := client.IAM().ListRoleBindings(ctx, roleIdentity, &iam.ListRoleBindingsRequest{})
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	// Find the specific binding
	var binding *iam.OrganisationRoleBinding
	for _, b := range bindings {
		if b.Identity == bindingIdentity {
			binding = &b
			break
		}
	}

	if binding == nil {
		d.SetId("")
		return nil
	}

	d.SetId(binding.Identity)
	d.Set("name", binding.Name)
	d.Set("description", binding.Description)
	d.Set("labels", binding.Labels)
	d.Set("annotations", binding.Annotations)
	d.Set("created_at", binding.CreatedAt.Format(time.RFC3339))
	d.Set("updated_at", binding.UpdatedAt.Format(time.RFC3339))

	// Set the appropriate identity field based on what's bound
	if binding.AppUser != nil {
		d.Set("user_id", binding.AppUser.Subject)
	}
	if binding.OrganisationTeam != nil {
		d.Set("team_id", binding.OrganisationTeam.Identity)
	}
	if binding.ServiceAccount != nil {
		d.Set("service_account_id", binding.ServiceAccount.Identity)
	}

	return nil
}

func resourceRoleBindingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Role bindings are immutable in the API, so we need to recreate them
	// This is a common pattern for resources that don't support updates
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Role binding updates not supported",
			Detail:   "Role bindings cannot be updated. Please delete and recreate the binding with the new configuration.",
		},
	}
}

func resourceRoleBindingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	roleIdentity := d.Get("role_id").(string)
	bindingIdentity := d.Get("id").(string)

	err = client.IAM().DeleteRoleBinding(ctx, roleIdentity, bindingIdentity)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
