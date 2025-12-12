package iaas

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	iaas "github.com/thalassa-cloud/client-go/iaas"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func resourceSnapshot() *schema.Resource {
	return &schema.Resource{
		Description: `
		Provides a Thalassa Cloud Snapshot resource. This can be used to create, manage, and delete snapshots of block volumes.
		`,
		CreateContext: resourceSnapshotCreate,
		ReadContext:   resourceSnapshotRead,
		UpdateContext: resourceSnapshotUpdate,
		DeleteContext: resourceSnapshotDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the Snapshot. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.StringLenBetween(1, 62),
				Description:  "Name of the snapshot",
			},
			"slug": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.StringLenBetween(0, 255),
				Description:  "A human readable description about the snapshot",
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Labels for the snapshot",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Annotations for the snapshot",
			},
			"delete_protection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Delete protection of the snapshot",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the snapshot",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Region of the snapshot",
			},
			"source_volume_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity of the source volume",
			},
			"size_gb": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Size of the snapshot in GB",
			},
			"snapshot_policy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity of the snapshot policy that created this snapshot",
			},
			"wait_until_available": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Wait until the snapshot is available",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceSnapshotCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	volumeIdentity := d.Get("volume_identity").(string)

	// Verify the volume exists
	_, err = client.IaaS().GetVolume(ctx, volumeIdentity)
	if err != nil {
		if tcclient.IsNotFound(err) {
			return diag.FromErr(fmt.Errorf("volume not found: %s", volumeIdentity))
		}
		return diag.FromErr(fmt.Errorf("failed to get volume: %w", err))
	}

	createSnapshot := iaas.CreateSnapshotRequest{
		Name:             d.Get("name").(string),
		Description:      d.Get("description").(string),
		Labels:           convert.ConvertToMap(d.Get("labels")),
		Annotations:      convert.ConvertToMap(d.Get("annotations")),
		VolumeIdentity:   volumeIdentity,
		DeleteProtection: d.Get("delete_protection").(bool),
	}

	snapshot, err := client.IaaS().CreateSnapshot(ctx, createSnapshot)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create snapshot: %w", err))
	}

	if snapshot != nil {
		d.SetId(snapshot.Identity)
		d.Set("slug", snapshot.Slug)
		d.Set("status", string(snapshot.Status))
	}

	if d.Get("wait_until_available").(bool) {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Minute)
		defer cancel()

		err = client.IaaS().WaitUntilSnapshotIsAvailable(ctxWithTimeout, snapshot.Identity)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to wait for snapshot to be available: %w", err))
		}
	}

	return resourceSnapshotRead(ctx, d, m)
}

func resourceSnapshotRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Thalassa client: %w", err))
	}

	identity := d.Id()
	snapshot, err := client.IaaS().GetSnapshot(ctx, identity)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error getting snapshot: %w", err))
	}
	if snapshot == nil {
		d.SetId("")
		return nil
	}

	d.SetId(snapshot.Identity)
	d.Set("name", snapshot.Name)
	d.Set("slug", snapshot.Slug)
	d.Set("description", snapshot.Description)
	d.Set("labels", snapshot.Labels)
	d.Set("annotations", snapshot.Annotations)
	d.Set("status", string(snapshot.Status))
	d.Set("delete_protection", snapshot.DeleteProtection)

	if snapshot.Region != nil {
		d.Set("region", snapshot.Region.Identity)
	}

	if snapshot.SourceVolumeId != nil {
		d.Set("source_volume_id", *snapshot.SourceVolumeId)
	}

	if snapshot.SizeGB != nil {
		d.Set("size_gb", *snapshot.SizeGB)
	}

	if snapshot.SnapshotPolicyId != nil {
		d.Set("snapshot_policy_id", *snapshot.SnapshotPolicyId)
	}

	return nil
}

func resourceSnapshotUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Thalassa client: %w", err))
	}

	updateSnapshot := iaas.UpdateSnapshotRequest{
		Name:             d.Get("name").(string),
		Description:      d.Get("description").(string),
		Labels:           convert.ConvertToMap(d.Get("labels")),
		Annotations:      convert.ConvertToMap(d.Get("annotations")),
		DeleteProtection: d.Get("delete_protection").(bool),
	}

	identity := d.Id()
	snapshot, err := client.IaaS().UpdateSnapshot(ctx, identity, updateSnapshot)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update snapshot: %w", err))
	}

	if snapshot != nil {
		d.Set("name", snapshot.Name)
		d.Set("description", snapshot.Description)
		d.Set("slug", snapshot.Slug)
		d.Set("status", string(snapshot.Status))
		d.Set("labels", snapshot.Labels)
		d.Set("annotations", snapshot.Annotations)
		d.Set("delete_protection", snapshot.DeleteProtection)
		return nil
	}

	return resourceSnapshotRead(ctx, d, m)
}

func resourceSnapshotDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Thalassa client: %w", err))
	}

	identity := d.Id()

	err = client.IaaS().DeleteSnapshot(ctx, identity)
	if err != nil {
		if !tcclient.IsNotFound(err) {
			return diag.FromErr(fmt.Errorf("failed to delete snapshot: %w", err))
		}
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	err = client.IaaS().WaitUntilSnapshotIsDeleted(ctxWithTimeout, identity)
	if err != nil {
		if !strings.Contains(err.Error(), "timeout") {
			return diag.FromErr(fmt.Errorf("failed to wait for snapshot deletion: %w", err))
		}
	}

	d.SetId("")
	return nil
}
