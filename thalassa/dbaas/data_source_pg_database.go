package dbaas

import (
	"context"
	"fmt"

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
		},
	}
}

func dataSourcePgDatabaseRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	dbClusterId := d.Get("db_cluster_id").(string)

	dbCluster, err := client.DbaaSAlphaV1().GetDbCluster(ctx, dbClusterId)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, database := range dbCluster.PostgresDatabases {
		if database.Name == d.Get("name").(string) {
			d.SetId(database.Identity)
			d.Set("name", database.Name)
			d.Set("db_cluster_id", dbClusterId)
			d.Set("owner_role_id", database.Owner)
			d.Set("connection_limit", database.ConnectionLimit)
			return nil
		}
	}

	return diag.FromErr(fmt.Errorf("database not found"))
}
