package iaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func dataSourceVpcPrivateConnector() *schema.Resource {
	return &schema.Resource{
		Description: "Get a VPC private connector by ID",
		ReadContext: dataSourceVpcPrivateConnectorRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the VPC private connector",
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the VPC private connector. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the VPC that the private connector belongs to",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the VPC private connector",
			},
			"attached_security_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of security group IDs attached to the VPC private connector",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Labels for the VPC private connector",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"annotations": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Annotations for the VPC private connector",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the VPC private connector",
			},
			"status_message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Additional information about the current status",
			},
			"endpoint_ipv4": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IPv4 endpoint address of the VPC private connector",
			},
			"endpoint_ipv6": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IPv6 endpoint address of the VPC private connector",
			},
		},
	}
}

func dataSourceVpcPrivateConnectorRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	identity := d.Get("id").(string)
	connector, err := client.IaaS().GetVpcPrivateConnector(ctx, identity)
	if err != nil {
		if tcclient.IsNotFound(err) {
			return diag.FromErr(fmt.Errorf("VPC private connector %s not found", identity))
		}
		return diag.FromErr(fmt.Errorf("error getting VPC private connector: %w", err))
	}
	if connector == nil {
		return diag.FromErr(fmt.Errorf("VPC private connector %s not found", identity))
	}

	d.SetId(connector.Identity)
	_ = d.Set("type", string(connector.Type))
	_ = d.Set("attached_security_groups", connector.AttachedSecurityGroups)
	_ = d.Set("labels", connector.Labels)
	_ = d.Set("annotations", connector.Annotations)
	_ = d.Set("status", string(connector.Status))
	_ = d.Set("status_message", convert.StringValue(connector.StatusMessage))
	_ = d.Set("endpoint_ipv4", convert.StringValue(connector.EndpointIPv4))
	_ = d.Set("endpoint_ipv6", convert.StringValue(connector.EndpointIPv6))

	if connector.Vpc != nil {
		_ = d.Set("vpc_id", connector.Vpc.Identity)
	}

	return nil
}
