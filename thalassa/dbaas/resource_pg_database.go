package dbaas

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thalassa-cloud/client-go/dbaas/dbaasalphav1"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func resourcePgDatabase() *schema.Resource {
	return &schema.Resource{
		Description:   "Create a PostgreSQL database",
		CreateContext: resourcePgDatabaseCreate,
		ReadContext:   resourcePgDatabaseRead,
		UpdateContext: resourcePgDatabaseUpdate,
		DeleteContext: resourcePgDatabaseDelete,
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
				Required:    true,
				Description: "The ID of the owner role",
			},
			"connection_limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     -1,
				Description: "The connection limit of the database",
			},
		},
	}
}

func resourcePgDatabaseCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	dbClusterId := d.Get("db_cluster_id").(string)
	var dbCluster *dbaasalphav1.DbCluster

	for {
		dbClusters, err := client.DbaaSAlphaV1().ListDbClusters(ctx, &dbaasalphav1.ListDbClustersRequest{})
		if err != nil {
			return diag.FromErr(err)
		}

		for _, cluster := range dbClusters {
			if cluster.Identity == dbClusterId {
				dbCluster = &cluster
				break
			}
		}
		if dbCluster.Status == dbaasalphav1.DbClusterStatusReady {
			break
		}
		time.Sleep(1 * time.Second)
	}

	dbCluster, err = client.DbaaSAlphaV1().GetDbCluster(ctx, dbClusterId)
	if err != nil {
		return diag.FromErr(err)
	}

	// check the owner exists
	ownerRoleId := d.Get("owner_role_id").(string)
	for _, role := range dbCluster.PostgresRoles {
		if role.Identity == ownerRoleId {
			break
		}
	}

	connectionLimit := d.Get("connection_limit").(int)

	ownerRoleName := ""
	for _, role := range dbCluster.PostgresRoles {
		if role.Identity == ownerRoleId {
			ownerRoleName = role.Name
			break
		}
	}
	if ownerRoleName == "" {
		return diag.FromErr(fmt.Errorf("owner role not found"))
	}

	// Check if the database already exists
	for _, database := range dbCluster.PostgresDatabases {
		if database.Name == d.Get("name").(string) {
			d.SetId(database.Identity)
			d.Set("name", database.Name)
			d.Set("db_cluster_id", dbClusterId)
			d.Set("owner_role_id", ownerRoleId)
			d.Set("connection_limit", database.ConnectionLimit)
			return resourcePgDatabaseRead(ctx, d, m)
		}
	}

	createDatabase := dbaasalphav1.CreatePgDatabaseRequest{
		Name:            d.Get("name").(string),
		Owner:           ownerRoleName,
		ConnectionLimit: &connectionLimit,
	}

	err = client.DbaaSAlphaV1().CreatePgDatabase(ctx, dbCluster.Identity, createDatabase)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get the database
	dbCluster, err = client.DbaaSAlphaV1().GetDbCluster(ctx, dbClusterId)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, database := range dbCluster.PostgresDatabases {
		if database.Name == createDatabase.Name {
			d.SetId(database.Identity)
			d.Set("name", database.Name)
			d.Set("db_cluster_id", dbClusterId)
			d.Set("owner_role_id", ownerRoleId)
			d.Set("connection_limit", database.ConnectionLimit)
			return resourcePgDatabaseRead(ctx, d, m)
		}
	}

	return diag.FromErr(fmt.Errorf("database not found"))
}

func resourcePgDatabaseRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	dbClusterId := d.Get("db_cluster_id").(string)
	var dbCluster *dbaasalphav1.DbCluster

	for {
		dbClusters, err := client.DbaaSAlphaV1().ListDbClusters(ctx, &dbaasalphav1.ListDbClustersRequest{})
		if err != nil {
			return diag.FromErr(err)
		}

		for _, cluster := range dbClusters {
			if cluster.Identity == dbClusterId {
				dbCluster = &cluster
				break
			}
		}
		if dbCluster.Status == dbaasalphav1.DbClusterStatusReady {
			break
		}
		time.Sleep(1 * time.Second)
	}

	dbCluster, err = client.DbaaSAlphaV1().GetDbCluster(ctx, dbClusterId)
	if err != nil {
		return diag.FromErr(err)
	}

	ownerRoleId := d.Get("owner_role_id").(string)

	for _, database := range dbCluster.PostgresDatabases {
		if database.Name == d.Get("name").(string) {
			d.SetId(database.Identity)
			d.Set("name", database.Name)
			d.Set("db_cluster_id", dbClusterId)
			d.Set("owner_role_id", ownerRoleId)
			d.Set("connection_limit", database.ConnectionLimit)
			return nil
		}
	}
	return diag.FromErr(fmt.Errorf("database not found"))
}

func resourcePgDatabaseUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	dbClusterId := d.Get("db_cluster_id").(string)
	var dbCluster *dbaasalphav1.DbCluster

	for {
		dbClusters, err := client.DbaaSAlphaV1().ListDbClusters(ctx, &dbaasalphav1.ListDbClustersRequest{})
		if err != nil {
			return diag.FromErr(err)
		}

		for _, cluster := range dbClusters {
			if cluster.Identity == dbClusterId {
				dbCluster = &cluster
				break
			}
		}
		if dbCluster.Status == dbaasalphav1.DbClusterStatusReady {
			break
		}
		time.Sleep(1 * time.Second)
	}

	connectionLimit := d.Get("connection_limit").(int)

	updateDatabase := dbaasalphav1.UpdatePgDatabaseRequest{
		ConnectionLimit: &connectionLimit,
	}

	err = client.DbaaSAlphaV1().UpdatePgDatabase(ctx, dbCluster.Identity, d.Get("id").(string), updateDatabase)
	if err != nil {
		return diag.FromErr(err)
	}

	ownerRoleId := d.Get("owner_role_id").(string)

	for _, database := range dbCluster.PostgresDatabases {
		if database.Name == d.Get("name").(string) {
			d.SetId(database.Identity)
			d.Set("name", database.Name)
			d.Set("db_cluster_id", dbClusterId)
			d.Set("owner_role_id", ownerRoleId)
			d.Set("connection_limit", database.ConnectionLimit)
			return resourcePgDatabaseRead(ctx, d, m)
		}
	}
	return diag.FromErr(fmt.Errorf("database not found"))
}

func resourcePgDatabaseDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	dbClusterId := d.Get("db_cluster_id").(string)
	var dbCluster *dbaasalphav1.DbCluster

	for {
		dbClusters, err := client.DbaaSAlphaV1().ListDbClusters(ctx, &dbaasalphav1.ListDbClustersRequest{})
		if err != nil {
			return diag.FromErr(err)
		}

		for _, cluster := range dbClusters {
			if cluster.Identity == dbClusterId {
				dbCluster = &cluster
				break
			}
		}
		if dbCluster.Status == dbaasalphav1.DbClusterStatusReady {
			break
		}
		time.Sleep(1 * time.Second)
	}

	// Check if the database is not already scheduled for deletion
	for _, database := range dbCluster.PostgresDatabases {
		if database.Name == d.Get("name").(string) {
			if database.DeleteScheduledAt != nil {
				d.SetId("")
				return nil
			}
		}
	}

	err = client.DbaaSAlphaV1().DeletePgDatabase(ctx, dbCluster.Identity, d.Get("id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
