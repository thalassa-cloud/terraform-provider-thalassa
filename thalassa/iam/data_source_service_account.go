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

func DataSourceServiceAccount() *schema.Resource {
	return &schema.Resource{
		Description: "Get a service account",
		ReadContext: dataSourceServiceAccountRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity of the service account",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.StringLenBetween(1, 255),
				Description:  "Name of the service account",
			},
			"slug": {
				Type:        schema.TypeString,
				Optional:    true,
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
				Computed:    true,
				Description: "Labels for the service account",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Computed:    true,
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
					},
				},
			},
		},
	}
}

func dataSourceServiceAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	slug := d.Get("slug").(string)
	name := d.Get("name").(string)

	if name == "" && slug == "" {
		return diag.FromErr(fmt.Errorf("either 'name' or 'slug' must be provided to look up a service account"))
	}

	accounts, err := client.IAM().ListServiceAccounts(ctx, &iam.ListServiceAccountsRequest{})
	if err != nil {
		return diag.FromErr(err)
	}

	var account *iam.ServiceAccount

	if name != "" && slug != "" {
		// Both name and slug provided: match both
		for _, a := range accounts {
			if a.Name == name && a.Slug == slug {
				account = &a
				break
			}
		}
		if account == nil {
			return diag.FromErr(fmt.Errorf("no service account found with name '%s' and slug '%s'", name, slug))
		}
	} else if name != "" {
		// Only name provided: match all by name
		var matchingAccounts []iam.ServiceAccount
		for _, a := range accounts {
			if a.Name == name {
				matchingAccounts = append(matchingAccounts, a)
			}
		}
		if len(matchingAccounts) == 0 {
			return diag.FromErr(fmt.Errorf("no service account found with name '%s'", name))
		}
		if len(matchingAccounts) > 1 {
			var slugs []string
			for _, a := range matchingAccounts {
				slugs = append(slugs, a.Slug)
			}
			return diag.FromErr(fmt.Errorf("multiple service accounts found with name '%s', please specify one of these slugs: %v", name, slugs))
		}
		account = &matchingAccounts[0]
	} else if slug != "" {
		// Only slug provided: match by slug
		for _, a := range accounts {
			if a.Slug == slug {
				account = &a
				break
			}
		}
		if account == nil {
			return diag.FromErr(fmt.Errorf("no service account found with slug '%s'", slug))
		}
	}

	// Set fields if found
	d.SetId(account.Identity)
	d.Set("id", account.Identity)
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

	return diag.Diagnostics{}
}
