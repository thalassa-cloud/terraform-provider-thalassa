package iaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func dataSourceCloudInitTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudInitTemplateRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"labels": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the Cloud Init Template. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"annotations": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"content": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"slug": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCloudInitTemplateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInitTemplate, err := client.IaaS().GetCloudInitTemplate(ctx, d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(cloudInitTemplate.Identity)
	d.Set("name", cloudInitTemplate.Name)
	d.Set("content", cloudInitTemplate.Content)
	d.Set("slug", cloudInitTemplate.Slug)

	if cloudInitTemplate.Labels != nil {
		d.Set("labels", cloudInitTemplate.Labels)
	}
	if cloudInitTemplate.Annotations != nil {
		d.Set("annotations", cloudInitTemplate.Annotations)
	}

	return nil
}
