package thalassa

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"

	iaas "github.com/thalassa-cloud/client-go/iaas"
)

const (
	ResourceVolumeAttachmentVirtualMachine = "cloud_virtual_machine"
)

func resourceBlockVolumeAttachment() *schema.Resource {
	return &schema.Resource{
		Description:   "Attach a block volume to a virtual machine. Volume must not be attached to another virtual machine.",
		CreateContext: resourceBlockVolumeAttachmentCreate,
		ReadContext:   resourceBlockVolumeAttachmentRead,
		DeleteContext: resourceBlockVolumeAttachmentDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the Volume Attachment. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"volume_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the volume to attach",
			},
			"vmi_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the virtual machine to attach the volume to",
			},
			"serial": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The device name to use for the volume attachment (e.g., /dev/sdb)",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validate.StringLenBetween(0, 255),
				Description:  "A human readable description about the volume attachment",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceBlockVolumeAttachmentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := getClient(getProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	volumeID := d.Get("volume_id").(string)
	vmiID := d.Get("vmi_id").(string)

	attachRequest := iaas.AttachVolumeRequest{
		ResourceType:     ResourceVolumeAttachmentVirtualMachine,
		ResourceIdentity: vmiID,
	}

	if v, ok := d.GetOk("description"); ok {
		attachRequest.Description = v.(string)
	}

	attachResponse, err := client.IaaS().AttachVolume(ctx, volumeID, attachRequest)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(attachResponse.Identity)
	return resourceBlockVolumeAttachmentRead(ctx, d, m)
}

func resourceBlockVolumeAttachmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := getClient(getProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	volumeID := d.Get("volume_id").(string)
	vmiID := d.Get("vmi_id").(string)

	// Get the volume to check if it's attached to the VMI
	volume, err := client.IaaS().GetVolume(ctx, volumeID)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error getting volume: %s", err))
	}

	// Check if the volume is attached to the specified VMI
	found := false
	if volume.Attachments != nil {
		for _, att := range volume.Attachments {
			if att.AttachedToResourceType == ResourceVolumeAttachmentVirtualMachine && att.AttachedToIdentity == vmiID {
				found = true
				d.Set("serial", att.Serial)
				d.Set("description", att.Description)
				break
			}
		}
	}

	if !found {
		d.SetId("")
		return nil
	}

	d.Set("volume_id", volumeID)
	d.Set("vmi_id", vmiID)

	return nil
}

func resourceBlockVolumeAttachmentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := getClient(getProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	volumeID := d.Get("volume_id").(string)
	vmiID := d.Get("vmi_id").(string)

	detachRequest := iaas.DetachVolumeRequest{
		ResourceType:     ResourceVolumeAttachmentVirtualMachine,
		ResourceIdentity: vmiID,
	}

	if err := client.IaaS().DetachVolume(ctx, volumeID, detachRequest); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
