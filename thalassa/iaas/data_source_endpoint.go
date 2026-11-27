package iaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func dataSourceEndpoint() *schema.Resource {
	return &schema.Resource{
		Description: "Get an endpoint by identity",
		ReadContext: dataSourceEndpointRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Identity of the endpoint",
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the endpoint. If not provided, the organisation of the (Terraform) provider will be used.",
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
	}
}

func dataSourceEndpointRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	identity := d.Get("id").(string)
	endpoint, err := client.IaaS().GetEndpoint(ctx, identity)
	if err != nil {
		if tcclient.IsNotFound(err) {
			return diag.FromErr(fmt.Errorf("endpoint %s not found", identity))
		}
		return diag.FromErr(fmt.Errorf("error getting endpoint: %w", err))
	}
	if endpoint == nil {
		return diag.FromErr(fmt.Errorf("endpoint %s not found", identity))
	}

	d.SetId(endpoint.Identity)
	_ = d.Set("name", endpoint.Name)
	_ = d.Set("endpoint_address", endpoint.EndpointAddress)
	_ = d.Set("endpoint_hostname", endpoint.EndpointHostname)
	_ = d.Set("endpoint_port", endpoint.EndpointPort)
	_ = d.Set("endpoint_type", string(endpoint.EndpointType))
	_ = d.Set("resource_type", endpoint.ResourceType)
	_ = d.Set("resource_identity", endpoint.ResourceIdentity)
	_ = d.Set("mac_address", endpoint.MacAddress)
	if endpoint.EndpointProtocol != nil {
		_ = d.Set("endpoint_protocol", string(*endpoint.EndpointProtocol))
	}

	if endpoint.Vpc != nil {
		_ = d.Set("vpc_id", endpoint.Vpc.Identity)
	}
	if endpoint.Subnet != nil {
		_ = d.Set("subnet_id", endpoint.Subnet.Identity)
	}

	return nil
}
