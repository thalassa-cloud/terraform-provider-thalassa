package objectstorage

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func DataSourceBucket() *schema.Resource {
	return &schema.Resource{
		Description: "Get an object storage bucket",
		ReadContext: dataSourceBucketRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity of the bucket",
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
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Region of the bucket",
			},
			"public": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the bucket is publicly accessible",
			},
			"policy": {
				Type:        schema.TypeString,
				Computed:    true,
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
		},
	}
}

func dataSourceBucketRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)
	region := d.Get("region").(string)

	bucket, err := client.ObjectStorage().GetBucket(ctx, name)
	if err != nil {
		return diag.FromErr(err)
	}

	if bucket == nil {
		return diag.FromErr(fmt.Errorf("bucket %s not found", name))
	}

	// Check if region filter is specified and matches
	if region != "" && bucket.Region != nil && bucket.Region.Identity != region {
		return diag.FromErr(fmt.Errorf("bucket %s not found in region %s", name, region))
	}

	d.SetId(bucket.Identity)
	d.Set("id", bucket.Identity)
	d.Set("name", bucket.Name)
	d.Set("public", bucket.Public)
	d.Set("status", bucket.Status)
	d.Set("endpoint", bucket.Endpoint)

	if bucket.Region != nil {
		d.Set("region", bucket.Region.Identity)
	}

	// Set policy as JSON string if available
	d.Set("policy", bucket.Policy)

	return diag.Diagnostics{}
}
