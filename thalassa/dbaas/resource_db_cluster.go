package dbaas

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/thalassa-cloud/client-go/dbaas"
	"github.com/thalassa-cloud/client-go/iaas"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
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
		CustomizeDiff: func(_ context.Context, diff *schema.ResourceDiff, _ any) error {
			raw, ok := diff.GetOk("restore_recovery_target")
			if !ok {
				return nil
			}

			blocks, ok := raw.([]any)
			if !ok || len(blocks) == 0 {
				return nil
			}

			block, ok := blocks[0].(map[string]any)
			if !ok {
				return nil
			}

			if err := validateRestoreRecoveryTargetBlock(block); err != nil {
				return err
			}

			restoreFromBackupID, hasRestoreFromBackupID := diff.GetOk("restore_from_backup_id")
			if !hasRestoreFromBackupID || strings.TrimSpace(restoreFromBackupID.(string)) == "" {
				return fmt.Errorf(
					"restore_recovery_target can only be used when restore_from_backup_id is set",
				)
			}

			return nil
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},
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
				ValidateFunc: func(val any, key string) (warns []string, errs []error) {
					if val == "" {
						errs = append(errs, fmt.Errorf("name is required"))
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
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the DB Cluster",
			},
			"labels": {
				Type:        schema.TypeMap,
				Default:     make(map[string]string),
				Optional:    true,
				Description: "Labels of the DB Cluster",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Default:     make(map[string]string),
				Optional:    true,
				Description: "Annotations of the DB Cluster",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Subnet of the DB Cluster",
				ValidateFunc: func(val any, key string) (warns []string, errs []error) {
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
				ValidateFunc: func(val any, key string) (warns []string, errs []error) {
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
				ValidateFunc: func(val any, key string) (warns []string, errs []error) {
					if val == "" {
						errs = append(errs, fmt.Errorf("engine version is required"))
					}
					warns = []string{}
					return
				},
			},
			"parameters": {
				Type:        schema.TypeMap,
				Default:     make(map[string]string),
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
				Required:    true,
				Description: "Storage type used to determine the size of the cluster storage",
				ValidateFunc: func(val any, key string) (warns []string, errs []error) {
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
				Default:     make(map[string]string),
				Optional:    true,
				Description: "Map of init db parameters",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"restore_recovery_target": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Recovery target for Point-In-Time Recovery (PITR). Only used when restore_from_backup_id is specified.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_time": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  barmanTargetTimeDescription,
							ValidateFunc: validateBarmanTargetTimeString,
						},
						"target_lsn": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "Log Sequence Number to restore to. Example: '0/1234567'",
							ValidateFunc: validateTargetLSNString,
						},
					},
				},
			},
			"auto_upgrade_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Auto upgrade policy for the cluster. Options: 'none', 'latest-version', 'latest-stable', 'latest-patch', 'latest-minor', 'latest-major'",
				ValidateFunc: func(val any, key string) (warns []string, errs []error) {
					validPolicies := []string{"none", "latest-version", "latest-stable", "latest-patch", "latest-minor", "latest-major"}
					policy := val.(string)
					valid := slices.Contains(validPolicies, policy)
					if !valid {
						errs = append(errs, fmt.Errorf("auto_upgrade_policy must be one of: %v", validPolicies))
					}
					warns = []string{}
					return
				},
			},
			"maintenance_day": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Day of the week for the maintenance window. 0 is Sunday, 6 is Saturday",
				ValidateFunc: func(val any, key string) (warns []string, errs []error) {
					day := val.(int)
					if day < 0 || day > 6 {
						errs = append(errs, fmt.Errorf("maintenance_day must be between 0 (Sunday) and 6 (Saturday)"))
					}
					warns = []string{}
					return
				},
			},
			"maintenance_start_at": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Start time of the maintenance window on the maintenance day in UTC. 0 is 00:00, 23 is 23:00",
				ValidateFunc: func(val any, key string) (warns []string, errs []error) {
					hour := val.(int)
					if hour < 0 || hour > 23 {
						errs = append(errs, fmt.Errorf("maintenance_start_at must be between 0 and 23"))
					}
					warns = []string{}
					return
				},
			},
			"restore_from_backup_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Identity of the DB object store used for barman backups (optional). Ignored if provision_db_object_store is true.",
			},
			"provision_db_object_store": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Flag to indicate if the DB object store should be provisioned for the cluster. If true, restore_from_backup_id will be ignored.",
			},
			"create_backup_before_destroy": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to create a backup before destroying the cluster. Only applies when the cluster is in ready status.",
			},
			"create_backup_before_destroy_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      30,
				Description:  "The timeout in minutes to wait for the pre-destroy backup to complete. Only used when create_backup_before_destroy is true.",
				ValidateFunc: validate.IntAtLeast(1),
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceDbClusterCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics { //nolint:gocyclo // create validates many optional DB cluster fields
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Safe type assertions with nil checks
	labels := make(map[string]string)
	if labelsRaw := d.Get("labels"); labelsRaw != nil {
		if labelsMap, ok := labelsRaw.(map[string]any); ok {
			for k, v := range labelsMap {
				if strVal, ok := v.(string); ok {
					labels[k] = strVal
				}
			}
		}
	}

	annotations := make(map[string]string)
	if annotationsRaw := d.Get("annotations"); annotationsRaw != nil {
		annotations = convert.ConvertToMap(annotationsRaw)
	}

	parameters := make(map[string]string)
	if parametersRaw := d.Get("parameters"); parametersRaw != nil {
		parameters = convert.ConvertToMap(parametersRaw)
	}

	subnetId := d.Get("subnet_id").(string)
	subnet, err := client.IaaS().GetSubnet(ctx, subnetId)
	if err != nil {
		if tcclient.IsNotFound(err) {
			return diag.FromErr(fmt.Errorf("subnet not found: %w", err))
		}
		return diag.FromErr(fmt.Errorf("failed to get subnet: %w", err))
	}

	databaseInstanceType := d.Get("database_instance_type").(string)
	databaseInstanceTypes, err := client.DBaaS().ListDatabaseInstanceTypes(ctx, &dbaas.ListDatabaseInstanceTypesRequest{})
	if err != nil {
		return diag.FromErr(fmt.Errorf("database instance type not found: %w", err))
	}
	foundInstanceType := false
	for _, instanceType := range databaseInstanceTypes {
		if strings.EqualFold(instanceType.Name, databaseInstanceType) || strings.EqualFold(instanceType.Identity, databaseInstanceType) || strings.EqualFold(instanceType.Slug, databaseInstanceType) {
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
		return diag.FromErr(fmt.Errorf("engine is required. Must be one of: 'postgres'"))
	}

	engineVersion := d.Get("engine_version").(string)
	tflog.Info(ctx, "engine", map[string]any{
		"engine":        engine,
		"engineVersion": engineVersion,
	})
	engineVersions, err := client.DBaaS().ListEngineVersions(ctx, dbaas.DbClusterDatabaseEngine(engine), &dbaas.ListEngineVersionsRequest{})
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

	createDbCluster := dbaas.CreateDbClusterRequest{
		Name:                         d.Get("name").(string),
		Description:                  d.Get("description").(string),
		Labels:                       dbaas.Labels(labels),
		Annotations:                  dbaas.Annotations(annotations),
		SubnetIdentity:               subnet.Identity,
		DeleteProtection:             d.Get("delete_protection").(bool),
		Engine:                       dbaas.DbClusterDatabaseEngine(engine),
		EngineVersion:                engineVersion,
		Parameters:                   parameters,
		AllocatedStorage:             uint64(d.Get("allocated_storage").(int)),
		DatabaseInstanceTypeIdentity: databaseInstanceType,
		AutoMinorVersionUpgrade:      d.Get("auto_minor_version_upgrade").(bool),
		Replicas:                     d.Get("replicas").(int),
	}

	foundVolumeTypeClass := false
	if volumeTypeClass := d.Get("volume_type_class"); volumeTypeClass != nil && volumeTypeClass != "" {
		volumeTypeClasses, err := client.IaaS().ListVolumeTypes(ctx, &iaas.ListVolumeTypesRequest{})
		if err != nil {
			if tcclient.IsNotFound(err) {
				return diag.FromErr(fmt.Errorf("volume type class not found: %w", err))
			}
			return diag.FromErr(fmt.Errorf("failed to get volume type class: %w", err))
		}
		volumeTypeClassStr := volumeTypeClass.(string)
		for _, class := range volumeTypeClasses {
			if strings.EqualFold(class.Name, volumeTypeClassStr) || strings.EqualFold(class.Identity, volumeTypeClassStr) {
				createDbCluster.VolumeTypeClassIdentity = class.Identity
				foundVolumeTypeClass = true
				break
			}
		}
		if !foundVolumeTypeClass {
			availableVolumeTypeClasses := []string{}
			for _, class := range volumeTypeClasses {
				availableVolumeTypeClasses = append(availableVolumeTypeClasses, class.Name)
			}
			return diag.FromErr(fmt.Errorf("volume type class not found: %s. Available volume type classes: %s", volumeTypeClassStr, strings.Join(availableVolumeTypeClasses, ", ")))
		}
	}

	if !foundVolumeTypeClass {
		return diag.FromErr(fmt.Errorf("volume type class not found: %s", d.Get("volume_type_class").(string)))
	}

	// Safe handling of security groups
	if securityGroupsRaw := d.Get("security_groups"); securityGroupsRaw != nil {
		createDbCluster.SecurityGroupAttachments = convert.ConvertToStringSlice(securityGroupsRaw)
	}

	// Safe handling of init_db
	if initDbRaw := d.Get("init_db"); initDbRaw != nil && len(initDbRaw.(map[string]any)) > 0 {
		if initDbMap, ok := initDbRaw.(map[string]any); ok {
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

			createDbCluster.PostgresInitDb = &dbaas.PostgresInitDb{
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

	if restoreFromBackupIdentity := d.Get("restore_from_backup_id"); restoreFromBackupIdentity != nil && restoreFromBackupIdentity != "" {
		if strVal, ok := restoreFromBackupIdentity.(string); ok {
			createDbCluster.RestoreFromBackupIdentity = convert.Ptr(strVal)
		}
	}

	// Handle restore recovery target
	if restoreRecoveryTargetRaw := d.Get("restore_recovery_target"); restoreRecoveryTargetRaw != nil {
		if restoreRecoveryTargetList, ok := restoreRecoveryTargetRaw.([]any); ok && len(restoreRecoveryTargetList) > 0 {
			if restoreRecoveryTargetMap, ok := restoreRecoveryTargetList[0].(map[string]any); ok {
				recoveryTarget := &dbaas.RestoreRecoveryTarget{}
				if targetTime, exists := restoreRecoveryTargetMap["target_time"]; exists && targetTime != nil {
					if strVal, ok := targetTime.(string); ok && strings.TrimSpace(strVal) != "" {
						parsedTargetTime, err := parseRecoveryTargetTime(strVal)
						if err != nil {
							return diag.FromErr(err)
						}
						recoveryTarget.TargetTime = convert.Ptr(formatBarmanTargetTime(parsedTargetTime))
					}
				}
				if targetLSN, exists := restoreRecoveryTargetMap["target_lsn"]; exists && targetLSN != nil {
					if strVal, ok := targetLSN.(string); ok && strVal != "" {
						recoveryTarget.TargetLSN = convert.Ptr(strVal)
					}
				}
				if recoveryTarget.TargetTime != nil || recoveryTarget.TargetLSN != nil {
					createDbCluster.RestoreRecoveryTarget = recoveryTarget
				}
			}
		}
	}

	// Handle auto upgrade policy
	if autoUpgradePolicy := d.Get("auto_upgrade_policy"); autoUpgradePolicy != nil && autoUpgradePolicy != "" {
		if strVal, ok := autoUpgradePolicy.(string); ok {
			policy := dbaas.DbClusterAutoUpgradePolicy(strVal)
			createDbCluster.AutoUpgradePolicy = &policy
		}
	}

	// Handle maintenance window
	if maintenanceDay := d.Get("maintenance_day"); maintenanceDay != nil {
		if intVal, ok := maintenanceDay.(int); ok {
			day := uint(intVal)
			createDbCluster.MaintenanceDay = &day
		}
	}
	if maintenanceStartAt := d.Get("maintenance_start_at"); maintenanceStartAt != nil {
		if intVal, ok := maintenanceStartAt.(int); ok {
			hour := uint(intVal)
			createDbCluster.MaintenanceStartAt = &hour
		}
	}

	// Handle DB object store
	if dbObjectStoreIdentity := d.Get("restore_from_backup_id"); dbObjectStoreIdentity != nil && dbObjectStoreIdentity != "" {
		if strVal, ok := dbObjectStoreIdentity.(string); ok {
			createDbCluster.DbObjectStoreIdentity = convert.Ptr(strVal)
		}
	}
	if provisionDbObjectStore := d.Get("provision_db_object_store"); provisionDbObjectStore != nil {
		if boolVal, ok := provisionDbObjectStore.(bool); ok {
			createDbCluster.ProvisionDbObjectStore = boolVal
		}
	}

	if d.Get("parameters") != nil && len(d.Get("parameters").(map[string]any)) > 0 {
		createDbCluster.Parameters = parameters
	}

	// Create the DbCluster
	createdDbCluster, err := client.DBaaS().CreateDbCluster(ctx, createDbCluster)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set the id to the name of the cluster
	d.SetId(createdDbCluster.Identity)

	var dbCluster *dbaas.DbCluster
	dbCluster, err = waitForReadyDbCluster(ctx, client, createdDbCluster.Identity)
	if err != nil {
		if tcclient.IsNotFound(err) {
			return diag.FromErr(fmt.Errorf("db cluster not found: %w", err))
		}
		return diag.FromErr(err)
	}

	d.SetId(dbCluster.Identity)
	_ = d.Set("name", dbCluster.Name)
	_ = d.Set("description", dbCluster.Description)
	_ = d.Set("labels", dbCluster.Labels)
	_ = d.Set("annotations", dbCluster.Annotations)
	_ = d.Set("delete_protection", dbCluster.DeleteProtection)
	_ = d.Set("replicas", dbCluster.Replicas)
	_ = d.Set("engine", dbCluster.Engine)
	_ = d.Set("engine_version", dbCluster.EngineVersion)
	_ = d.Set("parameters", dbCluster.Parameters)
	_ = d.Set("allocated_storage", dbCluster.AllocatedStorage)
	_ = d.Set("auto_minor_version_upgrade", dbCluster.AutoMinorVersionUpgrade)
	_ = d.Set("status", dbCluster.Status)
	_ = d.Set("endpoint_ipv4", dbCluster.EndpointIpv4)
	_ = d.Set("endpoint_ipv6", dbCluster.EndpointIpv6)
	_ = d.Set("port", dbCluster.Port)

	if dbCluster.Subnet != nil {
		_ = d.Set("subnet_id", dbCluster.Subnet.Identity)
	}
	if dbCluster.SecurityGroups != nil {
		_ = d.Set("security_groups", dbCluster.SecurityGroups)
	}

	return resourceDbClusterRead(ctx, d, m)
}

func resourceDbClusterRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	identity := d.Get("id").(string)
	var DbCluster *dbaas.DbCluster
	DbCluster, err = client.DBaaS().GetDbCluster(ctx, identity)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("failed to get db cluster: %w", err))
	}

	if DbCluster == nil {
		d.SetId("")
		return nil
	}

	d.SetId(DbCluster.Identity)
	_ = d.Set("name", DbCluster.Name)
	_ = d.Set("description", DbCluster.Description)
	_ = d.Set("labels", DbCluster.Labels)
	_ = d.Set("annotations", DbCluster.Annotations)
	_ = d.Set("delete_protection", DbCluster.DeleteProtection)
	_ = d.Set("replicas", DbCluster.Replicas)
	_ = d.Set("engine", DbCluster.Engine)
	_ = d.Set("engine_version", DbCluster.EngineVersion)
	_ = d.Set("parameters", DbCluster.Parameters)
	_ = d.Set("allocated_storage", DbCluster.AllocatedStorage)
	_ = d.Set("auto_minor_version_upgrade", DbCluster.AutoMinorVersionUpgrade)
	_ = d.Set("status", DbCluster.Status)
	_ = d.Set("endpoint_ipv4", DbCluster.EndpointIpv4)
	_ = d.Set("endpoint_ipv6", DbCluster.EndpointIpv6)
	_ = d.Set("port", DbCluster.Port)

	if DbCluster.Subnet != nil {
		_ = d.Set("subnet_id", DbCluster.Subnet.Identity)
	}
	if DbCluster.DatabaseInstanceType != nil {
		convert.SetReferenceField(d,
			"database_instance_type",
			DbCluster.DatabaseInstanceType.Identity,
			DbCluster.DatabaseInstanceType.Slug,
			DbCluster.DatabaseInstanceType.Name,
		)
	}
	if DbCluster.VolumeTypeClass != nil {
		convert.SetReferenceField(
			d,
			"volume_type_class",
			DbCluster.VolumeTypeClass.Identity,
			"",
			DbCluster.VolumeTypeClass.Name,
		)
	}
	if DbCluster.SecurityGroups != nil {
		securityGroupIds := make([]string, len(DbCluster.SecurityGroups))
		for i, sg := range DbCluster.SecurityGroups {
			securityGroupIds[i] = sg.Identity
		}
		_ = d.Set("security_groups", securityGroupIds)
	}

	return nil
}

func resourceDbClusterUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("id").(string)

	// Safe type assertions with nil checks
	labels := make(map[string]string)
	if labelsRaw := d.Get("labels"); labelsRaw != nil {
		labels = convert.ConvertToMap(labelsRaw)
	}

	annotations := make(map[string]string)
	if annotationsRaw := d.Get("annotations"); annotationsRaw != nil {
		annotations = convert.ConvertToMap(annotationsRaw)
	}

	parameters := make(map[string]string)
	if parametersRaw := d.Get("parameters"); parametersRaw != nil {
		parameters = convert.ConvertToMap(parametersRaw)
	}

	// Safe handling of security groups
	securityGroupAttachments := []string{}
	if securityGroupsRaw := d.Get("security_groups"); securityGroupsRaw != nil {
		securityGroupAttachments = convert.ConvertToStringSlice(securityGroupsRaw)
	}

	updateDbCluster := dbaas.UpdateDbClusterRequest{
		Name:                         d.Get("name").(string),
		Description:                  d.Get("description").(string),
		Labels:                       dbaas.Labels(labels),
		Annotations:                  dbaas.Annotations(annotations),
		SecurityGroupAttachments:     securityGroupAttachments,
		DeleteProtection:             d.Get("delete_protection").(bool),
		EngineVersion:                convert.Ptr(d.Get("engine_version").(string)),
		Parameters:                   parameters,
		AllocatedStorage:             uint64(d.Get("allocated_storage").(int)),
		Replicas:                     d.Get("replicas").(int),
		DatabaseInstanceTypeIdentity: convert.Ptr(d.Get("database_instance_type").(string)),
	}

	// Handle auto upgrade policy
	if autoUpgradePolicy := d.Get("auto_upgrade_policy"); autoUpgradePolicy != nil && autoUpgradePolicy != "" {
		if strVal, ok := autoUpgradePolicy.(string); ok {
			policy := dbaas.DbClusterAutoUpgradePolicy(strVal)
			updateDbCluster.AutoUpgradePolicy = &policy
		}
	}

	// Handle maintenance window
	if maintenanceDay := d.Get("maintenance_day"); maintenanceDay != nil {
		if intVal, ok := maintenanceDay.(int); ok {
			day := uint(intVal)
			updateDbCluster.MaintenanceDay = &day
		}
	}
	if maintenanceStartAt := d.Get("maintenance_start_at"); maintenanceStartAt != nil {
		if intVal, ok := maintenanceStartAt.(int); ok {
			hour := uint(intVal)
			updateDbCluster.MaintenanceStartAt = &hour
		}
	}

	// Handle DB object store
	if dbObjectStoreIdentity := d.Get("restore_from_backup_id"); dbObjectStoreIdentity != nil && dbObjectStoreIdentity != "" {
		if strVal, ok := dbObjectStoreIdentity.(string); ok {
			updateDbCluster.DbObjectStoreIdentity = convert.Ptr(strVal)
		}
	}

	_, err = client.DBaaS().UpdateDbCluster(ctx, id, updateDbCluster)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update db cluster: %w", err))
	}

	_, err = waitForReadyDbCluster(ctx, client, id)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceDbClusterRead(ctx, d, m)
}

