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

func resourceBlockVolume() *schema.Resource {
	return &schema.Resource{
		Description: `
		Provides a Thalassa Cloud Block Volume resource. This can be used to create, manage, and attach a detachable storage device to a virtual machine instance. 
		`,

		CreateContext: resourceBlockVolumeCreate,
		ReadContext:   resourceBlockVolumeRead,
		UpdateContext: resourceBlockVolumeUpdate,
		DeleteContext: resourceBlockVolumeDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Region of the Block Volume.",
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the Block Volume. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the Block Volume",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Block Volume",
			},
			"slug": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A human readable description about the Block Volume",
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Labels for the Block Volume",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Annotations for the Block Volume",
			},
			"volume_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StringIsNotWhiteSpace,
				Description:  "Volume type of the Block Volume",
			},
			"delete_protection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Delete protection of the Block Volume",
			},
			"size_gb": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validate.IntAtLeast(1),
				Description:  "Size of the Block Volume in GB",
			},
			"wait_until_ready": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Wait until the Block Volume is ready",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceBlockVolumeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	// lets find the volume type
	volumeType := d.Get("volume_type").(string)
	volumeTypeIdentity := ""
	volumeTypes, err := client.IaaS().ListVolumeTypes(ctx, &iaas.ListVolumeTypesRequest{})
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to list volume types: %w", err))
	}
	foundVolumeType := false
	for _, v := range volumeTypes {
		if v.Identity == volumeType || strings.EqualFold(v.Name, volumeType) {
			volumeTypeIdentity = v.Identity
			foundVolumeType = true
			break
		}
	}

	if !foundVolumeType {
		availableVolumeTypes := make([]string, len(volumeTypes))
		for i, v := range volumeTypes {
			availableVolumeTypes[i] = v.Name
		}
		return diag.FromErr(fmt.Errorf("volume type not found: %s. Available volume types: %v", volumeType, strings.Join(availableVolumeTypes, ", ")))
	}

	createBlockVolume := iaas.CreateVolume{
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		Labels:              convert.ConvertToMap(d.Get("labels")),
		Annotations:         convert.ConvertToMap(d.Get("annotations")),
		CloudRegionIdentity: region,
		VolumeTypeIdentity:  volumeTypeIdentity,
		Size:                d.Get("size_gb").(int),
		// DeleteProtection:          d.Get("delete_protection").(bool),
	}

	blockVolume, err := client.IaaS().CreateVolume(ctx, createBlockVolume)
	if err != nil {
		if tcclient.IsNotFound(err) {
			return diag.FromErr(fmt.Errorf("volume type not found: %w", err))
		}
		return diag.FromErr(fmt.Errorf("failed to create block volume: %w", err))
	}

	if blockVolume != nil {
		d.SetId(blockVolume.Identity)
		d.Set("slug", blockVolume.Slug)
		d.Set("status", blockVolume.Status)
	}

	if d.Get("wait_until_ready").(bool) {
		// wait until the volume is ready
		for {
			blockVolume, err := client.IaaS().GetVolume(ctx, blockVolume.Identity)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to check block volume status: %w", err))
			}
			if strings.EqualFold(blockVolume.Status, "available") || strings.EqualFold(blockVolume.Status, "ready") {
				break
			}
			time.Sleep(1 * time.Second)
		}
	}

	return resourceBlockVolumeRead(ctx, d, m)
}

func resourceBlockVolumeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Thalassa client: %w", err))
	}

	identity := d.Get("id").(string)
	blockVolume, err := client.IaaS().GetVolume(ctx, identity)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error getting blockVolume: %s", err))
	}
	if blockVolume == nil {
		d.SetId("")
		return nil
	}

	d.SetId(blockVolume.Identity)
	d.Set("name", blockVolume.Name)
	d.Set("slug", blockVolume.Slug)
	d.Set("description", blockVolume.Description)
	d.Set("labels", blockVolume.Labels)
	d.Set("annotations", blockVolume.Annotations)
	d.Set("status", blockVolume.Status)
	d.Set("size_gb", blockVolume.Size)

	if blockVolume.VolumeType != nil {
		if d.Get("volume_type").(string) != "" {
			d.Set("volume_type", d.Get("volume_type").(string))
		} else {
			d.Set("volume_type", blockVolume.VolumeType.Identity)
		}
	}

	if blockVolume.Region != nil {
		currentRegion := d.Get("region").(string)

		switch currentRegion {
		case "":
			d.Set("region", blockVolume.Region.Slug)
		case blockVolume.Region.Slug:
			d.Set("region", blockVolume.Region.Slug)
		case blockVolume.Region.Identity:
			d.Set("region", blockVolume.Region.Identity)
		default:
			d.Set("region", blockVolume.Region.Slug)
		}
	}

	return nil
}

func resourceBlockVolumeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Thalassa client: %w", err))
	}

	updateBlockVolume := iaas.UpdateVolume{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Labels:      convert.ConvertToMap(d.Get("labels")),
		Annotations: convert.ConvertToMap(d.Get("annotations")),
		Size:        d.Get("size_gb").(int),
	}
	identity := d.Get("id").(string)
	blockVolume, err := client.IaaS().UpdateVolume(ctx, identity, updateBlockVolume)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update block volume: %w", err))
	}
	if blockVolume != nil {
		d.Set("name", blockVolume.Name)
		d.Set("description", blockVolume.Description)
		d.Set("slug", blockVolume.Slug)
		d.Set("status", blockVolume.Status)
		d.Set("labels", blockVolume.Labels)
		d.Set("annotations", blockVolume.Annotations)
		if blockVolume.VolumeType != nil {
			if d.Get("volume_type").(string) != "" {
				d.Set("volume_type", d.Get("volume_type").(string))
			} else {
				d.Set("volume_type", blockVolume.VolumeType.Identity)
			}
		}
		d.Set("size_gb", blockVolume.Size)
		return nil
	}

	return resourceBlockVolumeRead(ctx, d, m)
}

func resourceBlockVolumeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Thalassa client: %w", err))
	}

	identity := d.Get("id").(string)

	err = client.IaaS().DeleteVolume(ctx, identity)
	if err != nil {
		if !tcclient.IsNotFound(err) {
			return diag.FromErr(fmt.Errorf("failed to delete block volume: %w", err))
		}
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()
	for {
		select {
		case <-ctxWithTimeout.Done():
			return diag.FromErr(fmt.Errorf("timeout waiting for block volume to be deleted"))
		case <-time.After(1 * time.Second):
		default:
		}

		blockVolume, err := client.IaaS().GetVolume(ctxWithTimeout, identity)
		if err != nil {
			if tcclient.IsNotFound(err) {
				break
			}
			return diag.FromErr(fmt.Errorf("failed to retrieve block volume: %w", err))
		}
		if blockVolume == nil {
			break
		}
		if strings.EqualFold(blockVolume.Status, "deleted") {
			break
		}
	}

	d.SetId("")
	return nil
}
