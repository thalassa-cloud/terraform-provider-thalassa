package thalassa

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	iaas "github.com/thalassa-cloud/client-go/iaas"
	kubernetes "github.com/thalassa-cloud/client-go/kubernetes"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

func resourceKubernetesCluster() *schema.Resource {
	return &schema.Resource{
		Description:   "Create an Kubernetes Cluster",
		CreateContext: resourceKubernetesClusterCreate,
		ReadContext:   resourceKubernetesClusterRead,
		UpdateContext: resourceKubernetesClusterUpdate,
		DeleteContext: resourceKubernetesClusterDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the Kubernetes Cluster. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Region of the Kubernetes Cluster. Required for hosted-control-plane clusters.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the Kubernetes Cluster",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "VPC of the Kubernetes Cluster. This is automatically set when a subnet is provided.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Subnet of the Kubernetes Cluster. Required for managed clusters.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the Kubernetes Cluster",
			},
			"slug": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A human readable description about the Kubernetes Cluster",
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Labels for the Kubernetes Cluster",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Annotations for the Kubernetes Cluster",
			},
			"cluster_version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cluster version of the Kubernetes Cluster",
			},
			"cluster_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "managed",
				ForceNew:     true,
				ValidateFunc: validate.StringInSlice([]string{"managed", "hosted-control-plane"}, false),
				Description:  "Cluster type of the Kubernetes Cluster",
			},
			"delete_protection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Delete protection of the Kubernetes Cluster",
			},
			"networking_cni": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StringInSlice([]string{"cilium", "custom"}, false),
				Description:  "CNI of the Kubernetes Cluster",
			},
			"networking_service_cidr": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IsCIDR,
				Description:  "Service CIDR of the Kubernetes Cluster",
			},
			"networking_pod_cidr": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IsCIDR,
				Description:  "Pod CIDR of the Kubernetes Cluster",
			},
			"pod_security_standards_profile": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "baseline",
				ValidateFunc: validate.StringInSlice([]string{"restricted", "baseline", "privileged"}, false),
				Description:  "Pod security standards profile of the Kubernetes Cluster",
			},
			"audit_log_profile": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "none",
				ValidateFunc: validate.StringInSlice([]string{"none", "basic", "advanced"}, false),
				Description:  "Audit log profile of the Kubernetes Cluster",
			},
			"default_network_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "deny-all",
				ValidateFunc: validate.StringInSlice([]string{"", "allow-all", "deny-all"}, false),
				Description:  "Default network policy of the Kubernetes Cluster",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceKubernetesClusterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := getClient(getProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.Get("cluster_type").(string) == "hosted-control-plane" {
		if d.Get("region").(string) == "" {
			return diag.FromErr(fmt.Errorf("region is required for hosted-control-plane clusters"))
		}
	}

	if d.Get("cluster_type").(string) == "managed" {
		if d.Get("subnet_id").(string) == "" {
			return diag.FromErr(fmt.Errorf("subnet is required for managed clusters"))
		}
	}

	version := d.Get("cluster_version").(string)

	// Get version from API
	kubernetesVersions, err := client.Kubernetes().ListKubernetesVersions(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, kv := range kubernetesVersions {
		if !kv.Enabled { // skip disabled versions
			continue
		}

		if kv.Identity == version || kv.Slug == version || kv.Name == version {
			version = kv.Identity
			break
		}
	}

	region := d.Get("region").(string)
	regions, err := client.IaaS().ListRegions(ctx, &iaas.ListRegionsRequest{})
	if err != nil {
		return diag.FromErr(err)
	}
	for _, r := range regions {
		if r.Identity == region || r.Slug == region || r.Name == region {
			region = r.Identity
			break
		}
	}

	createKubernetesCluster := kubernetes.CreateKubernetesCluster{
		Name:                      d.Get("name").(string),
		Description:               d.Get("description").(string),
		Labels:                    convertToMap(d.Get("labels")),
		Annotations:               convertToMap(d.Get("annotations")),
		Subnet:                    d.Get("subnet_id").(string),
		RegionIdentity:            region,
		DeleteProtection:          d.Get("delete_protection").(bool),
		ClusterType:               kubernetes.KubernetesClusterType(d.Get("cluster_type").(string)),
		KubernetesVersionIdentity: version,
		Networking: kubernetes.KubernetesClusterNetworking{
			CNI:         d.Get("networking_cni").(string),
			ServiceCIDR: d.Get("networking_service_cidr").(string),
			PodCIDR:     d.Get("networking_pod_cidr").(string),
		},
		PodSecurityStandardsProfile: kubernetes.KubernetesClusterPodSecurityStandards(d.Get("pod_security_standards_profile").(string)),
		AuditLogProfile:             kubernetes.KubernetesClusterAuditLoggingProfile(d.Get("audit_log_profile").(string)),
		DefaultNetworkPolicy:        kubernetes.KubernetesDefaultNetworkPolicies(d.Get("default_network_policy").(string)),
	}

	kubernetesCluster, err := client.Kubernetes().CreateKubernetesCluster(ctx, createKubernetesCluster)
	if err != nil {
		return diag.FromErr(err)
	}
	if kubernetesCluster != nil {
		d.SetId(kubernetesCluster.Identity)
		d.Set("slug", kubernetesCluster.Slug)
		d.Set("status", kubernetesCluster.Status)
	}

	// wait until the cluster is ready
	for {
		kubernetesCluster, err := client.Kubernetes().GetKubernetesCluster(ctx, kubernetesCluster.Identity)
		if err != nil {
			return diag.FromErr(err)
		}
		if strings.EqualFold(kubernetesCluster.Status, "ready") {
			break
		}
		time.Sleep(1 * time.Second)
	}

	return resourceKubernetesClusterRead(ctx, d, m)
}

func resourceKubernetesClusterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := getClient(getProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	slug := d.Get("id").(string)
	kubernetesCluster, err := client.Kubernetes().GetKubernetesCluster(ctx, slug)
	if err != nil && !tcclient.IsNotFound(err) {
		return diag.FromErr(fmt.Errorf("error getting kubernetesCluster: %s", err))
	}
	if kubernetesCluster == nil {
		return diag.FromErr(fmt.Errorf("kubernetesCluster was not found"))
	}

	currentlyConfiguredVersion := d.Get("cluster_version").(string)
	if !(kubernetesCluster.ClusterVersion.Name == currentlyConfiguredVersion || kubernetesCluster.ClusterVersion.Slug == currentlyConfiguredVersion || kubernetesCluster.ClusterVersion.Identity == currentlyConfiguredVersion) {
		d.Set("cluster_version", kubernetesCluster.ClusterVersion.Slug)
	}

	d.SetId(kubernetesCluster.Identity)
	d.Set("name", kubernetesCluster.Name)
	d.Set("slug", kubernetesCluster.Slug)
	d.Set("description", kubernetesCluster.Description)
	d.Set("labels", kubernetesCluster.Labels)
	d.Set("annotations", kubernetesCluster.Annotations)
	d.Set("cluster_type", kubernetesCluster.ClusterType)
	d.Set("delete_protection", kubernetesCluster.DeleteProtection)
	d.Set("networking_cni", kubernetesCluster.Configuration.Networking.CNI)
	d.Set("networking_service_cidr", kubernetesCluster.Configuration.Networking.ServiceCIDR)
	d.Set("networking_pod_cidr", kubernetesCluster.Configuration.Networking.PodCIDR)
	d.Set("pod_security_standards_profile", kubernetesCluster.PodSecurityStandardsProfile)
	d.Set("audit_log_profile", kubernetesCluster.AuditLogProfile)
	d.Set("default_network_policy", kubernetesCluster.DefaultNetworkPolicy)
	d.Set("status", kubernetesCluster.Status)

	if kubernetesCluster.VPC != nil {
		d.Set("vpc", kubernetesCluster.VPC.Identity)
	}
	if kubernetesCluster.Subnet != nil {
		d.Set("subnet_id", kubernetesCluster.Subnet.Identity)
	}
	if kubernetesCluster.Region != nil {
		d.Set("region", kubernetesCluster.Region.Identity)
	}

	return nil
}

func resourceKubernetesClusterUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := getClient(getProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	version := d.Get("cluster_version").(string)
	// Get version from API
	kubernetesVersions, err := client.Kubernetes().ListKubernetesVersions(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, kv := range kubernetesVersions {
		if !kv.Enabled { // skip disabled versions
			continue
		}

		if kv.Identity == version || kv.Slug == version || kv.Name == version {
			version = kv.Identity
			break
		}
	}

	updateKubernetesCluster := kubernetes.UpdateKubernetesCluster{
		Name:                        Ptr(d.Get("name").(string)),
		Description:                 Ptr(d.Get("description").(string)),
		Labels:                      convertToMap(d.Get("labels")),
		Annotations:                 convertToMap(d.Get("annotations")),
		DeleteProtection:            Ptr(d.Get("delete_protection").(bool)),
		KubernetesVersionIdentity:   Ptr(version),
		PodSecurityStandardsProfile: Ptr(kubernetes.KubernetesClusterPodSecurityStandards(d.Get("pod_security_standards_profile").(string))),
		AuditLogProfile:             Ptr(kubernetes.KubernetesClusterAuditLoggingProfile(d.Get("audit_log_profile").(string))),
		DefaultNetworkPolicy:        Ptr(kubernetes.KubernetesDefaultNetworkPolicies(d.Get("default_network_policy").(string))),
	}
	identity := d.Get("id").(string)
	kubernetesCluster, err := client.Kubernetes().UpdateKubernetesCluster(ctx, identity, updateKubernetesCluster)
	if err != nil {
		return diag.FromErr(err)
	}
	if kubernetesCluster != nil {
		d.Set("name", kubernetesCluster.Name)
		d.Set("description", kubernetesCluster.Description)
		d.Set("slug", kubernetesCluster.Slug)
		d.Set("status", kubernetesCluster.Status)
		if kubernetesCluster.VPC != nil {
			d.Set("vpc_id", kubernetesCluster.VPC.Identity)
		}
		if kubernetesCluster.Subnet != nil {
			d.Set("subnet_id", kubernetesCluster.Subnet.Identity)
		}

		d.Set("labels", kubernetesCluster.Labels)
		d.Set("annotations", kubernetesCluster.Annotations)
		d.Set("cluster_version", kubernetesCluster.ClusterVersion.Identity)
		d.Set("cluster_type", kubernetesCluster.ClusterType)
		d.Set("delete_protection", kubernetesCluster.DeleteProtection)
		d.Set("networking_cni", kubernetesCluster.Configuration.Networking.CNI)
		d.Set("networking_service_cidr", kubernetesCluster.Configuration.Networking.ServiceCIDR)
		d.Set("networking_pod_cidr", kubernetesCluster.Configuration.Networking.PodCIDR)
		d.Set("pod_security_standards_profile", kubernetesCluster.PodSecurityStandardsProfile)
		d.Set("audit_log_profile", kubernetesCluster.AuditLogProfile)
		d.Set("default_network_policy", kubernetesCluster.DefaultNetworkPolicy)

		return nil
	}

	return resourceKubernetesClusterRead(ctx, d, m)
}

func resourceKubernetesClusterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := getClient(getProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	identity := d.Get("id").(string)

	err = client.Kubernetes().DeleteKubernetesCluster(ctx, identity)
	if err != nil {
		if !tcclient.IsNotFound(err) {
			return diag.FromErr(err)
		}
	}

	// wait until the cluster is deleted
	for {
		kubernetesCluster, err := client.Kubernetes().GetKubernetesCluster(ctx, identity)
		if err != nil {
			if tcclient.IsNotFound(err) {
				break
			}
			return diag.FromErr(err)
		}
		if kubernetesCluster == nil {
			break
		}
		if strings.EqualFold(kubernetesCluster.Status, "deleted") {
			break
		}
		time.Sleep(1 * time.Second)
	}

	d.SetId("")
	return nil
}
