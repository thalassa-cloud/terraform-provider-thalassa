package iaas

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	iaas "github.com/thalassa-cloud/client-go/iaas"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func DataSourceVolumeType() *schema.Resource {
	return &schema.Resource{
		Description: "Get an volume type by name. Volume Types are used to create block volumes. The matching name is case insensitive.",
		ReadContext: dataSourceVolumeTypeRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the volume type.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the volume type.",
			},
			"storage_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The storage type of the volume type. For example: 'block'.",
			},
			"allow_resize": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the volume type allows resizing. If false, the volume size cannot be changed after creation.",
			},
		},
	}
}

func dataSourceVolumeTypeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}
	name, ok := d.Get("name").(string)
	if !ok {
		return diag.FromErr(fmt.Errorf("name is not a string"))
	}

	machineTypes, err := client.IaaS().ListVolumeTypes(ctx, &iaas.ListVolumeTypesRequest{})
	if err != nil {
		return diag.FromErr(err)
	}

	for _, machineType := range machineTypes {
		if strings.EqualFold(machineType.Name, name) {
			d.SetId(machineType.Identity)
			d.Set("id", machineType.Identity)
			d.Set("name", machineType.Name)
			d.Set("description", machineType.Description)
			d.Set("storage_type", machineType.StorageType)
			d.Set("allow_resize", machineType.AllowResize)
			return diag.Diagnostics{}
		}
	}
	return diag.FromErr(fmt.Errorf("not found"))
}
