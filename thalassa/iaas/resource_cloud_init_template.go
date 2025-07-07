package iaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	iaas "github.com/thalassa-cloud/client-go/iaas"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"

	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
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
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The identity of the cloud init template",
			},
			"slug": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The slug of the cloud init template",
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the Machine Type. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Labels to add to the cloud init template",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Annotations to add to the cloud init template",
			},
			"name": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The name of the cloud init template",
			},
			"content": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The content of the cloud init template",
			},
		},
	}

}

func resourceCloudInitTemplateCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating client: %s", err))
	}

	create := iaas.CreateCloudInitTemplateRequest{
		Name:    d.Get("name").(string),
		Content: d.Get("content").(string),
	}

	if labels, ok := d.GetOk("labels"); ok {
		create.Labels = convert.ConvertToMap(labels.(map[string]interface{}))
	}
	if annotations, ok := d.GetOk("annotations"); ok {
		create.Annotations = convert.ConvertToMap(annotations.(map[string]interface{}))
	}

	cloudInitTemplate, err := client.IaaS().CreateCloudInitTemplate(ctx, create)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating cloud init template: %s", err))
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
		return diag.FromErr(fmt.Errorf("error creating client: %s", err))
	}

	cloudInitTemplate, err := client.IaaS().GetCloudInitTemplate(ctx, d.Id())
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error reading cloud init template: %s", err))
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
		return diag.FromErr(fmt.Errorf("error creating client: %s", err))
	}

	err = client.IaaS().DeleteCloudInitTemplate(ctx, d.Id())
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error deleting cloud init template: %s", err))
	}

	d.SetId("")
	return nil
}
