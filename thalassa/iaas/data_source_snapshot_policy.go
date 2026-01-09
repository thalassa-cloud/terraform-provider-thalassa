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

func DataSourceSnapshotPolicy() *schema.Resource {
	return &schema.Resource{
		Description: "Get a snapshot policy",
		ReadContext: dataSourceSnapshotPolicyRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity of the snapshot policy",
			},
			"identity": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Identity of the snapshot policy",
				ExactlyOneOf: []string{"identity", "name"},
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Name of the snapshot policy",
				ExactlyOneOf: []string{"identity", "name"},
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the Snapshot Policy. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"slug": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Slug of the snapshot policy",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the snapshot policy",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Region of the snapshot policy",
			},
			"ttl": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Time to live for snapshots created by this policy (e.g., 168h, 7d)",
			},
			"keep_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Maximum number of snapshots to retain",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the snapshot policy is enabled",
			},
			"schedule": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cron schedule for the snapshot policy",
			},
			"timezone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timezone for the snapshot policy schedule",
			},
			"target": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Target configuration for the snapshot policy",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of target (selector or explicit)",
						},
						"selector": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Label selector for volumes (when type is selector)",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"volume_identities": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of volume identities (when type is explicit)",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"next_snapshot_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When the next snapshot will be created",
			},
			"last_snapshot_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When the last snapshot was created",
			},
			"labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Labels for the snapshot policy",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Annotations for the snapshot policy",
			},
		},
	}
}

func dataSourceSnapshotPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	var policy *iaas.SnapshotPolicy

	if identity, ok := d.GetOk("identity"); ok {
		// Look up by identity
		policy, err = client.IaaS().GetSnapshotPolicy(ctx, identity.(string))
		if err != nil {
			if tcclient.IsNotFound(err) {
				return diag.FromErr(fmt.Errorf("snapshot policy not found: %s", identity.(string)))
			}
			return diag.FromErr(err)
		}
	} else if name, ok := d.GetOk("name"); ok {
		// Look up by name
		policies, err := client.IaaS().ListSnapshotPolicies(ctx, &iaas.ListSnapshotPoliciesRequest{})
		if err != nil {
			return diag.FromErr(err)
		}

		// Find the policy with the matching name
		for _, p := range policies {
			if p.Name == name.(string) {
				policy = &p
				break
			}
		}

		if policy == nil {
			return diag.FromErr(fmt.Errorf("snapshot policy with name %s not found", name.(string)))
		}
	}

	// Set the ID and other attributes
	d.SetId(policy.Identity)
	d.Set("id", policy.Identity)
	d.Set("name", policy.Name)
	d.Set("slug", policy.Slug)
	d.Set("description", policy.Description)
	d.Set("enabled", policy.Enabled)
	d.Set("schedule", policy.Schedule)
	d.Set("timezone", policy.Timezone)
	d.Set("labels", policy.Labels)
	d.Set("annotations", policy.Annotations)

	// Convert TTL duration to string
	d.Set("ttl", policy.Ttl.String())

	if policy.KeepCount != nil {
		d.Set("keep_count", *policy.KeepCount)
	}

	if policy.Region != nil {
		d.Set("region", policy.Region.Identity)
	}

	if policy.NextSnapshotAt != nil {
		d.Set("next_snapshot_at", policy.NextSnapshotAt.Format("2006-01-02T15:04:05Z07:00"))
	}

	if policy.LastSnapshotAt != nil {
		d.Set("last_snapshot_at", policy.LastSnapshotAt.Format("2006-01-02T15:04:05Z07:00"))
	}

	// Set target
	target := map[string]interface{}{
		"type": string(policy.Target.Type),
	}

	switch policy.Target.Type {
	case iaas.SnapshotPolicyTargetTypeSelector:
		target["selector"] = policy.Target.Selector
		target["volume_identities"] = []string{}
	case iaas.SnapshotPolicyTargetTypeExplicit:
		target["selector"] = map[string]string{}
		target["volume_identities"] = policy.Target.VolumeIdentities
	}

	d.Set("target", []interface{}{target})

	return nil
}
