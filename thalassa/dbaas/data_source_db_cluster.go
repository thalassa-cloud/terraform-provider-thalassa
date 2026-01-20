package dbaas

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/thalassa-cloud/client-go/dbaas"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func dataSourceDbCluster() *schema.Resource {
	return &schema.Resource{
		Description: "Get an DB Cluster",
		ReadContext: dataSourceDbClusterRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity of the DB Cluster",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the DB Cluster",
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the Db Cluster. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"slug": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Slug of the DB Cluster",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the DB Cluster",
			},
			"labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Labels of the DB Cluster",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Annotations of the DB Cluster",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Subnet of the DB Cluster",
			},
			"database_instance_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Database instance type of the DB Cluster",
			},
			"replicas": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of instances in the cluster",
			},
			"engine": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Database engine of the cluster",
			},
			"engine_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Version of the database engine",
			},
			"parameters": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Map of parameter name to database engine specific parameter value",
			},
			"allocated_storage": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Amount of storage allocated to the cluster in GB",
			},
			"volume_type_class": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Storage type used to determine the size of the cluster storage",
			},
			"auto_minor_version_upgrade": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag indicating if the cluster should automatically upgrade to the latest minor version",
			},
			"database_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the database on the cluster",
			},
			"delete_protection": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag indicating if the cluster should be protected from deletion",
			},
			"security_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of security groups associated with the cluster",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the cluster",
			},
			"endpoint_ipv4": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IPv4 address of the cluster endpoint",
			},
			"endpoint_ipv6": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IPv6 address of the cluster endpoint",
			},
			"port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Port of the cluster endpoint",
			},
		},
	}
}

func dataSourceDbClusterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("id").(string)
	name := d.Get("name").(string)
	slug := d.Get("slug").(string)
	var DbCluster *dbaas.DbCluster
	dbClusters, err := client.DBaaS().ListDbClusters(ctx, &dbaas.ListDbClustersRequest{})
	if err != nil {
		return diag.FromErr(err)
	}

	// Track if we found multiple matches to provide better error messages
	var matchingClusters []*dbaas.DbCluster

	for _, dbCluster := range dbClusters {
		if dbCluster.Identity == id {
			DbCluster = &dbCluster
			break
		}
		if strings.EqualFold(dbCluster.Name, name) {
			// If name matches, also check slug to ensure uniqueness
			if slug != "" && strings.EqualFold(dbCluster.Slug, slug) {
				DbCluster = &dbCluster
				break
			} else if slug == "" {
				// If no slug provided, collect all name matches
				matchingClusters = append(matchingClusters, &dbCluster)
			}
			continue
		}
		if strings.EqualFold(dbCluster.Slug, slug) {
			DbCluster = &dbCluster
			break
		}
	}

	// Handle multiple matches case
	if DbCluster == nil && len(matchingClusters) > 1 {
		var names []string
		for _, cluster := range matchingClusters {
			names = append(names, fmt.Sprintf("%s (slug: %s)", cluster.Name, cluster.Slug))
		}
		return diag.Errorf("Multiple clusters found with name '%s'. Please specify a slug to uniquely identify the cluster. Matching clusters: %v", name, names)
	} else if len(matchingClusters) == 1 {
		DbCluster = matchingClusters[0]
	}

	if DbCluster == nil {
		return diag.FromErr(fmt.Errorf("DbCluster was not found"))
	}

	d.SetId(DbCluster.Identity)
	d.Set("name", DbCluster.Name)
	d.Set("slug", DbCluster.Slug)
	d.Set("description", DbCluster.Description)
	d.Set("labels", DbCluster.Labels)
	d.Set("annotations", DbCluster.Annotations)
	d.Set("delete_protection", DbCluster.DeleteProtection)
	d.Set("replicas", DbCluster.Replicas)
	d.Set("engine", DbCluster.Engine)
	d.Set("engine_version", DbCluster.EngineVersion)
	d.Set("parameters", DbCluster.Parameters)
	d.Set("allocated_storage", DbCluster.AllocatedStorage)
	d.Set("auto_minor_version_upgrade", DbCluster.AutoMinorVersionUpgrade)
	d.Set("status", DbCluster.Status)
	d.Set("endpoint_ipv4", DbCluster.EndpointIpv4)
	d.Set("endpoint_ipv6", DbCluster.EndpointIpv6)
	d.Set("port", DbCluster.Port)

	// Handle optional fields
	if DbCluster.DatabaseName != nil {
		d.Set("database_name", *DbCluster.DatabaseName)
	}
	if DbCluster.Subnet != nil {
		d.Set("subnet_id", DbCluster.Subnet.Identity)
	}
	if DbCluster.DatabaseInstanceType != nil {
		d.Set("database_instance_type", DbCluster.DatabaseInstanceType.Identity)
	}
	if DbCluster.VolumeTypeClass != nil {
		d.Set("volume_type_class", DbCluster.VolumeTypeClass.Identity)
	}
	if DbCluster.SecurityGroups != nil {
		securityGroupIds := make([]string, len(DbCluster.SecurityGroups))
		for i, sg := range DbCluster.SecurityGroups {
			securityGroupIds[i] = sg.Identity
		}
		d.Set("security_groups", securityGroupIds)
	}

	return nil
}
