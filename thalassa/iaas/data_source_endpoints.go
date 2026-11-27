package iaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	iaas "github.com/thalassa-cloud/client-go/iaas"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func dataSourceEndpoints() *schema.Resource {
	return &schema.Resource{
		Description: "List endpoints. Optionally filter by VPC identity or resource identity.",
		ReadContext: dataSourceEndpointsRead,
		Schema: map[string]*schema.Schema{
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the endpoints. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter endpoints by VPC identity",
			},
			"resource_identity": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter endpoints by the identity of the attached resource",
			},
			"endpoints": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of endpoints that match the filter criteria",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identity of the endpoint",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the endpoint",
						},
						"endpoint_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Address of the endpoint",
						},
						"endpoint_hostname": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Hostname of the endpoint",
						},
						"endpoint_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Port of the endpoint",
						},
						"endpoint_protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Protocol of the endpoint",
						},
						"endpoint_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the endpoint (vpc or platform)",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the resource that the endpoint is attached to",
						},
						"resource_identity": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identity of the resource that the endpoint is attached to",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identity of the VPC that the endpoint belongs to",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identity of the subnet that the endpoint belongs to",
						},
						"mac_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "MAC address of the endpoint",
						},
					},
				},
			},
		},
	}
}

func dataSourceEndpointsRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	endpoints, err := client.IaaS().ListEndpoints(ctx, &iaas.ListEndpointsRequest{})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error listing endpoints: %w", err))
	}

	vpcID := d.Get("vpc_id").(string)
	resourceIdentity := d.Get("resource_identity").(string)

	result := make([]map[string]any, 0, len(endpoints))
	for _, endpoint := range endpoints {
		if vpcID != "" {
			if endpoint.Vpc == nil || endpoint.Vpc.Identity != vpcID {
				continue
			}
		}
		if resourceIdentity != "" && endpoint.ResourceIdentity != resourceIdentity {
			continue
		}

		item := map[string]any{
			"id":                endpoint.Identity,
			"name":              endpoint.Name,
			"endpoint_address":  endpoint.EndpointAddress,
			"endpoint_hostname": endpoint.EndpointHostname,
			"endpoint_port":     endpoint.EndpointPort,
			"endpoint_type":     string(endpoint.EndpointType),
			"resource_type":     endpoint.ResourceType,
			"resource_identity": endpoint.ResourceIdentity,
			"mac_address":       endpoint.MacAddress,
		}
		if endpoint.EndpointProtocol != nil {
			item["endpoint_protocol"] = string(*endpoint.EndpointProtocol)
		}
		if endpoint.Vpc != nil {
			item["vpc_id"] = endpoint.Vpc.Identity
		}
		if endpoint.Subnet != nil {
			item["subnet_id"] = endpoint.Subnet.Identity
		}
		result = append(result, item)
	}

	id := "endpoints"
	if vpcID != "" {
		id = fmt.Sprintf("endpoints-vpc-%s", vpcID)
	}
	if resourceIdentity != "" {
		id = fmt.Sprintf("%s-resource-%s", id, resourceIdentity)
	}

	d.SetId(id)
	_ = d.Set("endpoints", result)
	return nil
}
