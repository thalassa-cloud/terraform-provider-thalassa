package iaas

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"

	iaas "github.com/thalassa-cloud/client-go/iaas"
)

// validateCIDRorIP validates that a string is either a valid CIDR block or IP address
func validateCIDRorIP(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return warnings, errors
	}

	// Check if it's a CIDR block
	if strings.Contains(v, "/") {
		_, _, err := net.ParseCIDR(v)
		if err != nil {
			errors = append(errors, fmt.Errorf("expected %s to be a valid CIDR block, got: %s", k, v))
		}
		return warnings, errors
	}

	// Check if it's a valid IP address
	ip := net.ParseIP(v)
	if ip == nil {
		errors = append(errors, fmt.Errorf("expected %s to be a valid IP address or CIDR block, got: %s", k, v))
	}
	return warnings, errors
}

func resourceLoadBalancerListener() *schema.Resource {
	return &schema.Resource{
		Description:   "Create a listener for a loadbalancer",
		CreateContext: resourceLoadBalancerListenerCreate,
		ReadContext:   resourceLoadBalancerListenerRead,
		UpdateContext: resourceLoadBalancerListenerUpdate,
		DeleteContext: resourceLoadBalancerListenerDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the Loadbalancer Listener. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"loadbalancer_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the loadbalancer to create the listener on",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.StringLenBetween(1, 62),
				Description:  "Name of the Loadbalancer Listener",
			},
			"slug": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.StringLenBetween(0, 255),
				Description:  "A human readable description about the loadbalancer listener",
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Labels for the Loadbalancer Listener",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Annotations for the Loadbalancer Listener",
			},
			"port": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validate.IntBetween(1, 65535),
				Description:  "The port the listener is listening on",
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.StringInSlice([]string{"http", "https", "tcp", "udp", "grpc", "quic"}, false),
				Description:  "The protocol the listener is using",
			},
			"target_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the target group to attach to the listener",
			},
			"max_connections": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The maximum number of connections that the listener can handle",
			},
			"connection_idle_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The amount of seconds used for configuring the idle connection timeout on a listener",
			},
			"allowed_sources": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateCIDRorIP,
				},
				Description: "A list of CIDR blocks or IP addresses that are allowed to connect to the listener",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceLoadBalancerListenerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	loadbalancerID := d.Get("loadbalancer_id").(string)
	createListener := iaas.CreateListener{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Labels:      convert.ConvertToMap(d.Get("labels")),
		Annotations: convert.ConvertToMap(d.Get("annotations")),
		Port:        d.Get("port").(int),
		Protocol:    iaas.LoadbalancerProtocol(d.Get("protocol").(string)),
		TargetGroup: d.Get("target_group_id").(string),
	}

	if v, ok := d.GetOk("max_connections"); ok {
		maxConn := uint32(v.(int))
		createListener.MaxConnections = &maxConn
	}

	if v, ok := d.GetOk("connection_idle_timeout"); ok {
		timeout := uint32(v.(int))
		createListener.ConnectionIdleTimeout = &timeout
	}

	if v, ok := d.GetOk("allowed_sources"); ok {
		sources := make([]string, len(v.([]interface{})))
		for i, source := range v.([]interface{}) {
			sources[i] = source.(string)
		}
		createListener.AllowedSources = sources
	}

	listener, err := client.IaaS().CreateListener(ctx, loadbalancerID, createListener)
	if err != nil {
		return diag.FromErr(err)
	}

	if listener != nil {
		d.SetId(listener.Identity)
		d.Set("slug", listener.Slug)
		return nil
	}

	return resourceLoadBalancerListenerRead(ctx, d, m)
}

func resourceLoadBalancerListenerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	loadbalancerID := d.Get("loadbalancer_id").(string)
	listenerID := d.Id()

	listener, err := client.IaaS().GetListener(ctx, iaas.GetLoadbalancerListenerRequest{
		Loadbalancer: loadbalancerID,
		Listener:     listenerID,
	})
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error getting loadbalancer listener: %s", err))
	}

	if listener == nil {
		d.SetId("")
		return nil
	}

	d.SetId(listener.Identity)
	d.Set("name", listener.Name)
	d.Set("slug", listener.Slug)
	d.Set("description", listener.Description)
	d.Set("labels", listener.Labels)
	d.Set("annotations", listener.Annotations)
	d.Set("port", listener.Port)
	d.Set("protocol", listener.Protocol)
	if listener.TargetGroup != nil {
		d.Set("target_group_id", listener.TargetGroup.Identity)
	}
	if listener.MaxConnections != nil {
		d.Set("max_connections", *listener.MaxConnections)
	}
	if listener.ConnectionIdleTimeout != nil {
		d.Set("connection_idle_timeout", *listener.ConnectionIdleTimeout)
	}
	d.Set("allowed_sources", listener.AllowedSources)

	return nil
}

func resourceLoadBalancerListenerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	loadbalancerID := d.Get("loadbalancer_id").(string)
	listenerID := d.Id()

	updateListener := iaas.UpdateListener{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Labels:      convert.ConvertToMap(d.Get("labels")),
		Annotations: convert.ConvertToMap(d.Get("annotations")),
		Port:        d.Get("port").(int),
		Protocol:    iaas.LoadbalancerProtocol(d.Get("protocol").(string)),
		TargetGroup: d.Get("target_group_id").(string),
	}

	if v, ok := d.GetOk("max_connections"); ok {
		maxConn := uint32(v.(int))
		updateListener.MaxConnections = &maxConn
	}

	if v, ok := d.GetOk("connection_idle_timeout"); ok {
		timeout := uint32(v.(int))
		updateListener.ConnectionIdleTimeout = &timeout
	}

	if v, ok := d.GetOk("allowed_sources"); ok {
		sources := make([]string, len(v.([]interface{})))
		for i, source := range v.([]interface{}) {
			sources[i] = source.(string)
		}
		updateListener.AllowedSources = sources
	}

	listener, err := client.IaaS().UpdateListener(ctx, loadbalancerID, listenerID, updateListener)
	if err != nil {
		return diag.FromErr(err)
	}

	if listener != nil {
		d.Set("name", listener.Name)
		d.Set("description", listener.Description)
		d.Set("labels", listener.Labels)
		d.Set("annotations", listener.Annotations)
		d.Set("port", listener.Port)
		d.Set("protocol", listener.Protocol)
		if listener.TargetGroup != nil {
			d.Set("target_group_id", listener.TargetGroup.Identity)
		}
		if listener.MaxConnections != nil {
			d.Set("max_connections", *listener.MaxConnections)
		}
		if listener.ConnectionIdleTimeout != nil {
			d.Set("connection_idle_timeout", *listener.ConnectionIdleTimeout)
		}
		d.Set("allowed_sources", listener.AllowedSources)
		return nil
	}

	return resourceLoadBalancerListenerRead(ctx, d, m)
}

func resourceLoadBalancerListenerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	loadbalancerID := d.Get("loadbalancer_id").(string)
	listenerID := d.Id()

	if err := client.IaaS().DeleteListener(ctx, loadbalancerID, listenerID); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
