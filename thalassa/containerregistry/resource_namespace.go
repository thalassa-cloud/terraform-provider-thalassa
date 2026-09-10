package containerregistry

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	tcregistry "github.com/thalassa-cloud/client-go/containerregistry"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func ResourceNamespace() *schema.Resource {
	return &schema.Resource{
		Description:   "Create and manage a container registry namespace in Thalassa Cloud",
		CreateContext: resourceNamespaceCreate,
		ReadContext:   resourceNamespaceRead,
		UpdateContext: resourceNamespaceUpdate,
		DeleteContext: resourceNamespaceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity of the container registry namespace",
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Organisation of the namespace. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"region": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StringLenBetween(1, 255),
				Description:  "Region slug where the namespace is created (e.g. nl-01)",
			},
			"namespace": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StringLenBetween(1, 255),
				Description:  "Name of the container registry namespace",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				ValidateFunc: validate.StringLenBetween(0, 1024),
				Description:  "Human-readable description of the namespace",
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Labels for the namespace",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Annotations for the namespace",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation timestamp (RFC3339)",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last update timestamp (RFC3339)",
			},
			"object_version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Object version of the namespace",
			},
			"total_size_bytes": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total size of images in the namespace in bytes",
			},
		},
	}
}

func resourceNamespaceCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	createReq := tcregistry.CreateContainerRegistryNamespaceRequest{
		Region:      d.Get("region").(string),
		Namespace:   d.Get("namespace").(string),
		Description: d.Get("description").(string),
		Labels:      tcregistry.Labels(convert.ConvertToMap(d.Get("labels"))),
		Annotations: tcregistry.Annotations(convert.ConvertToMap(d.Get("annotations"))),
	}

	ns, err := client.ContainerRegistry().CreateContainerRegistryNamespace(ctx, createReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating container registry namespace: %w", err))
	}
	if ns == nil {
		return diag.FromErr(fmt.Errorf("error creating container registry namespace: empty response"))
	}

	d.SetId(ns.Identity)
	return resourceNamespaceRead(ctx, d, m)
}

func resourceNamespaceRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	ns, err := client.ContainerRegistry().GetContainerRegistryNamespace(ctx, d.Id())
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error getting container registry namespace: %w", err))
	}
	if ns == nil {
		d.SetId("")
		return nil
	}

	return setNamespaceResourceState(d, ns)
}

func setNamespaceResourceState(d *schema.ResourceData, ns *tcregistry.ContainerRegistryNamespace) diag.Diagnostics {
	d.SetId(ns.Identity)
	_ = d.Set("namespace", ns.Namespace)
	_ = d.Set("description", ns.Description)
	_ = d.Set("labels", ns.Labels)
	_ = d.Set("annotations", ns.Annotations)
	_ = d.Set("created_at", ns.CreatedAt.Format(TimeFormatRFC3339))
	_ = d.Set("updated_at", ns.UpdatedAt.Format(TimeFormatRFC3339))
	_ = d.Set("object_version", ns.ObjectVersion)
	_ = d.Set("total_size_bytes", ns.TotalSizeBytes)
	if ns.Region != nil && ns.Region.Slug != "" {
		_ = d.Set("region", ns.Region.Slug)
	}
	return nil
}

func resourceNamespaceUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	updateReq := tcregistry.UpdateContainerRegistryNamespaceRequest{
		Description: d.Get("description").(string),
		Labels:      tcregistry.Labels(convert.ConvertToMap(d.Get("labels"))),
		Annotations: tcregistry.Annotations(convert.ConvertToMap(d.Get("annotations"))),
	}

	ns, err := client.ContainerRegistry().UpdateContainerRegistryNamespace(ctx, d.Id(), updateReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating container registry namespace: %w", err))
	}
	if ns != nil {
		if diags := setNamespaceResourceState(d, ns); diags != nil {
			return diags
		}
	}

	return resourceNamespaceRead(ctx, d, m)
}

func resourceNamespaceDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := client.ContainerRegistry().DeleteContainerRegistryNamespace(ctx, d.Id()); err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error deleting container registry namespace: %w", err))
	}

	d.SetId("")
	return nil
}
