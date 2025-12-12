package iaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	iaas "github.com/thalassa-cloud/client-go/iaas"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func DataSourceSnapshot() *schema.Resource {
	return &schema.Resource{
		Description: "Get a snapshot",
		ReadContext: dataSourceSnapshotRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity of the snapshot",
			},
			"identity": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Identity of the snapshot",
				ExactlyOneOf: []string{"identity", "name"},
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Name of the snapshot",
				ExactlyOneOf: []string{"identity", "name"},
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the Snapshot. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"slug": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Slug of the snapshot",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the snapshot",
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
			"delete_protection": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Delete protection of the snapshot",
			},
			"snapshot_policy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity of the snapshot policy that created this snapshot",
			},
			"labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Labels for the snapshot",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Annotations for the snapshot",
			},
		},
	}
}

func dataSourceSnapshotRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	var snapshot *iaas.Snapshot

	if identity, ok := d.GetOk("identity"); ok {
		// Look up by identity
		snapshot, err = client.IaaS().GetSnapshot(ctx, identity.(string))
		if err != nil {
			if tcclient.IsNotFound(err) {
				return diag.FromErr(fmt.Errorf("snapshot not found: %s", identity.(string)))
			}
			return diag.FromErr(err)
		}
	} else if name, ok := d.GetOk("name"); ok {
		// Look up by name
		snapshots, err := client.IaaS().ListSnapshots(ctx, &iaas.ListSnapshotsRequest{})
		if err != nil {
			return diag.FromErr(err)
		}

		// Find the snapshot with the matching name
		for _, s := range snapshots {
			if s.Name == name.(string) {
				snapshot = &s
				break
			}
		}

		if snapshot == nil {
			return diag.FromErr(fmt.Errorf("snapshot with name %s not found", name.(string)))
		}
	}

	// Set the ID and other attributes
	d.SetId(snapshot.Identity)
	d.Set("id", snapshot.Identity)
	d.Set("name", snapshot.Name)
	d.Set("slug", snapshot.Slug)
	d.Set("description", snapshot.Description)
	d.Set("status", string(snapshot.Status))
	d.Set("delete_protection", snapshot.DeleteProtection)
	d.Set("labels", snapshot.Labels)
	d.Set("annotations", snapshot.Annotations)

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
