package kubernetes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thalassa-cloud/client-go/kubernetes"

	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func DataSourceKubernetesCluster() *schema.Resource {
	return &schema.Resource{
		Description: "Get an Kubernetes cluster",
		ReadContext: dataSourceKubernetesClusterRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Kubernetes version.",
			},
			"slug": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The slug of the Kubernetes version.",
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Organisation of the Kubernetes Cluster",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Region of the Kubernetes Cluster. Required for hosted-control-plane clusters.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Subnet of the Kubernetes Cluster.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "VPC of the Kubernetes Cluster.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A human readable description about the Kubernetes Cluster",
			},
			"labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Labels for the Kubernetes Cluster",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Annotations for the Kubernetes Cluster",
			},
			"cluster_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cluster version of the Kubernetes Cluster",
			},
			"cluster_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cluster type of the Kubernetes Cluster",
			},
			"delete_protection": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Delete protection of the Kubernetes Cluster",
			},
			"networking_cni": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CNI of the Kubernetes Cluster",
			},
			"networking_service_cidr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Service CIDR of the Kubernetes Cluster",
			},
			"networking_pod_cidr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Pod CIDR of the Kubernetes Cluster",
			},
			"pod_security_standards_profile": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Pod security standards profile of the Kubernetes Cluster",
			},
			"audit_log_profile": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Audit log profile of the Kubernetes Cluster",
			},
			"default_network_policy": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Default network policy of the Kubernetes Cluster",
			},
			"kubernetes_api_server_endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Kubernetes API server endpoint of the Kubernetes Cluster",
			},
			"kubernetes_api_server_ca_certificate": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Kubernetes API server CA certificate of the Kubernetes Cluster",
			},
			"api_server_acls": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "API server ACLs for the Kubernetes Cluster",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allowed_cidrs": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of allowed CIDRs for API server access",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"auto_upgrade_policy": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Auto upgrade policy of the Kubernetes Cluster",
			},
			"maintenance_day": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Day of the week when the cluster will be upgraded (0-6, where 0 is Sunday)",
			},
			"maintenance_start_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Time of day when the cluster will be upgraded in minutes from midnight",
			},
		},
	}
}

func dataSourceKubernetesClusterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := provider.GetProvider(m)
	slug := d.Get("slug").(string)

	clusters, err := provider.Client.Kubernetes().ListKubernetesClusters(ctx, &kubernetes.ListKubernetesClustersRequest{})
	if err != nil {
		return diag.FromErr(err)
	}

	for _, cluster := range clusters {
		if slug != "" && cluster.Slug == slug {
			d.SetId(cluster.Identity)
			d.Set("id", cluster.Identity)
			d.Set("name", cluster.Name)
			d.Set("slug", cluster.Slug)
			d.Set("description", cluster.Description)
			d.Set("cluster_version", cluster.ClusterVersion.Name)
			d.Set("cluster_type", cluster.ClusterType)
			d.Set("delete_protection", cluster.DeleteProtection)
			d.Set("networking_cni", cluster.Configuration.Networking.CNI)
			d.Set("networking_service_cidr", cluster.Configuration.Networking.ServiceCIDR)
			d.Set("networking_pod_cidr", cluster.Configuration.Networking.PodCIDR)
			d.Set("pod_security_standards_profile", cluster.PodSecurityStandardsProfile)
			d.Set("audit_log_profile", cluster.AuditLogProfile)
			d.Set("default_network_policy", cluster.DefaultNetworkPolicy)
			if cluster.Region != nil {
				d.Set("region", cluster.Region.Identity)
			}
			if cluster.Subnet != nil {
				d.Set("subnet_id", cluster.Subnet.Identity)
			}
			if cluster.VPC != nil {
				d.Set("vpc_id", cluster.VPC.Identity)
			}
			d.Set("kubernetes_api_server_endpoint", cluster.APIServerURL)
			d.Set("kubernetes_api_server_ca_certificate", cluster.APIServerCA)

			// Set API server ACLs
			if cluster.ApiServerACLs.AllowedCIDRs != nil {
				apiServerACLs := map[string]interface{}{
					"allowed_cidrs": cluster.ApiServerACLs.AllowedCIDRs,
				}
				if err := d.Set("api_server_acls", []interface{}{apiServerACLs}); err != nil {
					return diag.FromErr(fmt.Errorf("error setting api_server_acls: %s", err))
				}
			}

			// Set auto upgrade policy
			d.Set("auto_upgrade_policy", cluster.AutoUpgradePolicy)

			// Set maintenance settings
			if cluster.MaintenanceDay != nil {
				d.Set("maintenance_day", int(*cluster.MaintenanceDay))
			}
			if cluster.MaintenanceStartAt != nil {
				d.Set("maintenance_start_at", int(*cluster.MaintenanceStartAt))
			}

			// Set labels and annotations directly
			if err := d.Set("labels", cluster.Labels); err != nil {
				return diag.FromErr(fmt.Errorf("error setting labels: %s", err))
			}

			if err := d.Set("annotations", cluster.Annotations); err != nil {
				return diag.FromErr(fmt.Errorf("error setting annotations: %s", err))
			}

			return diag.Diagnostics{}
		}
	}
	return diag.FromErr(fmt.Errorf("not found"))
}
