package iaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	iaas "github.com/thalassa-cloud/client-go/iaas"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func DataSourceVpc() *schema.Resource {
	return &schema.Resource{
		Description: "Get an vpc",
		ReadContext: dataSourceVpcRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity of the Vpc",
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the Vpc. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.StringLenBetween(1, 62),
				Description:  "Name of the Vpc",
			},
			"slug": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Slug of the Vpc",
			},
			"cidrs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of CIDRs for the Vpc",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validate.IsCIDR,
				},
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.StringLenBetween(0, 255),
				Description:  "A human readable description about the vpc",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Region of the Vpc. Provide the identity of the region. Can only be set on creation.",
			},
			"labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Labels for the Vpc",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Annotations for the Vpc",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the Vpc",
			},
		},
	}
}

func dataSourceVpcRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := provider.GetProvider(m)
	slug := d.Get("slug").(string)
	name := d.Get("name").(string)
	region := d.Get("region").(string)

	vpcs, err := provider.Client.IaaS().ListVpcs(ctx, &iaas.ListVpcsRequest{})
	if err != nil {
		return diag.FromErr(err)
	}

	var vpc *iaas.Vpc
	// First find all VPCs matching name and region
	var matchingVpcs []iaas.Vpc
	for _, v := range vpcs {
		if v.CloudRegion != nil && region != "" && v.CloudRegion.Identity != region {
			continue
		}
		if v.Name == name {
			matchingVpcs = append(matchingVpcs, v)
		}
	}

	if len(matchingVpcs) > 1 {
		// Multiple VPCs found with same name - require slug
		if slug == "" {
			var slugs []string
			for _, v := range matchingVpcs {
				slugs = append(slugs, v.Slug)
			}
			return diag.FromErr(fmt.Errorf("multiple VPCs found with name '%s', please specify one of these slugs: %v", name, slugs))
		}

		// Find exact match using slug
		for _, v := range matchingVpcs {
			if v.Slug == slug {
				vpc = &v
				break
			}
		}
	} else if len(matchingVpcs) == 1 {
		// Single match found
		vpc = &matchingVpcs[0]
	} else if slug != "" {
		// No matches by name, try finding by slug
		for _, v := range vpcs {
			if v.CloudRegion != nil && region != "" && v.CloudRegion.Identity != region {
				continue
			}
			if v.Slug == slug {
				vpc = &v
				break
			}
		}
	}

	if vpc != nil {

		d.SetId(vpc.Identity)
		d.Set("id", vpc.Identity)
		d.Set("name", vpc.Name)
		d.Set("slug", vpc.Slug)
		d.Set("description", vpc.Description)
		d.Set("status", vpc.Status)
		d.Set("labels", vpc.Labels)
		d.Set("annotations", vpc.Annotations)
		if vpc.CloudRegion != nil {
			d.Set("region", vpc.CloudRegion.Identity)
		}
		d.Set("cidrs", vpc.CIDRs)
		return diag.Diagnostics{}
	}

	return diag.FromErr(fmt.Errorf("vpc %s not found", name))
}
