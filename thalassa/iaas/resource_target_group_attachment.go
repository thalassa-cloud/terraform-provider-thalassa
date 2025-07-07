package iaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"

	iaas "github.com/thalassa-cloud/client-go/iaas"
)

func resourceTargetGroupAttachment() *schema.Resource {
	return &schema.Resource{
		Description:   "Attach a Virtual Machine Instance to a target group",
		CreateContext: resourceTargetGroupAttachmentCreate,
		ReadContext:   resourceTargetGroupAttachmentRead,
		DeleteContext: resourceTargetGroupAttachmentDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the Target Group Attachment. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"target_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the target group to attach to",
			},
			"vmi_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the Virtual Machine Instance to attach",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTargetGroupAttachmentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	targetGroupID := d.Get("target_group_id").(string)
	vmiID := d.Get("vmi_id").(string)

	// lets try and find the target group and vmi
	_, err = client.IaaS().GetTargetGroup(ctx, iaas.GetTargetGroupRequest{Identity: targetGroupID})
	if err != nil {
		if tcclient.IsNotFound(err) {
			return diag.FromErr(fmt.Errorf("target group %s was not found", targetGroupID))
		}
		return diag.FromErr(fmt.Errorf("error getting target group: %w", err))
	}
	_, err = client.IaaS().GetMachine(ctx, vmiID)
	if err != nil {
		if tcclient.IsNotFound(err) {
			return diag.FromErr(fmt.Errorf("virtual machine instance %s was not found", vmiID))
		}
		return diag.FromErr(fmt.Errorf("error getting virtual machine instance: %w", err))
	}

	// Create a single target
	attachmentRequest := iaas.AttachTargetGroupRequest{
		TargetGroupID: targetGroupID,
		AttachTarget: iaas.AttachTarget{
			ServerIdentity: vmiID,
		},
	}
	attachResponse, err := client.IaaS().AttachServerToTargetGroup(ctx, attachmentRequest)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting client: %w", err))
	}
	if attachResponse != nil {
		d.SetId(createTargetGroupAttachmentID(vmiID, targetGroupID))
		d.Set("target_group_id", targetGroupID)
		d.Set("vmi_id", vmiID)
	}
	return resourceTargetGroupAttachmentRead(ctx, d, m)
}

func createTargetGroupAttachmentID(vmiID, targetGroupID string) string {
	return fmt.Sprintf("%s:%s", vmiID, targetGroupID)
}

func resourceTargetGroupAttachmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting client: %w", err))
	}

	targetGroupID := d.Get("target_group_id").(string)
	vmiID := d.Get("vmi_id").(string)

	// Check if the VMI is attached to the target group
	found := false

	// Get the target group to check if the attachment exists
	tg, err := client.IaaS().GetTargetGroup(ctx, iaas.GetTargetGroupRequest{Identity: targetGroupID})
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error getting target group: %s", err))
	}

	if tg.LoadbalancerTargetGroupAttachments != nil {
		for _, att := range tg.LoadbalancerTargetGroupAttachments {
			if vmiID != "" && att.VirtualMachineInstance != nil && att.VirtualMachineInstance.Identity == vmiID {
				found = true
				d.Set("vmi_id", att.VirtualMachineInstance.Identity)
				break
			}
		}
	}

	if !found {
		d.SetId("")
		return nil
	}

	d.SetId(createTargetGroupAttachmentID(vmiID, targetGroupID))
	d.Set("target_group_id", targetGroupID)
	d.Set("vmi_id", vmiID)

	return nil
}

func resourceTargetGroupAttachmentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting client: %w", err))
	}

	targetGroupID := d.Get("target_group_id").(string)
	vmiID := d.Get("vmi_id").(string)
	if err := client.IaaS().DetachServerFromTargetGroup(ctx, iaas.DetachTargetRequest{
		TargetGroupID: targetGroupID,
		AttachmentID:  vmiID,
	}); err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error detaching server from target group: %w", err))
	}

	d.SetId("")
	return nil
}
