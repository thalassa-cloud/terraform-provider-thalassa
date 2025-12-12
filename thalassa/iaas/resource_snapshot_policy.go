package iaas

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
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

func resourceSnapshotPolicy() *schema.Resource {
	return &schema.Resource{
		Description: `
		Provides a Thalassa Cloud Snapshot Policy resource. This can be used to create, manage, and delete snapshot policies that automatically create snapshots of volumes based on a schedule.
		`,
		CreateContext: resourceSnapshotPolicyCreate,
		ReadContext:   resourceSnapshotPolicyRead,
		UpdateContext: resourceSnapshotPolicyUpdate,
		DeleteContext: resourceSnapshotPolicyDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the Snapshot Policy. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.StringLenBetween(1, 62),
				Description:  "Name of the snapshot policy",
			},
			"slug": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.StringLenBetween(0, 255),
				Description:  "A human readable description about the snapshot policy",
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Labels for the snapshot policy",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Annotations for the snapshot policy",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Region of the snapshot policy",
			},
			"ttl": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Time to live for snapshots created by this policy. Supports formats like '168h' (hours), '7d' (days), '1w' (weeks). Examples: '24h', '7d', '30d'",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					ttlStr := val.(string)
					_, err := parseDuration(ttlStr)
					if err != nil {
						errs = append(errs, fmt.Errorf("ttl must be a valid duration (e.g., '24h', '7d', '1w'): %w", err))
					}
					return
				},
			},
			"keep_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Maximum number of snapshots to retain. When this limit is reached, the oldest snapshots will be deleted.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the snapshot policy is enabled",
			},
			"schedule": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cron schedule for the snapshot policy (e.g., '0 2 * * *' for daily at 2 AM)",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					schedule := val.(string)
					if !regexp.MustCompile(`^[0-9,\-\*]+ [0-9,\-\*]+ [0-9,\-\*]+ [0-9,\-\*]+ [0-9,\-\*]+$`).MatchString(schedule) {
						errs = append(errs, fmt.Errorf("schedule must be in valid cron format (e.g., '0 2 * * *')"))
					}
					return
				},
			},
			"timezone": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Timezone for the snapshot policy schedule (e.g., 'UTC', 'America/New_York')",
			},
			"target": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Target configuration for the snapshot policy",
				MaxItems:    1,
				MinItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "Type of target: 'selector' to target volumes based on labels, or 'explicit' to target specific volumes",
							ValidateFunc: validate.StringInSlice([]string{"selector", "explicit"}, false),
						},
						"selector": {
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "Label selector for volumes (required when type is 'selector')",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"volume_identities": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of volume identities (required when type is 'explicit')",
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
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// parseDuration parses a duration string that supports 'd' for days and 'w' for weeks
// in addition to standard Go duration formats
func parseDuration(s string) (time.Duration, error) {
	// Handle days (d)
	if strings.HasSuffix(s, "d") {
		daysStr := strings.TrimSuffix(s, "d")
		days, err := strconv.Atoi(daysStr)
		if err != nil {
			return 0, fmt.Errorf("invalid days format: %w", err)
		}
		return time.Duration(days) * 24 * time.Hour, nil
	}

	// Handle weeks (w)
	if strings.HasSuffix(s, "w") {
		weeksStr := strings.TrimSuffix(s, "w")
		weeks, err := strconv.Atoi(weeksStr)
		if err != nil {
			return 0, fmt.Errorf("invalid weeks format: %w", err)
		}
		return time.Duration(weeks) * 7 * 24 * time.Hour, nil
	}

	// Use standard Go duration parsing for other formats (h, m, s, etc.)
	return time.ParseDuration(s)
}

func resourceSnapshotPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	region := d.Get("region").(string)
	regions, err := client.IaaS().ListRegions(ctx, &iaas.ListRegionsRequest{})
	if err != nil {
		if tcclient.IsNotFound(err) {
			return diag.FromErr(fmt.Errorf("region not found: %w", err))
		}
		return diag.FromErr(fmt.Errorf("failed to find region: %w", err))
	}
	foundRegion := false
	for _, r := range regions {
		if r.Identity == region || r.Slug == region || r.Name == region {
			region = r.Identity
			foundRegion = true
			break
		}
	}
	if !foundRegion {
		availableRegions := make([]string, len(regions))
		for i, r := range regions {
			availableRegions[i] = r.Slug
		}
		return diag.FromErr(fmt.Errorf("region not found: %s. Available regions: %v", region, strings.Join(availableRegions, ", ")))
	}

	// Parse TTL
	ttlStr := d.Get("ttl").(string)
	ttl, err := parseDuration(ttlStr)
	if err != nil {
		return diag.FromErr(fmt.Errorf("invalid ttl format: %w", err))
	}

	// Parse target
	targetList := d.Get("target").([]interface{})
	if len(targetList) != 1 {
		return diag.FromErr(fmt.Errorf("target must have exactly one element"))
	}
	targetMap := targetList[0].(map[string]interface{})
	targetType := targetMap["type"].(string)

	var target iaas.SnapshotPolicyTarget
	if targetType == "selector" {
		selector := convert.ConvertToMap(targetMap["selector"])
		if len(selector) == 0 {
			return diag.FromErr(fmt.Errorf("selector is required when target type is 'selector'"))
		}
		target = iaas.SnapshotPolicyTarget{
			Type:     iaas.SnapshotPolicyTargetTypeSelector,
			Selector: selector,
		}
	} else if targetType == "explicit" {
		volumeIdentities := convert.ConvertToStringSlice(targetMap["volume_identities"])
		if len(volumeIdentities) == 0 {
			return diag.FromErr(fmt.Errorf("volume_identities is required when target type is 'explicit'"))
		}
		target = iaas.SnapshotPolicyTarget{
			Type:             iaas.SnapshotPolicyTargetTypeExplicit,
			VolumeIdentities: volumeIdentities,
		}
	}

	keepCount := (*int)(nil)
	if v, ok := d.GetOk("keep_count"); ok {
		kc := v.(int)
		keepCount = &kc
	}

	createPolicy := iaas.CreateSnapshotPolicyRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Labels:      convert.ConvertToMap(d.Get("labels")),
		Annotations: convert.ConvertToMap(d.Get("annotations")),
		Region:      region,
		Ttl:         ttl,
		KeepCount:   keepCount,
		Enabled:     d.Get("enabled").(bool),
		Schedule:    d.Get("schedule").(string),
		Timezone:    d.Get("timezone").(string),
		Target:      target,
	}

	policy, err := client.IaaS().CreateSnapshotPolicy(ctx, createPolicy)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create snapshot policy: %w", err))
	}

	if policy != nil {
		d.SetId(policy.Identity)
		d.Set("slug", policy.Slug)
	}

	return resourceSnapshotPolicyRead(ctx, d, m)
}

func resourceSnapshotPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Thalassa client: %w", err))
	}

	identity := d.Id()
	policy, err := client.IaaS().GetSnapshotPolicy(ctx, identity)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error getting snapshot policy: %w", err))
	}
	if policy == nil {
		d.SetId("")
		return nil
	}

	d.SetId(policy.Identity)
	d.Set("name", policy.Name)
	d.Set("slug", policy.Slug)
	d.Set("description", policy.Description)
	d.Set("labels", policy.Labels)
	d.Set("annotations", policy.Annotations)
	d.Set("enabled", policy.Enabled)
	d.Set("schedule", policy.Schedule)
	d.Set("timezone", policy.Timezone)

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

	if policy.Target.Type == iaas.SnapshotPolicyTargetTypeSelector {
		target["selector"] = policy.Target.Selector
		target["volume_identities"] = []string{}
	} else if policy.Target.Type == iaas.SnapshotPolicyTargetTypeExplicit {
		target["selector"] = map[string]string{}
		target["volume_identities"] = policy.Target.VolumeIdentities
	}

	d.Set("target", []interface{}{target})

	return nil
}

func resourceSnapshotPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Thalassa client: %w", err))
	}

	// Parse TTL
	ttlStr := d.Get("ttl").(string)
	ttl, err := parseDuration(ttlStr)
	if err != nil {
		return diag.FromErr(fmt.Errorf("invalid ttl format: %w", err))
	}

	// Parse target
	targetList := d.Get("target").([]interface{})
	if len(targetList) != 1 {
		return diag.FromErr(fmt.Errorf("target must have exactly one element"))
	}
	targetMap := targetList[0].(map[string]interface{})
	targetType := targetMap["type"].(string)

	var target iaas.SnapshotPolicyTarget
	if targetType == "selector" {
		selector := convert.ConvertToMap(targetMap["selector"])
		if len(selector) == 0 {
			return diag.FromErr(fmt.Errorf("selector is required when target type is 'selector'"))
		}
		target = iaas.SnapshotPolicyTarget{
			Type:     iaas.SnapshotPolicyTargetTypeSelector,
			Selector: selector,
		}
	} else if targetType == "explicit" {
		volumeIdentities := convert.ConvertToStringSlice(targetMap["volume_identities"])
		if len(volumeIdentities) == 0 {
			return diag.FromErr(fmt.Errorf("volume_identities is required when target type is 'explicit'"))
		}
		target = iaas.SnapshotPolicyTarget{
			Type:             iaas.SnapshotPolicyTargetTypeExplicit,
			VolumeIdentities: volumeIdentities,
		}
	}

	keepCount := (*int)(nil)
	if v, ok := d.GetOk("keep_count"); ok {
		kc := v.(int)
		keepCount = &kc
	}

	updatePolicy := iaas.UpdateSnapshotPolicyRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Labels:      convert.ConvertToMap(d.Get("labels")),
		Annotations: convert.ConvertToMap(d.Get("annotations")),
		Ttl:         ttl,
		KeepCount:   keepCount,
		Enabled:     d.Get("enabled").(bool),
		Schedule:    d.Get("schedule").(string),
		Timezone:    d.Get("timezone").(string),
		Target:      target,
	}

	identity := d.Id()
	policy, err := client.IaaS().UpdateSnapshotPolicy(ctx, identity, updatePolicy)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update snapshot policy: %w", err))
	}

	if policy != nil {
		d.Set("name", policy.Name)
		d.Set("description", policy.Description)
		d.Set("slug", policy.Slug)
		d.Set("enabled", policy.Enabled)
		d.Set("schedule", policy.Schedule)
		d.Set("timezone", policy.Timezone)
		d.Set("labels", policy.Labels)
		d.Set("annotations", policy.Annotations)
		d.Set("ttl", policy.Ttl.String())
		if policy.KeepCount != nil {
			d.Set("keep_count", *policy.KeepCount)
		}
		return nil
	}

	return resourceSnapshotPolicyRead(ctx, d, m)
}

func resourceSnapshotPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Thalassa client: %w", err))
	}

	identity := d.Id()

	err = client.IaaS().DeleteSnapshotPolicy(ctx, identity)
	if err != nil {
		if !tcclient.IsNotFound(err) {
			return diag.FromErr(fmt.Errorf("failed to delete snapshot policy: %w", err))
		}
	}

	d.SetId("")
	return nil
}
