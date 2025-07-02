package iaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	iaas "github.com/thalassa-cloud/client-go/iaas"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func resourceCloudInitTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudInitTemplateCreate,
		ReadContext:   resourceCloudInitTemplateRead,
		// UpdateContext: resourceCloudInitTemplateUpdate,
		DeleteContext: resourceCloudInitTemplateDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"slug": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the Machine Type. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"annotations": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"content": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
		},
	}

}

func resourceCloudInitTemplateCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	create := iaas.CreateCloudInitTemplateRequest{
		Name:    d.Get("name").(string),
		Content: d.Get("content").(string),
	}

	if labels, ok := d.GetOk("labels"); ok {
		create.Labels = convertLabels(labels.(map[string]interface{}))
	}
	if annotations, ok := d.GetOk("annotations"); ok {
		create.Annotations = convertAnnotations(annotations.(map[string]interface{}))
	}

	cloudInitTemplate, err := client.IaaS().CreateCloudInitTemplate(ctx, create)
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

	return resourceCloudInitTemplateRead(ctx, d, m)
}

func resourceCloudInitTemplateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInitTemplate, err := client.IaaS().GetCloudInitTemplate(ctx, d.Id())
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

func resourceCloudInitTemplateDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.IaaS().DeleteCloudInitTemplate(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func convertLabels(labels map[string]interface{}) map[string]string {
	convertedLabels := make(map[string]string)
	for key, value := range labels {
		convertedLabels[key] = value.(string)
	}
	return convertedLabels
}

func convertAnnotations(annotations map[string]interface{}) map[string]string {
	convertedAnnotations := make(map[string]string)
	for key, value := range annotations {
		convertedAnnotations[key] = value.(string)
	}
	return convertedAnnotations
}
