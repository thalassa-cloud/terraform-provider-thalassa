package iam

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	iam "github.com/thalassa-cloud/client-go/iam"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func DataSourceIamResourceAccessBindings() *schema.Resource {
	return &schema.Resource{
		Description: "List principals with IAM access to a resource identity in Thalassa Cloud. " +
			"Requires projects to be enabled for the organisation (currently private beta).",
		ReadContext: dataSourceIamResourceAccessBindingsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Data source identity (resource_identity)",
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Organisation for the query. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"resource_identity": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Identity of the resource to list access bindings for",
			},
			"resource_types": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				Description: "Resource types to include in the access binding query",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"permissions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional permission filter (create, read, update, delete, list, *)",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"binding_types": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional binding type filter",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"bindings": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Access bindings for the resource",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource type of the binding",
						},
						"binding_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the binding",
						},
						"binding_identity": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identity of the binding",
						},
						"policy_identity": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identity of the IAM policy that grants access",
						},
						"policy_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the IAM policy that grants access",
						},
						"principal_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Principal type (user or service_account)",
						},
						"principal_identity": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identity of the principal",
						},
						"principal_display_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Display name of the principal",
						},
						"permissions": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Permissions granted by the binding",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"expires_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Expiry timestamp of the binding, if set",
						},
					},
				},
			},
		},
	}
}

func dataSourceIamResourceAccessBindingsRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceIdentity := d.Get("resource_identity").(string)

	resourceTypes := make([]string, 0)
	if rtList, ok := d.Get("resource_types").([]any); ok {
		for _, rt := range rtList {
			resourceTypes = append(resourceTypes, rt.(string))
		}
	}

	req := iam.ListResourceAccessBindingsRequest{
		ResourceIdentity: resourceIdentity,
		ResourceTypes:    resourceTypes,
	}

	if pList, ok := d.Get("permissions").([]any); ok && len(pList) > 0 {
		permissions := make([]string, 0, len(pList))
		for _, p := range pList {
			permissions = append(permissions, p.(string))
		}
		req.Permissions = convertStringsToPermissions(permissions)
	}

	if btList, ok := d.Get("binding_types").([]any); ok && len(btList) > 0 {
		bindingTypes := make([]string, 0, len(btList))
		for _, bt := range btList {
			bindingTypes = append(bindingTypes, bt.(string))
		}
		req.BindingTypes = bindingTypes
	}

	result, err := client.IAM().ListResourceAccessBindings(ctx, req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error listing resource access bindings: %s", err))
	}
	if result == nil {
		return diag.FromErr(fmt.Errorf("empty response when listing resource access bindings"))
	}

	bindings := make([]map[string]any, 0, len(result.Bindings))
	for _, item := range result.Bindings {
		binding := map[string]any{
			"resource_type":          item.ResourceType,
			"binding_type":           item.BindingType,
			"binding_identity":       item.BindingIdentity,
			"policy_identity":        item.PolicyIdentity,
			"policy_name":            item.PolicyName,
			"principal_type":         string(item.PrincipalType),
			"principal_identity":     item.PrincipalIdentity,
			"principal_display_name": item.PrincipalDisplayName,
			"permissions":            toListOfInterfaces(item.Permissions),
		}
		if item.ExpiresAt != nil {
			binding["expires_at"] = item.ExpiresAt.Format(TimeFormatRFC3339)
		}
		bindings = append(bindings, binding)
	}

	d.SetId(resourceIdentity)
	_ = d.Set("bindings", bindings)

	return nil
}
