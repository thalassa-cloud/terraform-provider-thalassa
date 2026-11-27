package dbaas

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/thalassa-cloud/client-go/dbaas"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
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
				ValidateFunc: func(val any, key string) (warns []string, errs []error) {
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
			"extensions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "PostgreSQL extensions to install on the database",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the extension",
						},
						"version": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Version of the extension. When empty, the image default is used.",
						},
					},
				},
			},
			"auto_update_extensions": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.StringInSlice([]string{"none", "customerManaged", "all"}, false),
				Description:  "Automatic extension upgrade policy. Must be none, customerManaged, or all.",
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

func resourcePgDatabaseCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	dbClusterId := d.Get("db_cluster_id").(string)

	dbCluster, err := waitForReadyDbCluster(ctx, client, dbClusterId)
	if err != nil {
		if tcclient.IsNotFound(err) {
			return diag.FromErr(fmt.Errorf("db cluster not found: %w", err))
		}
		return diag.FromErr(fmt.Errorf("error getting db cluster: %w", err))
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
			_ = d.Set("name", database.Name)
			_ = d.Set("db_cluster_id", dbClusterId)
			_ = d.Set("owner_role_id", ownerRoleId)
			_ = d.Set("connection_limit", database.ConnectionLimit)
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
	if extensions := convertPgDatabaseExtensions(d.Get("extensions")); extensions != nil {
		createDatabase.Extensions = extensions
	}
	if autoUpdateExtensions, ok := d.GetOk("auto_update_extensions"); ok {
		policy := dbaas.DbClusterPostgresDatabaseAutoUpdateExtensions(autoUpdateExtensions.(string))
		createDatabase.AutoUpdateExtensions = policy
	}

	_, err = client.DBaaS().CreatePgDatabase(ctx, dbCluster.Identity, createDatabase)
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
			_ = d.Set("name", database.Name)
			_ = d.Set("db_cluster_id", dbClusterId)
			_ = d.Set("owner_role_id", ownerRoleId)
			_ = d.Set("connection_limit", database.ConnectionLimit)
			return resourcePgDatabaseRead(ctx, d, m)
		}
	}

	return diag.FromErr(fmt.Errorf("database not found"))
}

func resourcePgDatabaseRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {

	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	dbClusterId := d.Get("db_cluster_id").(string)

	dbCluster, err := waitForReadyDbCluster(ctx, client, dbClusterId)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil // a deleted db cluster means the pg database is also deleted
		}
		return diag.FromErr(err)
	}

	ownerRoleId := d.Get("owner_role_id").(string)

	for _, database := range dbCluster.PostgresDatabases {
		if strings.EqualFold(database.Name, d.Get("name").(string)) {
			return setPgDatabaseState(d, dbClusterId, ownerRoleId, database)
		}
	}
	return diag.FromErr(fmt.Errorf("database not found"))
}

func resourcePgDatabaseUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	dbClusterId := d.Get("db_cluster_id").(string)

	dbCluster, err := waitForReadyDbCluster(ctx, client, dbClusterId)
	if err != nil {
		if tcclient.IsNotFound(err) {
			return diag.FromErr(fmt.Errorf("db cluster not found: %w", err))
		}
		return diag.FromErr(fmt.Errorf("error getting db cluster: %w", err))
	}

	connectionLimit := d.Get("connection_limit").(int)

	updateDatabase := dbaas.UpdatePgDatabaseRequest{
		ConnectionLimit: &connectionLimit,
	}

	if allowConnections, ok := d.GetOk("allow_connections"); ok {
		allowConnectionsBool := allowConnections.(bool)
		updateDatabase.AllowConnections = &allowConnectionsBool
	}
	if d.HasChange("extensions") {
		updateDatabase.Extensions = convertPgDatabaseExtensions(d.Get("extensions"))
		if updateDatabase.Extensions == nil {
			updateDatabase.Extensions = &dbaas.PgDatabaseExtensions{Extensions: []dbaas.DatabasePostgresExtension{}}
		}
	}
	if d.HasChange("auto_update_extensions") {
		policy := dbaas.DbClusterPostgresDatabaseAutoUpdateExtensions(d.Get("auto_update_extensions").(string))
		updateDatabase.AutoUpdateExtensions = convert.Ptr(policy)
	}

	_, err = client.DBaaS().UpdatePgDatabase(ctx, dbCluster.Identity, d.Get("id").(string), updateDatabase)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating pg database: %w", err))
	}

	ownerRoleId := d.Get("owner_role_id").(string)

	for _, database := range dbCluster.PostgresDatabases {
		if strings.EqualFold(database.Name, d.Get("name").(string)) {
			d.SetId(database.Identity)
			_ = d.Set("name", database.Name)
			_ = d.Set("db_cluster_id", dbClusterId)
			_ = d.Set("owner_role_id", ownerRoleId)
			_ = d.Set("connection_limit", database.ConnectionLimit)
			return resourcePgDatabaseRead(ctx, d, m)
		}
	}
	return diag.FromErr(fmt.Errorf("database not found"))
}

func resourcePgDatabaseDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	dbClusterId := d.Get("db_cluster_id").(string)

	dbCluster, err := waitForReadyDbCluster(ctx, client, dbClusterId)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil // a deleted db cluster means the pg database is also deleted
		}
		return diag.FromErr(fmt.Errorf("error getting db cluster: %w", err))
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

func convertPgDatabaseExtensions(raw any) *dbaas.PgDatabaseExtensions {
	if raw == nil {
		return nil
	}
	list, ok := raw.([]any)
	if !ok || len(list) == 0 {
		return nil
	}

	extensions := make([]dbaas.DatabasePostgresExtension, 0, len(list))
	for _, item := range list {
		if item == nil {
			continue
		}
		cfg, ok := item.(map[string]any)
		if !ok || cfg == nil {
			continue
		}
		ext := dbaas.DatabasePostgresExtension{}
		if name, exists := cfg["name"]; exists && name != nil {
			if strVal, ok := name.(string); ok {
				ext.Name = strVal
			}
		}
		if version, exists := cfg["version"]; exists && version != nil {
			if strVal, ok := version.(string); ok {
				ext.Version = strVal
			}
		}
		if ext.Name == "" {
			continue
		}
		extensions = append(extensions, ext)
	}
	if len(extensions) == 0 {
		return nil
	}
	return &dbaas.PgDatabaseExtensions{Extensions: extensions}
}

func flattenPgDatabaseExtensions(extensions *dbaas.PgDatabaseExtensions) []any {
	if extensions == nil || len(extensions.Extensions) == 0 {
		return nil
	}
	result := make([]any, 0, len(extensions.Extensions))
	for _, ext := range extensions.Extensions {
		result = append(result, map[string]any{
			"name":    ext.Name,
			"version": ext.Version,
		})
	}
	return result
}

func flattenPgDatabaseExtensionStatus(statuses []dbaas.DatabasePostgresExtensionStatus) []any {
	if len(statuses) == 0 {
		return nil
	}
	result := make([]any, 0, len(statuses))
	for _, status := range statuses {
		result = append(result, map[string]any{
			"name":               status.Name,
			"installed_version":  status.InstalledVersion,
			"default_version":    status.DefaultVersion,
			"available_versions": status.AvailableVersions,
			"upgrade_available":  status.UpgradeAvailable,
		})
	}
	return result
}

func setPgDatabaseState(d *schema.ResourceData, dbClusterId, ownerRoleId string, database dbaas.DbClusterPostgresDatabase) diag.Diagnostics {
	d.SetId(database.Identity)
	_ = d.Set("name", database.Name)
	_ = d.Set("db_cluster_id", dbClusterId)
	_ = d.Set("owner_role_id", ownerRoleId)
	_ = d.Set("connection_limit", database.ConnectionLimit)
	_ = d.Set("extensions", flattenPgDatabaseExtensions(database.Extensions))
	_ = d.Set("auto_update_extensions", string(database.AutoUpdateExtensions))
	_ = d.Set("extension_status", flattenPgDatabaseExtensionStatus(database.ExtensionStatus))
	return nil
}