func isBackupComplete(backup *dbaas.DbClusterBackup) bool {
	if backup == nil {
		return false
	}
	switch backup.Status {
	case dbaas.ObjectStatusReady, dbaas.ObjectStatus("completed"):
		return true
	default:
		return false
	}
}

func createBackupBeforeDestroy(ctx context.Context, dbaasClient *dbaas.Client, dbClusterIdentity string, timeoutMinutes int) diag.Diagnostics {
	backupName := fmt.Sprintf("terraform-pre-destroy-%d", time.Now().Unix())
	createBackup := dbaas.CreateDbClusterBackupRequest{
		Name:   backupName,
		Labels: dbaas.Labels{"terraform": "true", "purpose": "pre-destroy"},
	}

	backup, err := dbaasClient.CreateDbBackup(ctx, dbClusterIdentity, createBackup)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create pre-destroy backup: %w", err))
	}

	tflog.Info(ctx, "waiting for pre-destroy backup to complete", map[string]any{
		"backup_id": backup.Identity,
	})

	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Duration(timeoutMinutes)*time.Minute)
	defer cancel()

	for {
		select {
		case <-ctxWithTimeout.Done():
			if ctxWithTimeout.Err() == context.Canceled {
				return diag.FromErr(fmt.Errorf("pre-destroy backup wait cancelled"))
			}
			return diag.FromErr(fmt.Errorf("timeout waiting for pre-destroy backup %q to complete (last status: %s)", backup.Identity, backup.Status))
		default:
		}

		backup, err = dbaasClient.GetDbBackup(ctxWithTimeout, backup.Identity)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to retrieve pre-destroy backup: %w", err))
		}

		if isBackupComplete(backup) {
			tflog.Info(ctx, "pre-destroy backup completed", map[string]any{
				"backup_id": backup.Identity,
				"status":    backup.Status,
			})
			return nil
		}
		if backup.Status == dbaas.ObjectStatusFailed {
			return diag.FromErr(fmt.Errorf("pre-destroy backup failed: %s", backup.StatusMessage))
		}

		tflog.Debug(ctx, "pre-destroy backup still in progress", map[string]any{
			"backup_id": backup.Identity,
			"status":    backup.Status,
		})
		time.Sleep(1 * time.Second)
	}
}

func resourceDbClusterDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("id").(string)
	// Get the cluster
	dbCluster, err := client.DBaaS().GetDbCluster(ctx, id)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("failed to retrieve db cluster: %w", err))
	}

	if d.Get("create_backup_before_destroy").(bool) && dbCluster.Status == dbaas.DbClusterStatusReady && dbCluster.DbObjectStore != nil {
		timeout := d.Get("create_backup_before_destroy_timeout").(int)
		if diags := createBackupBeforeDestroy(ctx, client.DBaaS(), dbCluster.Identity, timeout); diags != nil {
			return diags
		}
	}

	err = client.DBaaS().DeleteDbCluster(ctx, dbCluster.Identity)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("failed to delete db cluster: %w", err))
	}

	if err := waitForDeletedDbCluster(ctx, client, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
