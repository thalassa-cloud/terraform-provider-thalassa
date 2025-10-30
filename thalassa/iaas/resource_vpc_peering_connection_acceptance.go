package iaas

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"

	iaas "github.com/thalassa-cloud/client-go/iaas"
)

func resourceVpcPeeringConnectionAcceptance() *schema.Resource {
	return &schema.Resource{
		Description:   "Accept a VPC peering connection",
		CreateContext: resourceVpcPeeringConnectionAcceptanceCreate,
		ReadContext:   resourceVpcPeeringConnectionAcceptanceRead,
		UpdateContext: updateVpcPeeringConnectionAcceptance,
		DeleteContext: resourceVpcPeeringConnectionAcceptanceDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"peering_connection_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Identity of the VPC peering connection to accept",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current status of the VPC peering connection after the action",
			},
			"status_message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Additional information about the current status",
			},
			"wait_for_active": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,

				Description: "Whether to wait for the VPC peering connection to be active",
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

func resourceVpcPeeringConnectionAcceptanceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	peeringConnectionId := d.Get("peering_connection_id").(string)

	// ensure the peering connection is in the pending state
	peeringConnection, err := client.IaaS().GetVpcPeeringConnection(ctx, peeringConnectionId)
	if err != nil {
		return diag.FromErr(err)
	}
	if peeringConnection == nil {
		return diag.FromErr(fmt.Errorf("VPC peering connection %s was not found", peeringConnectionId))
	}
	if strings.EqualFold(string(peeringConnection.Status), "active") {
		return setVpcPeeringConnectionAcceptanceData(d, peeringConnection)
	}
	// else, accept the peering connection

	acceptRequest := iaas.AcceptVpcPeeringConnectionRequest{}
	peeringConnection, err = client.IaaS().AcceptVpcPeeringConnection(ctx, peeringConnectionId, acceptRequest)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to accept VPC peering connection %s: %w", peeringConnectionId, err))
	}

	d.SetId(peeringConnectionId)

	if d.Get("wait_for_active").(bool) {
		timeout := d.Get("wait_for_active_timeout").(int)
		ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Minute)
		defer cancel()
		if diag := waitForVpcPeeringConnectionActive(ctxWithTimeout, client.IaaS(), peeringConnectionId); diag != nil {
			return diag
		}
	}
	return setVpcPeeringConnectionAcceptanceData(d, peeringConnection)
}

func updateVpcPeeringConnectionAcceptance(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	// check if the peering connection is already accepted
	peeringConnection, err := client.IaaS().GetVpcPeeringConnection(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if peeringConnection == nil {
		return diag.FromErr(fmt.Errorf("VPC peering connection %s was not found", d.Id()))
	}
	if strings.EqualFold(string(peeringConnection.Status), "active") {
		return setVpcPeeringConnectionAcceptanceData(d, peeringConnection)
	}

	peeringConnectionId := d.Get("peering_connection_id").(string)
	acceptRequest := iaas.AcceptVpcPeeringConnectionRequest{}
	peeringConnection, err = client.IaaS().AcceptVpcPeeringConnection(ctx, peeringConnectionId, acceptRequest)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to accept VPC peering connection %s: %w", peeringConnectionId, err))
	}

	return setVpcPeeringConnectionAcceptanceData(d, peeringConnection)
}

func resourceVpcPeeringConnectionAcceptanceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	peeringConnectionId := d.Get("peering_connection_id").(string)
	peeringConnection, err := client.IaaS().GetVpcPeeringConnection(ctx, peeringConnectionId)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	return setVpcPeeringConnectionAcceptanceData(d, peeringConnection)
}

func resourceVpcPeeringConnectionAcceptanceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to get client: %w", err))
	}

	peeringConnectionId := d.Get("peering_connection_id").(string)
	if err := client.IaaS().DeleteVpcPeeringConnection(ctx, peeringConnectionId); err != nil {
		if !tcclient.IsNotFound(err) {
			return diag.FromErr(fmt.Errorf("failed to reject VPC peering connection %s: %w", peeringConnectionId, err))
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
			_, err := client.IaaS().GetVpcPeeringConnection(ctxWithTimeout, peeringConnectionId)
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

func setVpcPeeringConnectionAcceptanceData(d *schema.ResourceData, connection *iaas.VpcPeeringConnection) diag.Diagnostics {
	d.Set("peering_connection_id", connection.Identity)
	d.Set("status", connection.Status)

	if connection.StatusMessage != nil {
		d.Set("status_message", *connection.StatusMessage)
	}
	return nil
}

func waitForVpcPeeringConnectionActive(ctx context.Context, client *iaas.Client, peeringConnectionId string) diag.Diagnostics {
	for {
		select {
		case <-ctx.Done():
			return diag.FromErr(fmt.Errorf("timeout waiting for VPC peering connection to be active"))
		case <-time.After(1 * time.Second):
		}
		peeringConnection, err := client.GetVpcPeeringConnection(ctx, peeringConnectionId)
		if err != nil {
			if tcclient.IsNotFound(err) {
				return diag.FromErr(fmt.Errorf("VPC peering connection %s was not found", peeringConnectionId))
			}
			return diag.FromErr(fmt.Errorf("failed to get VPC peering connection %s: %w", peeringConnectionId, err))
		}
		if peeringConnection == nil {
			return diag.FromErr(fmt.Errorf("failed to get VPC peering connection %s: %w", peeringConnectionId, err))
		}
		if strings.EqualFold(string(peeringConnection.Status), "active") {
			break
		}
		if strings.EqualFold(string(peeringConnection.Status), "rejected") {
			return diag.FromErr(fmt.Errorf("VPC peering connection %s was rejected", peeringConnectionId))
		}
		if strings.EqualFold(string(peeringConnection.Status), "failed") {
			return diag.FromErr(fmt.Errorf("VPC peering connection %s failed", peeringConnectionId))
		}
	}
	return nil
}
