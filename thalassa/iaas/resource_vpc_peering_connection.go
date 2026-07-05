package iaas

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"

	iaas "github.com/thalassa-cloud/client-go/iaas"
)

func resourceVpcPeeringConnection() *schema.Resource {
	return &schema.Resource{
		Description:   "Create and manage a VPC peering connection",
		CreateContext: resourceVpcPeeringConnectionCreate,
		ReadContext:   resourceVpcPeeringConnectionRead,
		UpdateContext: resourceVpcPeeringConnectionUpdate,
		DeleteContext: resourceVpcPeeringConnectionDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.StringLenBetween(1, 63),
				Description:  "Name of the VPC peering connection. Must be between 1 and 63 characters and contain only ASCII characters.",
			},
			"description": {
				Type:         schema.TypeString,
				Default:      "",
				Optional:     true,
				ValidateFunc: validate.StringLenBetween(0, 500),
				Description:  "Description of the VPC peering connection. Must be at most 500 characters and contain only ASCII characters.",
			},
			"requester_vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the VPC that will initiate the peering request",
			},
			"accepter_vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the VPC that will accept or deny the peering request",
			},
			"accepter_organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "ID of the organisation that owns the accepter VPC",
			},
			"auto_accept": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Whether the peering connection should be automatically accepted. Only allowed if requester and accepter VPCs are in the same region and owned by the same organisation.",
			},
			"labels": {
				Type:        schema.TypeMap,
				Default:     make(map[string]string),
				Optional:    true,
				Description: "Labels for the VPC peering connection",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"annotations": {
				Type:        schema.TypeMap,
				Default:     make(map[string]string),
				Optional:    true,
				Description: "Annotations for the VPC peering connection",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
			"requester_vpc": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information about the requester VPC",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"identity": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the VPC",
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
							Description: "ID of the VPC",
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
							Description: "ID of the organisation",
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

			"wait_for_active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to wait for the VPC peering connection to be active (requires acceptance by the accepter VPC owner). If false, the resource will be marked as created, but the peering connection may not be active yet.",
			},
			"wait_for_active_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
				Description: "The timeout in minutes to wait for the VPC peering connection to be active",
			},
			"wait_for_deleted_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     5,
				Description: "The timeout in minutes to wait for the VPC peering connection to be deleted. Set to 0 to disable waiting.",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceVpcPeeringConnectionCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	createRequest := iaas.CreateVpcPeeringConnectionRequest{
		Name:                         d.Get("name").(string),
		Description:                  d.Get("description").(string),
		Labels:                       convert.ConvertToMap(d.Get("labels")),
		Annotations:                  convert.ConvertToMap(d.Get("annotations")),
		RequesterVpcIdentity:         d.Get("requester_vpc_id").(string),
		AccepterVpcIdentity:          d.Get("accepter_vpc_id").(string),
		AccepterOrganisationIdentity: d.Get("accepter_organisation_id").(string),
		AutoAccept:                   d.Get("auto_accept").(bool),
	}

	peeringConnection, err := client.IaaS().CreateVpcPeeringConnection(ctx, createRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(peeringConnection.Identity)

	if d.Get("wait_for_active").(bool) {
		timeout := d.Get("wait_for_active_timeout").(int)
		ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Minute)
		defer cancel()
		if diag := waitForVpcPeeringConnectionActive(ctxWithTimeout, client.IaaS(), peeringConnection.Identity); diag != nil {
			return diag
		}
	}

	return resourceVpcPeeringConnectionRead(ctx, d, m)
}

func resourceVpcPeeringConnectionRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	peeringConnection, err := client.IaaS().GetVpcPeeringConnection(ctx, d.Id())
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	return setVpcPeeringConnectionData(d, peeringConnection)
}

func resourceVpcPeeringConnectionUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	// get the peering connection
	peeringConnection, err := client.IaaS().GetVpcPeeringConnection(ctx, d.Id())
	if err != nil {
		if tcclient.IsNotFound(err) {
			return diag.FromErr(fmt.Errorf("VPC peering connection %s was not found", d.Id()))
		}
		return diag.FromErr(err)
	}

	if peeringConnection == nil {
		return diag.FromErr(fmt.Errorf("VPC peering connection %s was not found", d.Id()))
	}

	// check if an update is needed
	if vpcPeeringConnectionMetadataMatchesState(peeringConnection, d) {
		return setVpcPeeringConnectionData(d, peeringConnection)
	}

	updateRequest := iaas.UpdateVpcPeeringConnectionRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Labels:      convert.ConvertToMap(d.Get("labels")),
		Annotations: convert.ConvertToMap(d.Get("annotations")),
	}

	peeringConnection, err = client.IaaS().UpdateVpcPeeringConnection(ctx, d.Id(), updateRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	return setVpcPeeringConnectionData(d, peeringConnection)
}

func resourceVpcPeeringConnectionDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.IaaS().DeleteVpcPeeringConnection(ctx, d.Id())
	if err != nil {
		if !tcclient.IsNotFound(err) {
			return diag.FromErr(err)
		}
	}

	// wait until the peering connection is deleted
	timeout := d.Get("wait_for_deleted_timeout").(int)
	if timeout > 0 {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Minute)
		defer cancel()
		for {
			select {
			case <-ctxWithTimeout.Done():
				return diag.FromErr(fmt.Errorf("timeout while waiting for peering connection to be deleted"))
			case <-time.After(1 * time.Second):
			}
			_, err := client.IaaS().GetVpcPeeringConnection(ctxWithTimeout, d.Id())
			if err != nil {
				if tcclient.IsNotFound(err) {
					break
				}
				return diag.FromErr(err)
			}
		}
	}
	d.SetId("")
	return nil
}

func vpcPeeringConnectionMetadataMatchesState(connection *iaas.VpcPeeringConnection, d *schema.ResourceData) bool {
	return connection.Name == d.Get("name").(string) &&
		connection.Description == d.Get("description").(string) &&
		reflect.DeepEqual(connection.Labels, convert.ConvertToMap(d.Get("labels"))) &&
		reflect.DeepEqual(connection.Annotations, convert.ConvertToMap(d.Get("annotations")))
}

func setVpcPeeringConnectionData(d *schema.ResourceData, connection *iaas.VpcPeeringConnection) diag.Diagnostics {
	d.SetId(connection.Identity)
	_ = d.Set("name", connection.Name)
	_ = d.Set("description", connection.Description)
	_ = d.Set("status", connection.Status)
	_ = d.Set("created_at", connection.CreatedAt.Format(time.RFC3339))
	_ = d.Set("updated_at", connection.UpdatedAt.Format(time.RFC3339))

	if connection.StatusMessage != nil {
		_ = d.Set("status_message", *connection.StatusMessage)
	}
	if connection.ExpiresAt != nil {
		_ = d.Set("expires_at", connection.ExpiresAt.Format(time.RFC3339))
	}
	if connection.RequesterNextHopIP != nil {
		_ = d.Set("requester_next_hop_ip", *connection.RequesterNextHopIP)
	}
	if connection.AccepterNextHopIP != nil {
		_ = d.Set("accepter_next_hop_ip", *connection.AccepterNextHopIP)
	}

	// Set labels and annotations
	if connection.Labels != nil {
		_ = d.Set("labels", connection.Labels)
	}
	if connection.Annotations != nil {
		_ = d.Set("annotations", connection.Annotations)
	}

	if connection.RequesterVpc != nil {
		_ = d.Set("requester_vpc_id", connection.RequesterVpc.Identity)
		requesterVpc := []map[string]any{
			{
				"identity": connection.RequesterVpc.Identity,
				"name":     connection.RequesterVpc.Name,
			},
		}
		_ = d.Set("requester_vpc", requesterVpc)
	}

	// Set accepter VPC information
	if connection.AccepterVpc != nil {
		_ = d.Set("accepter_vpc_id", connection.AccepterVpc.Identity)
		accepterVpc := []map[string]any{
			{
				"identity": connection.AccepterVpc.Identity,
				"name":     connection.AccepterVpc.Name,
			},
		}
		_ = d.Set("accepter_vpc", accepterVpc)
	}

	// Set requester organisation information
	if connection.RequesterOrganisation != nil {
		requesterOrg := []map[string]any{
			{
				"identity": connection.RequesterOrganisation.Identity,
				"name":     connection.RequesterOrganisation.Name,
			},
		}
		_ = d.Set("requester_organisation", requesterOrg)
	}

	// Set accepter organisation information
	if connection.AccepterOrganisation != nil {
		accepterOrg := []map[string]any{
			{
				"identity": connection.AccepterOrganisation.Identity,
				"name":     connection.AccepterOrganisation.Name,
			},
		}
		_ = d.Set("accepter_organisation", accepterOrg)
	}

	return nil
}
