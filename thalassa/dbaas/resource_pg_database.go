package dbaas

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thalassa-cloud/client-go/dbaas"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
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
				ForceNew:    true,
				Description: "The name of the database",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					if strings.TrimSpace(val.(string)) == "" {
						errs = append(errs, fmt.Errorf("database name is required"))
					}
					warns = []string{}
					return
				},
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
			"allow_connections": {
				Type:        schema.TypeBool,
				Default:     true,
				Optional:    true,
				Description: "If false then no one can connect to this database. Defaults to true.",
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
	var dbCluster *dbaas.DbCluster

	for {
		select {
		case <-ctx.Done():
			return diag.FromErr(ctx.Err())
		default:
			time.Sleep(1 * time.Second)
		}
		dbCluster, err = client.DBaaS().GetDbCluster(ctx, dbClusterId)
		if err != nil {
			if tcclient.IsNotFound(err) {
				return diag.FromErr(fmt.Errorf("db cluster not found: %w", err))
			}
			return diag.FromErr(fmt.Errorf("error getting db cluster: %w", err))
		}
		if dbCluster == nil {
			return diag.FromErr(fmt.Errorf("db cluster not found"))
		}
		if dbCluster.Status == dbaas.DbClusterStatusReady {
			break
		}
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
		return diag.FromErr(fmt.Errorf("owner role %s not found for database cluster", ownerRoleId))
	}

	// Check if the database already exists
	for _, database := range dbCluster.PostgresDatabases {
		if strings.EqualFold(database.Name, d.Get("name").(string)) {
			d.SetId(database.Identity)
			d.Set("name", database.Name)
			d.Set("db_cluster_id", dbClusterId)
			d.Set("owner_role_id", ownerRoleId)
			d.Set("connection_limit", database.ConnectionLimit)
			return resourcePgDatabaseRead(ctx, d, m)
		}
	}

	createDatabase := dbaas.CreatePgDatabaseRequest{
		Name:            d.Get("name").(string),
		Owner:           ownerRoleName,
		ConnectionLimit: &connectionLimit,
	}

	if allowConnections, ok := d.GetOk("allow_connections"); ok {
		allowConnectionsBool := allowConnections.(bool)
		createDatabase.AllowConnections = &allowConnectionsBool
	}

	err = client.DBaaS().CreatePgDatabase(ctx, dbCluster.Identity, createDatabase)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating pg database: %w", err))
	}

	// Get the database
	dbCluster, err = client.DBaaS().GetDbCluster(ctx, dbClusterId)
	if err != nil {
		if tcclient.IsNotFound(err) {
			return diag.FromErr(fmt.Errorf("db cluster not found: %w", err))
		}
		return diag.FromErr(fmt.Errorf("error getting db cluster: %w", err))
	}

	for _, database := range dbCluster.PostgresDatabases {
		if strings.EqualFold(database.Name, createDatabase.Name) {
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
	var dbCluster *dbaas.DbCluster

	for {
		select {
		case <-ctx.Done():
			return diag.FromErr(ctx.Err())
		default:
			time.Sleep(1 * time.Second)
		}
		dbCluster, err = client.DBaaS().GetDbCluster(ctx, dbClusterId)
		if err != nil {
			if tcclient.IsNotFound(err) {
				d.SetId("")
				return nil // a deleted db cluster means the pg database is also deleted
			}
			return diag.FromErr(err)
		}
		if dbCluster == nil {
			return diag.FromErr(fmt.Errorf("db cluster not found"))
		}

		if dbCluster.Status == dbaas.DbClusterStatusReady {
			break
		}
	}

	ownerRoleId := d.Get("owner_role_id").(string)

	for _, database := range dbCluster.PostgresDatabases {
		if strings.EqualFold(database.Name, d.Get("name").(string)) {
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
	var dbCluster *dbaas.DbCluster

	for {
		select {
		case <-ctx.Done():
			return diag.FromErr(ctx.Err())
		default:
			time.Sleep(1 * time.Second)
		}
		dbCluster, err = client.DBaaS().GetDbCluster(ctx, dbClusterId)
		if err != nil {
			if tcclient.IsNotFound(err) {
				return diag.FromErr(fmt.Errorf("db cluster not found: %w", err))
			}
			return diag.FromErr(fmt.Errorf("error getting db cluster: %w", err))
		}

		if dbCluster == nil {
			return diag.FromErr(fmt.Errorf("db cluster not found"))
		}

		if dbCluster.Status == dbaas.DbClusterStatusReady {
			break
		}
	}

	connectionLimit := d.Get("connection_limit").(int)

	updateDatabase := dbaas.UpdatePgDatabaseRequest{
		ConnectionLimit: &connectionLimit,
	}

	if allowConnections, ok := d.GetOk("allow_connections"); ok {
		allowConnectionsBool := allowConnections.(bool)
		updateDatabase.AllowConnections = &allowConnectionsBool
	}

	err = client.DBaaS().UpdatePgDatabase(ctx, dbCluster.Identity, d.Get("id").(string), updateDatabase)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating pg database: %w", err))
	}

	ownerRoleId := d.Get("owner_role_id").(string)

	for _, database := range dbCluster.PostgresDatabases {
		if strings.EqualFold(database.Name, d.Get("name").(string)) {
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
	var dbCluster *dbaas.DbCluster

	for {
		select {
		case <-ctx.Done():
			return diag.FromErr(ctx.Err())
		default:
			time.Sleep(1 * time.Second)
		}
		dbCluster, err = client.DBaaS().GetDbCluster(ctx, dbClusterId)
		if err != nil {
			if tcclient.IsNotFound(err) {
				d.SetId("")
				return nil // a deleted db cluster means the pg database is also deleted
			}
			return diag.FromErr(fmt.Errorf("error getting db cluster: %w", err))
		}

		if dbCluster == nil {
			return diag.FromErr(fmt.Errorf("db cluster not found"))
		}

		if dbCluster.Status == dbaas.DbClusterStatusReady {
			break
		}
	}

	// Check if the database is not already scheduled for deletion
	for _, database := range dbCluster.PostgresDatabases {
		if strings.EqualFold(database.Name, d.Get("name").(string)) {
			if database.DeleteScheduledAt != nil {
				d.SetId("")
				return nil
			}
		}
	}

	err = client.DBaaS().DeletePgDatabase(ctx, dbCluster.Identity, d.Get("id").(string), true)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil // a deleted db cluster means the pg database is also deleted
		}
		return diag.FromErr(fmt.Errorf("error deleting pg database: %w", err))
	}

	d.SetId("")
	return nil
}
