package iaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	iaas "github.com/thalassa-cloud/client-go/iaas"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func dataSourceSubnet() *schema.Resource {
	return &schema.Resource{
		Description: "Get a subnet by name",
		ReadContext: dataSourceSubnetRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"slug": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the Subnet. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSubnetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	subnets, err := client.IaaS().ListSubnets(ctx, &iaas.ListSubnetsRequest{})
	if err != nil {
		return diag.FromErr(err)
	}

	// First count how many subnets exist with the given name in the VPC
	var matchingSubnets []iaas.Subnet
	for _, subnet := range subnets {
		if subnet.Name == d.Get("name").(string) && subnet.Vpc.Identity == d.Get("vpc_id").(string) {
			matchingSubnets = append(matchingSubnets, subnet)
		}
	}

	if len(matchingSubnets) > 1 {
		// Multiple subnets found with same name in VPC - require slug to be set
		slug := d.Get("slug").(string)
		if slug == "" {
			return diag.FromErr(fmt.Errorf("multiple subnets found with name '%s' in VPC, please specify slug", d.Get("name").(string)))
		}

		// Find exact match using slug
		for _, subnet := range matchingSubnets {
			if subnet.Slug == slug {
				d.SetId(subnet.Identity)
				d.Set("name", subnet.Name)
				d.Set("vpc_id", subnet.Vpc.Identity)
				d.Set("cidr", subnet.Cidr)
				d.Set("slug", subnet.Slug)
				return nil
			}
		}
		return diag.FromErr(fmt.Errorf("no subnet found with name '%s' and slug '%s' in VPC", d.Get("name").(string), slug))
	} else if len(matchingSubnets) == 1 {
		// Single match found
		subnet := matchingSubnets[0]
		d.SetId(subnet.Identity)
		d.Set("name", subnet.Name)
		d.Set("vpc_id", subnet.Vpc.Identity)
		d.Set("cidr", subnet.Cidr)
		d.Set("slug", subnet.Slug)
		return nil
	}

	return diag.FromErr(fmt.Errorf("subnet not found"))
}
