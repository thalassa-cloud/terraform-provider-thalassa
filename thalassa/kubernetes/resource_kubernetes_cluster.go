package kubernetes

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
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func resourceKubernetesCluster() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages a Kubernetes cluster in the Thalassa cloud platform. This resource supports both managed clusters and hosted control plane clusters, allowing you to deploy production-ready Kubernetes environments with configurable networking, security policies, and auto-upgrade capabilities. The cluster can be customized with specific CNI plugins (Cilium or custom), network CIDRs, pod security standards, audit logging, and API server access controls.",
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
				Optional:    true,
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
				Optional:    true,
				Description: "Cluster version of the Kubernetes Cluster, can be a name, slug or identity of the Kubernetes version. If not provided, the latest stable version will be used for provisioning.",
			},
			"cluster_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "managed",
				ForceNew:     true,
				ValidateFunc: validate.StringInSlice([]string{"managed", "hosted-control-plane"}, false),
				Description:  "Cluster type of the Kubernetes Cluster. Must be one of: managed, hosted-control-plane. Default: managed.",
			},
			"delete_protection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Delete protection of the Kubernetes Cluster",
			},
			"networking_cni": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "cilium", // Default to Cilium
				ValidateFunc: validate.StringInSlice([]string{"cilium", "custom"}, false),
				Description:  "CNI of the Kubernetes Cluster",
			},
			"networking_service_cidr": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validate.IsCIDR,
				Default:      "172.16.0.0/18",
				Description:  "Service CIDR of the Kubernetes Cluster. Must be a valid CIDR block.",
			},
			"networking_pod_cidr": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validate.IsCIDR,
				Default:      "192.168.0.0/16",
				Description:  "Pod CIDR of the Kubernetes Cluster. Must be a valid CIDR block.",
			},
			"autoscaler_config": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Configuration for the cluster autoscaler. These values can also be configured using annotations on a KubernetesNodePool object.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scale_down_disabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Flag to disable the scale down of node pools by the cluster autoscaler",
						},
						"scale_down_delay_after_add": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Delay after adding a node to the node pool by the cluster autoscaler",
						},
						"estimator": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Estimator to use for the cluster autoscaler. Available values: binpacking",
						},
						"expander": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Expander to use for the cluster autoscaler",
						},
						"ignore_daemonsets_utilization": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Flag to ignore the utilization of daemonsets by the cluster autoscaler",
						},
						"balance_similar_node_groups": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Flag to balance the utilization of similar node groups by the cluster autoscaler",
						},
						"expendable_pods_priority_cutoff": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Priority cutoff for the expendable pods by the cluster autoscaler",
						},
						"scale_down_unneeded_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Time after which a node can be scaled down by the cluster autoscaler",
						},
						"scale_down_utilization_threshold": {
							Type:        schema.TypeFloat,
							Optional:    true,
							Description: "Utilization threshold for the cluster autoscaler. The autoscaler might scale down non-empty nodes with utilization below a threshold. To prevent this behavior, set the utilization threshold to 0",
						},
						"max_graceful_termination_sec": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Maximum graceful termination time for the cluster autoscaler. If the pod is not stopped within this time then the node is terminated anyway.",
						},
						"enable_proactive_scale_up": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Flag to enable the proactive scale up of the cluster autoscaler. Whether to enable/disable proactive scale-ups, defaults to false",
						},
					},
				},
			},
			"pod_security_standards_profile": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "baseline",
				ValidateFunc: validate.StringInSlice([]string{"restricted", "baseline", "privileged"}, false),
				Description:  "Pod security standards profile of the Kubernetes Cluster. Must be one of: restricted, baseline, privileged. Default: baseline.",
			},
			"audit_log_profile": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "none",
				ValidateFunc: validate.StringInSlice([]string{"none", "basic", "advanced"}, false),
				Description:  "Audit log profile of the Kubernetes Cluster. Must be one of: none, basic, advanced. Default: none.",
			},
			"default_network_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "deny-all",
				ValidateFunc: validate.StringInSlice([]string{"", "allow-all", "deny-all"}, false),
				Description:  "Default network policy of the Kubernetes Cluster. Must be one of: allow-all, deny-all. Default: deny-all.",
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
			"internal_endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "VPC-internal endpoint for the Kubernetes Cluster",
			},
			"advertise_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Advertise port for the Kubernetes Cluster within the VPC",
			},
			"konnectivity_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Konnectivity port for the Kubernetes Cluster within the VPC",
			},
			"api_server_acls": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    10,
				Description: "API server ACLs for the Kubernetes Cluster",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allowed_cidrs": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of allowed CIDRs for API server access",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"auto_upgrade_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "none",
				ValidateFunc: validate.StringInSlice([]string{"none", "latest-version", "latest-stable"}, false),
				Description:  "Auto upgrade policy of the Kubernetes Cluster. Must be one of: none, latest-version, latest-stable. Default: none.",
			},
			"maintenance_day": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validate.IntBetween(0, 6),
				Description:  "Day of the week when the cluster will be upgraded (0-6, where 0 is Sunday)",
			},
			"maintenance_start_at": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validate.IntBetween(0, 1439),
				Description:  "Time of day when the cluster will be upgraded in minutes from midnight (0-1439)",
			},
			"security_group_attachments": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List identities of security group that will be attached to the Kubernetes Cluster",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"disable_public_endpoint": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Disable public endpoint of the Kubernetes Cluster. When set to true, the Kubernetes Cluster will only be accessible via the private VPC endpoint and the user will need to provide a solution to access the Kubernetes API server.",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceKubernetesClusterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
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
		Labels:                    convert.ConvertToMap(d.Get("labels")),
		Annotations:               convert.ConvertToMap(d.Get("annotations")),
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
		ApiServerACLs:               convertApiServerACLs(d.Get("api_server_acls")),
		AutoUpgradePolicy:           kubernetes.KubernetesClusterAutoUpgradePolicy(d.Get("auto_upgrade_policy").(string)),
		DisablePublicEndpoint:       d.Get("disable_public_endpoint").(bool),
	}

	if securityGroupAttachments, ok := d.GetOk("security_group_attachments"); ok {
		createKubernetesCluster.SecurityGroupAttachments = convert.ConvertToStringSlice(securityGroupAttachments)
	}

	// Set autoscaler config if provided
	if autoscalerConfig, ok := d.GetOk("autoscaler_config"); ok {
		config := convertAutoscalerConfig(autoscalerConfig)
		createKubernetesCluster.AutoscalerConfig = config
	}

	// Set maintenance settings if provided
	if maintenanceDay, ok := d.GetOk("maintenance_day"); ok {
		day := uint(maintenanceDay.(int))
		createKubernetesCluster.MaintenanceDay = &day
	}
	if maintenanceStartAt, ok := d.GetOk("maintenance_start_at"); ok {
		startAt := uint(maintenanceStartAt.(int))
		createKubernetesCluster.MaintenanceStartAt = &startAt
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

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 20*time.Minute)
	defer cancel()

	// wait until the cluster is ready
	for {
		select {
		case <-ctxWithTimeout.Done():
			return diag.FromErr(fmt.Errorf("timeout waiting for cluster to be ready"))
		default:
		}
		kubernetesCluster, err := client.Kubernetes().GetKubernetesCluster(ctxWithTimeout, kubernetesCluster.Identity)
		if err != nil {
			return diag.FromErr(err)
		}
		if strings.EqualFold(kubernetesCluster.Status, "error") {
			return diag.FromErr(fmt.Errorf("cluster is in error state: %s", kubernetesCluster.StatusMessage))
		}
		if strings.EqualFold(kubernetesCluster.Status, "ready") {
			d.Set("kubernetes_api_server_endpoint", kubernetesCluster.APIServerURL)
			d.Set("kubernetes_api_server_ca_certificate", kubernetesCluster.APIServerCA)
			break
		}
		time.Sleep(1 * time.Second)
	}

	return resourceKubernetesClusterRead(ctx, d, m)
}

func resourceKubernetesClusterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	slug := d.Get("id").(string)
	kubernetesCluster, err := client.Kubernetes().GetKubernetesCluster(ctx, slug)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error getting kubernetesCluster: %s", err))
	}
	if kubernetesCluster == nil {
		return diag.FromErr(fmt.Errorf("kubernetesCluster was not found"))
	}

	// Only set cluster_version in state if it was defined in the configuration.
	// This avoids introducing a value into state that the user did not specify.
	if _, hasVersion := d.GetOk("cluster_version"); hasVersion {
		currentlyConfiguredVersion := d.Get("cluster_version").(string)
		if !(kubernetesCluster.ClusterVersion.Name == currentlyConfiguredVersion || kubernetesCluster.ClusterVersion.Slug == currentlyConfiguredVersion || kubernetesCluster.ClusterVersion.Identity == currentlyConfiguredVersion) {
			d.Set("cluster_version", kubernetesCluster.ClusterVersion.Slug)
		} else {
			d.Set("cluster_version", currentlyConfiguredVersion)
		}
	} else {
		d.Set("cluster_version", nil)
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
	d.Set("kubernetes_api_server_endpoint", kubernetesCluster.APIServerURL)
	d.Set("kubernetes_api_server_ca_certificate", kubernetesCluster.APIServerCA)
	if kubernetesCluster.InternalEndpoint != nil {
		d.Set("internal_endpoint", *kubernetesCluster.InternalEndpoint)
	}
	if kubernetesCluster.AdvertisePort != nil {
		d.Set("advertise_port", *kubernetesCluster.AdvertisePort)
	}
	if kubernetesCluster.KonnectivityPort != nil {
		d.Set("konnectivity_port", *kubernetesCluster.KonnectivityPort)
	}
	d.Set("disable_public_endpoint", kubernetesCluster.DisablePublicEndpoint)

	// Set API server ACLs
	if len(kubernetesCluster.ApiServerACLs.AllowedCIDRs) > 0 {
		apiServerACLs := map[string]interface{}{
			"allowed_cidrs": kubernetesCluster.ApiServerACLs.AllowedCIDRs,
		}
		if err := d.Set("api_server_acls", []interface{}{apiServerACLs}); err != nil {
			return diag.FromErr(fmt.Errorf("error setting api_server_acls: %s", err))
		}
	}

	securityGroupAttachments := []string{}
	for _, sg := range kubernetesCluster.SecurityGroups {
		securityGroupAttachments = append(securityGroupAttachments, sg.Identity)
	}
	d.Set("security_group_attachments", securityGroupAttachments)

	// Set auto upgrade policy
	d.Set("auto_upgrade_policy", kubernetesCluster.AutoUpgradePolicy)

	// Set maintenance settings
	if kubernetesCluster.MaintenanceDay != nil {
		d.Set("maintenance_day", int(*kubernetesCluster.MaintenanceDay))
	}
	if kubernetesCluster.MaintenanceStartAt != nil {
		d.Set("maintenance_start_at", int(*kubernetesCluster.MaintenanceStartAt))
	}

	// Set autoscaler config
	if kubernetesCluster.AutoscalerConfig != nil {
		autoscalerConfig := map[string]interface{}{
			"scale_down_disabled":              kubernetesCluster.AutoscalerConfig.ScaleDownDisabled,
			"scale_down_delay_after_add":       kubernetesCluster.AutoscalerConfig.ScaleDownDelayAfterAdd,
			"estimator":                        kubernetesCluster.AutoscalerConfig.Estimator,
			"expander":                         kubernetesCluster.AutoscalerConfig.Expander,
			"ignore_daemonsets_utilization":    kubernetesCluster.AutoscalerConfig.IgnoreDaemonsetsUtilization,
			"balance_similar_node_groups":      kubernetesCluster.AutoscalerConfig.BalanceSimilarNodeGroups,
			"expendable_pods_priority_cutoff":  kubernetesCluster.AutoscalerConfig.ExpendablePodsPriorityCutoff,
			"scale_down_unneeded_time":         kubernetesCluster.AutoscalerConfig.ScaleDownUnneededTime,
			"scale_down_utilization_threshold": kubernetesCluster.AutoscalerConfig.ScaleDownUtilizationThreshold,
			"max_graceful_termination_sec":     kubernetesCluster.AutoscalerConfig.MaxGracefulTerminationSec,
			"enable_proactive_scale_up":        kubernetesCluster.AutoscalerConfig.EnableProactiveScaleUp,
		}
		if err := d.Set("autoscaler_config", []interface{}{autoscalerConfig}); err != nil {
			return diag.FromErr(fmt.Errorf("error setting autoscaler_config: %s", err))
		}
	}

	if kubernetesCluster.VPC != nil {
		d.Set("vpc_id", kubernetesCluster.VPC.Identity)
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
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	var versionPtr *string
	if version, ok := d.GetOk("cluster_version"); ok && d.HasChange("cluster_version") {
		versionStr := version.(string)
		if versionStr != "" {
			kubernetesVersions, err := client.Kubernetes().ListKubernetesVersions(ctx)
			if err != nil {
				return diag.FromErr(err)
			}
			for _, kv := range kubernetesVersions {
				if !kv.Enabled { // skip disabled versions
					continue
				}
				if kv.Identity == versionStr || kv.Slug == versionStr || kv.Name == versionStr {
					resolved := kv.Identity
					versionPtr = &resolved
					break
				}
			}
			if versionPtr == nil {
				// If we couldn't resolve, still honor the provided value
				versionPtr = &versionStr
			}
		}
	}
	// If cluster_version is not defined or hasn't changed, versionPtr remains nil
	// and no version update will be sent to the API

	updateKubernetesCluster := kubernetes.UpdateKubernetesCluster{
		Name:                        convert.Ptr(d.Get("name").(string)),
		Description:                 convert.Ptr(d.Get("description").(string)),
		Labels:                      convert.ConvertToMap(d.Get("labels")),
		Annotations:                 convert.ConvertToMap(d.Get("annotations")),
		DeleteProtection:            convert.Ptr(d.Get("delete_protection").(bool)),
		KubernetesVersionIdentity:   versionPtr,
		PodSecurityStandardsProfile: convert.Ptr(kubernetes.KubernetesClusterPodSecurityStandards(d.Get("pod_security_standards_profile").(string))),
		AuditLogProfile:             convert.Ptr(kubernetes.KubernetesClusterAuditLoggingProfile(d.Get("audit_log_profile").(string))),
		DefaultNetworkPolicy:        convert.Ptr(kubernetes.KubernetesDefaultNetworkPolicies(d.Get("default_network_policy").(string))),
		ApiServerACLs:               convertApiServerACLs(d.Get("api_server_acls")),
		AutoUpgradePolicy:           kubernetes.KubernetesClusterAutoUpgradePolicy(d.Get("auto_upgrade_policy").(string)),
		DisablePublicEndpoint:       convert.Ptr(d.Get("disable_public_endpoint").(bool)),
	}

	if securityGroupAttachments, ok := d.GetOk("security_group_attachments"); ok {
		updateKubernetesCluster.SecurityGroupAttachments = convert.ConvertToStringSlice(securityGroupAttachments)
	}

	// Set autoscaler config if provided
	if autoscalerConfig, ok := d.GetOk("autoscaler_config"); ok {
		config := convertAutoscalerConfig(autoscalerConfig)
		updateKubernetesCluster.AutoscalerConfig = config
	}

	// Set maintenance settings if provided
	if maintenanceDay, ok := d.GetOk("maintenance_day"); ok {
		day := uint(maintenanceDay.(int))
		updateKubernetesCluster.MaintenanceDay = &day
	}
	if maintenanceStartAt, ok := d.GetOk("maintenance_start_at"); ok {
		startAt := uint(maintenanceStartAt.(int))
		updateKubernetesCluster.MaintenanceStartAt = &startAt
	}

	identity := d.Get("id").(string)
	kubernetesCluster, err := client.Kubernetes().UpdateKubernetesCluster(ctx, identity, updateKubernetesCluster)
	if err != nil {
		return diag.FromErr(err)
	}
	if kubernetesCluster != nil {

		currentlyConfiguredVersionInt, ok := d.GetOk("cluster_version")
		if ok {
			currentlyConfiguredVersion := currentlyConfiguredVersionInt.(string)
			if !(kubernetesCluster.ClusterVersion.Name == currentlyConfiguredVersion || kubernetesCluster.ClusterVersion.Slug == currentlyConfiguredVersion || kubernetesCluster.ClusterVersion.Identity == currentlyConfiguredVersion) {
				d.Set("cluster_version", kubernetesCluster.ClusterVersion.Slug)
			} else {
				d.Set("cluster_version", currentlyConfiguredVersion)
			}
		}

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
		d.Set("cluster_type", kubernetesCluster.ClusterType)
		d.Set("delete_protection", kubernetesCluster.DeleteProtection)
		d.Set("networking_cni", kubernetesCluster.Configuration.Networking.CNI)
		d.Set("networking_service_cidr", kubernetesCluster.Configuration.Networking.ServiceCIDR)
		d.Set("networking_pod_cidr", kubernetesCluster.Configuration.Networking.PodCIDR)
		d.Set("pod_security_standards_profile", kubernetesCluster.PodSecurityStandardsProfile)
		d.Set("audit_log_profile", kubernetesCluster.AuditLogProfile)
		d.Set("default_network_policy", kubernetesCluster.DefaultNetworkPolicy)
		d.Set("disable_public_endpoint", kubernetesCluster.DisablePublicEndpoint)

		// Set API server ACLs
		if len(kubernetesCluster.ApiServerACLs.AllowedCIDRs) > 0 {
			apiServerACLs := map[string]interface{}{
				"allowed_cidrs": kubernetesCluster.ApiServerACLs.AllowedCIDRs,
			}
			if err := d.Set("api_server_acls", []interface{}{apiServerACLs}); err != nil {
				return diag.FromErr(fmt.Errorf("error setting api_server_acls: %s", err))
			}
		}

		// Set auto upgrade policy
		d.Set("auto_upgrade_policy", kubernetesCluster.AutoUpgradePolicy)

		// Set maintenance settings
		if kubernetesCluster.MaintenanceDay != nil {
			d.Set("maintenance_day", int(*kubernetesCluster.MaintenanceDay))
		}
		if kubernetesCluster.MaintenanceStartAt != nil {
			d.Set("maintenance_start_at", int(*kubernetesCluster.MaintenanceStartAt))
		}

		// Set autoscaler config
		if kubernetesCluster.AutoscalerConfig != nil {
			autoscalerConfig := map[string]interface{}{
				"scale_down_disabled":              kubernetesCluster.AutoscalerConfig.ScaleDownDisabled,
				"scale_down_delay_after_add":       kubernetesCluster.AutoscalerConfig.ScaleDownDelayAfterAdd,
				"estimator":                        kubernetesCluster.AutoscalerConfig.Estimator,
				"expander":                         kubernetesCluster.AutoscalerConfig.Expander,
				"ignore_daemonsets_utilization":    kubernetesCluster.AutoscalerConfig.IgnoreDaemonsetsUtilization,
				"balance_similar_node_groups":      kubernetesCluster.AutoscalerConfig.BalanceSimilarNodeGroups,
				"expendable_pods_priority_cutoff":  kubernetesCluster.AutoscalerConfig.ExpendablePodsPriorityCutoff,
				"scale_down_unneeded_time":         kubernetesCluster.AutoscalerConfig.ScaleDownUnneededTime,
				"scale_down_utilization_threshold": kubernetesCluster.AutoscalerConfig.ScaleDownUtilizationThreshold,
				"max_graceful_termination_sec":     kubernetesCluster.AutoscalerConfig.MaxGracefulTerminationSec,
				"enable_proactive_scale_up":        kubernetesCluster.AutoscalerConfig.EnableProactiveScaleUp,
			}
			if err := d.Set("autoscaler_config", []interface{}{autoscalerConfig}); err != nil {
				return diag.FromErr(fmt.Errorf("error setting autoscaler_config: %s", err))
			}
		}

		// Set computed fields
		if kubernetesCluster.InternalEndpoint != nil {
			d.Set("internal_endpoint", *kubernetesCluster.InternalEndpoint)
		}
		if kubernetesCluster.AdvertisePort != nil {
			d.Set("advertise_port", *kubernetesCluster.AdvertisePort)
		}
		if kubernetesCluster.KonnectivityPort != nil {
			d.Set("konnectivity_port", *kubernetesCluster.KonnectivityPort)
		}

		return nil
	}

	return resourceKubernetesClusterRead(ctx, d, m)
}

func resourceKubernetesClusterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
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
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 20*time.Minute)
	defer cancel()
	for {
		select {
		case <-ctxWithTimeout.Done():
			return diag.FromErr(fmt.Errorf("timeout waiting for cluster to be deleted"))
		default:
		}
		kubernetesCluster, err := client.Kubernetes().GetKubernetesCluster(ctxWithTimeout, identity)
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

// convertApiServerACLs converts the API server ACLs from Terraform schema to the API format
func convertApiServerACLs(acls interface{}) kubernetes.KubernetesApiServerACLs {
	if acls == nil {
		return kubernetes.KubernetesApiServerACLs{}
	}

	aclsList, ok := acls.([]interface{})
	if !ok || len(aclsList) == 0 {
		return kubernetes.KubernetesApiServerACLs{}
	}

	first := aclsList[0]
	if first == nil {
		return kubernetes.KubernetesApiServerACLs{}
	}

	acl, ok := first.(map[string]interface{})
	if !ok || acl == nil {
		return kubernetes.KubernetesApiServerACLs{}
	}

	var allowedCIDRs []string
	if v, exists := acl["allowed_cidrs"]; exists && v != nil {
		allowedCIDRs = convert.ConvertToStringSlice(v)
	}

	return kubernetes.KubernetesApiServerACLs{
		AllowedCIDRs: allowedCIDRs,
	}
}

// convertAutoscalerConfig converts the autoscaler config from Terraform schema to the API format
func convertAutoscalerConfig(config interface{}) *kubernetes.AutoscalerConfig {
	if config == nil {
		return nil
	}

	configList, ok := config.([]interface{})
	if !ok || len(configList) == 0 {
		return nil
	}

	first := configList[0]
	if first == nil {
		return nil
	}

	cfg, ok := first.(map[string]interface{})
	if !ok || cfg == nil {
		return nil
	}

	result := &kubernetes.AutoscalerConfig{}

	if v, exists := cfg["scale_down_disabled"]; exists && v != nil {
		result.ScaleDownDisabled = v.(bool)
	}
	if v, exists := cfg["scale_down_delay_after_add"]; exists && v != nil {
		result.ScaleDownDelayAfterAdd = v.(string)
	}
	if v, exists := cfg["estimator"]; exists && v != nil {
		result.Estimator = v.(string)
	}
	if v, exists := cfg["expander"]; exists && v != nil {
		result.Expander = v.(string)
	}
	if v, exists := cfg["ignore_daemonsets_utilization"]; exists && v != nil {
		result.IgnoreDaemonsetsUtilization = v.(bool)
	}
	if v, exists := cfg["balance_similar_node_groups"]; exists && v != nil {
		result.BalanceSimilarNodeGroups = v.(bool)
	}
	if v, exists := cfg["expendable_pods_priority_cutoff"]; exists && v != nil {
		result.ExpendablePodsPriorityCutoff = v.(int)
	}
	if v, exists := cfg["scale_down_unneeded_time"]; exists && v != nil {
		result.ScaleDownUnneededTime = v.(string)
	}
	if v, exists := cfg["scale_down_utilization_threshold"]; exists && v != nil {
		result.ScaleDownUtilizationThreshold = v.(float64)
	}
	if v, exists := cfg["max_graceful_termination_sec"]; exists && v != nil {
		result.MaxGracefulTerminationSec = v.(int)
	}
	if v, exists := cfg["enable_proactive_scale_up"]; exists && v != nil {
		result.EnableProactiveScaleUp = v.(bool)
	}

	return result
}
