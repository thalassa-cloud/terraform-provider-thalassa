package containerregistry

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	tcregistry "github.com/thalassa-cloud/client-go/containerregistry"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func DataSourceNamespace() *schema.Resource {
	return &schema.Resource{
		Description: "Get a container registry namespace by identity or by region and namespace name",
		ReadContext: dataSourceNamespaceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
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
				Optional:     true,
				ValidateFunc: validate.StringLenBetween(1, 255),
				Description:  "Region slug of the namespace (required when looking up by namespace name)",
			},
			"namespace": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.StringLenBetween(1, 255),
				Description:  "Name of the container registry namespace",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Human-readable description of the namespace",
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

func dataSourceNamespaceRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	identity := d.Get("id").(string)
	region := d.Get("region").(string)
	namespace := d.Get("namespace").(string)

	if identity == "" && (region == "" || namespace == "") {
		return diag.FromErr(fmt.Errorf("either 'id' or both 'region' and 'namespace' must be provided to look up a container registry namespace"))
	}

	if identity != "" {
		ns, err := client.ContainerRegistry().GetContainerRegistryNamespace(ctx, identity)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error getting container registry namespace by id: %w", err))
		}
		if ns == nil {
			return diag.FromErr(fmt.Errorf("no container registry namespace found with id '%s'", identity))
		}
		return setNamespaceDataSourceState(d, ns)
	}

	namespaces, err := client.ContainerRegistry().ListContainerRegistryNamespaces(ctx, &tcregistry.ListContainerRegistryNamespacesRequest{})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error listing container registry namespaces: %w", err))
	}

	var match *tcregistry.ContainerRegistryNamespace
	for i := range namespaces {
		ns := &namespaces[i]
		if ns.Namespace != namespace {
			continue
		}
		if ns.Region == nil || ns.Region.Slug != region {
			continue
		}
		if match != nil {
			return diag.FromErr(fmt.Errorf("multiple container registry namespaces found for region '%s' and namespace '%s'", region, namespace))
		}
		match = ns
	}
	if match == nil {
		return diag.FromErr(fmt.Errorf("no container registry namespace found for region '%s' and namespace '%s'", region, namespace))
	}

	return setNamespaceDataSourceState(d, match)
}

func setNamespaceDataSourceState(d *schema.ResourceData, ns *tcregistry.ContainerRegistryNamespace) diag.Diagnostics {
	d.SetId(ns.Identity)
	_ = d.Set("id", ns.Identity)
	_ = d.Set("namespace", ns.Namespace)
	_ = d.Set("description", ns.Description)
	_ = d.Set("created_at", ns.CreatedAt.Format(TimeFormatRFC3339))
	_ = d.Set("updated_at", ns.UpdatedAt.Format(TimeFormatRFC3339))
	_ = d.Set("object_version", ns.ObjectVersion)
	_ = d.Set("total_size_bytes", ns.TotalSizeBytes)
	if ns.Region != nil {
		switch {
		case ns.Region.Slug != "":
			_ = d.Set("region", ns.Region.Slug)
		case ns.Region.Name != "":
			_ = d.Set("region", ns.Region.Name)
		}
	}
	return nil
}
