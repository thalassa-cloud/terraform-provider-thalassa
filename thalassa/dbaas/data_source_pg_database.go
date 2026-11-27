package dbaas

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func dataSourcePgDatabase() *schema.Resource {

	return &schema.Resource{
		ReadContext: dataSourcePgDatabaseRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the database",
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the Db Cluster. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"db_cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the database cluster",
			},
			"owner_role_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the owner role",
			},
			"connection_limit": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The connection limit of the database",
			},
			"extensions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "PostgreSQL extensions installed on the database",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the extension",
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Version of the extension",
						},
					},
				},
			},
			"auto_update_extensions": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Automatic extension upgrade policy",
			},
			"extension_status": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Observed PostgreSQL extension state",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the installed extension",
						},
						"installed_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Installed version of the extension",
						},
						"default_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Default version of the extension",
						},
						"available_versions": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Available versions of the extension",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"upgrade_available": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether an upgrade is available for the extension",
						},
					},
				},
			},
		},
	}
}

func dataSourcePgDatabaseRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {

	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	dbClusterId := d.Get("db_cluster_id").(string)

	dbCluster, err := client.DBaaS().GetDbCluster(ctx, dbClusterId)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, database := range dbCluster.PostgresDatabases {
		if strings.EqualFold(database.Name, d.Get("name").(string)) {
			d.SetId(database.Identity)
			_ = d.Set("name", database.Name)
			_ = d.Set("db_cluster_id", dbClusterId)
			_ = d.Set("owner_role_id", database.Owner)
			_ = d.Set("connection_limit", database.ConnectionLimit)
			_ = d.Set("extensions", flattenPgDatabaseExtensions(database.Extensions))
			_ = d.Set("auto_update_extensions", string(database.AutoUpdateExtensions))
			_ = d.Set("extension_status", flattenPgDatabaseExtensionStatus(database.ExtensionStatus))
			return nil
		}
	}

	return diag.FromErr(fmt.Errorf("database not found"))
}
