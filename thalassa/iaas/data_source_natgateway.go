package iaas

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/thalassa-cloud/client-go/filters"
	iaas "github.com/thalassa-cloud/client-go/iaas"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func DataSourceNatGateway() *schema.Resource {
	return &schema.Resource{
		Description: "Get a NAT Gateway by name, slug, or other attributes",
		ReadContext: dataSourceNatGatewayRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the NAT Gateway",
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the NAT Gateway. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.StringLenBetween(1, 62),
				Description:  "The name of the NAT Gateway to look up",
			},
			"slug": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The slug of the NAT Gateway to look up",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the VPC that contains the NAT Gateway",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the subnet that contains the NAT Gateway",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The region of the NAT Gateway",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The status of the NAT Gateway to filter by",
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Labels to filter NAT Gateways by",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			// Computed fields
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A human readable description about the NAT Gateway",
			},
			"endpoint_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Endpoint IP of the NAT Gateway",
			},
			"v4_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "V4 IP of the NAT Gateway",
			},
			"v6_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "V6 IP of the NAT Gateway",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Annotations for the NAT Gateway",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp when the NAT Gateway was created",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp when the NAT Gateway was last updated",
			},
			"security_group_attachments": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of security group IDs attached to the NAT Gateway",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceNatGatewayRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	natGateways, err := provider.Client.IaaS().ListNatGateways(ctx, &iaas.ListNatGatewaysRequest{
		Filters: requestFilters,
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error listing NAT Gateways: %w", err))
	}

	var matchingNatGateways []iaas.VpcNatGateway
	name := d.Get("name").(string)
	slug := d.Get("slug").(string)
	for _, ng := range natGateways {
		// Filter by name
		if name != "" && ng.Name != name {
			continue
		}

		// Filter by slug
		if slug != "" && ng.Slug != slug {
			continue
		}

		matchingNatGateways = append(matchingNatGateways, ng)
	}

	if len(matchingNatGateways) == 0 {
		return diag.Errorf("no NAT Gateway found matching the specified criteria")
	}

	if len(matchingNatGateways) > 1 {
		var names []string
		for _, ng := range matchingNatGateways {
			names = append(names, ng.Name)
		}
		return diag.Errorf("multiple NAT Gateways found matching the criteria: %v. Please provide more specific filters", names)
	}

	natGateway := matchingNatGateways[0]

	// Set the resource data
	d.SetId(natGateway.Identity)
	d.Set("id", natGateway.Identity)
	d.Set("name", natGateway.Name)
	d.Set("slug", natGateway.Slug)
	d.Set("description", natGateway.Description)
	d.Set("endpoint_ip", natGateway.EndpointIP)
	d.Set("status", natGateway.Status)
	d.Set("v4_ip", natGateway.V4IP)
	d.Set("v6_ip", natGateway.V6IP)
	d.Set("labels", natGateway.Labels)
	d.Set("annotations", natGateway.Annotations)
	d.Set("created_at", natGateway.CreatedAt.Format(time.RFC3339))
	if !natGateway.UpdatedAt.IsZero() {
		d.Set("updated_at", natGateway.UpdatedAt.Format(time.RFC3339))
	}

	// Set VPC information
	if natGateway.Vpc != nil {
		d.Set("vpc_id", natGateway.Vpc.Identity)
	}

	// Set subnet information
	if natGateway.Subnet != nil {
		d.Set("subnet_id", natGateway.Subnet.Identity)
	}
	// Set security group attachments
	if len(natGateway.SecurityGroups) > 0 {
		securityGroupIDs := make([]string, len(natGateway.SecurityGroups))
		for i, sg := range natGateway.SecurityGroups {
			securityGroupIDs[i] = sg.Identity
		}
		d.Set("security_group_attachments", securityGroupIDs)
	}

	return nil
}
