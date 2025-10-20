package iaas

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/thalassa-cloud/client-go/filters"
	iaas "github.com/thalassa-cloud/client-go/iaas"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func DataSourceLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Description: "Get a Load Balancer by name, slug, or other attributes",
		ReadContext: dataSourceLoadBalancerRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the Load Balancer",
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the Load Balancer. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.StringLenBetween(1, 62),
				Description:  "The name of the Load Balancer to look up",
			},
			"slug": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The slug of the Load Balancer to look up",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the VPC that contains the Load Balancer",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the subnet that contains the Load Balancer",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The region of the Load Balancer",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The status of the Load Balancer to filter by",
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Labels to filter Load Balancers by",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			// Computed fields
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A human readable description about the Load Balancer",
			},
			"ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IP address of the Load Balancer",
			},
			"external_ip_addresses": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The external IP addresses of the Load Balancer",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"delete_protection": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Delete protection for the Load Balancer",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Annotations for the Load Balancer",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp when the Load Balancer was created",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp when the Load Balancer was last updated",
			},
			"security_group_attachments": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of security group IDs attached to the Load Balancer",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"listeners": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of listeners attached to the Load Balancer",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier of the listener",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the listener",
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The port of the listener",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The protocol of the listener",
						},
						"target_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The target group ID associated with the listener",
						},
					},
				},
			},
		},
	}
}

func dataSourceLoadBalancerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := provider.GetProvider(m)

	// Build filters based on provided criteria
	var requestFilters []filters.Filter
	vpcID := d.Get("vpc_id").(string)
	subnetID := d.Get("subnet_id").(string)
	region := d.Get("region").(string)
	status := d.Get("status").(string)
	labels := d.Get("labels").(map[string]interface{})

	// Add VPC filter
	if vpcID != "" {
		requestFilters = append(requestFilters, &filters.FilterKeyValue{
			Key:   filters.FilterVpcIdentity,
			Value: vpcID,
		})
	}

	// Add subnet filter
	if subnetID != "" {
		requestFilters = append(requestFilters, &filters.FilterKeyValue{
			Key:   filters.FilterSubnetIdentity,
			Value: subnetID,
		})
	}

	// Add region filter
	if region != "" {
		requestFilters = append(requestFilters, &filters.FilterKeyValue{
			Key:   filters.FilterRegion,
			Value: region,
		})
	}

	// Add status filter
	if status != "" {
		requestFilters = append(requestFilters, &filters.FilterKeyValue{
			Key:   "status",
			Value: status,
		})
	}

	if len(labels) > 0 {
		requestFilters = append(requestFilters, &filters.LabelFilter{
			MatchLabels: convert.ConvertToMap(labels),
		})
	}

	loadBalancers, err := provider.Client.IaaS().ListLoadbalancers(ctx, &iaas.ListLoadbalancersRequest{
		Filters: requestFilters,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	var matchingLoadBalancers []iaas.VpcLoadbalancer
	name := d.Get("name").(string)
	slug := d.Get("slug").(string)

	// Filter by name and slug (these are not supported by API filters)
	for _, lb := range loadBalancers {
		// Filter by name
		if name != "" && lb.Name != name {
			continue
		}

		// Filter by slug
		if slug != "" && lb.Slug != slug {
			continue
		}

		matchingLoadBalancers = append(matchingLoadBalancers, lb)
	}

	if len(matchingLoadBalancers) == 0 {
		return diag.Errorf("no Load Balancer found matching the specified criteria")
	}

	if len(matchingLoadBalancers) > 1 {
		var names []string
		for _, lb := range matchingLoadBalancers {
			names = append(names, lb.Name)
		}
		return diag.Errorf("multiple Load Balancers found matching the criteria: %v. Please provide more specific filters", names)
	}

	loadBalancer := matchingLoadBalancers[0]

	// Set the resource data
	d.SetId(loadBalancer.Identity)
	d.Set("id", loadBalancer.Identity)
	d.Set("name", loadBalancer.Name)
	d.Set("slug", loadBalancer.Slug)
	d.Set("description", loadBalancer.Description)
	if len(loadBalancer.ExternalIpAddresses) > 0 {
		d.Set("ip_address", loadBalancer.ExternalIpAddresses[0])
	}
	d.Set("external_ip_addresses", loadBalancer.ExternalIpAddresses)
	d.Set("delete_protection", loadBalancer.DeleteProtection)
	// Note: Internal field is not available on VpcLoadbalancer type
	d.Set("status", loadBalancer.Status)
	d.Set("labels", loadBalancer.Labels)
	d.Set("annotations", loadBalancer.Annotations)
	d.Set("created_at", loadBalancer.CreatedAt.Format(time.RFC3339))
	if !loadBalancer.UpdatedAt.IsZero() {
		d.Set("updated_at", loadBalancer.UpdatedAt.Format(time.RFC3339))
	}

	// Set VPC information
	if loadBalancer.Vpc != nil {
		d.Set("vpc_id", loadBalancer.Vpc.Identity)
	}

	// Set subnet information
	if loadBalancer.Subnet != nil {
		d.Set("subnet_id", loadBalancer.Subnet.Identity)
	}

	// Note: Region information is not directly available on VpcLoadbalancer

	// Set security group attachments
	if len(loadBalancer.SecurityGroups) > 0 {
		securityGroupIDs := make([]string, len(loadBalancer.SecurityGroups))
		for i, sg := range loadBalancer.SecurityGroups {
			securityGroupIDs[i] = sg.Identity
		}
		d.Set("security_group_attachments", securityGroupIDs)
	}

	// Note: Listeners information would need separate API call to fetch

	return nil
}
