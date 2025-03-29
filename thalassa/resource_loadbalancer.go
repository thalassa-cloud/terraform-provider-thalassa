package thalassa

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"

	iaas "github.com/thalassa-cloud/client-go/pkg/iaas"
)

func resourceLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Description:   "Create an loadbalancer",
		CreateContext: resourceLoadBalancerCreate,
		ReadContext:   resourceLoadBalancerRead,
		UpdateContext: resourceLoadBalancerUpdate,
		DeleteContext: resourceLoadBalancerDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organisation": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the Loadbalancer. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"vpc": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "VPC of the Loadbalancer",
			},
			"subnet": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Subnet of the Loadbalancer",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StringLenBetween(1, 62),
				Description:  "Name of the Loadbalancer",
			},
			"slug": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.StringLenBetween(0, 255),
				Description:  "A human readable description about the loadbalancer",
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Labels for the Loadbalancer",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Annotations for the Loadbalancer",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceLoadBalancerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := getClient(getProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	createLoadbalancer := iaas.CreateLoadbalancer{
		Name:                 d.Get("name").(string),
		Description:          d.Get("description").(string),
		Labels:               convertToMap(d.Get("labels")),
		Annotations:          convertToMap(d.Get("annotations")),
		Subnet:               d.Get("subnet").(string),
		DeleteProtection:     d.Get("delete_protection").(bool),
		InternalLoadbalancer: d.Get("internal").(bool),
	}

	loadbalancer, err := client.IaaS().CreateLoadbalancer(ctx, createLoadbalancer)

	if err != nil {
		return diag.FromErr(err)
	}
	if loadbalancer != nil {
		d.SetId(loadbalancer.Identity)
		d.Set("slug", loadbalancer.Slug)
		return nil
	}
	return resourceLoadBalancerRead(ctx, d, m)
}

func resourceLoadBalancerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := getClient(getProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	slug := d.Get("slug").(string)
	loadbalancer, err := client.IaaS().GetLoadbalancer(ctx, slug)
	if err != nil && !tcclient.IsNotFound(err) {
		return diag.FromErr(fmt.Errorf("error getting loadbalancer: %s", err))
	}
	if loadbalancer == nil {
		return diag.FromErr(fmt.Errorf("loadbalancer was not found"))
	}

	d.SetId(loadbalancer.Identity)
	d.Set("name", loadbalancer.Name)
	d.Set("slug", loadbalancer.Slug)
	d.Set("description", loadbalancer.Description)
	d.Set("labels", loadbalancer.Labels)
	d.Set("annotations", loadbalancer.Annotations)
	d.Set("subnet", loadbalancer.Subnet.Identity)
	d.Set("vpc", loadbalancer.Subnet.Vpc.Identity)
	// d.Set("delete_protection", loadbalancer.DeleteProtection)
	// d.Set("internal", loadbalancer.InternalLoadbalancer)

	return nil
}

func resourceLoadBalancerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := getClient(getProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	updateLoadbalancer := iaas.UpdateLoadbalancer{
		Name:             d.Get("name").(string),
		Description:      d.Get("description").(string),
		Labels:           convertToMap(d.Get("labels")),
		Annotations:      convertToMap(d.Get("annotations")),
		DeleteProtection: d.Get("delete_protection").(bool),
	}

	slug := d.Get("slug").(string)

	loadbalancer, err := client.IaaS().UpdateLoadbalancer(ctx, slug, updateLoadbalancer)
	if err != nil {
		return diag.FromErr(err)
	}
	if loadbalancer != nil {
		d.Set("name", loadbalancer.Name)
		d.Set("description", loadbalancer.Description)
		d.Set("slug", loadbalancer.Slug)
		d.Set("labels", loadbalancer.Labels)
		d.Set("annotations", loadbalancer.Annotations)
		// d.Set("delete_protection", loadbalancer.DeleteProtection)
		return nil
	}

	return resourceLoadBalancerRead(ctx, d, m)
}

func resourceLoadBalancerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := getClient(getProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("id").(string)

	error := client.IaaS().DeleteLoadbalancer(ctx, id)
	if error != nil {
		return diag.FromErr(error)
	}

	d.SetId("")

	return nil
}
