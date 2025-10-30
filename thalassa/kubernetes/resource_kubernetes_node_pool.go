package kubernetes

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	kubernetes "github.com/thalassa-cloud/client-go/kubernetes"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func resourceKubernetesNodePool() *schema.Resource {
	return &schema.Resource{
		Description:   "Create an Kubernetes Node Pool for a Kubernetes Cluster. This resource is only available for managed Kubernetes Clusters. A Node Pool is a group of nodes that are identically configured and are automatically joined to the Kubernetes Cluster. Node Pools can be scaled up and down as needed.",
		CreateContext: resourceKubernetesNodePoolCreate,
		ReadContext:   resourceKubernetesNodePoolRead,
		UpdateContext: resourceKubernetesNodePoolUpdate,
		DeleteContext: resourceKubernetesNodePoolDelete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
			// If autoscaling is enabled, replicas must be unset by the user
			if d.Get("enable_autoscaling").(bool) {
				if _, ok := d.GetOk("replicas"); ok {
					return fmt.Errorf("replicas must be unset when enable_autoscaling is true")
				}
				if _, ok := d.GetOk("min_replicas"); !ok {
					return fmt.Errorf("min_replicas must be set when enable_autoscaling is true")
				}
				if _, ok := d.GetOk("max_replicas"); !ok {
					return fmt.Errorf("max_replicas must be set when enable_autoscaling is true")
				}
			} else {
				// replicas must be set
				if _, ok := d.GetOk("replicas"); !ok {
					return fmt.Errorf("replicas must be set when enable_autoscaling is false")
				}
			}
			return nil
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the Kubernetes Node Pool. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StringLenBetween(1, 62),
				Description:  "Name of the Kubernetes Node Pool",
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Labels for the Kubernetes Node Pool. Optional. These labels are used for filtering and grouping resources in the Thalassa Console. Labels are not applied to the Kubernetes nodes created for this Node Pool, please use node_labels instead.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Annotations for the Kubernetes Node Pool. Optional. These annotations are used for additional metadata and configuration. Annotations are not applied to the Kubernetes nodes created for this Node Pool, please use node_annotations instead.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Subnet ID where the Kubernetes node pool nodes will be deployed. This subnet must be in the same VPC as the Kubernetes cluster.",
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Kubernetes Cluster of the Kubernetes Node Pool",
			},
			"slug": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Slug of the Kubernetes Node Pool",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the Kubernetes Node Pool",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.StringLenBetween(0, 255),
				Description:  "A human readable description about the Kubernetes Node Pool",
			},
			"kubernetes_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Kubernetes version for the node pool nodes. Optional - if not specified, the cluster's version will be used. Can be specified as version name, slug, or identity. Must be an enabled version.",
			},
			"upgrade_strategy": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(kubernetes.KubernetesNodePoolUpgradeStrategyAlways),
				ValidateFunc: validate.StringInSlice([]string{
					string(kubernetes.KubernetesNodePoolUpgradeStrategyManual),
					string(kubernetes.KubernetesNodePoolUpgradeStrategyAuto),
					// Legacy options. Provided for backward compatibility.
					string(kubernetes.KubernetesNodePoolUpgradeStrategyAlways),
					string(kubernetes.KubernetesNodePoolUpgradeStrategyOnDelete),
					string(kubernetes.KubernetesNodePoolUpgradeStrategyInplace),
					string(kubernetes.KubernetesNodePoolUpgradeStrategyNever),
				}, false),
				Description: "Upgrade strategy for the Kubernetes Node Pool",
			},
			"enable_autoscaling": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable autoscaling for the Kubernetes Node Pool",
			},
			"enable_autohealing": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable autohealing for the Kubernetes Node Pool",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Availability zone for the Kubernetes Node Pool",
			},
			"replicas": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Number of replicas for the Kubernetes Node Pool. Do not set this when enable_autoscaling is true.",
			},
			"min_replicas": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Minimum number of replicas for the Kubernetes Node Pool. May only be set when enable_autoscaling is true.",
			},
			"max_replicas": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Maximum number of replicas for the Kubernetes Node Pool. May only be set when enable_autoscaling is true.",
			},
			"machine_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Machine type for the Kubernetes Node Pool",
			},
			"node_taints": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Taints for the Kubernetes Node Pool",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"effect": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Effect of the taint",
							ValidateFunc: validate.StringInSlice([]string{
								"NoSchedule",
								"NoExecute",
								"PreferNoSchedule",
							}, false),
						},
						"key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Key of the taint",
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								if _, ok := v.(string); !ok {
									errors = append(errors, fmt.Errorf("expected key to be a string"))
								}
								// may not contain whitespace
								if strings.Contains(v.(string), " ") || strings.Contains(v.(string), ".") {
									errors = append(errors, fmt.Errorf("key may not contain whitespace or dots"))
								}
								return
							},
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Value of the taint. Optional.",
						},
						"operator": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Operator of the taint",
							ValidateFunc: validate.StringInSlice([]string{
								"Equal",
								"Exists",
							}, false),
						},
					},
				},
			},
			"node_labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Labels for the Kubernetes Nodes within this Node Pool. Optional. These labels are applied to the Kubernetes nodes created for this Node Pool. Labels must match the same constraints as Kubernetes labels.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"node_annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Annotations for the Kubernetes Nodes within this Node Pool. Optional. These annotations are applied to the Kubernetes nodes created for this Node Pool. Annotations must match the same constraints as Kubernetes annotations.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"security_group_attachments": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List identities of security group that will be attached to the machines in the Node Pool",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceKubernetesNodePoolCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	var subnetIdentity *string
	if subnetId, ok := d.GetOk("subnet_id"); ok {
		subnetIdentity = convert.Ptr(subnetId.(string))
	}

	var kubernetesVersionIdentity *string
	if kubernetesVersion, ok := d.GetOk("kubernetes_version"); ok {
		// Get version from API
		kubernetesVersions, err := client.Kubernetes().ListKubernetesVersions(ctx)
		if err != nil {
			return diag.FromErr(err)
		}
		for _, kv := range kubernetesVersions {
			if !kv.Enabled { // skip disabled versions
				continue
			}

			if kv.Identity == kubernetesVersion || kv.Slug == kubernetesVersion || kv.Name == kubernetesVersion {
				kubernetesVersionIdentity = convert.Ptr(kv.Identity)
				break
			}
		}
		if kubernetesVersionIdentity == nil {
			return diag.FromErr(fmt.Errorf("kubernetes version '%s' not found or not enabled. Please check available versions and ensure the version is enabled", kubernetesVersion))
		}
	} else {
		// fetch the cluster's version
		kubernetesClusterIdentity := d.Get("cluster_id").(string)
		kubernetesCluster, err := client.Kubernetes().GetKubernetesCluster(ctx, kubernetesClusterIdentity)
		if err != nil {
			return diag.FromErr(err)
		}
		kubernetesVersionIdentity = convert.Ptr(kubernetesCluster.ClusterVersion.Identity)
		if kubernetesVersionIdentity == nil {
			return diag.FromErr(fmt.Errorf("kubernetes version not found for cluster '%s'", kubernetesClusterIdentity))
		}
	}
	// If kubernetes_version is not provided, kubernetesVersionIdentity will be nil
	// and the cluster's version will be used automatically

	enableAutoscaling := d.Get("enable_autoscaling").(bool)
	var replicas int

	if enableAutoscaling {
		minReplicas := d.Get("min_replicas").(int)
		maxReplicas := d.Get("max_replicas").(int)
		if minReplicas > maxReplicas {
			return diag.FromErr(fmt.Errorf("autoscaling configuration error: min_replicas (%d) cannot be greater than max_replicas (%d). Please ensure min_replicas <= max_replicas", minReplicas, maxReplicas))
		}
		if minReplicas < 0 {
			return diag.FromErr(fmt.Errorf("autoscaling configuration error: min_replicas must be at least 0, got %d", minReplicas))
		}
		// When autoscaling is enabled, start with min_replicas
		replicas = minReplicas
	} else {
		// When autoscaling is disabled, replicas is required
		if replicasVal, ok := d.GetOk("replicas"); ok {
			replicas = replicasVal.(int)
		} else {
			return diag.FromErr(fmt.Errorf("replicas is required when enable_autoscaling is false. Set replicas to the desired number of nodes for this node pool"))
		}
	}

	createKubernetesNodePool := kubernetes.CreateKubernetesNodePool{
		Name:              d.Get("name").(string),
		MachineType:       d.Get("machine_type").(string), // TODO: check if machine type is valid
		Replicas:          replicas,
		Description:       d.Get("description").(string),
		Labels:            convert.ConvertToMap(d.Get("labels")),
		Annotations:       convert.ConvertToMap(d.Get("annotations")),
		EnableAutoscaling: enableAutoscaling,
		AvailabilityZone:  d.Get("availability_zone").(string),
		MinReplicas:       d.Get("min_replicas").(int),
		MaxReplicas:       d.Get("max_replicas").(int),
		NodeSettings: kubernetes.KubernetesNodeSettings{
			Annotations: convert.ConvertToMap(d.Get("node_annotations")),
			Labels:      convert.ConvertToMap(d.Get("node_labels")),
			Taints:      convertToNodeTaints(d.Get("node_taints").([]interface{})),
		},
		EnableAutoHealing:         d.Get("enable_autohealing").(bool),
		UpgradeStrategy:           convert.Ptr(kubernetes.KubernetesNodePoolUpgradeStrategy(d.Get("upgrade_strategy").(string))),
		SubnetIdentity:            subnetIdentity,
		KubernetesVersionIdentity: kubernetesVersionIdentity,
	}

	if securityGroupAttachments, ok := d.GetOk("security_group_attachments"); ok {
		createKubernetesNodePool.SecurityGroupAttachments = convert.ConvertToStringSlice(securityGroupAttachments)
	}

	kubernetesClusterIdentity := d.Get("cluster_id").(string)
	kubernetesNodePool, err := client.Kubernetes().CreateKubernetesNodePool(ctx, kubernetesClusterIdentity, createKubernetesNodePool)
	if err != nil {
		return diag.FromErr(err)
	}
	if kubernetesNodePool != nil {
		d.SetId(kubernetesNodePool.Identity)
		d.Set("slug", kubernetesNodePool.Slug)
		d.Set("status", kubernetesNodePool.Status)

		for {
			kubernetesNodePool, err := client.Kubernetes().GetKubernetesNodePool(ctx, kubernetesClusterIdentity, kubernetesNodePool.Identity)
			if err != nil {
				return diag.FromErr(err)
			}
			if kubernetesNodePool.Status == kubernetes.KubernetesNodePoolStatusReady {
				break
			}
			time.Sleep(1 * time.Second)
		}
		return nil
	}
	return resourceKubernetesNodePoolRead(ctx, d, m)
}

func resourceKubernetesNodePoolRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	identity := d.Get("id").(string)
	kubernetesClusterIdentity := d.Get("cluster_id").(string)
	kubernetesNodePool, err := client.Kubernetes().GetKubernetesNodePool(ctx, kubernetesClusterIdentity, identity)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error getting kubernetesNodePool: %s", err))
	}

	if kubernetesNodePool == nil {
		d.SetId("")
		return nil
	}

	currentlyConfiguredVersion, ok := d.GetOk("kubernetes_version")
	if ok {
		if kubernetesNodePool.KubernetesVersion != nil {
			if !(kubernetesNodePool.KubernetesVersion.Name == currentlyConfiguredVersion || kubernetesNodePool.KubernetesVersion.Slug == currentlyConfiguredVersion || kubernetesNodePool.KubernetesVersion.Identity == currentlyConfiguredVersion) {
				d.Set("kubernetes_version", kubernetesNodePool.KubernetesVersion.Slug)
			}
		}
	}

	d.SetId(kubernetesNodePool.Identity)
	d.Set("name", kubernetesNodePool.Name)
	d.Set("slug", kubernetesNodePool.Slug)
	d.Set("description", kubernetesNodePool.Description)
	d.Set("labels", convertFromNodeLabels(kubernetesNodePool.Labels))
	d.Set("annotations", convertFromNodeLabels(kubernetesNodePool.Annotations))
	d.Set("status", kubernetesNodePool.Status)

	// if replicas is set, set it in state
	if _, ok := d.GetOk("replicas"); ok {
		d.Set("replicas", kubernetesNodePool.Replicas)
	}
	d.Set("availability_zone", kubernetesNodePool.AvailabilityZone)
	if _, ok := d.GetOk("min_replicas"); ok {
		d.Set("min_replicas", kubernetesNodePool.MinReplicas)
	}
	if _, ok := d.GetOk("max_replicas"); ok {
		d.Set("max_replicas", kubernetesNodePool.MaxReplicas)
	}
	d.Set("machine_type", kubernetesNodePool.MachineType)
	d.Set("enable_autoscaling", kubernetesNodePool.EnableAutoscaling)
	d.Set("enable_autohealing", kubernetesNodePool.EnableAutoHealing)
	d.Set("node_taints", convertFromNodeTaints(kubernetesNodePool.NodeSettings.Taints))
	d.Set("node_labels", convertFromNodeLabels(kubernetesNodePool.NodeSettings.Labels))
	d.Set("node_annotations", convertFromNodeLabels(kubernetesNodePool.NodeSettings.Annotations))

	if kubernetesNodePool.Subnet != nil {
		d.Set("subnet_id", kubernetesNodePool.Subnet.Identity)
	}

	securityGroupAttachments := []string{}
	for _, sg := range kubernetesNodePool.SecurityGroups {
		securityGroupAttachments = append(securityGroupAttachments, sg.Identity)
	}
	d.Set("security_group_attachments", securityGroupAttachments)

	return nil
}

func resourceKubernetesNodePoolUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}
	nodePoolIdentity := d.Get("id").(string)
	kubernetesClusterIdentity := d.Get("cluster_id").(string)
	if kubernetesClusterIdentity == "" {
		return diag.FromErr(fmt.Errorf("kubernetes cluster identity is required"))
	}

	var kubernetesVersionIdentity *string
	if kubernetesVersion, ok := d.GetOk("kubernetes_version"); ok {
		// Get version from API
		kubernetesVersions, err := client.Kubernetes().ListKubernetesVersions(ctx)
		if err != nil {
			return diag.FromErr(err)
		}
		for _, kv := range kubernetesVersions {
			if !kv.Enabled { // skip disabled versions
				continue
			}

			if kv.Identity == kubernetesVersion || kv.Slug == kubernetesVersion || kv.Name == kubernetesVersion {
				kubernetesVersionIdentity = convert.Ptr(kv.Identity)
				break
			}
		}
		if kubernetesVersionIdentity == nil {
			return diag.FromErr(fmt.Errorf("kubernetes version '%s' not found or not enabled. Please check available versions and ensure the version is enabled", kubernetesVersion))
		}
	} else {
		// fetch the cluster's version
		kubernetesCluster, err := client.Kubernetes().GetKubernetesCluster(ctx, kubernetesClusterIdentity)
		if err != nil {
			return diag.FromErr(err)
		}
		kubernetesVersionIdentity = convert.Ptr(kubernetesCluster.ClusterVersion.Identity)
		if kubernetesVersionIdentity == nil {
			return diag.FromErr(fmt.Errorf("kubernetes version not found for cluster '%s'", kubernetesClusterIdentity))
		}
	}
	// If kubernetes_version is not provided, kubernetesVersionIdentity will be nil
	// and the cluster's version will be used automatically

	enableAutoscaling := d.Get("enable_autoscaling").(bool)
	var replicas *int

	if enableAutoscaling {
		minReplicas := d.Get("min_replicas").(int)
		maxReplicas := d.Get("max_replicas").(int)
		if minReplicas > maxReplicas {
			return diag.FromErr(fmt.Errorf("autoscaling configuration error: min_replicas (%d) cannot be greater than max_replicas (%d). Please ensure min_replicas <= max_replicas", minReplicas, maxReplicas))
		}
		if minReplicas < 0 {
			return diag.FromErr(fmt.Errorf("autoscaling configuration error: min_replicas must be at least 1, got %d", minReplicas))
		}

		currentNodePool, err := client.Kubernetes().GetKubernetesNodePool(ctx, kubernetesClusterIdentity, nodePoolIdentity)
		if err != nil {
			return diag.FromErr(err)
		}

		// When autoscaling is enabled, set replicas to the current value, with min_replicas as the minimum and max_replicas as the maximum
		replicas = &currentNodePool.Replicas
		if *replicas < minReplicas {
			*replicas = minReplicas
		}
		if *replicas > maxReplicas {
			*replicas = maxReplicas
		}
	} else {
		// When autoscaling is disabled, replicas is required
		if replicasVal, ok := d.GetOk("replicas"); ok {
			replicasInt := replicasVal.(int)
			replicas = &replicasInt
		} else {
			return diag.FromErr(fmt.Errorf("replicas is required when enable_autoscaling is false. Set replicas to the desired number of nodes for this node pool"))
		}
	}

	updateKubernetesNodePool := kubernetes.UpdateKubernetesNodePool{
		Description:               d.Get("description").(string),
		Labels:                    convert.ConvertToMap(d.Get("labels")),
		Annotations:               convert.ConvertToMap(d.Get("annotations")),
		MachineType:               d.Get("machine_type").(string),
		Replicas:                  replicas,
		AvailabilityZone:          d.Get("availability_zone").(string),
		EnableAutoscaling:         convert.Ptr(enableAutoscaling),
		MinReplicas:               convert.Ptr(d.Get("min_replicas").(int)),
		MaxReplicas:               convert.Ptr(d.Get("max_replicas").(int)),
		EnableAutoHealing:         convert.Ptr(d.Get("enable_autohealing").(bool)),
		UpgradeStrategy:           convert.Ptr(kubernetes.KubernetesNodePoolUpgradeStrategy(d.Get("upgrade_strategy").(string))),
		KubernetesVersionIdentity: kubernetesVersionIdentity,
		NodeSettings: &kubernetes.KubernetesNodeSettings{
			Annotations: convert.ConvertToMap(d.Get("node_annotations")),
			Labels:      convert.ConvertToMap(d.Get("node_labels")),
			Taints:      convertToNodeTaints(d.Get("node_taints").([]interface{})),
		},
	}
	if securityGroupAttachments, ok := d.GetOk("security_group_attachments"); ok {
		updateKubernetesNodePool.SecurityGroupAttachments = convert.ConvertToStringSlice(securityGroupAttachments)
	}

	kubernetesNodePool, err := client.Kubernetes().UpdateKubernetesNodePool(ctx, kubernetesClusterIdentity, nodePoolIdentity, updateKubernetesNodePool)
	if err != nil {
		return diag.FromErr(err)
	}

	if kubernetesNodePool != nil {
		d.Set("slug", kubernetesNodePool.Slug)
		d.Set("status", kubernetesNodePool.Status)

		for {
			kubernetesNodePool, err := client.Kubernetes().GetKubernetesNodePool(ctx, kubernetesClusterIdentity, nodePoolIdentity)
			if err != nil {
				return diag.FromErr(err)
			}
			if kubernetesNodePool.Status == kubernetes.KubernetesNodePoolStatusReady {
				break
			}
			time.Sleep(1 * time.Second)
		}
		if _, ok := d.GetOk("replicas"); ok {
			d.Set("replicas", kubernetesNodePool.Replicas)
		}

		if _, ok := d.GetOk("min_replicas"); ok {
			d.Set("min_replicas", kubernetesNodePool.MinReplicas)
		}
		if _, ok := d.GetOk("max_replicas"); ok {
			d.Set("max_replicas", kubernetesNodePool.MaxReplicas)
		}
		d.Set("machine_type", kubernetesNodePool.MachineType)
		d.Set("labels", kubernetesNodePool.Labels)
		d.Set("annotations", kubernetesNodePool.Annotations)
		d.Set("enable_autoscaling", kubernetesNodePool.EnableAutoscaling)
		d.Set("enable_autohealing", kubernetesNodePool.EnableAutoHealing)
		d.Set("node_taints", convertFromNodeTaints(kubernetesNodePool.NodeSettings.Taints))
		d.Set("node_labels", convertFromNodeLabels(kubernetesNodePool.NodeSettings.Labels))
		d.Set("node_annotations", convertFromNodeLabels(kubernetesNodePool.NodeSettings.Annotations))
	}

	return resourceKubernetesNodePoolRead(ctx, d, m)
}

func resourceKubernetesNodePoolDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	nodePoolIdentity := d.Get("id").(string)
	kubernetesClusterIdentity := d.Get("cluster_id").(string)
	err = client.Kubernetes().DeleteKubernetesNodePool(ctx, kubernetesClusterIdentity, nodePoolIdentity)
	if err != nil {
		if !tcclient.IsNotFound(err) {
			return diag.FromErr(err)
		}
	}

	for {
		kubernetesNodePool, err := client.Kubernetes().GetKubernetesNodePool(ctx, kubernetesClusterIdentity, nodePoolIdentity)
		if err != nil {
			return diag.FromErr(err)
		}
		if kubernetesNodePool == nil {
			break
		}
		if kubernetesNodePool.Status == kubernetes.KubernetesNodePoolStatusDeleted {
			break
		}
		time.Sleep(1 * time.Second)
	}

	d.SetId("")
	return nil
}

func convertToNodeTaints(taints []interface{}) []kubernetes.NodeTaint {
	nodeTaints := make([]kubernetes.NodeTaint, len(taints))
	for i, taint := range taints {
		taintMap := taint.(map[string]interface{})
		nodeTaints[i] = kubernetes.NodeTaint{
			Key:      taintMap["key"].(string),
			Value:    taintMap["value"].(string),
			Operator: taintMap["operator"].(string),
			Effect:   taintMap["effect"].(string),
		}
	}
	return nodeTaints
}

func convertFromNodeTaints(taints []kubernetes.NodeTaint) []interface{} {
	nodeTaints := make([]interface{}, len(taints))
	for i, taint := range taints {
		nodeTaints[i] = map[string]interface{}{
			"key":      taint.Key,
			"value":    taint.Value,
			"operator": taint.Operator,
			"effect":   taint.Effect,
		}
	}
	return nodeTaints
}

func convertToNodeLabels(labels map[string]interface{}) map[string]string {
	nodeLabels := make(map[string]string)
	for key, value := range labels {
		nodeLabels[key] = value.(string)
	}
	return nodeLabels
}

func convertFromNodeLabels(labels map[string]string) map[string]interface{} {
	nodeLabels := make(map[string]interface{})
	for key, value := range labels {
		nodeLabels[key] = value
	}
	return nodeLabels
}
