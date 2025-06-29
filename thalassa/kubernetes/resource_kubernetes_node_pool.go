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
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
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
				Optional:    true,
				ForceNew:    true,
				Description: "Subnet of the Kubernetes Cluster. Required for managed Kubernetes Clusters.",
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
				Description: "Kubernetes version for the Kubernetes Node Pool. Optional. Will use the Kubernetes Cluster version if not set.",
			},
			"upgrade_strategy": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(kubernetes.KubernetesNodePoolUpgradeStrategyAlways),
				ValidateFunc: validate.StringInSlice([]string{
					string(kubernetes.KubernetesNodePoolUpgradeStrategyAlways),
					string(kubernetes.KubernetesNodePoolUpgradeStrategyOnDelete),
					string(kubernetes.KubernetesNodePoolUpgradeStrategyInplace),
					string(kubernetes.KubernetesNodePoolUpgradeStrategyNever),
				}, false),
				Description: "Upgrade strategy for the Kubernetes Node Pool",
			},
			// TODO: add these back in when the API is updated
			// "enable_autoscaling": {
			// 	Type:        schema.TypeBool,
			// 	Optional:    true,
			// 	Default:     false,
			// 	Description: "Enable autoscaling for the Kubernetes Node Pool",
			// },
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
				Default:     1,
				Description: "Number of replicas for the Kubernetes Node Pool",
			},
			"min_replicas": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Minimum number of replicas for the Kubernetes Node Pool",
			},
			"max_replicas": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Maximum number of replicas for the Kubernetes Node Pool",
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
			return diag.FromErr(fmt.Errorf("kubernetes version not found"))
		}
	}

	createKubernetesNodePool := kubernetes.CreateKubernetesNodePool{
		Name:        d.Get("name").(string),
		MachineType: d.Get("machine_type").(string), // TODO: check if machine type is valid
		Replicas:    d.Get("replicas").(int),
		Description: d.Get("description").(string),
		Labels:      convert.ConvertToMap(d.Get("labels")),
		Annotations: convert.ConvertToMap(d.Get("annotations")),
		// EnableAutoscaling: d.Get("enable_autoscaling").(bool),
		AvailabilityZone: d.Get("availability_zone").(string),
		MinReplicas:      d.Get("min_replicas").(int),
		MaxReplicas:      d.Get("max_replicas").(int),
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
	if kubernetesNodePool.KubernetesVersion != nil {
		if !ok || !(kubernetesNodePool.KubernetesVersion.Name == currentlyConfiguredVersion || kubernetesNodePool.KubernetesVersion.Slug == currentlyConfiguredVersion || kubernetesNodePool.KubernetesVersion.Identity == currentlyConfiguredVersion) {
			d.Set("kubernetes_version", kubernetesNodePool.KubernetesVersion.Slug)
		}
	}

	d.SetId(kubernetesNodePool.Identity)
	d.Set("name", kubernetesNodePool.Name)
	d.Set("slug", kubernetesNodePool.Slug)
	d.Set("description", kubernetesNodePool.Description)
	d.Set("labels", convertFromNodeLabels(kubernetesNodePool.Labels))
	d.Set("annotations", convertFromNodeLabels(kubernetesNodePool.Annotations))
	d.Set("status", kubernetesNodePool.Status)
	d.Set("replicas", kubernetesNodePool.Replicas)
	d.Set("availability_zone", kubernetesNodePool.AvailabilityZone)
	d.Set("min_replicas", kubernetesNodePool.MinReplicas)
	d.Set("max_replicas", kubernetesNodePool.MaxReplicas)
	d.Set("machine_type", kubernetesNodePool.MachineType)
	// d.Set("enable_autoscaling", kubernetesNodePool.EnableAutoscaling)
	d.Set("enable_autohealing", kubernetesNodePool.EnableAutoHealing)
	d.Set("node_taints", convertFromNodeTaints(kubernetesNodePool.NodeSettings.Taints))
	d.Set("node_labels", convertFromNodeLabels(kubernetesNodePool.NodeSettings.Labels))
	d.Set("node_annotations", convertFromNodeLabels(kubernetesNodePool.NodeSettings.Annotations))

	if kubernetesNodePool.Subnet != nil {
		d.Set("subnet_id", kubernetesNodePool.Subnet.Identity)
	}

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
			return diag.FromErr(fmt.Errorf("kubernetes version not found"))
		}
	}

	updateKubernetesNodePool := kubernetes.UpdateKubernetesNodePool{
		Description:      d.Get("description").(string),
		Labels:           convert.ConvertToMap(d.Get("labels")),
		Annotations:      convert.ConvertToMap(d.Get("annotations")),
		MachineType:      d.Get("machine_type").(string),
		Replicas:         convert.Ptr(d.Get("replicas").(int)),
		AvailabilityZone: d.Get("availability_zone").(string),
		// EnableAutoscaling:         convert.Ptr(d.Get("enable_autoscaling").(bool)),
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
		d.Set("replicas", kubernetesNodePool.Replicas)
		d.Set("min_replicas", kubernetesNodePool.MinReplicas)
		d.Set("max_replicas", kubernetesNodePool.MaxReplicas)
		d.Set("machine_type", kubernetesNodePool.MachineType)
		// d.Set("enable_autoscaling", kubernetesNodePool.EnableAutoscaling)
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
			Key:    taintMap["key"].(string),
			Value:  taintMap["value"].(string),
			Effect: taintMap["effect"].(string),
		}
	}
	return nodeTaints
}

func convertFromNodeTaints(taints []kubernetes.NodeTaint) []interface{} {
	nodeTaints := make([]interface{}, len(taints))
	for i, taint := range taints {
		nodeTaints[i] = map[string]interface{}{
			"key":    taint.Key,
			"value":  taint.Value,
			"effect": taint.Effect,
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
