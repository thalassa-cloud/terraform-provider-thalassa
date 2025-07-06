package iaas

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"

	iaas "github.com/thalassa-cloud/client-go/iaas"
)

func resourceVirtualMachineInstance() *schema.Resource {
	return &schema.Resource{
		Description:   "Create an virtual machine instance within a subnet on the Thalassa Cloud platform",
		CreateContext: resourceVirtualMachineInstanceCreate,
		ReadContext:   resourceVirtualMachineInstanceRead,
		UpdateContext: resourceVirtualMachineInstanceUpdate,
		DeleteContext: resourceVirtualMachineInstanceDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Subnet of the Virtual Machine Instance",
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the Machine Type. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.StringLenBetween(1, 62),
				ForceNew:     true,
				Description:  "Name of the Virtual Machine Instance",
			},
			"slug": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Slug of the Virtual Machine Instance",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.StringLenBetween(0, 255),
				Description:  "A human readable description about the virtual machine instance",
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Labels for the virtual machine instance",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Annotations for the virtual machine instance",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Availability zone of the virtual machine instance",
			},
			"machine_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Machine type of the virtual machine instance",
			},
			"machine_image": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Machine image of the virtual machine instance",
			},
			"delete_protection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Delete protection of the virtual machine instance",
			},
			"cloud_init": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Cloud init of the virtual machine instance",
			},
			"cloud_init_template_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Cloud init template id of the virtual machine instance",
			},
			"root_volume_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Root volume id of the virtual machine instance. Must be provided if root_volume_type is not set.",
			},
			"root_volume_size_gb": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Root volume size of the virtual machine instance. Must be provided if root_volume_id is not set.",
			},
			"root_volume_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Root volume type of the virtual machine instance. Must be provided if root_volume_id is not set.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the virtual machine instance",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Desired state of the virtual machine instance. Can be 'running', 'stopped', 'deleted'",
			},
			"ip_addresses": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "IP addresses of the virtual machine instance",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"attached_volume_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Attached volume ids of the virtual machine instance",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"security_group_attachments": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Security group attached to the virtual machine instance",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CustomizeDiff: func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
			// Get all values from the diff
			rootVolumeID := diff.Get("root_volume_id")
			rootVolumeSize := diff.Get("root_volume_size_gb")
			rootVolumeType := diff.Get("root_volume_type")

			tflog.Debug(ctx, "Validating root volume combination",
				map[string]interface{}{
					"root_volume_id":      rootVolumeID,
					"root_volume_size_gb": rootVolumeSize,
					"root_volume_type":    rootVolumeType,
				})

			// If root_volume_id is set, we're good
			if rootVolumeID != nil && rootVolumeID.(string) != "" {
				return nil
			}

			// If root_volume_id is not set, both root_volume_size_gb and root_volume_type must be set
			if rootVolumeSize == nil || rootVolumeType == nil {
				return fmt.Errorf("either root_volume_id must be provided, or both root_volume_size_gb and root_volume_type must be provided")
			}

			// Additional validation for root_volume_size_gb if needed
			if rootVolumeSize.(int) <= 0 {
				return fmt.Errorf("root_volume_size_gb must be greater than 0")
			}

			return nil
		},
	}
}

func resourceVirtualMachineInstanceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Thalassa client: %w", err))
	}

	rootVolume := iaas.CreateMachineVolume{
		Name: convert.Ptr(d.Get("name").(string)),
	}

	if rootVolumeSize, ok := d.GetOk("root_volume_size_gb"); ok {
		rootVolume.Size = rootVolumeSize.(int)
	}

	if rootVolumeType, ok := d.GetOk("root_volume_type"); ok {
		volumeType, err := lookupVolumeType(ctx, client.IaaS(), rootVolumeType.(string))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to lookup volume type: %w", err))
		}
		rootVolume.VolumeTypeIdentity = volumeType.Identity
	}

	if rootVolumeId, ok := d.GetOk("root_volume_id"); ok {
		d.Set("root_volume_id", rootVolumeId.(string))
	}

	createVirtualMachineInstance := iaas.CreateMachine{
		Name:                     d.Get("name").(string),
		Description:              d.Get("description").(string),
		Labels:                   convert.ConvertToMap(d.Get("labels")),
		Annotations:              convert.ConvertToMap(d.Get("annotations")),
		Subnet:                   d.Get("subnet_id").(string),
		MachineType:              d.Get("machine_type").(string),
		MachineImage:             d.Get("machine_image").(string),
		DeleteProtection:         d.Get("delete_protection").(bool),
		CloudInit:                d.Get("cloud_init").(string),
		RootVolume:               rootVolume,
		SecurityGroupAttachments: convertToStrList(d.Get("security_group_attachments")),
	}

	if cloudInitTemplateId, ok := d.GetOk("cloud_init_template_id"); ok {
		cloudInitTemplate, err := client.IaaS().GetCloudInitTemplate(ctx, cloudInitTemplateId.(string))
		if err != nil {
			if tcclient.IsNotFound(err) {
				return diag.FromErr(fmt.Errorf("cloud init template not found: %w", err))
			}
			return diag.FromErr(fmt.Errorf("failed to get cloud init template: %w", err))
		}
		d.Set("cloud_init_template_id", cloudInitTemplate.Identity)
		createVirtualMachineInstance.CloudInit = cloudInitTemplate.Content
	}

	if availabilityZone, ok := d.GetOk("availability_zone"); ok {
		createVirtualMachineInstance.AvailabilityZone = convert.Ptr(availabilityZone.(string))
	}

	virtualMachineInstance, err := client.IaaS().CreateMachine(ctx, createVirtualMachineInstance)

	if err != nil {
		if tcclient.IsNotFound(err) {
			return diag.FromErr(fmt.Errorf("used resource for creating virtual machine instance not found: %w", err))
		}
		return diag.FromErr(fmt.Errorf("failed to create virtual machine instance: %w", err))
	}
	if virtualMachineInstance != nil {
		identity := virtualMachineInstance.Identity
		d.SetId(identity)
		d.Set("slug", virtualMachineInstance.Slug)

		// wait until the virtual machine instance is ready
		ctxWithTimeout, cancel := context.WithTimeout(ctx, 20*time.Minute)
		defer cancel()
		for {
			select {
			case <-ctxWithTimeout.Done():
				return diag.FromErr(fmt.Errorf("timeout while waiting for virtual machine instance to be ready"))
			case <-time.After(1 * time.Second):
				// continue
				virtualMachineInstance, err = client.IaaS().GetMachine(ctxWithTimeout, identity)
				if err != nil {
					if tcclient.IsNotFound(err) {
						return diag.FromErr(fmt.Errorf("virtual machine instance %s was not found after creation", identity))
					}
					return diag.FromErr(err)
				}
				if virtualMachineInstance == nil {
					return diag.FromErr(fmt.Errorf("virtual machine instance %s was not found after creation", identity))
				}

				if strings.EqualFold(virtualMachineInstance.Status.Status, "ready") || strings.EqualFold(virtualMachineInstance.Status.Status, "running") {
					d.Set("ip_addresses", getIPAddresses(virtualMachineInstance))
					d.Set("attached_volume_ids", getAttachedVolumeIds(virtualMachineInstance))
					d.Set("status", virtualMachineInstance.Status.Status)
					d.Set("state", virtualMachineInstance.State)
					d.Set("security_group_attachments", virtualMachineInstance.SecurityGroupAttachments)
					if virtualMachineInstance.AvailabilityZone != nil {
						d.Set("availability_zone", *virtualMachineInstance.AvailabilityZone)
					} else {
						d.Set("availability_zone", "")
					}
					return nil
				} else if strings.EqualFold(virtualMachineInstance.Status.Status, "error") || strings.EqualFold(virtualMachineInstance.Status.Status, "failed") {
					return diag.FromErr(fmt.Errorf("virtual machine instance is in error state: %s", virtualMachineInstance.Status.StatusMessage))
				}
			}
		}
	}
	return resourceVirtualMachineInstanceRead(ctx, d, m)
}

func resourceVirtualMachineInstanceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Thalassa client: %w", err))
	}

	identity := d.Get("id").(string)
	cloudInitTemplateId := d.Get("cloud_init_template_id").(string)
	virtualMachineInstance, err := client.IaaS().GetMachine(ctx, identity)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error getting virtual machine instance: %s", err))
	}
	if virtualMachineInstance == nil {
		d.SetId("")
		return nil
	}

	d.SetId(virtualMachineInstance.Identity)
	d.Set("name", virtualMachineInstance.Name)
	d.Set("slug", virtualMachineInstance.Slug)
	d.Set("description", virtualMachineInstance.Description)
	d.Set("labels", virtualMachineInstance.Labels)
	d.Set("annotations", virtualMachineInstance.Annotations)
	d.Set("status", virtualMachineInstance.Status.Status)
	d.Set("state", virtualMachineInstance.State)
	d.Set("ip_addresses", getIPAddresses(virtualMachineInstance))
	d.Set("attached_volume_ids", getAttachedVolumeIds(virtualMachineInstance))
	d.Set("cloud_init_template_id", cloudInitTemplateId)

	if d.Get("machine_type").(string) != "" {
		d.Set("machine_type", d.Get("machine_type").(string))
	} else {
		d.Set("machine_type", virtualMachineInstance.MachineType.Identity)
	}
	if d.Get("machine_image").(string) != "" {
		d.Set("machine_image", d.Get("machine_image").(string))
	} else {
		d.Set("machine_image", virtualMachineInstance.MachineImage.Identity)
	}

	d.Set("subnet_id", virtualMachineInstance.Subnet.Identity)
	d.Set("delete_protection", virtualMachineInstance.DeleteProtection)
	d.Set("cloud_init", virtualMachineInstance.CloudInit)
	d.Set("security_group_attachments", virtualMachineInstance.SecurityGroupAttachments)
	if virtualMachineInstance.AvailabilityZone != nil {
		d.Set("availability_zone", *virtualMachineInstance.AvailabilityZone)
	} else {
		d.Set("availability_zone", "")
	}

	if virtualMachineInstance.PersistentVolume != nil {
		d.Set("root_volume_size_gb", virtualMachineInstance.PersistentVolume.Size)
		d.Set("root_volume_id", virtualMachineInstance.PersistentVolume.Identity)

		if virtualMachineInstance.PersistentVolume.VolumeType != nil {
			if d.Get("root_volume_type").(string) != "" {
				d.Set("root_volume_type", d.Get("root_volume_type").(string))
			} else {
				d.Set("root_volume_type", virtualMachineInstance.PersistentVolume.VolumeType.Identity)
			}
		}
	}

	return nil
}

func resourceVirtualMachineInstanceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Thalassa client: %w", err))
	}

	subnetId := d.Get("subnet_id").(string)

	state := iaas.MachineState(d.Get("state").(string))
	availabilityZone := d.Get("availability_zone").(string)
	machineType := d.Get("machine_type").(string)
	deleteProtection := d.Get("delete_protection").(bool)
	cloudInitTemplateId := d.Get("cloud_init_template_id").(string)

	updateVirtualMachineInstance := iaas.UpdateMachine{
		Name:             d.Get("name").(string),
		Description:      d.Get("description").(string),
		Labels:           convert.ConvertToMap(d.Get("labels")),
		Annotations:      convert.ConvertToMap(d.Get("annotations")),
		Subnet:           &subnetId,
		State:            &state,
		AvailabilityZone: &availabilityZone,
		MachineType:      &machineType,
		DeleteProtection: &deleteProtection,
		// SecurityGroupAttachments: securityGroupAttachments,
	}

	identity := d.Get("id").(string)

	virtualMachineInstance, err := client.IaaS().UpdateMachine(ctx, identity, updateVirtualMachineInstance)
	if err != nil {
		if tcclient.IsNotFound(err) {
			return diag.FromErr(fmt.Errorf("used resource for updating virtual machine instance not found: %w", err))
		}
		return diag.FromErr(fmt.Errorf("failed to update virtual machine instance: %w", err))
	}
	if virtualMachineInstance != nil {
		d.Set("name", virtualMachineInstance.Name)
		d.Set("description", virtualMachineInstance.Description)
		d.Set("slug", virtualMachineInstance.Slug)
		d.Set("labels", virtualMachineInstance.Labels)
		d.Set("annotations", virtualMachineInstance.Annotations)
		d.Set("status", virtualMachineInstance.Status.Status)
		d.Set("state", virtualMachineInstance.State)

		d.Set("ip_addresses", getIPAddresses(virtualMachineInstance))
		d.Set("attached_volume_ids", getAttachedVolumeIds(virtualMachineInstance))

		if d.Get("machine_type").(string) != "" {
			if virtualMachineInstance.MachineType != nil {
				// check if the machine type is the same as the one in the diff
				stateReference := d.Get("machine_type").(string)
				switch stateReference {
				case virtualMachineInstance.MachineType.Identity:
					d.Set("machine_type", stateReference)
				case virtualMachineInstance.MachineType.Slug:
					d.Set("machine_type", stateReference)
				default:
					d.Set("machine_type", virtualMachineInstance.MachineType.Identity)
				}
			}
		} else {
			d.Set("machine_type", virtualMachineInstance.MachineType.Identity)
		}

		if d.Get("machine_image").(string) != "" {
			if virtualMachineInstance.MachineImage != nil {
				// check if the machine image is the same as the one in the diff
				stateReference := d.Get("machine_image").(string)
				switch stateReference {
				case virtualMachineInstance.MachineImage.Identity:
					d.Set("machine_image", stateReference)
				case virtualMachineInstance.MachineImage.Slug:
					d.Set("machine_image", stateReference)
				default:
					d.Set("machine_image", virtualMachineInstance.MachineImage.Identity)
				}
			}
		} else {
			d.Set("machine_image", virtualMachineInstance.MachineImage.Identity)
		}

		d.Set("subnet_id", virtualMachineInstance.Subnet.Identity)
		d.Set("delete_protection", virtualMachineInstance.DeleteProtection)
		d.Set("cloud_init", virtualMachineInstance.CloudInit)
		d.Set("cloud_init_template_id", cloudInitTemplateId)
		d.Set("security_group_attachments", virtualMachineInstance.SecurityGroupAttachments)
		if virtualMachineInstance.AvailabilityZone != nil {
			d.Set("availability_zone", *virtualMachineInstance.AvailabilityZone)
		} else {
			d.Set("availability_zone", "")
		}

		if virtualMachineInstance.PersistentVolume != nil {
			d.Set("root_volume_size_gb", virtualMachineInstance.PersistentVolume.Size)
			d.Set("root_volume_id", virtualMachineInstance.PersistentVolume.Identity)
			if virtualMachineInstance.PersistentVolume.VolumeType != nil {
				if d.Get("root_volume_type").(string) != "" {
					d.Set("root_volume_type", d.Get("root_volume_type").(string))
				} else {
					d.Set("root_volume_type", virtualMachineInstance.PersistentVolume.VolumeType.Identity)
				}
			}
		}

		return nil
	}

	return resourceVirtualMachineInstanceRead(ctx, d, m)
}

func resourceVirtualMachineInstanceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Thalassa client: %w", err))
	}

	id := d.Get("id").(string)

	err = client.IaaS().DeleteMachine(ctx, id)
	if err != nil && !tcclient.IsNotFound(err) {
		return diag.FromErr(fmt.Errorf("failed to delete virtual machine instance: %w", err))
	}

	// wait until the virtual machine instance is deleted
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 20*time.Minute)
	defer cancel()
	for {
		select {
		case <-ctxWithTimeout.Done():
			return diag.FromErr(fmt.Errorf("timeout while waiting for virtual machine instance to be deleted"))
		case <-time.After(2 * time.Second):
			m, err := client.IaaS().GetMachine(ctxWithTimeout, id)
			if err != nil {
				if tcclient.IsNotFound(err) {
					d.SetId("")
					return nil
				}
				return diag.FromErr(fmt.Errorf("error getting virtual machine instance: %s", err))
			}

			if strings.EqualFold(m.Status.Status, "deleted") {
				d.SetId("")
				return nil
			}
		}
	}
}

func getIPAddresses(virtualMachineInstance *iaas.Machine) []string {
	ipAddresses := []string{}
	for _, interf := range virtualMachineInstance.Interfaces {
		ipAddresses = append(ipAddresses, interf.IPAddresses...)
	}
	return ipAddresses
}

func getAttachedVolumeIds(virtualMachineInstance *iaas.Machine) []string {
	attachedVolumeIds := []string{}
	for _, volumeAttachment := range virtualMachineInstance.VolumeAttachments {
		attachedVolumeIds = append(attachedVolumeIds, volumeAttachment.PersistentVolume.Identity)
	}
	return attachedVolumeIds
}

func getInterfaces(virtualMachineInstance *iaas.Machine) []map[string]interface{} {
	interfaces := []map[string]interface{}{}
	for _, interf := range virtualMachineInstance.Interfaces {
		interfaces = append(interfaces, map[string]interface{}{
			"name":         interf.Name,
			"mac_address":  interf.MacAddress,
			"ip_addresses": interf.IPAddresses,
		})
	}
	return interfaces
}

func getVolumeAttachments(virtualMachineInstance *iaas.Machine) []map[string]interface{} {
	volumeAttachments := []map[string]interface{}{}
	for _, volumeAttachment := range virtualMachineInstance.VolumeAttachments {
		v := map[string]interface{}{
			"serial": volumeAttachment.Serial,
		}
		if volumeAttachment.PersistentVolume != nil {
			v["size_gb"] = volumeAttachment.PersistentVolume.Size
			v["volume_type"] = volumeAttachment.PersistentVolume.VolumeType.Identity
			v["volume_id"] = volumeAttachment.PersistentVolume.Identity
		}
		volumeAttachments = append(volumeAttachments, v)
	}
	return volumeAttachments
}

func convertToStrList(v interface{}) []string {
	if v == nil {
		return []string{}
	}
	values := []string{}
	for _, v := range v.([]interface{}) {
		values = append(values, v.(string))
	}
	return values
}
