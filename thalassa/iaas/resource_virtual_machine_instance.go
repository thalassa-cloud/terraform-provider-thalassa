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
				Default:      "",
				Optional:     true,
				ValidateFunc: validate.StringLenBetween(0, 255),
				Description:  "A human readable description about the virtual machine instance",
			},
			"labels": {
				Type:        schema.TypeMap,
				Default:     make(map[string]string),
				Optional:    true,
				Description: "Labels for the virtual machine instance",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Default:     make(map[string]string),
				Optional:    true,
				Description: "Annotations for the virtual machine instance",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
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
				Description: "Machine image for the virtual machine instance. You may pass the image identity, slug, or name (name match is case-insensitive)",
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
				ForceNew:    true,
				Description: "Cloud init template id of the virtual machine instance. If provided, the cloud init will be set to the content of the template.",
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
				Description: "List identities of security group that will be attached to the Virtual Machine Instance",
				Elem: &schema.Schema{
					Type:        schema.TypeString,
					Description: "The identity of the security group that will be attached to the Virtual Machine Instance",
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CustomizeDiff: func(ctx context.Context, diff *schema.ResourceDiff, meta any) error {
			// Get all values from the diff
			rootVolumeID := diff.Get("root_volume_id")
			rootVolumeSize := diff.Get("root_volume_size_gb")
			rootVolumeType := diff.Get("root_volume_type")

			tflog.Debug(ctx, "Validating root volume combination",
				map[string]any{
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

func resourceVirtualMachineInstanceCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
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
		_ = d.Set("root_volume_id", rootVolumeId.(string))
	}

	machineImageRef := d.Get("machine_image").(string)
	machineImageIdentity, err := lookupMachineImageIdentity(ctx, client.IaaS(), machineImageRef)
	if err != nil {
		return diag.FromErr(err)
	}

	createVirtualMachineInstance := iaas.CreateMachine{
		Name:                     d.Get("name").(string),
		Description:              d.Get("description").(string),
		Labels:                   convert.ConvertToMap(d.Get("labels")),
		Annotations:              convert.ConvertToMap(d.Get("annotations")),
		Subnet:                   d.Get("subnet_id").(string),
		MachineType:              d.Get("machine_type").(string),
		MachineImage:             machineImageIdentity,
		DeleteProtection:         d.Get("delete_protection").(bool),
		CloudInit:                d.Get("cloud_init").(string),
		RootVolume:               rootVolume,
		SecurityGroupAttachments: convert.ConvertToStringSlice(d.Get("security_group_attachments")),
	}

	if cloudInitTemplateId, ok := d.GetOk("cloud_init_template_id"); ok {
		cloudInitTemplate, err := client.IaaS().GetCloudInitTemplate(ctx, cloudInitTemplateId.(string))
		if err != nil {
			if tcclient.IsNotFound(err) {
				return diag.FromErr(fmt.Errorf("cloud init template not found: %w", err))
			}
			return diag.FromErr(fmt.Errorf("failed to get cloud init template: %w", err))
		}
		_ = d.Set("cloud_init_template_id", cloudInitTemplate.Identity)
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
		_ = d.Set("slug", virtualMachineInstance.Slug)

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
					_ = d.Set("ip_addresses", getIPAddresses(virtualMachineInstance))
					_ = d.Set("attached_volume_ids", getAttachedVolumeIds(virtualMachineInstance))
					_ = d.Set("status", virtualMachineInstance.Status.Status)
					_ = d.Set("state", virtualMachineInstance.State)

					securityGroupAttachments := make([]string, len(virtualMachineInstance.SecurityGroups))
					for i, securityGroup := range virtualMachineInstance.SecurityGroups {
						securityGroupAttachments[i] = securityGroup.Identity
					}
					_ = d.Set("security_group_attachments", securityGroupAttachments)

					if virtualMachineInstance.AvailabilityZone != nil {
						_ = d.Set("availability_zone", *virtualMachineInstance.AvailabilityZone)
					} else {
						_ = d.Set("availability_zone", "")
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

func resourceVirtualMachineInstanceRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Thalassa client: %w", err))
	}

	identity := d.Get("id").(string)
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
	_ = d.Set("name", virtualMachineInstance.Name)
	_ = d.Set("slug", virtualMachineInstance.Slug)
	_ = d.Set("description", virtualMachineInstance.Description)
	_ = d.Set("labels", virtualMachineInstance.Labels)
	_ = d.Set("annotations", virtualMachineInstance.Annotations)
	_ = d.Set("status", virtualMachineInstance.Status.Status)
	_ = d.Set("state", virtualMachineInstance.State)
	_ = d.Set("ip_addresses", getIPAddresses(virtualMachineInstance))
	_ = d.Set("attached_volume_ids", getAttachedVolumeIds(virtualMachineInstance))

	cloudInitTemplateId := d.Get("cloud_init_template_id").(string)
	_ = d.Set("cloud_init_template_id", cloudInitTemplateId)

	if virtualMachineInstance.MachineType != nil {
		setMachineTypeField(d, virtualMachineInstance.MachineType)
	}
	if virtualMachineInstance.MachineImage != nil {
		setMachineImageField(d, virtualMachineInstance.MachineImage)
	}

	_ = d.Set("subnet_id", virtualMachineInstance.Subnet.Identity)
	_ = d.Set("delete_protection", virtualMachineInstance.DeleteProtection)
	if virtualMachineInstance.CloudInit != nil && cloudInitTemplateId == "" {
		_ = d.Set("cloud_init", *virtualMachineInstance.CloudInit)
	} else {
		_ = d.Set("cloud_init", "")
	}

	securityGroupAttachments := make([]string, len(virtualMachineInstance.SecurityGroups))
	for i, securityGroup := range virtualMachineInstance.SecurityGroups {
		securityGroupAttachments[i] = securityGroup.Identity
	}
	_ = d.Set("security_group_attachments", securityGroupAttachments)

	if virtualMachineInstance.AvailabilityZone != nil {
		_ = d.Set("availability_zone", *virtualMachineInstance.AvailabilityZone)
	} else {
		_ = d.Set("availability_zone", "")
	}

	if virtualMachineInstance.PersistentVolume != nil {
		_ = d.Set("root_volume_size_gb", virtualMachineInstance.PersistentVolume.Size)
		_ = d.Set("root_volume_id", virtualMachineInstance.PersistentVolume.Identity)

		if virtualMachineInstance.PersistentVolume.VolumeType != nil {
			setRootVolumeTypeField(d, virtualMachineInstance.PersistentVolume.VolumeType)
		}
	}

	return nil
}

func resourceVirtualMachineInstanceUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Thalassa client: %w", err))
	}

	subnetId := d.Get("subnet_id").(string)

	currentMachine, err := client.IaaS().GetMachine(ctx, d.Get("id").(string))
	if err != nil {
		if tcclient.IsNotFound(err) {
			return diag.FromErr(fmt.Errorf("virtual machine instance not found: %w", err))
		}
		return diag.FromErr(fmt.Errorf("failed to get virtual machine instance: %w", err))
	}

	state := iaas.MachineState(d.Get("state").(string))
	var availabilityZone *string
	if d.Get("availability_zone").(string) != "" {
		availabilityZone = convert.Ptr(d.Get("availability_zone").(string))
	} else {
		availabilityZone = currentMachine.AvailabilityZone
	}
	machineType := d.Get("machine_type").(string)
	deleteProtection := d.Get("delete_protection").(bool)
	cloudInitTemplateId := d.Get("cloud_init_template_id").(string)

	updateVirtualMachineInstance := iaas.UpdateMachine{
		Name:                     d.Get("name").(string),
		Description:              d.Get("description").(string),
		Labels:                   convert.ConvertToMap(d.Get("labels")),
		Annotations:              convert.ConvertToMap(d.Get("annotations")),
		Subnet:                   &subnetId,
		State:                    &state,
		AvailabilityZone:         availabilityZone,
		MachineType:              &machineType,
		DeleteProtection:         &deleteProtection,
		SecurityGroupAttachments: convert.ConvertToStringSlice(d.Get("security_group_attachments")),
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
		_ = d.Set("name", virtualMachineInstance.Name)
		_ = d.Set("description", virtualMachineInstance.Description)
		_ = d.Set("slug", virtualMachineInstance.Slug)
		_ = d.Set("labels", virtualMachineInstance.Labels)
		_ = d.Set("annotations", virtualMachineInstance.Annotations)
		_ = d.Set("status", virtualMachineInstance.Status.Status)
		_ = d.Set("state", virtualMachineInstance.State)

		_ = d.Set("ip_addresses", getIPAddresses(virtualMachineInstance))
		_ = d.Set("attached_volume_ids", getAttachedVolumeIds(virtualMachineInstance))

		if virtualMachineInstance.MachineType != nil {
			setMachineTypeField(d, virtualMachineInstance.MachineType)
		}
		if virtualMachineInstance.MachineImage != nil {
			setMachineImageField(d, virtualMachineInstance.MachineImage)
		}

		_ = d.Set("subnet_id", virtualMachineInstance.Subnet.Identity)
		_ = d.Set("delete_protection", virtualMachineInstance.DeleteProtection)
		_ = d.Set("cloud_init", virtualMachineInstance.CloudInit)
		_ = d.Set("cloud_init_template_id", cloudInitTemplateId)

		securityGroupAttachments := make([]string, len(virtualMachineInstance.SecurityGroups))
		for i, securityGroup := range virtualMachineInstance.SecurityGroups {
			securityGroupAttachments[i] = securityGroup.Identity
		}
		_ = d.Set("security_group_attachments", securityGroupAttachments)

		if virtualMachineInstance.AvailabilityZone != nil {
			_ = d.Set("availability_zone", *virtualMachineInstance.AvailabilityZone)
		} else if currentMachine.AvailabilityZone != nil {
			_ = d.Set("availability_zone", *currentMachine.AvailabilityZone)
		} else {
			_ = d.Set("availability_zone", "")
		}

		if virtualMachineInstance.PersistentVolume != nil {
			_ = d.Set("root_volume_size_gb", virtualMachineInstance.PersistentVolume.Size)
			_ = d.Set("root_volume_id", virtualMachineInstance.PersistentVolume.Identity)
			if virtualMachineInstance.PersistentVolume.VolumeType != nil {
				setRootVolumeTypeField(d, virtualMachineInstance.PersistentVolume.VolumeType)
			}
		}

		return nil
	}

	return resourceVirtualMachineInstanceRead(ctx, d, m)
}

func resourceVirtualMachineInstanceDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
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

// setMachineTypeField keeps the user's reference (identity, slug, or name) when it still matches the API value.
func setMachineTypeField(d *schema.ResourceData, mt *iaas.MachineType) {
	if mt == nil {
		return
	}
	convert.SetReferenceField(d, "machine_type", mt.Identity, mt.Slug, mt.Name)
}

func setRootVolumeTypeField(d *schema.ResourceData, vt *iaas.VolumeType) {
	if vt == nil {
		return
	}
	convert.SetReferenceField(d, "root_volume_type", vt.Identity, "", vt.Name)
}

// setMachineImageField mirrors region handling on block_volume: the API returns identity, but Terraform
// may use identity, slug, or name; we keep the user's value when it still matches the resolved image.
func setMachineImageField(d *schema.ResourceData, mi *iaas.MachineImage) {
	if mi == nil {
		return
	}
	current := d.Get("machine_image").(string)
	switch {
	case current == "":
		_ = d.Set("machine_image", mi.Identity)
	case current == mi.Identity:
		_ = d.Set("machine_image", current)
	case current == mi.Slug:
		_ = d.Set("machine_image", current)
	case strings.EqualFold(current, mi.Name):
		_ = d.Set("machine_image", current)
	default:
		_ = d.Set("machine_image", mi.Identity)
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
	attachedVolumeIds := make([]string, 0, len(virtualMachineInstance.VolumeAttachments))
	for _, volumeAttachment := range virtualMachineInstance.VolumeAttachments {
		attachedVolumeIds = append(attachedVolumeIds, volumeAttachment.PersistentVolume.Identity)
	}
	return attachedVolumeIds
}
