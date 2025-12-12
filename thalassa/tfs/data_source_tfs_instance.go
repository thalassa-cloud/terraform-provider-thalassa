package tfs

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
	tfs "github.com/thalassa-cloud/client-go/tfs"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func DataSourceTfsInstance() *schema.Resource {
	return &schema.Resource{
		Description: "Get a TFS (Thalassa Filesystem Service) instance",
		ReadContext: dataSourceTfsInstanceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity of the TFS instance",
			},
			"identity": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Identity of the TFS instance",
				ExactlyOneOf: []string{"identity", "name"},
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Name of the TFS instance",
				ExactlyOneOf: []string{"identity", "name"},
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the TFS Instance. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"slug": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Slug of the TFS instance",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the TFS instance",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the TFS instance",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Region of the TFS instance",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity of the VPC",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity of the subnet",
			},
			"size_gb": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Size of the TFS instance in GB",
			},
			"delete_protection": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Delete protection of the TFS instance",
			},
			"endpoints": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of endpoints for the TFS instance",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"identity": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identity of the endpoint",
						},
						"address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP address of the endpoint",
						},
						"hostname": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Hostname of the endpoint",
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Port of the endpoint (defaults to 2049 for NFS)",
						},
					},
				},
			},
			"security_group_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of security group identities attached to the TFS instance",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Labels for the TFS instance",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Annotations for the TFS instance",
			},
		},
	}
}

func dataSourceTfsInstanceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	var tfsInstance *tfs.TfsInstance

	if identity, ok := d.GetOk("identity"); ok {
		// Look up by identity
		tfsInstance, err = client.Tfs().GetTfsInstance(ctx, identity.(string))
		if err != nil {
			if tcclient.IsNotFound(err) {
				return diag.FromErr(fmt.Errorf("TFS instance not found: %s", identity.(string)))
			}
			return diag.FromErr(err)
		}
	} else if name, ok := d.GetOk("name"); ok {
		// Look up by name
		tfsInstances, err := client.Tfs().ListTfsInstances(ctx, &tfs.ListTfsInstancesRequest{})
		if err != nil {
			return diag.FromErr(err)
		}

		// Find the TFS instance with the matching name
		for _, t := range tfsInstances {
			if t.Name == name.(string) {
				tfsInstance = &t
				break
			}
		}

		if tfsInstance == nil {
			return diag.FromErr(fmt.Errorf("TFS instance with name %s not found", name.(string)))
		}
	}

	// Set the ID and other attributes
	d.SetId(tfsInstance.Identity)
	d.Set("id", tfsInstance.Identity)
	d.Set("name", tfsInstance.Name)
	d.Set("slug", tfsInstance.Slug)
	d.Set("status", string(tfsInstance.Status))
	d.Set("delete_protection", tfsInstance.DeleteProtection)
	d.Set("labels", tfsInstance.Labels)
	d.Set("annotations", tfsInstance.Annotations)

	if tfsInstance.Description != nil {
		d.Set("description", *tfsInstance.Description)
	}

	if tfsInstance.Region != nil {
		d.Set("region", tfsInstance.Region.Identity)
	}

	if tfsInstance.Vpc != nil {
		d.Set("vpc_id", tfsInstance.Vpc.Identity)
	}

	if tfsInstance.Subnet != nil {
		d.Set("subnet_id", tfsInstance.Subnet.Identity)
	}

	d.Set("size_gb", tfsInstance.SizeGB)

	// Set endpoints
	if len(tfsInstance.Endpoints) > 0 {
		endpoints := make([]map[string]interface{}, len(tfsInstance.Endpoints))
		for i, endpoint := range tfsInstance.Endpoints {
			endpointMap := map[string]interface{}{
				"identity": endpoint.Identity,
				"address":  endpoint.EndpointAddress,
				"hostname": endpoint.EndpointHostname,
				"port":     2049, // Standard NFS port
			}
			endpoints[i] = endpointMap
		}
		d.Set("endpoints", endpoints)
	} else {
		d.Set("endpoints", []map[string]interface{}{})
	}

	// Set security group IDs
	if len(tfsInstance.SecurityGroups) > 0 {
		securityGroupIds := make([]string, len(tfsInstance.SecurityGroups))
		for i, sg := range tfsInstance.SecurityGroups {
			securityGroupIds[i] = sg.Identity
		}
		d.Set("security_group_ids", securityGroupIds)
	} else {
		d.Set("security_group_ids", []string{})
	}

	return nil
}
