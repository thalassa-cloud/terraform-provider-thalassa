package iaas

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"

	iaas "github.com/thalassa-cloud/client-go/iaas"
)

var (
	targetGroupProtocols = []string{
		"tcp", "udp", "http", "https", "grpc", "quic",
	}
	healthCheckProtocols = []string{
		"tcp", "udp", "http", "https",
	}
	loadbalancingPolicies = []string{"ROUND_ROBIN", "RANDOM", "MAGLEV"}
)

func validateOptionalStringInSlice(valid []string) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, es []error) {
		s, ok := v.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected string at %s", k))
			return
		}
		if s == "" {
			return
		}
		return validate.StringInSlice(valid, false)(v, k)
	}
}

func resourceTargetGroup() *schema.Resource {
	return &schema.Resource{
		Description:   "Create a target group for a load balancer",
		CreateContext: resourceTargetGroupCreate,
		ReadContext:   resourceTargetGroupRead,
		UpdateContext: resourceTargetGroupUpdate,
		DeleteContext: resourceTargetGroupDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the Target Group. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StringLenBetween(1, 62),
				Description:  "Name of the Target Group",
			},
			"slug": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:         schema.TypeString,
				Default:      "",
				Optional:     true,
				ValidateFunc: validate.StringLenBetween(0, 255),
				Description:  "A human readable description about the target group",
			},
			"labels": {
				Type:        schema.TypeMap,
				Default:     make(map[string]string),
				Optional:    true,
				Description: "Labels for the Target Group",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Default:     make(map[string]string),
				Optional:    true,
				Description: "Annotations for the Target Group",
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StringInSlice(targetGroupProtocols, false),
				Description:  "Protocol for routing traffic to targets (tcp, udp, http, https, grpc, quic).",
			},
			"port": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IntBetween(1, 65535),
				Description:  "The port on which the targets receive traffic",
			},
			"enable_proxy_protocol": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "When true, the load balancer uses PROXY protocol toward backends in this target group. All targets must support PROXY protocol.",
			},
			"loadbalancing_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateOptionalStringInSlice(loadbalancingPolicies),
				Description:  "Load balancing algorithm: ROUND_ROBIN (default), RANDOM, or MAGLEV.",
			},
			"target_selector": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Label selector for automatic target membership; when set, targets matching these labels join the group.",
			},
			"health_check_path": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.StringLenBetween(0, 255),
				Description:  "HTTP(S) health check path; leave empty for TCP/UDP checks or when not using HTTP health checks.",
			},
			"health_check_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validate.IntBetween(1, 65535),
				Description:  "Port for health checks; if omitted but other health check settings are set, defaults to the target group port.",
			},
			"health_check_protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateOptionalStringInSlice(healthCheckProtocols),
				Description:  "Health check protocol (tcp, udp, http, https). If omitted when configuring a health check, defaults to tcp.",
			},
			"health_check_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      30,
				ValidateFunc: validate.IntBetween(5, 300),
				Description:  "Seconds between health checks of each target (periodSeconds).",
			},
			"health_check_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      5,
				ValidateFunc: validate.IntBetween(1, 300),
				Description:  "Seconds to wait for a health check response before failure (timeoutSeconds).",
			},
			"healthy_threshold": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3,
				ValidateFunc: validate.IntBetween(1, 10),
				Description:  "Consecutive successes required to mark a target healthy.",
			},
			"unhealthy_threshold": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3,
				ValidateFunc: validate.IntBetween(1, 10),
				Description:  "Consecutive failures required to mark a target unhealthy.",
			},
			"attachments": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "The targets to attach to the target group. If provided, the targets will be attached to the target group when the resource is created. Overwrites the target group attachment resource.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ID of the target (e.g. instance ID)",
						},
					},
				},
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The VPC this target group belongs to",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTargetGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating client: %w", err))
	}

	enableProxy := d.Get("enable_proxy_protocol").(bool)
	createTargetGroup := iaas.CreateTargetGroup{
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		Labels:              convert.ConvertToMap(d.Get("labels")),
		Annotations:         convert.ConvertToMap(d.Get("annotations")),
		Vpc:                 d.Get("vpc_id").(string),
		TargetPort:          d.Get("port").(int),
		Protocol:            iaas.LoadbalancerProtocol(d.Get("protocol").(string)),
		EnableProxyProtocol: &enableProxy,
	}
	if ts := convert.ConvertToMap(d.Get("target_selector")); len(ts) > 0 {
		createTargetGroup.TargetSelector = ts
	}
	if v := strings.TrimSpace(d.Get("loadbalancing_policy").(string)); v != "" {
		p := iaas.LoadbalancingPolicy(v)
		createTargetGroup.LoadbalancingPolicy = &p
	}
	if hc := expandTargetGroupHealthCheck(d); hc != nil {
		createTargetGroup.HealthCheck = hc
	}

	tg, err := client.IaaS().CreateTargetGroup(ctx, createTargetGroup)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating target group: %w", err))
	}
	if tg != nil {
		d.SetId(tg.Identity)
		d.Set("slug", tg.Slug)
		// Attach targets if specified
		if attachments, ok := d.GetOk("attachments"); ok {
			attachmentList := attachments.([]interface{})
			attach := make([]iaas.AttachTarget, len(attachmentList))
			for i, a := range attachmentList {
				row := a.(map[string]interface{})
				attach[i] = iaas.AttachTarget{ServerIdentity: row["id"].(string)}
			}
			batch := iaas.TargetGroupAttachmentsBatch{
				TargetGroupID: tg.Identity,
				Attachments:   attach,
			}
			if err := client.IaaS().SetTargetGroupServerAttachments(ctx, batch); err != nil {
				return diag.FromErr(fmt.Errorf("error setting target group server attachments: %w", err))
			}
		}
		return nil
	}
	return resourceTargetGroupRead(ctx, d, m)
}

func resourceTargetGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating client: %w", err))
	}

	id := d.Get("id").(string)
	tg, err := client.IaaS().GetTargetGroup(ctx, iaas.GetTargetGroupRequest{Identity: id})
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error getting target group: %w", err))
	}
	if tg == nil {
		d.SetId("")
		return nil
	}

	d.SetId(tg.Identity)
	d.Set("name", tg.Name)
	d.Set("slug", tg.Slug)
	d.Set("description", tg.Description)
	d.Set("labels", tg.Labels)
	d.Set("annotations", tg.Annotations)
	d.Set("vpc_id", tg.Vpc.Identity)
	d.Set("protocol", string(tg.Protocol))
	d.Set("port", tg.TargetPort)

	if tg.EnableProxyProtocol != nil {
		d.Set("enable_proxy_protocol", *tg.EnableProxyProtocol)
	} else {
		d.Set("enable_proxy_protocol", false)
	}
	if tg.LoadbalancingPolicy != nil {
		d.Set("loadbalancing_policy", string(*tg.LoadbalancingPolicy))
	}
	if tg.TargetSelector != nil {
		d.Set("target_selector", tg.TargetSelector)
	} else {
		d.Set("target_selector", map[string]string{})
	}

	if tg.HealthCheck != nil {
		d.Set("health_check_path", tg.HealthCheck.Path)
		d.Set("health_check_port", tg.HealthCheck.Port)
		d.Set("health_check_protocol", string(tg.HealthCheck.Protocol))
		d.Set("health_check_interval", tg.HealthCheck.PeriodSeconds)
		d.Set("health_check_timeout", tg.HealthCheck.TimeoutSeconds)
		d.Set("healthy_threshold", tg.HealthCheck.HealthyThreshold)
		d.Set("unhealthy_threshold", tg.HealthCheck.UnhealthyThreshold)
	}

	// Set targets from attachments
	if tg.LoadbalancerTargetGroupAttachments != nil {
		targets := make([]map[string]interface{}, len(tg.LoadbalancerTargetGroupAttachments))
		for i, att := range tg.LoadbalancerTargetGroupAttachments {
			if att.VirtualMachineInstance != nil {
				targets[i] = map[string]interface{}{
					"id": att.VirtualMachineInstance.Identity,
				}
			}
		}
		d.Set("attachments", targets)
	}

	return nil
}

func resourceTargetGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	enableProxy := d.Get("enable_proxy_protocol").(bool)
	updateTargetGroup := iaas.UpdateTargetGroup{
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		Labels:              convert.ConvertToMap(d.Get("labels")),
		Annotations:         convert.ConvertToMap(d.Get("annotations")),
		TargetPort:          d.Get("port").(int),
		Protocol:            iaas.LoadbalancerProtocol(d.Get("protocol").(string)),
		EnableProxyProtocol: &enableProxy,
	}
	if d.HasChange("target_selector") {
		updateTargetGroup.TargetSelector = convert.ConvertToMap(d.Get("target_selector"))
	}
	if d.HasChange("loadbalancing_policy") {
		if v := strings.TrimSpace(d.Get("loadbalancing_policy").(string)); v != "" {
			p := iaas.LoadbalancingPolicy(v)
			updateTargetGroup.LoadbalancingPolicy = &p
		}
	}
	if hc := expandTargetGroupHealthCheck(d); hc != nil {
		updateTargetGroup.HealthCheck = hc
	}

	id := d.Get("id").(string)
	tg, err := client.IaaS().UpdateTargetGroup(ctx, iaas.UpdateTargetGroupRequest{
		Identity:          id,
		UpdateTargetGroup: updateTargetGroup,
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating target group: %w", err))
	}
	if tg != nil {
		// Attach targets if specified
		if attachments, ok := d.GetOk("attachments"); ok {
			attachmentList := attachments.([]interface{})
			attachments := make([]iaas.AttachTarget, len(attachmentList))
			for i, a := range attachmentList {
				attachment := a.(map[string]interface{})
				attachments[i] = iaas.AttachTarget{ServerIdentity: attachment["id"].(string)}
			}
			batch := iaas.TargetGroupAttachmentsBatch{
				TargetGroupID: tg.Identity,
				Attachments:   attachments,
			}
			if err := client.IaaS().SetTargetGroupServerAttachments(ctx, batch); err != nil {
				return diag.FromErr(fmt.Errorf("error setting target group server attachments: %w", err))
			}
		}
		return nil
	}
	return resourceTargetGroupRead(ctx, d, m)
}

func resourceTargetGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating client: %w", err))
	}

	id := d.Get("id").(string)
	err = client.IaaS().DeleteTargetGroup(ctx, iaas.DeleteTargetGroupRequest{Identity: id})
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error deleting target group: %w", err))
	}
	d.SetId("")
	return nil
}

func expandTargetGroupHealthCheck(d *schema.ResourceData) *iaas.BackendHealthCheck {
	hcPort := d.Get("health_check_port").(int)
	hcProto := strings.TrimSpace(d.Get("health_check_protocol").(string))
	hcPath := d.Get("health_check_path").(string)
	if hcPort == 0 && hcProto == "" && hcPath == "" {
		return nil
	}
	if hcPort == 0 {
		hcPort = d.Get("port").(int)
	}
	if hcProto == "" {
		hcProto = "tcp"
	}
	return &iaas.BackendHealthCheck{
		Protocol:           iaas.LoadbalancerProtocol(hcProto),
		Port:               int32(hcPort),
		Path:               hcPath,
		PeriodSeconds:      d.Get("health_check_interval").(int),
		TimeoutSeconds:     d.Get("health_check_timeout").(int),
		HealthyThreshold:   int32(d.Get("healthy_threshold").(int)),
		UnhealthyThreshold: int32(d.Get("unhealthy_threshold").(int)),
	}
}
