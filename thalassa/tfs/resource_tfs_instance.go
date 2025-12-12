package tfs

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
	tfs "github.com/thalassa-cloud/client-go/tfs"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func resourceTfsInstance() *schema.Resource {
	return &schema.Resource{
		Description: `
		Provides a Thalassa Cloud TFS (Thalassa Filesystem Service) instance resource. 
		TFS provides a high-availability, multi-availability zone Network File System (NFS) service 
		for shared storage across your infrastructure. TFS supports NFSv4 and NFSv4.1 protocols.
		`,
		CreateContext: resourceTfsInstanceCreate,
		ReadContext:   resourceTfsInstanceRead,
		UpdateContext: resourceTfsInstanceUpdate,
		DeleteContext: resourceTfsInstanceDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the TFS Instance. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.StringLenBetween(1, 62),
				Description:  "Name of the TFS instance",
			},
			"slug": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.StringLenBetween(0, 255),
				Description:  "A human readable description about the TFS instance",
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Labels for the TFS instance",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Annotations for the TFS instance",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Region of the TFS instance",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Identity of the VPC to create the TFS instance in",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Identity of the subnet to create the TFS instance in",
			},
			"size_gb": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validate.IntAtLeast(1),
				Description:  "Size of the TFS instance in GB",
			},
			"security_group_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of security group identities to attach to the TFS instance",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"delete_protection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Delete protection of the TFS instance",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the TFS instance",
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
			"wait_until_available": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Wait until the TFS instance is available",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTfsInstanceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Validate and resolve region
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

	// Validate VPC exists
	vpcId := d.Get("vpc_id").(string)
	_, err = client.IaaS().GetVpc(ctx, vpcId)
	if err != nil {
		if tcclient.IsNotFound(err) {
			return diag.FromErr(fmt.Errorf("VPC not found: %s", vpcId))
		}
		return diag.FromErr(fmt.Errorf("failed to get VPC: %w", err))
	}

	// Validate subnet exists
	subnetId := d.Get("subnet_id").(string)
	_, err = client.IaaS().GetSubnet(ctx, subnetId)
	if err != nil {
		if tcclient.IsNotFound(err) {
			return diag.FromErr(fmt.Errorf("subnet not found: %s", subnetId))
		}
		return diag.FromErr(fmt.Errorf("failed to get subnet: %w", err))
	}

	// Get security group IDs
	var securityGroupIds []string
	if v, ok := d.GetOk("security_group_ids"); ok && v != nil {
		securityGroupIds = convert.ConvertToStringSlice(v)
	}

	createTfsInstance := tfs.CreateTfsInstanceRequest{
		Name:                     d.Get("name").(string),
		Description:              d.Get("description").(string),
		Labels:                   convert.ConvertToMap(d.Get("labels")),
		Annotations:              convert.ConvertToMap(d.Get("annotations")),
		CloudRegionIdentity:      region,
		VpcIdentity:              vpcId,
		SubnetIdentity:           subnetId,
		SizeGB:                   d.Get("size_gb").(int),
		SecurityGroupAttachments: securityGroupIds,
		DeleteProtection:         d.Get("delete_protection").(bool),
	}

	tfsInstance, err := client.Tfs().CreateTfsInstance(ctx, createTfsInstance)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create TFS instance: %w", err))
	}

	if tfsInstance != nil {
		d.SetId(tfsInstance.Identity)
		d.Set("slug", tfsInstance.Slug)
		d.Set("status", string(tfsInstance.Status))
	}

	if d.Get("wait_until_available").(bool) {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Minute)
		defer cancel()

		err = client.Tfs().WaitUntilTfsInstanceIsAvailable(ctxWithTimeout, tfsInstance.Identity)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to wait for TFS instance to be available: %w", err))
		}
	}

	return resourceTfsInstanceRead(ctx, d, m)
}

func resourceTfsInstanceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Thalassa client: %w", err))
	}

	identity := d.Id()
	tfsInstance, err := client.Tfs().GetTfsInstance(ctx, identity)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error getting TFS instance: %w", err))
	}
	if tfsInstance == nil {
		d.SetId("")
		return nil
	}

	d.SetId(tfsInstance.Identity)
	d.Set("name", tfsInstance.Name)
	d.Set("slug", tfsInstance.Slug)
	d.Set("status", string(tfsInstance.Status))
	d.Set("delete_protection", tfsInstance.DeleteProtection)
	d.Set("labels", tfsInstance.Labels)
	d.Set("annotations", tfsInstance.Annotations)
	d.Set("size_gb", tfsInstance.SizeGB)

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

func resourceTfsInstanceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Thalassa client: %w", err))
	}

	// Get security group IDs
	var securityGroupIds []string
	if v, ok := d.GetOk("security_group_ids"); ok && v != nil {
		securityGroupIds = convert.ConvertToStringSlice(v)
	}

	updateTfsInstance := tfs.UpdateTfsInstanceRequest{
		Name:                     d.Get("name").(string),
		Description:              d.Get("description").(string),
		Labels:                   convert.ConvertToMap(d.Get("labels")),
		Annotations:              convert.ConvertToMap(d.Get("annotations")),
		SizeGB:                   d.Get("size_gb").(int),
		SecurityGroupAttachments: securityGroupIds,
		DeleteProtection:         d.Get("delete_protection").(bool),
	}

	identity := d.Id()
	tfsInstance, err := client.Tfs().UpdateTfsInstance(ctx, identity, updateTfsInstance)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update TFS instance: %w", err))
	}

	if tfsInstance != nil {
		d.Set("name", tfsInstance.Name)
		if tfsInstance.Description != nil {
			d.Set("description", *tfsInstance.Description)
		}
		d.Set("slug", tfsInstance.Slug)
		d.Set("status", string(tfsInstance.Status))
		d.Set("labels", tfsInstance.Labels)
		d.Set("annotations", tfsInstance.Annotations)
		d.Set("delete_protection", tfsInstance.DeleteProtection)
		d.Set("size_gb", tfsInstance.SizeGB)

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

		return nil
	}

	return resourceTfsInstanceRead(ctx, d, m)
}

func resourceTfsInstanceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Thalassa client: %w", err))
	}

	identity := d.Id()

	err = client.Tfs().DeleteTfsInstance(ctx, identity)
	if err != nil {
		if !tcclient.IsNotFound(err) {
			return diag.FromErr(fmt.Errorf("failed to delete TFS instance: %w", err))
		}
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	err = client.Tfs().WaitUntilTfsInstanceIsDeleted(ctxWithTimeout, identity)
	if err != nil {
		if !strings.Contains(err.Error(), "timeout") {
			return diag.FromErr(fmt.Errorf("failed to wait for TFS instance deletion: %w", err))
		}
	}

	d.SetId("")
	return nil
}
