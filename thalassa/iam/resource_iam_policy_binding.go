package iam

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	iam "github.com/thalassa-cloud/client-go/iam"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func ResourceIamPolicyBinding() *schema.Resource {
	return &schema.Resource{
		Description: "Manage a binding of a principal to a project IAM policy in Thalassa Cloud. " +
			"Requires projects to be enabled for the organisation (currently private beta).",
		CreateContext: resourceIamPolicyBindingCreate,
		ReadContext:   resourceIamPolicyBindingRead,
		UpdateContext: resourceIamPolicyBindingUpdate,
		DeleteContext: resourceIamPolicyBindingDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity of the IAM policy binding",
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Organisation of the IAM policy binding. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Identity of the IAM policy to bind",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StringLenBetween(1, 255),
				Description:  "Name of the IAM policy binding",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				ValidateFunc: validate.StringLenBetween(0, 255),
				Description:  "Description of the IAM policy binding",
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Labels for the IAM policy binding",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Annotations for the IAM policy binding",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"user_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"user_id", "service_account_id"},
				Description:  "Identity of the user to bind to this IAM policy",
			},
			"service_account_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"user_id", "service_account_id"},
				Description:  "Identity of the service account to bind to this IAM policy",
			},
			"expires_at": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"duration"},
				Description:   "RFC3339 timestamp when the binding expires. Conflicts with duration.",
			},
			"duration": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"expires_at"},
				Description:   "Relative duration until the binding expires (for example 24h). Conflicts with expires_at.",
			},
			"slug": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Slug of the IAM policy binding",
			},
			"object_version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Object version of the IAM policy binding",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation timestamp of the IAM policy binding",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last update timestamp of the IAM policy binding",
			},
		},
	}
}

func resourceIamPolicyBindingCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	policyIdentity := d.Get("policy_id").(string)

	createReq := iam.CreateIamPolicyBindingRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Labels:      convert.ConvertToMap(d.Get("labels")),
		Annotations: convert.ConvertToMap(d.Get("annotations")),
	}

	if userIdentity, ok := d.Get("user_id").(string); ok && userIdentity != "" {
		createReq.UserIdentity = &userIdentity
	}
	if serviceAccountIdentity, ok := d.Get("service_account_id").(string); ok && serviceAccountIdentity != "" {
		createReq.ServiceAccountIdentity = &serviceAccountIdentity
	}

	if expiresAt, ok := d.Get("expires_at").(string); ok && expiresAt != "" {
		parsed, err := time.Parse(time.RFC3339, expiresAt)
		if err != nil {
			return diag.FromErr(fmt.Errorf("expires_at must be a valid RFC3339 timestamp: %s", err))
		}
		createReq.ExpiresAt = &parsed
	}
	if duration, ok := d.Get("duration").(string); ok && duration != "" {
		createReq.Duration = duration
	}

	binding, err := client.IAM().CreateIamPolicyBinding(ctx, policyIdentity, createReq)
	if err != nil {
		return diag.FromErr(err)
	}
	if binding == nil {
		return diag.FromErr(fmt.Errorf("failed to create IAM policy binding"))
	}

	d.SetId(binding.Identity)
	return resourceIamPolicyBindingRead(ctx, d, m)
}

func resourceIamPolicyBindingRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	policyIdentity := d.Get("policy_id").(string)
	bindingIdentity := d.Id()

	bindings, err := client.IAM().ListIamPolicyBindings(ctx, policyIdentity, &iam.ListIamPolicyBindingsRequest{})
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	var binding *iam.IamPolicyBinding
	for i := range bindings {
		if bindings[i].Identity == bindingIdentity {
			binding = &bindings[i]
			break
		}
	}
	if binding == nil {
		d.SetId("")
		return nil
	}

	d.SetId(binding.Identity)
	_ = d.Set("name", binding.Name)
	_ = d.Set("slug", binding.Slug)
	_ = d.Set("description", binding.Description)
	_ = d.Set("labels", binding.Labels)
	_ = d.Set("annotations", binding.Annotations)
	_ = d.Set("object_version", binding.ObjectVersion)
	_ = d.Set("created_at", binding.CreatedAt.Format(TimeFormatRFC3339))
	if binding.UpdatedAt != nil {
		_ = d.Set("updated_at", binding.UpdatedAt.Format(TimeFormatRFC3339))
	}
	if binding.ExpiresAt != nil {
		_ = d.Set("expires_at", binding.ExpiresAt.Format(TimeFormatRFC3339))
	}

	if binding.User != nil {
		_ = d.Set("user_id", binding.User.Subject)
	}
	if binding.ServiceAccount != nil {
		_ = d.Set("service_account_id", binding.ServiceAccount.Identity)
	}

	return nil
}

func resourceIamPolicyBindingUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	policyIdentity := d.Get("policy_id").(string)
	bindingIdentity := d.Id()

	updateReq := iam.UpdateIamPolicyBindingRequest{
		Description: d.Get("description").(string),
		Labels:      convert.ConvertToMap(d.Get("labels")),
		Annotations: convert.ConvertToMap(d.Get("annotations")),
	}

	_, err = client.IAM().UpdateIamPolicyBinding(ctx, policyIdentity, bindingIdentity, updateReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating IAM policy binding: %s", err))
	}

	return resourceIamPolicyBindingRead(ctx, d, m)
}

func resourceIamPolicyBindingDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	policyIdentity := d.Get("policy_id").(string)
	bindingIdentity := d.Id()

	err = client.IAM().DeleteIamPolicyBinding(ctx, policyIdentity, bindingIdentity)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
