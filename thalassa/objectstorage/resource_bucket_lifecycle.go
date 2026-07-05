package objectstorage

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/thalassa-cloud/client-go/objectstorage"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"

	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func resourceBucketLifecycle() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage object storage bucket lifecycle rules. Each apply replaces the full rule set.",
		CreateContext: resourceBucketLifecycleCreate,
		ReadContext:   resourceBucketLifecycleRead,
		UpdateContext: resourceBucketLifecycleUpdate,
		DeleteContext: resourceBucketLifecycleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organisation_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"bucket_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the bucket.",
			},
			"rule": lifecycleRuleSchema(),
		},
	}
}

func resourceBucketLifecycleCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	return resourceBucketLifecycleApply(ctx, d, m)
}

func resourceBucketLifecycleUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	return resourceBucketLifecycleApply(ctx, d, m)
}

func resourceBucketLifecycleApply(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	bucketName := d.Get("bucket_name").(string)
	rules := expandLifecycleRules(d.Get("rule").([]any))

	if lifecycleHasNoncurrentRules(rules) {
		bucket, err := client.ObjectStorage().GetBucket(ctx, bucketName)
		if err != nil {
			return diag.FromErr(fmt.Errorf("reading bucket for lifecycle validation: %w", err))
		}
		if bucket.Versioning != objectstorage.ObjectStorageBucketVersioningEnabled {
			return diag.Errorf("noncurrent version lifecycle rules require bucket versioning to be Enabled on %q", bucketName)
		}
	}

	_, err = client.ObjectStorage().SetBucketLifecycle(ctx, bucketName, objectstorage.SetBucketLifecycleRequest{
		Rules: rules,
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("setting bucket lifecycle: %w", err))
	}

	d.SetId(bucketName)
	return resourceBucketLifecycleRead(ctx, d, m)
}

func resourceBucketLifecycleRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	bucketName := d.Id()
	if bucketName == "" {
		bucketName = d.Get("bucket_name").(string)
	}

	lifecycle, err := client.ObjectStorage().GetBucketLifecycle(ctx, bucketName)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("reading bucket lifecycle: %w", err))
	}

	d.SetId(bucketName)
	_ = d.Set("bucket_name", bucketName)
	_ = d.Set("rule", flattenLifecycleRules(lifecycle.Rules))

	return nil
}

func resourceBucketLifecycleDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	bucketName := d.Get("bucket_name").(string)
	if err := client.ObjectStorage().DeleteBucketLifecycle(ctx, bucketName); err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("deleting bucket lifecycle: %w", err))
	}

	d.SetId("")
	return nil
}
