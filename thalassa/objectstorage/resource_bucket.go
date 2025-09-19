package objectstorage

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/thalassa-cloud/client-go/objectstorage"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func resourceBucket() *schema.Resource {
	return &schema.Resource{
		Description:   "Create and manage an object storage bucket",
		CreateContext: resourceBucketCreate,
		ReadContext:   resourceBucketRead,
		UpdateContext: resourceBucketUpdate,
		DeleteContext: resourceBucketDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the bucket. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.StringLenBetween(1, 63),
				Description:  "Name of the bucket",
				ForceNew:     true,
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Region of the bucket",
				ForceNew:    true,
			},
			"public": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the bucket is publicly accessible",
			},
			"policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The bucket policy as a JSON string",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the bucket",
			},
			"endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The endpoint URL for the bucket",
			},
			"versioning": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether the bucket is versioned",
			},
			"object_lock_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether the bucket has object lock enabled",
			},
		},
	}
}

func resourceBucketCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	bucketName := d.Get("name").(string)
	region := d.Get("region").(string)
	public := d.Get("public").(bool)
	versioning := d.Get("versioning").(bool)
	objectLockEnabled := d.Get("object_lock_enabled").(bool)

	var policyDoc *objectstorage.PolicyDocument
	if v, ok := d.GetOk("policy"); ok && v.(string) != "" {
		var doc objectstorage.PolicyDocument
		if err := json.Unmarshal([]byte(v.(string)), &doc); err != nil {
			return diag.FromErr(fmt.Errorf("invalid policy JSON: %w", err))
		}
		policyDoc = &doc
	}

	var bucketVersioning objectstorage.ObjectStorageBucketVersioning
	if versioning {
		bucketVersioning = objectstorage.ObjectStorageBucketVersioningEnabled
	} else {
		bucketVersioning = objectstorage.ObjectStorageBucketVersioningDisabled
	}

	createReq := objectstorage.CreateBucketRequest{
		BucketName:        bucketName,
		Public:            public,
		Region:            region,
		PolicyDocument:    policyDoc,
		Versioning:        bucketVersioning,
		ObjectLockEnabled: objectLockEnabled,
	}

	bucket, err := client.ObjectStorage().CreateBucket(ctx, createReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating bucket: %w", err))
	}

	d.SetId(bucket.Identity)
	return resourceBucketRead(ctx, d, m)
}

func resourceBucketRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting client: %w", err))
	}

	name := d.Get("name").(string)
	bucket, err := client.ObjectStorage().GetBucket(ctx, name)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting bucket: %w", err))
	}
	if bucket == nil {
		d.SetId("")
		return nil
	}

	d.SetId(bucket.Identity)
	d.Set("name", bucket.Name)
	d.Set("public", bucket.Public)
	d.Set("status", bucket.Status)
	d.Set("endpoint", bucket.Endpoint)
	if bucket.Region != nil {
		d.Set("region", bucket.Region.Identity)
	}
	d.Set("policy", bucket.Policy)
	d.Set("versioning", bucket.Versioning == objectstorage.ObjectStorageBucketVersioningEnabled)
	d.Set("object_lock_enabled", bucket.ObjectLockEnabled)

	return nil
}

func resourceBucketUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting client: %w", err))
	}
	name := d.Get("name").(string)

	// get the bucket
	bucket, err := client.ObjectStorage().GetBucket(ctx, name)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting bucket: %w", err))
	}
	if bucket == nil {
		return diag.FromErr(fmt.Errorf("bucket not found"))
	}

	updateReq := objectstorage.UpdateBucketRequest{}
	if d.HasChange("public") {
		updateReq.Public = d.Get("public").(bool)
	}
	if d.HasChange("policy") {
		if v, ok := d.GetOk("policy"); ok && v.(string) != "" {
			var doc objectstorage.PolicyDocument
			if err := json.Unmarshal([]byte(v.(string)), &doc); err != nil {
				return diag.FromErr(fmt.Errorf("invalid policy JSON: %w", err))
			}
			updateReq.PolicyDocument = &doc
		}
	}

	if d.HasChange("versioning") {
		switch d.Get("versioning").(bool) {
		case true:
			updateReq.Versioning = objectstorage.ObjectStorageBucketVersioningEnabled
		case false:
			if bucket.Versioning == objectstorage.ObjectStorageBucketVersioningEnabled {
				updateReq.Versioning = objectstorage.ObjectStorageBucketVersioningSuspended
			} else {
				updateReq.Versioning = objectstorage.ObjectStorageBucketVersioningDisabled
			}
		}
	}

	if d.HasChange("object_lock_enabled") {
		updateReq.ObjectLockEnabled = convert.Ptr(d.Get("object_lock_enabled").(bool))
	}

	_, err = client.ObjectStorage().UpdateBucket(ctx, name, updateReq)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceBucketRead(ctx, d, m)
}

func resourceBucketDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting bucket: %w", err))
	}

	name := d.Get("name").(string)
	if err := client.ObjectStorage().DeleteBucket(ctx, name); err != nil {
		return diag.FromErr(fmt.Errorf("error deleting bucket: %w", err))
	}
	d.SetId("")
	return nil
}
