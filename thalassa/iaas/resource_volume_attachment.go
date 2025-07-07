package iaas

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"

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
		UpdateContext: resourceBlockVolumeAttachmentUpdate,
		DeleteContext: resourceBlockVolumeAttachmentDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
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
			"wait_for_attached": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    false,
				Default:     true,
				Description: "Wait for the volume to be attached to the virtual machine. If false, the volume will be attached and the resource will be marked as created, but the volume may not be attached to the virtual machine yet.",
			},
			"wait_for_attached_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     5,
				ForceNew:    false,
				Description: "The timeout in minutes to wait for the volume to be attached to the virtual machine. Only used if wait_for_attached is true. If not provided, the default timeout of 5 minutes will be used.",
			},
			"wait_for_detached_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    false,
				Default:     5,
				Description: "The timeout in minutes to wait for the volume to be detached from the virtual machine. Only used if wait_for_detached is true. If not provided, the default timeout of 5 minutes will be used.",
			},
			"wait_for_detached": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    false,
				Default:     true,
				Description: "Wait for the volume to be detached from the virtual machine. If false, the volume will be detached and the resource will be marked as deleted, but the volume may not be detached from the virtual machine yet.",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceBlockVolumeAttachmentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
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

	// check volume status
	volume, err := client.IaaS().GetVolume(ctx, volumeID)
	if err != nil {
		if tcclient.IsNotFound(err) {
			return diag.FromErr(fmt.Errorf("volume not found: %w", err))
		}
		return diag.FromErr(fmt.Errorf("error getting volume: %s", err))
	}

	switch volume.Status {
	case "Detaching":
		// wait until detaching is complete
		timeout := d.Get("wait_for_detached_timeout").(int)
		ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Minute)
		defer cancel()
		for {
			select {
			case <-ctxWithTimeout.Done():
				return diag.FromErr(fmt.Errorf("timeout waiting for volume to be available"))
			default:
			}
			volume, err = client.IaaS().GetVolume(ctxWithTimeout, volumeID)
			if err != nil {
				return diag.FromErr(fmt.Errorf("error getting volume: %s", err))
			}
			if volume.Status == "Available" {
				break
			} else if volume.Status == "Attaching" {
				return diag.FromErr(fmt.Errorf("volume is already being attached to a different virtual machine"))
			} else if volume.Status == "Attached" {
				return diag.FromErr(fmt.Errorf("volume is already attached to a different virtual machine"))
			}
			time.Sleep(1 * time.Second)
		}
	case "Attached":
		// check if the volume is already attached to the VMI
		if volume.Attachments != nil {
			for _, att := range volume.Attachments {
				if att.AttachedToResourceType == ResourceVolumeAttachmentVirtualMachine && att.AttachedToIdentity == vmiID {
					d.SetId(att.Identity)
					return resourceBlockVolumeAttachmentRead(ctx, d, m)
				}
			}
		}
		return diag.FromErr(fmt.Errorf("volume is already attached to a different virtual machine"))
	}

	attachResponse, err := client.IaaS().AttachVolume(ctx, volumeID, attachRequest)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(attachResponse.Identity)

	if d.Get("wait_for_attached").(bool) {
		// wait until the volume is attached
		timeout := d.Get("wait_for_attached_timeout").(int)
		ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Minute)
		defer cancel()
		for {
			select {
			case <-ctxWithTimeout.Done():
				return diag.FromErr(fmt.Errorf("timeout waiting for volume to be attached"))
			default:
			}
			volume, err = client.IaaS().GetVolume(ctxWithTimeout, volumeID)
			if err != nil {
				return diag.FromErr(fmt.Errorf("error getting volume: %s", err))
			}
			if volume.Status == "Attached" {
				break
			}
			time.Sleep(1 * time.Second)
		}
	}

	return resourceBlockVolumeAttachmentRead(ctx, d, m)
}

func resourceBlockVolumeAttachmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
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

func resourceBlockVolumeAttachmentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	volumeID := d.Get("volume_id").(string)
	vmiID := d.Get("vmi_id").(string)

	volume, err := client.IaaS().GetVolume(ctx, volumeID)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error getting volume: %s", err))
	}

	// get the volume attachment
	vmi, err := client.IaaS().GetMachine(ctx, vmiID)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error getting machine instance: %s", err))
	}

	if volume.Status != "Attached" {
		return diag.FromErr(fmt.Errorf("volume is not attached to the virtual machine"))
	}

	// check if the volume is stil attached to the VMI
	found := false
	if vmi.VolumeAttachments != nil {
		for _, att := range vmi.VolumeAttachments {
			if att.VirtualMachine.Identity == vmi.Identity {
				found = true
				break
			}
		}
	}
	if !found {
		return diag.FromErr(fmt.Errorf("volume is not attached to the virtual machine"))
	}

	return nil
}

func resourceBlockVolumeAttachmentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
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
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if d.Get("wait_for_detached").(bool) {
		// wait until the volume is detached
		timeout := d.Get("wait_for_detached_timeout").(int)
		ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Minute)
		defer cancel()
		for {
			select {
			case <-ctxWithTimeout.Done():
				return diag.FromErr(fmt.Errorf("timeout waiting for volume to be detached"))
			default:
			}
			volume, err := client.IaaS().GetVolume(ctxWithTimeout, volumeID)
			if err != nil {
				if tcclient.IsNotFound(err) {
					d.SetId("")
					return nil
				}
				return diag.FromErr(fmt.Errorf("error getting volume: %s", err))
			}
			if volume.Status == "Available" {
				break
			}
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}
