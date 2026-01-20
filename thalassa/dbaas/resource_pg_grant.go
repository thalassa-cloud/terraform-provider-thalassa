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
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func resourcePgGrant() *schema.Resource {
	return &schema.Resource{
		Description:   "Create a PostgreSQL grant for a role on a database",
		CreateContext: resourcePgGrantCreate,
		ReadContext:   resourcePgGrantRead,
		UpdateContext: resourcePgGrantUpdate,
		DeleteContext: resourcePgGrantDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the PostgreSQL grant (grant name)",
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the Db Cluster. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the grant",
			},
			"db_cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the database cluster",
			},
			"role_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the role to grant permissions to",
			},
			"database_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the database to grant permissions on",
			},
			"read": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether the role can read from the database",
			},
			"write": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether the role can write to the database",
			},
		},
	}
}

func resourcePgGrantCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	dbClusterId := d.Get("db_cluster_id").(string)
	dbCluster, err := client.DBaaS().GetDbCluster(ctx, dbClusterId)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil // a deleted db cluster means the pg grant is also deleted
		}
		return diag.FromErr(err)
	}

	// Verify the role exists
	roleName := d.Get("role_name").(string)
	roleFound := false
	for _, role := range dbCluster.PostgresRoles {
		if strings.EqualFold(role.Name, roleName) {
			roleFound = true
			break
		}
	}
	if !roleFound {
		return diag.FromErr(fmt.Errorf("role %s not found in database cluster", roleName))
	}

	// Verify the database exists
	databaseName := d.Get("database_name").(string)
	databaseFound := false
	for _, database := range dbCluster.PostgresDatabases {
		if strings.EqualFold(database.Name, databaseName) {
			databaseFound = true
			break
		}
	}
	if !databaseFound {
		return diag.FromErr(fmt.Errorf("database %s not found in database cluster", databaseName))
	}

	createGrant := dbaas.CreatePgGrantRequest{
		Name:         d.Get("name").(string),
		RoleName:     roleName,
		DatabaseName: databaseName,
		Read:         d.Get("read").(bool),
		Write:        d.Get("write").(bool),
	}

	err = client.DBaaS().CreatePgGrant(ctx, dbCluster.Identity, createGrant)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating pg grant: %w", err))
	}

	// Use the grant name as the ID
	d.SetId(createGrant.Name)
	d.Set("name", createGrant.Name)
	d.Set("db_cluster_id", dbClusterId)
	d.Set("role_name", roleName)
	d.Set("database_name", databaseName)
	d.Set("read", createGrant.Read)
	d.Set("write", createGrant.Write)

	return resourcePgGrantRead(ctx, d, m)
}

func resourcePgGrantRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Since there's no GetPgGrant method, we can't directly read the grant.
	// We'll verify the cluster exists and that the role and database still exist.
	// The grant itself is managed by the API and we trust the state.
	dbClusterId := d.Get("db_cluster_id").(string)
	dbCluster, err := client.DBaaS().GetDbCluster(ctx, dbClusterId)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil // a deleted db cluster means the pg grant is also deleted
		}
		return diag.FromErr(err)
	}

	// Verify the role still exists
	roleName := d.Get("role_name").(string)
	roleFound := false
	for _, role := range dbCluster.PostgresRoles {
		if strings.EqualFold(role.Name, roleName) {
			roleFound = true
			break
		}
	}
	if !roleFound {
		// Role was deleted, grant is likely gone too
		d.SetId("")
		return nil
	}

	// Verify the database still exists
	databaseName := d.Get("database_name").(string)
	databaseFound := false
	for _, database := range dbCluster.PostgresDatabases {
		if strings.EqualFold(database.Name, databaseName) {
			databaseFound = true
			break
		}
	}
	if !databaseFound {
		// Database was deleted, grant is likely gone too
		d.SetId("")
		return nil
	}

	// Grant is assumed to exist if role and database exist
	// We can't verify the actual grant without a GetPgGrant method
	return nil
}

func resourcePgGrantUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	dbClusterId := d.Get("db_cluster_id").(string)
	dbCluster, err := client.DBaaS().GetDbCluster(ctx, dbClusterId)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil // a deleted db cluster means the pg grant is also deleted
		}
		return diag.FromErr(err)
	}

	grantName := d.Get("name").(string)
	updateGrant := dbaas.UpdatePgGrantRequest{
		Read:  convert.Ptr(d.Get("read").(bool)),
		Write: convert.Ptr(d.Get("write").(bool)),
	}

	err = client.DBaaS().UpdatePgGrant(ctx, dbCluster.Identity, grantName, updateGrant)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil // grant was deleted
		}
		return diag.FromErr(fmt.Errorf("error updating pg grant: %w", err))
	}

	return resourcePgGrantRead(ctx, d, m)
}

func resourcePgGrantDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
				return nil // a deleted db cluster means the pg grant is also deleted
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

	grantName := d.Get("name").(string)
	err = client.DBaaS().DeletePgGrant(ctx, dbCluster.Identity, grantName)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil // grant was already deleted
		}
		return diag.FromErr(fmt.Errorf("error deleting pg grant: %w", err))
	}

	d.SetId("")
	return nil
}
