package kubernetes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thalassa-cloud/client-go/kubernetes"

	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func DataSourceKubernetesVersion() *schema.Resource {
	return &schema.Resource{
		Description: "Get an Kubernetes version",
		ReadContext: dataSourceKubernetesVersionRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the Kubernetes version.",
			},
			"slug": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The slug of the Kubernetes version.",
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the Kubernetes Version. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"kubernetes_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Kubernetes version",
			},
			"containerd_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The containerd version.",
			},
			"cni_plugins_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CNI plugins version.",
			},
			"crictl_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crictl version.",
			},
			"runc_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The runc version.",
			},
			"cilium_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The cilium version.",
			},
			"cloud_controller_manager_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The cloud controller manager version.",
			},
			"istio_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The istio version.",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "The annotations of the Kubernetes version.",
			},
		},
	}
}

func dataSourceKubernetesVersionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := provider.GetProvider(m)
	name := d.Get("name").(string)
	slug := d.Get("slug").(string)
	id := d.Get("id").(string)

	versions, err := provider.Client.Kubernetes().ListKubernetesVersions(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	var version *kubernetes.KubernetesVersion
	for _, v := range versions {
		if name != "" && v.Name == name {
			version = &v
		}
		if slug != "" && v.Slug == slug {
			version = &v
		}
		if id != "" && v.Identity == id {
			version = &v
		}
	}

	if version == nil {
		return diag.FromErr(fmt.Errorf("version not found"))
	}

	d.SetId(version.Identity)
	d.Set("id", version.Identity)
	d.Set("kubernetes_version", version.KubernetesVersion)
	d.Set("containerd_version", version.ContainerdVersion)
	d.Set("cni_plugins_version", version.CNIPluginsVersion)
	d.Set("crictl_version", version.CrictlVersion)
	d.Set("runc_version", version.RuncVersion)
	d.Set("cilium_version", version.CiliumVersion)
	d.Set("cloud_controller_manager_version", version.CloudControllerManagerVersion)
	d.Set("istio_version", version.IstioVersion)

	if err := d.Set("annotations", version.Annotations); err != nil {
		return diag.FromErr(fmt.Errorf("error setting annotations: %s", err))
	}

	return diag.Diagnostics{}
}
