package iaas

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"

	iaas "github.com/thalassa-cloud/client-go/iaas"
)

func resourceLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Description:   "Create an loadbalancer within a VPC",
		CreateContext: resourceLoadBalancerCreate,
		ReadContext:   resourceLoadBalancerRead,
		UpdateContext: resourceLoadBalancerUpdate,
		DeleteContext: resourceLoadBalancerDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the Loadbalancer. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Region of the Loadbalancer",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "VPC of the Loadbalancer",
			},
			"subnet_id": {
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
			"delete_protection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Delete protection for the Loadbalancer",
			},
			"internal": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Internal loadbalancer",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceLoadBalancerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	createLoadbalancer := iaas.CreateLoadbalancer{
		Name:                 d.Get("name").(string),
		Description:          d.Get("description").(string),
		Labels:               convert.ConvertToMap(d.Get("labels")),
		Annotations:          convert.ConvertToMap(d.Get("annotations")),
		Subnet:               d.Get("subnet_id").(string),
		DeleteProtection:     d.Get("delete_protection").(bool),
		InternalLoadbalancer: d.Get("internal").(bool),
	}

	if deleteProtection := d.Get("delete_protection").(bool); deleteProtection {
		createLoadbalancer.DeleteProtection = deleteProtection
	}
	if internal := d.Get("internal").(bool); internal {
		createLoadbalancer.InternalLoadbalancer = internal
	}

	loadbalancer, err := client.IaaS().CreateLoadbalancer(ctx, createLoadbalancer)

	if err != nil {
		return diag.FromErr(err)
	}

	if loadbalancer == nil {
		return diag.FromErr(fmt.Errorf("error creating loadbalancer: %s", err))
	}
	d.SetId(loadbalancer.Identity)
	d.Set("name", loadbalancer.Name)
	d.Set("slug", loadbalancer.Slug)
	d.Set("description", loadbalancer.Description)
	d.Set("labels", loadbalancer.Labels)
	d.Set("annotations", loadbalancer.Annotations)
	d.Set("subnet_id", loadbalancer.Subnet.Identity)
	d.Set("vpc_id", loadbalancer.Vpc.Identity)
	if d.Get("delete_protection") != nil {
		d.Set("delete_protection", d.Get("delete_protection").(bool))
	}
	if d.Get("internal") != nil {
		d.Set("internal", d.Get("internal").(bool))
	}

	// wait until the loadbalancer is ready
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 20*time.Minute)
	defer cancel()
	for {
		select {
		case <-ctxWithTimeout.Done():
			return diag.FromErr(fmt.Errorf("timeout while waiting for loadbalancer to be ready"))
		case <-time.After(1 * time.Second):
		}
		loadbalancer, err = client.IaaS().GetLoadbalancer(ctxWithTimeout, loadbalancer.Identity)
		if err != nil {
			return diag.FromErr(err)
		}
		if loadbalancer.Status == "ready" {
			break
		}
	}

	return resourceLoadBalancerRead(ctx, d, m)
}

func resourceLoadBalancerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	slug := d.Get("id").(string)
	loadbalancer, err := client.IaaS().GetLoadbalancer(ctx, slug)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error getting loadbalancer: %s", err))
	}
	if loadbalancer == nil {
		d.SetId("")
		return nil
	}

	d.SetId(loadbalancer.Identity)
	d.Set("name", loadbalancer.Name)
	d.Set("slug", loadbalancer.Slug)
	d.Set("description", loadbalancer.Description)
	d.Set("labels", loadbalancer.Labels)
	d.Set("annotations", loadbalancer.Annotations)
	d.Set("subnet_id", loadbalancer.Subnet.Identity)
	d.Set("vpc_id", loadbalancer.Vpc.Identity)
	// d.Set("delete_protection", loadbalancer.DeleteProtection)
	// d.Set("internal", loadbalancer.InternalLoadbalancer)

	return nil
}

func resourceLoadBalancerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	updateLoadbalancer := iaas.UpdateLoadbalancer{
		Name:             d.Get("name").(string),
		Description:      d.Get("description").(string),
		Labels:           convert.ConvertToMap(d.Get("labels")),
		Annotations:      convert.ConvertToMap(d.Get("annotations")),
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

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 20*time.Minute)
	defer cancel()
	for {
		select {
		case <-ctxWithTimeout.Done():
			return diag.FromErr(fmt.Errorf("timeout while waiting for loadbalancer to be ready"))
		case <-time.After(1 * time.Second):
		}
		// continue
		loadbalancer, err = client.IaaS().GetLoadbalancer(ctxWithTimeout, slug)
		if err != nil {
			if tcclient.IsNotFound(err) {
				return diag.FromErr(fmt.Errorf("loadbalancer %s was not found after update", slug))
			}
			return diag.FromErr(err)
		}
		if loadbalancer.Status == "ready" {
			break
		}
	}

	return resourceLoadBalancerRead(ctx, d, m)
}

func resourceLoadBalancerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("id").(string)
	if err := client.IaaS().DeleteLoadbalancer(ctx, id); err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	// wait until the loadbalancer is deleted
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 20*time.Minute)
	defer cancel()

	for {
		select {
		case <-ctxWithTimeout.Done():
			return diag.FromErr(fmt.Errorf("timeout while waiting for loadbalancer to be deleted"))
		case <-time.After(1 * time.Second):
		}
		loadbalancer, err := client.IaaS().GetLoadbalancer(ctxWithTimeout, id)
		if err != nil {
			if tcclient.IsNotFound(err) {
				d.SetId("")
				return nil
			}
			return diag.FromErr(err)
		}
		if loadbalancer.Status == "deleted" {
			break
		}
	}

	d.SetId("")
	return nil
}
