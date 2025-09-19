package iam

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"

	iam "github.com/thalassa-cloud/client-go/iam"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

func ResourceServiceAccount() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a service account in Thalassa Cloud",
		CreateContext: resourceServiceAccountCreate,
		ReadContext:   resourceServiceAccountRead,
		UpdateContext: resourceServiceAccountUpdate,
		DeleteContext: resourceServiceAccountDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.StringLenBetween(1, 255),
				Description:  "Name of the service account",
			},
			"slug": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Slug of the service account",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.StringLenBetween(0, 255),
				Description:  "Description of the service account",
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Labels for the service account",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Annotations for the service account",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation timestamp of the service account",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last update timestamp of the service account",
			},
			"object_version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Object version of the service account",
			},
			"role_bindings": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Role bindings for the service account",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"identity": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identity of the role binding",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the role binding",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the role binding",
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
						"role_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identity of the role binding",
						},
					},
				},
			},
		},
	}
}

func resourceServiceAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	createReq := iam.CreateServiceAccountRequest{
		Name:        d.Get("name").(string),
		Description: convert.Ptr(d.Get("description").(string)),
		Labels:      convert.ConvertToMap(d.Get("labels")),
		Annotations: convert.ConvertToMap(d.Get("annotations")),
	}

	account, err := client.IAM().CreateServiceAccount(ctx, createReq)
	if err != nil {
		return diag.FromErr(err)
	}

	if account != nil {
		d.SetId(account.Identity)
		return resourceServiceAccountRead(ctx, d, m)
	}

	return diag.FromErr(fmt.Errorf("failed to create service account"))
}

func resourceServiceAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	identity := d.Get("id").(string)
	account, err := client.IAM().GetServiceAccount(ctx, identity)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if account == nil {
		d.SetId("")
		return nil
	}

	d.SetId(account.Identity)
	d.Set("name", account.Name)
	d.Set("slug", account.Slug)
	d.Set("description", account.Description)
	d.Set("labels", account.Labels)
	d.Set("annotations", account.Annotations)
	d.Set("created_at", account.CreatedAt.Format(TimeFormatRFC3339))
	d.Set("object_version", account.ObjectVersion)

	if account.UpdatedAt != nil {
		d.Set("updated_at", account.UpdatedAt.Format(TimeFormatRFC3339))
	}

	// Set role bindings
	roleBindingsList := make([]map[string]interface{}, len(account.RoleBindings))
	for i, binding := range account.RoleBindings {
		roleId := ""
		if binding.OrganisationRole != nil {
			roleId = binding.OrganisationRole.Identity
		}
		bindingMap := map[string]interface{}{
			"identity":    binding.Identity,
			"name":        binding.Name,
			"description": binding.Description,
			"created_at":  binding.CreatedAt.Format(TimeFormatRFC3339),
			"updated_at":  binding.UpdatedAt.Format(TimeFormatRFC3339),
			"role_id":     roleId,
		}
		roleBindingsList[i] = bindingMap
	}
	d.Set("role_bindings", roleBindingsList)

	return nil
}

func resourceServiceAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	identity := d.Get("id").(string)

	updateReq := iam.UpdateServiceAccountRequest{
		Name:        convert.Ptr(d.Get("name").(string)),
		Description: convert.Ptr(d.Get("description").(string)),
		Labels:      convert.ConvertToMap(d.Get("labels")),
		Annotations: convert.ConvertToMap(d.Get("annotations")),
	}

	account, err := client.IAM().UpdateServiceAccount(ctx, identity, updateReq)
	if err != nil {
		return diag.FromErr(err)
	}

	if account != nil {
		return resourceServiceAccountRead(ctx, d, m)
	}

	return diag.FromErr(fmt.Errorf("failed to update service account"))
}

func resourceServiceAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	identity := d.Get("id").(string)
	err = client.IAM().DeleteServiceAccount(ctx, identity)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
