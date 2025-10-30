package iaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	iaas "github.com/thalassa-cloud/client-go/iaas"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func DataSourceVpcPeeringConnection() *schema.Resource {
	return &schema.Resource{
		Description: "Get a VPC peering connection",
		ReadContext: dataSourceVpcPeeringConnectionRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity of the VPC peering connection",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.StringLenBetween(1, 63),
				Description:  "Name of the VPC peering connection",
			},
			"identity": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Identity of the VPC peering connection",
			},
			"requester_vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Identity of the requester VPC",
			},
			"accepter_vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Identity of the accepter VPC",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Status of the VPC peering connection",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the VPC peering connection",
			},
			"auto_accept": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the peering connection is set to auto-accept",
			},
			"labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Labels for the VPC peering connection",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Annotations for the VPC peering connection",
			},
			"status_message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Additional information about the current status",
			},
			"expires_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Time when the peering request expires if not accepted",
			},
			"requester_next_hop_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Next hop IP address for the requester VPC",
			},
			"accepter_next_hop_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Next hop IP address for the accepter VPC",
			},
			"requester_vpc": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information about the requester VPC",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"identity": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identity of the VPC",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the VPC",
						},
					},
				},
			},
			"accepter_vpc": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information about the accepter VPC",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"identity": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identity of the VPC",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the VPC",
						},
					},
				},
			},
			"requester_organisation": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information about the requester organisation",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"identity": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identity of the organisation",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the organisation",
						},
					},
				},
			},
			"accepter_organisation": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information about the accepter organisation",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"identity": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identity of the organisation",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the organisation",
						},
					},
				},
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Time when the VPC peering connection was created",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Time when the VPC peering connection was last updated",
			},
		},
	}
}

func dataSourceVpcPeeringConnectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := provider.GetProvider(m)
	identity := d.Get("identity").(string)
	name := d.Get("name").(string)
	requesterVpcId := d.Get("requester_vpc_id").(string)
	accepterVpcId := d.Get("accepter_vpc_id").(string)
	status := d.Get("status").(string)

	// If identity is provided, get the specific peering connection
	if identity != "" {
		peeringConnection, err := provider.Client.IaaS().GetVpcPeeringConnection(ctx, identity)
		if err != nil {
			return diag.FromErr(err)
		}
		return setVpcPeeringConnectionData(d, peeringConnection)
	}

	// Otherwise, list peering connections and find by criteria
	peeringConnections, err := provider.Client.IaaS().ListVpcPeeringConnections(ctx, &iaas.ListVpcPeeringConnectionsRequest{})
	if err != nil {
		return diag.FromErr(err)
	}

	var peeringConnection *iaas.VpcPeeringConnection
	for _, conn := range peeringConnections {
		// Filter by name if provided
		if name != "" && conn.Name != name {
			continue
		}
		// Filter by requester VPC if provided
		if requesterVpcId != "" && (conn.RequesterVpc == nil || conn.RequesterVpc.Identity != requesterVpcId) {
			continue
		}
		// Filter by accepter VPC if provided
		if accepterVpcId != "" && (conn.AccepterVpc == nil || conn.AccepterVpc.Identity != accepterVpcId) {
			continue
		}
		// Filter by status if provided
		if status != "" && conn.Status != iaas.VpcPeeringConnectionStatus(status) {
			continue
		}

		peeringConnection = &conn
		break
	}

	if peeringConnection == nil {
		return diag.FromErr(fmt.Errorf("VPC peering connection not found with the given criteria"))
	}

	return setVpcPeeringConnectionData(d, peeringConnection)
}
