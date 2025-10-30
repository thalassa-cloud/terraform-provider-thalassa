package iaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	iaas "github.com/thalassa-cloud/client-go/iaas"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func DataSourceVpcPeeringConnections() *schema.Resource {
	return &schema.Resource{
		Description: "Get all VPC peering connections for the current organisation",
		ReadContext: dataSourceVpcPeeringConnectionsRead,
		Schema: map[string]*schema.Schema{
			"peering_connections": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of VPC peering connections",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identity of the VPC peering connection",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the VPC peering connection",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the VPC peering connection",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Current status of the VPC peering connection",
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
				},
			},
		},
	}
}

func dataSourceVpcPeeringConnectionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := provider.GetProvider(m)

	peeringConnections, err := provider.Client.IaaS().ListVpcPeeringConnections(ctx, &iaas.ListVpcPeeringConnectionsRequest{})
	if err != nil {
		return diag.FromErr(err)
	}

	// Convert peering connections to the expected format
	connections := make([]map[string]interface{}, len(peeringConnections))
	for i, conn := range peeringConnections {
		connMap := map[string]interface{}{
			"id":          conn.Identity,
			"name":        conn.Name,
			"description": conn.Description,
			"status":      conn.Status,
			"created_at":  conn.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			"updated_at":  conn.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}

		// Set optional fields
		if conn.StatusMessage != nil {
			connMap["status_message"] = *conn.StatusMessage
		}
		if conn.ExpiresAt != nil {
			connMap["expires_at"] = conn.ExpiresAt.Format("2006-01-02T15:04:05Z07:00")
		}
		if conn.RequesterNextHopIP != nil {
			connMap["requester_next_hop_ip"] = *conn.RequesterNextHopIP
		}
		if conn.AccepterNextHopIP != nil {
			connMap["accepter_next_hop_ip"] = *conn.AccepterNextHopIP
		}

		// Set labels and annotations
		if conn.Labels != nil {
			connMap["labels"] = conn.Labels
		}
		if conn.Annotations != nil {
			connMap["annotations"] = conn.Annotations
		}

		// Set requester VPC information
		if conn.RequesterVpc != nil {
			requesterVpc := []map[string]interface{}{
				{
					"identity": conn.RequesterVpc.Identity,
					"name":     conn.RequesterVpc.Name,
				},
			}
			connMap["requester_vpc"] = requesterVpc
		}

		// Set accepter VPC information
		if conn.AccepterVpc != nil {
			accepterVpc := []map[string]interface{}{
				{
					"identity": conn.AccepterVpc.Identity,
					"name":     conn.AccepterVpc.Name,
				},
			}
			connMap["accepter_vpc"] = accepterVpc
		}

		// Set requester organisation information
		if conn.RequesterOrganisation != nil {
			requesterOrg := []map[string]interface{}{
				{
					"identity": conn.RequesterOrganisation.Identity,
					"name":     conn.RequesterOrganisation.Name,
				},
			}
			connMap["requester_organisation"] = requesterOrg
		}

		// Set accepter organisation information
		if conn.AccepterOrganisation != nil {
			accepterOrg := []map[string]interface{}{
				{
					"identity": conn.AccepterOrganisation.Identity,
					"name":     conn.AccepterOrganisation.Name,
				},
			}
			connMap["accepter_organisation"] = accepterOrg
		}

		connections[i] = connMap
	}

	d.SetId(fmt.Sprintf("vpc-peering-connections-%d", len(connections)))
	d.Set("peering_connections", connections)

	return nil
}
