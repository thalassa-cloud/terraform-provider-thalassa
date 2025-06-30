package dbaas

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thalassa-cloud/client-go/dbaas/dbaasalphav1"
	"github.com/thalassa-cloud/client-go/iaas"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func resourceDbCluster() *schema.Resource {
	return &schema.Resource{
		Description:   "Create an DB Cluster",
		CreateContext: resourceDbClusterCreate,
		ReadContext:   resourceDbClusterRead,
		UpdateContext: resourceDbClusterUpdate,
		DeleteContext: resourceDbClusterDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the DB Cluster",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					if val == "" {
						errs = append(errs, fmt.Errorf("name is required"))
					}
					warns = []string{}
					return
				},
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the DB Cluster",
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Labels of the DB Cluster",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Annotations of the DB Cluster",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Subnet of the DB Cluster",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					if val == "" {
						errs = append(errs, fmt.Errorf("subnet is required"))
					}
					warns = []string{}
					return
				},
			},
			"database_instance_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Database instance type of the DB Cluster",
			},
			"replicas": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Number of instances in the cluster",
			},
			"engine": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Database engine of the cluster",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					// Must be one of: postgres; PostgresInitDb is only supported for Postgres clusters
					if val != "postgres" {
						errs = append(errs, fmt.Errorf("invalid engine: %s", val))
					}
					warns = []string{}
					return
				},
			},
			"engine_version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Version of the database engine",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					if val == "" {
						errs = append(errs, fmt.Errorf("engine version is required"))
					}
					warns = []string{}
					return
				},
			},
			"parameters": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Map of parameter name to database engine specific parameter value",
			},
			"allocated_storage": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Amount of storage allocated to the cluster in GB",
			},
			"volume_type_class": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Storage type used to determine the size of the cluster storage",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					if val == "" {
						errs = append(errs, fmt.Errorf("volume type class is required"))
					}
					warns = []string{}
					return
				},
			},
			"auto_minor_version_upgrade": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Flag indicating if the cluster should automatically upgrade to the latest minor version",
			},
			"database_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the database on the cluster",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					if val == "" {
						errs = append(errs, fmt.Errorf("database name is required"))
					}
					warns = []string{}
					return
				},
			},
			"delete_protection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Flag indicating if the cluster should be protected from deletion",
			},
			"security_groups": {
				Type:        schema.TypeList,
				Optional:    true,
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
			"init_db": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Map of init db parameters",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"restore_from_backup_identity": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Identity of the backup to restore from",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceDbClusterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("id").(string)

	// Safe type assertions with nil checks
	labels := make(map[string]string)
	if labelsRaw := d.Get("labels"); labelsRaw != nil {
		if labelsMap, ok := labelsRaw.(map[string]interface{}); ok {
			for k, v := range labelsMap {
				if strVal, ok := v.(string); ok {
					labels[k] = strVal
				}
			}
		}
	}

	annotations := make(map[string]string)
	if annotationsRaw := d.Get("annotations"); annotationsRaw != nil {
		if annotationsMap, ok := annotationsRaw.(map[string]interface{}); ok {
			for k, v := range annotationsMap {
				if strVal, ok := v.(string); ok {
					annotations[k] = strVal
				}
			}
		}
	}

	parameters := make(map[string]string)
	if parametersRaw := d.Get("parameters"); parametersRaw != nil {
		if parametersMap, ok := parametersRaw.(map[string]interface{}); ok {
			for k, v := range parametersMap {
				if strVal, ok := v.(string); ok {
					parameters[k] = strVal
				}
			}
		}
	}

	subnetId := d.Get("subnet_id").(string)
	subnet, err := client.IaaS().GetSubnet(ctx, subnetId)
	if err != nil {
		return diag.FromErr(fmt.Errorf("subnet not found: %w", err))
	}

	databaseInstanceType := d.Get("database_instance_type").(string)
	databaseInstanceTypes, err := client.DbaaSAlphaV1().ListDatabaseInstanceTypes(ctx, &dbaasalphav1.ListDatabaseInstanceTypesRequest{})
	if err != nil {
		return diag.FromErr(fmt.Errorf("database instance type not found: %w", err))
	}
	foundInstanceType := false
	for _, instanceType := range databaseInstanceTypes {
		if instanceType.Name == databaseInstanceType {
			databaseInstanceType = instanceType.Identity
			foundInstanceType = true
			break
		}
	}
	if !foundInstanceType {
		return diag.FromErr(fmt.Errorf("database instance type not found: %s", databaseInstanceType))
	}

	engine := d.Get("engine").(string)
	if engine == "" {
		return diag.FromErr(fmt.Errorf("engine is required"))
	}

	engineVersion := d.Get("engine_version").(string)
	tflog.Info(ctx, "engine", map[string]interface{}{
		"engine":        engine,
		"engineVersion": engineVersion,
	})
	engineVersions, err := client.DbaaSAlphaV1().ListEngineVersions(ctx, dbaasalphav1.DbClusterDatabaseEngine(engine), &dbaasalphav1.ListEngineVersionsRequest{})
	if err != nil {
		return diag.FromErr(fmt.Errorf("engine version not found: %w", err))
	}
	foundEngineVersion := false
	for _, version := range engineVersions {
		if version.EngineVersion == engineVersion {
			foundEngineVersion = true
		}
	}
	if !foundEngineVersion {
		return diag.FromErr(fmt.Errorf("engine version not found: %s", engineVersion))
	}

	createDbCluster := dbaasalphav1.CreateDbClusterRequest{
		Name:                         d.Get("name").(string),
		Description:                  d.Get("description").(string),
		Labels:                       dbaasalphav1.Labels(labels),
		Annotations:                  dbaasalphav1.Annotations(annotations),
		SubnetIdentity:               subnet.Identity,
		DeleteProtection:             d.Get("delete_protection").(bool),
		Engine:                       dbaasalphav1.DbClusterDatabaseEngine(d.Get("engine").(string)),
		EngineVersion:                engineVersion,
		Parameters:                   parameters,
		AllocatedStorage:             uint64(d.Get("allocated_storage").(int)),
		DatabaseInstanceTypeIdentity: databaseInstanceType,
		AutoMinorVersionUpgrade:      d.Get("auto_minor_version_upgrade").(bool),
		Instances:                    d.Get("replicas").(int),
	}

	if volumeTypeClass := d.Get("volume_type_class"); volumeTypeClass != nil && volumeTypeClass != "" {
		volumeTypeClasses, err := client.IaaS().ListVolumeTypes(ctx, &iaas.ListVolumeTypesRequest{})
		if err != nil {
			return diag.FromErr(err)
		}
		for _, class := range volumeTypeClasses {
			if strings.ToLower(class.Name) == volumeTypeClass {
				createDbCluster.VolumeTypeClassIdentity = class.Identity
				break
			}
		}
	}

	// Safe handling of security groups
	if securityGroupsRaw := d.Get("security_groups"); securityGroupsRaw != nil {
		if securityGroupsList, ok := securityGroupsRaw.([]interface{}); ok {
			securityGroups := make([]string, len(securityGroupsList))
			for i, sg := range securityGroupsList {
				if strVal, ok := sg.(string); ok {
					securityGroups[i] = strVal
				}
			}
			createDbCluster.SecurityGroupAttachments = securityGroups
		}
	}

	// Safe handling of init_db
	if initDbRaw := d.Get("init_db"); initDbRaw != nil && len(initDbRaw.(map[string]interface{})) > 0 {
		if initDbMap, ok := initDbRaw.(map[string]interface{}); ok {
			// Safe extraction of init_db values with defaults
			dataChecksums := false
			if val, exists := initDbMap["data_checksums"]; exists && val != nil {
				if strVal, ok := val.(string); ok {
					dataChecksums = strVal == "true"
				}
			}

			encoding := ""
			if val, exists := initDbMap["encoding"]; exists && val != nil {
				if strVal, ok := val.(string); ok {
					encoding = strVal
				}
			}

			locale := ""
			if val, exists := initDbMap["locale"]; exists && val != nil {
				if strVal, ok := val.(string); ok {
					locale = strVal
				}
			}

			localeProvider := ""
			if val, exists := initDbMap["locale_provider"]; exists && val != nil {
				if strVal, ok := val.(string); ok {
					localeProvider = strVal
				}
			}

			lcCollate := ""
			if val, exists := initDbMap["lc_collate"]; exists && val != nil {
				if strVal, ok := val.(string); ok {
					lcCollate = strVal
				}
			}

			lcCtype := ""
			if val, exists := initDbMap["lc_ctype"]; exists && val != nil {
				if strVal, ok := val.(string); ok {
					lcCtype = strVal
				}
			}

			icuLocale := ""
			if val, exists := initDbMap["icu_locale"]; exists && val != nil {
				if strVal, ok := val.(string); ok {
					icuLocale = strVal
				}
			}

			icuRules := ""
			if val, exists := initDbMap["icu_rules"]; exists && val != nil {
				if strVal, ok := val.(string); ok {
					icuRules = strVal
				}
			}

			builtinLocale := ""
			if val, exists := initDbMap["builtin_locale"]; exists && val != nil {
				if strVal, ok := val.(string); ok {
					builtinLocale = strVal
				}
			}

			walSegmentSize := 0
			if val, exists := initDbMap["wal_segment_size"]; exists && val != nil {
				if strVal, ok := val.(string); ok {
					if intVal, err := strconv.Atoi(strVal); err == nil {
						walSegmentSize = intVal
					}
				}
			}

			createDbCluster.PostgresInitDb = &dbaasalphav1.PostgresInitDb{
				DataChecksums:  dataChecksums,
				Encoding:       encoding,
				Locale:         locale,
				LocaleProvider: localeProvider,
				LcCollate:      lcCollate,
				LcCtype:        lcCtype,
				IcuLocale:      icuLocale,
				IcuRules:       icuRules,
				BuiltinLocale:  builtinLocale,
				WalSegmentSize: walSegmentSize,
			}
		}
	}

	if restoreFromBackupIdentity := d.Get("restore_from_backup_identity"); restoreFromBackupIdentity != nil && restoreFromBackupIdentity != "" {
		if strVal, ok := restoreFromBackupIdentity.(string); ok {
			createDbCluster.RestoreFromBackupIdentity = convert.Ptr(strVal)
		}
	}

	if d.Get("parameters") != nil && len(d.Get("parameters").(map[string]interface{})) > 0 {
		createDbCluster.Parameters = parameters
	}

	// Create the DbCluster
	createdDbCluster, err := client.DbaaSAlphaV1().CreateDbCluster(ctx, createDbCluster)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set the id to the name of the cluster
	d.SetId(createdDbCluster.Identity)

	var dbCluster *dbaasalphav1.DbCluster
	// Wait for the cluster to be ready
	for {
		dbclusters, err := client.DbaaSAlphaV1().ListDbClusters(ctx, &dbaasalphav1.ListDbClustersRequest{})
		if err != nil {

			return diag.FromErr(err)
		}

		var foundCluster *dbaasalphav1.DbCluster
		for _, dbCluster := range dbclusters {
			if dbCluster.Identity == id || dbCluster.Name == dbCluster.Name {
				foundCluster = &dbCluster
				break
			}
		}

		if foundCluster != nil && foundCluster.Status == dbaasalphav1.DbClusterStatusReady {
			dbCluster = foundCluster
			break
		}

		time.Sleep(1 * time.Second)
	}

	d.SetId(dbCluster.Identity)
	d.Set("name", dbCluster.Name)
	d.Set("description", dbCluster.Description)
	d.Set("labels", dbCluster.Labels)
	d.Set("annotations", dbCluster.Annotations)
	d.Set("delete_protection", dbCluster.DeleteProtection)
	d.Set("replicas", dbCluster.Replicas)
	d.Set("engine", dbCluster.Engine)
	d.Set("engine_version", dbCluster.EngineVersion)
	d.Set("parameters", dbCluster.Parameters)
	d.Set("allocated_storage", dbCluster.AllocatedStorage)
	d.Set("auto_minor_version_upgrade", dbCluster.AutoMinorVersionUpgrade)
	d.Set("status", dbCluster.Status)
	d.Set("endpoint_ipv4", dbCluster.EndpointIpv4)
	d.Set("endpoint_ipv6", dbCluster.EndpointIpv6)
	d.Set("port", dbCluster.Port)

	// Handle optional fields
	if dbCluster.DatabaseName != nil {
		d.Set("database_name", *dbCluster.DatabaseName)
	}
	if dbCluster.Subnet != nil {
		d.Set("subnet_id", dbCluster.Subnet.Identity)
	}
	if dbCluster.DatabaseInstanceType != nil {
		d.Set("database_instance_type", dbCluster.DatabaseInstanceType.Identity)
	}
	if dbCluster.VolumeTypeClass != nil {
		d.Set("volume_type_class", dbCluster.VolumeTypeClass.Identity)
	}
	if dbCluster.SecurityGroups != nil {
		d.Set("security_groups", dbCluster.SecurityGroups)
	}

	return resourceDbClusterRead(ctx, d, m)
}

func resourceDbClusterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	slug := d.Get("id").(string)
	var DbCluster *dbaasalphav1.DbCluster
	dbClusters, err := client.DbaaSAlphaV1().ListDbClusters(ctx, &dbaasalphav1.ListDbClustersRequest{})
	if err != nil {
		return diag.FromErr(err)
	}

	for _, dbCluster := range dbClusters {
		if dbCluster.Identity == slug {
			DbCluster = &dbCluster
			break
		}
	}

	if DbCluster == nil {
		return diag.FromErr(fmt.Errorf("DbCluster was not found"))
	}

	d.SetId(DbCluster.Identity)
	d.Set("name", DbCluster.Name)
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

func resourceDbClusterUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("id").(string)

	// Safe type assertions with nil checks
	labels := make(map[string]string)
	if labelsRaw := d.Get("labels"); labelsRaw != nil {
		if labelsMap, ok := labelsRaw.(map[string]interface{}); ok {
			for k, v := range labelsMap {
				if strVal, ok := v.(string); ok {
					labels[k] = strVal
				}
			}
		}
	}

	annotations := make(map[string]string)
	if annotationsRaw := d.Get("annotations"); annotationsRaw != nil {
		if annotationsMap, ok := annotationsRaw.(map[string]interface{}); ok {
			for k, v := range annotationsMap {
				if strVal, ok := v.(string); ok {
					annotations[k] = strVal
				}
			}
		}
	}

	parameters := make(map[string]string)
	if parametersRaw := d.Get("parameters"); parametersRaw != nil {
		if parametersMap, ok := parametersRaw.(map[string]interface{}); ok {
			for k, v := range parametersMap {
				if strVal, ok := v.(string); ok {
					parameters[k] = strVal
				}
			}
		}
	}

	// Safe handling of security groups
	securityGroupAttachments := []string{}
	if securityGroupsRaw := d.Get("security_groups"); securityGroupsRaw != nil {
		if securityGroupsList, ok := securityGroupsRaw.([]interface{}); ok {
			securityGroupAttachments = make([]string, len(securityGroupsList))
			for i, sg := range securityGroupsList {
				if strVal, ok := sg.(string); ok {
					securityGroupAttachments[i] = strVal
				}
			}
		}
	}

	updateDbCluster := dbaasalphav1.UpdateDbClusterRequest{
		Name:                         d.Get("name").(string),
		Description:                  d.Get("description").(string),
		Labels:                       dbaasalphav1.Labels(labels),
		Annotations:                  dbaasalphav1.Annotations(annotations),
		SecurityGroupAttachments:     securityGroupAttachments,
		DeleteProtection:             d.Get("delete_protection").(bool),
		EngineVersion:                convert.Ptr(d.Get("engine_version").(string)),
		Parameters:                   parameters,
		AllocatedStorage:             uint64(d.Get("allocated_storage").(int)),
		AutoMinorVersionUpgrade:      d.Get("auto_minor_version_upgrade").(bool),
		Replicas:                     d.Get("replicas").(int),
		DatabaseInstanceTypeIdentity: convert.Ptr(d.Get("database_instance_type").(string)),
	}

	// Safe handling of optional database_name
	if databaseName := d.Get("database_name"); databaseName != nil && databaseName != "" {
		if strVal, ok := databaseName.(string); ok {
			updateDbCluster.DatabaseName = convert.Ptr(strVal)
		}
	}

	updatedDbCluster, err := client.DbaaSAlphaV1().UpdateDbCluster(ctx, id, updateDbCluster)
	if err != nil {
		return diag.FromErr(err)
	}

	for {
		updatedDbCluster, err = client.DbaaSAlphaV1().GetDbCluster(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}
		if updatedDbCluster.Status == dbaasalphav1.DbClusterStatusReady {
			break
		}
		time.Sleep(1 * time.Second)
	}

	if updatedDbCluster == nil {
		return diag.FromErr(fmt.Errorf("DbCluster was not found"))
	}

	d.SetId(updatedDbCluster.Identity)

	return resourceDbClusterRead(ctx, d, m)
}

func resourceDbClusterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("id").(string)
	// Get the cluster
	dbCluster, err := client.DbaaSAlphaV1().GetDbCluster(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DbaaSAlphaV1().DeleteDbCluster(ctx, dbCluster.Identity)
	if err != nil {
		return diag.FromErr(err)
	}

	for {
		dbClusters, err := client.DbaaSAlphaV1().ListDbClusters(ctx, &dbaasalphav1.ListDbClustersRequest{})
		if err != nil {
			return diag.FromErr(err)
		}
		var foundCluster *dbaasalphav1.DbCluster
		for _, dbCluster := range dbClusters {
			if dbCluster.Identity == id {
				foundCluster = &dbCluster
				break
			}
		}
		if foundCluster != nil && foundCluster.Status == dbaasalphav1.DbClusterStatusDeleted {
			break
		}
		if foundCluster == nil {
			// Assume the cluster is deleted
			break
		}
		time.Sleep(1 * time.Second)
	}

	d.SetId("")
	return nil
}
