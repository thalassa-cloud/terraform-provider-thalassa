package iam

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	iam "github.com/thalassa-cloud/client-go/iam"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func ResourceIamPolicy() *schema.Resource {
	return &schema.Resource{
		Description: "Manage a project IAM policy in Thalassa Cloud. " +
			"Requires projects to be enabled for the organisation (currently private beta).",
		CreateContext: resourceIamPolicyCreate,
		ReadContext:   resourceIamPolicyRead,
		UpdateContext: resourceIamPolicyUpdate,
		DeleteContext: resourceIamPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity of the IAM policy",
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Organisation of the IAM policy. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StringLenBetween(1, 255),
				Description:  "Name of the IAM policy",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				ValidateFunc: validate.StringLenBetween(0, 255),
				Description:  "Description of the IAM policy",
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Labels for the IAM policy",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Annotations for the IAM policy",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"replicate_to_children": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "When true, replicate this policy to child projects",
			},
			"slug": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Slug of the IAM policy",
			},
			"system": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the IAM policy is a system policy",
			},
			"is_read_only": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the IAM policy is read-only and cannot be modified",
			},
			"object_version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Object version of the IAM policy",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation timestamp of the IAM policy",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last update timestamp of the IAM policy",
			},
		},
	}
}

func resourceIamPolicyCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	createReq := iam.CreateIamPolicyRequest{
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		Labels:              convert.ConvertToMap(d.Get("labels")),
		Annotations:         convert.ConvertToMap(d.Get("annotations")),
		ReplicateToChildren: d.Get("replicate_to_children").(bool),
	}

	policy, err := client.IAM().CreateIamPolicy(ctx, createReq)
	if err != nil {
		return diag.FromErr(err)
	}
	if policy == nil {
		return diag.FromErr(fmt.Errorf("failed to create IAM policy"))
	}

	d.SetId(policy.Identity)
	return resourceIamPolicyRead(ctx, d, m)
}

func resourceIamPolicyRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	identity := d.Id()
	policy, err := client.IAM().GetIamPolicy(ctx, identity)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	if policy == nil {
		d.SetId("")
		return nil
	}

	d.SetId(policy.Identity)
	_ = d.Set("name", policy.Name)
	_ = d.Set("slug", policy.Slug)
	_ = d.Set("description", policy.Description)
	_ = d.Set("labels", policy.Labels)
	_ = d.Set("annotations", policy.Annotations)
	_ = d.Set("replicate_to_children", policy.ReplicateToChildren)
	_ = d.Set("system", policy.System)
	_ = d.Set("is_read_only", policy.IsReadOnly)
	_ = d.Set("object_version", policy.ObjectVersion)
	_ = d.Set("created_at", policy.CreatedAt.Format(TimeFormatRFC3339))
	if policy.UpdatedAt != nil {
		_ = d.Set("updated_at", policy.UpdatedAt.Format(TimeFormatRFC3339))
	}

	return nil
}

func resourceIamPolicyUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	identity := d.Id()
	updateReq := iam.UpdateIamPolicyRequest{
		Description: d.Get("description").(string),
		Labels:      convert.ConvertToMap(d.Get("labels")),
		Annotations: convert.ConvertToMap(d.Get("annotations")),
	}

	_, err = client.IAM().UpdateIamPolicy(ctx, identity, updateReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating IAM policy: %s", err))
	}

	return resourceIamPolicyRead(ctx, d, m)
}

func resourceIamPolicyDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	identity := d.Id()
	err = client.IAM().DeleteIamPolicy(ctx, identity)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error deleting IAM policy: %s", err))
	}

	d.SetId("")
	return nil
}
