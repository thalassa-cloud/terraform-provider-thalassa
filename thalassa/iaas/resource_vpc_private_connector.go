package iaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	iaas "github.com/thalassa-cloud/client-go/iaas"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func resourceVpcPrivateConnector() *schema.Resource {
	return &schema.Resource{
		Description: "Create a VPC private connector for cloud services within a VPC. " +
			"The connector provides private connectivity to platform cloud services.",
		CreateContext: resourceVpcPrivateConnectorCreate,
		ReadContext:   resourceVpcPrivateConnectorRead,
		UpdateContext: resourceVpcPrivateConnectorUpdate,
		DeleteContext: resourceVpcPrivateConnectorDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity of the VPC private connector",
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the VPC private connector. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Identity of the VPC that the private connector belongs to",
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StringInSlice([]string{string(iaas.VpcPrivateConnectorTypeCloudServices)}, false),
				Description:  "Type of the VPC private connector. Currently only cloud_services is supported.",
			},
			"attached_security_groups": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of security group identities attached to the VPC private connector",
				Elem: &schema.Schema{
					Type:        schema.TypeString,
					Description: "Identity of a security group",
				},
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Labels for the VPC private connector",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
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
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceVpcPrivateConnectorCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	createReq := iaas.CreateVpcPrivateConnectorRequest{
		VpcIdentity:            d.Get("vpc_id").(string),
		Type:                   iaas.VpcPrivateConnectorType(d.Get("type").(string)),
		AttachedSecurityGroups: convert.ConvertToStringSlice(d.Get("attached_security_groups")),
		Labels:                 convert.ConvertToMap(d.Get("labels")),
		Annotations:            convert.ConvertToMap(d.Get("annotations")),
	}

	connector, err := client.IaaS().CreateVpcPrivateConnector(ctx, createReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating VPC private connector: %w", err))
	}
	if connector == nil {
		return diag.FromErr(fmt.Errorf("VPC private connector was not returned after creation"))
	}

	d.SetId(connector.Identity)
	return resourceVpcPrivateConnectorRead(ctx, d, m)
}

func resourceVpcPrivateConnectorRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	identity := d.Id()
	connector, err := client.IaaS().GetVpcPrivateConnector(ctx, identity)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error getting VPC private connector: %w", err))
	}
	if connector == nil {
		d.SetId("")
		return nil
	}

	return setVpcPrivateConnectorState(d, connector)
}

func resourceVpcPrivateConnectorUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	updateReq := iaas.UpdateVpcPrivateConnectorRequest{
		AttachedSecurityGroups: convert.ConvertToStringSlice(d.Get("attached_security_groups")),
		Labels:                 convert.ConvertToMap(d.Get("labels")),
		Annotations:            convert.ConvertToMap(d.Get("annotations")),
	}

	_, err = client.IaaS().UpdateVpcPrivateConnector(ctx, d.Id(), updateReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating VPC private connector: %w", err))
	}

	return resourceVpcPrivateConnectorRead(ctx, d, m)
}

func resourceVpcPrivateConnectorDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	if id == "" {
		return nil
	}

	err = client.IaaS().DeleteVpcPrivateConnector(ctx, id)
	if err != nil && !tcclient.IsNotFound(err) {
		return diag.FromErr(fmt.Errorf("error deleting VPC private connector: %w", err))
	}

	d.SetId("")
	return nil
}

func setVpcPrivateConnectorState(d *schema.ResourceData, connector *iaas.VpcPrivateConnector) diag.Diagnostics {
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
