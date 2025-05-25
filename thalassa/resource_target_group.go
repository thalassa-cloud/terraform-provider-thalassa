package thalassa

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"

	iaas "github.com/thalassa-cloud/client-go/iaas"
)

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
				Required:    true,
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
				Optional:     true,
				ValidateFunc: validate.StringLenBetween(0, 255),
				Description:  "A human readable description about the target group",
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Labels for the Target Group",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Annotations for the Target Group",
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StringInSlice([]string{"tcp", "udp", "http", "https", "grpc", "quic"}, false),
				Description:  "The protocol to use for routing traffic to the targets. Must be one of: tcp, udp, http, https, grpc, quic.",
			},
			"port": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IntBetween(1, 65535),
				Description:  "The port on which the targets receive traffic",
			},
			"health_check_path": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.StringLenBetween(1, 255),
				Description:  "The path to use for health checks (only for HTTP/HTTPS)",
			},
			"health_check_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validate.IntBetween(1, 65535),
				Description:  "The port to use for health checks",
			},
			"health_check_protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.StringInSlice([]string{"tcp", "https"}, false),
				Description:  "The protocol to use for health checks. Must be one of: tcp, http.",
			},
			"health_check_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      30,
				ValidateFunc: validate.IntBetween(5, 300),
				Description:  "The approximate amount of time, in seconds, between health checks of an individual target",
			},
			"health_check_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      5,
				ValidateFunc: validate.IntBetween(2, 60),
				Description:  "The amount of time, in seconds, during which no response means a failed health check",
			},
			"healthy_threshold": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3,
				ValidateFunc: validate.IntBetween(2, 10),
				Description:  "The number of consecutive health checks successes required before considering an unhealthy target healthy",
			},
			"unhealthy_threshold": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3,
				ValidateFunc: validate.IntBetween(2, 10),
				Description:  "The number of consecutive health check failures required before considering a target unhealthy",
			},
			"attachments": {
				Type:        schema.TypeList,
				Optional:    true,
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
	client, err := getClient(getProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	healthCheck := &iaas.BackendHealthCheck{
		Protocol:           iaas.LoadbalancerProtocol(d.Get("health_check_protocol").(string)),
		Port:               int32(d.Get("health_check_port").(int)),
		Path:               d.Get("health_check_path").(string),
		PeriodSeconds:      d.Get("health_check_interval").(int),
		TimeoutSeconds:     d.Get("health_check_timeout").(int),
		HealthyThreshold:   int32(d.Get("healthy_threshold").(int)),
		UnhealthyThreshold: int32(d.Get("unhealthy_threshold").(int)),
	}

	createTargetGroup := iaas.CreateTargetGroup{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Labels:      convertToMap(d.Get("labels")),
		Annotations: convertToMap(d.Get("annotations")),
		Vpc:         d.Get("vpc_id").(string),
		TargetPort:  d.Get("port").(int),
		Protocol:    iaas.LoadbalancerProtocol(d.Get("protocol").(string)),
		HealthCheck: healthCheck,
	}

	tg, err := client.IaaS().CreateTargetGroup(ctx, createTargetGroup)
	if err != nil {
		return diag.FromErr(err)
	}
	if tg != nil {
		d.SetId(tg.Identity)
		d.Set("slug", tg.Slug)
		// Attach targets if specified
		if targets, ok := d.GetOk("targets"); ok {
			targetList := targets.([]interface{})
			attachments := make([]iaas.AttachTarget, len(targetList))
			for i, t := range targetList {
				target := t.(map[string]interface{})
				attachments[i] = iaas.AttachTarget{ServerIdentity: target["id"].(string)}
			}
			batch := iaas.TargetGroupAttachmentsBatch{
				TargetGroupID: tg.Identity,
				Attachments:   attachments,
			}
			if err := client.IaaS().SetTargetGroupServerAttachments(ctx, batch); err != nil {
				return diag.FromErr(err)
			}
		}
		return nil
	}
	return resourceTargetGroupRead(ctx, d, m)
}

func resourceTargetGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := getClient(getProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("id").(string)
	tg, err := client.IaaS().GetTargetGroup(ctx, iaas.GetTargetGroupRequest{Identity: id})
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error getting target group: %s", err))
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
	d.Set("protocol", tg.Protocol)
	d.Set("port", tg.TargetPort)

	if tg.HealthCheck != nil {
		d.Set("health_check_path", tg.HealthCheck.Path)
		d.Set("health_check_port", tg.HealthCheck.Port)
		d.Set("health_check_protocol", tg.HealthCheck.Protocol)
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
	client, err := getClient(getProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	healthCheck := &iaas.BackendHealthCheck{
		Protocol:           iaas.LoadbalancerProtocol(d.Get("health_check_protocol").(string)),
		Port:               int32(d.Get("health_check_port").(int)),
		Path:               d.Get("health_check_path").(string),
		PeriodSeconds:      d.Get("health_check_interval").(int),
		TimeoutSeconds:     d.Get("health_check_timeout").(int),
		HealthyThreshold:   int32(d.Get("healthy_threshold").(int)),
		UnhealthyThreshold: int32(d.Get("unhealthy_threshold").(int)),
	}

	updateTargetGroup := iaas.UpdateTargetGroup{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Labels:      convertToMap(d.Get("labels")),
		Annotations: convertToMap(d.Get("annotations")),
		TargetPort:  d.Get("port").(int),
		Protocol:    iaas.LoadbalancerProtocol(d.Get("protocol").(string)),
		HealthCheck: healthCheck,
	}

	id := d.Get("id").(string)
	tg, err := client.IaaS().UpdateTargetGroup(ctx, iaas.UpdateTargetGroupRequest{
		Identity:          id,
		UpdateTargetGroup: updateTargetGroup,
	})
	if err != nil {
		return diag.FromErr(err)
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
				return diag.FromErr(err)
			}
		}
		return nil
	}
	return resourceTargetGroupRead(ctx, d, m)
}

func resourceTargetGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := getClient(getProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("id").(string)
	err = client.IaaS().DeleteTargetGroup(ctx, iaas.DeleteTargetGroupRequest{Identity: id})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
