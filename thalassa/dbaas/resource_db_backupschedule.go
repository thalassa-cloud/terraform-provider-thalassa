package dbaas

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thalassa-cloud/client-go/dbaas/dbaasalphav1"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func resourceDbBackupSchedule() *schema.Resource {
	return &schema.Resource{
		Description:   "Create a database backup schedule",
		CreateContext: resourceDbBackupScheduleCreate,
		ReadContext:   resourceDbBackupScheduleRead,
		UpdateContext: resourceDbBackupScheduleUpdate,
		DeleteContext: resourceDbBackupScheduleDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the database backup schedule",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the database backup schedule",
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the Db Backup Schedule. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			// TODO: missing api implementation
			// "description": {
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// 	Description: "The description of the database backup schedule",
			// },
			// TODO: missing api implementation
			// "labels": {
			// 	Type:        schema.TypeMap,
			// 	Optional:    true,
			// 	Description: "The labels of the database backup schedule",
			// 	Default:     map[string]string{},
			// 	Elem: &schema.Schema{
			// 		Type: schema.TypeString,
			// 	},
			// },
			// "annotations": {
			// 	Type:        schema.TypeMap,
			// 	Optional:    true,
			// 	Description: "The annotations of the database backup schedule",
			// 	Default:     map[string]string{},
			// 	Elem: &schema.Schema{
			// 		Type: schema.TypeString,
			// 	},
			// },
			"db_cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the database cluster",
			},
			"retention_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "7d",
				Description: "The retention policy of the database backup schedule (7d, 14d, 30d, 90d, 180d, 365d, 730d)",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					if !regexp.MustCompile(`^[1-9][0-9]*d$`).MatchString(val.(string)) {
						errs = append(errs, fmt.Errorf("retention_policy must be in the format of <number>d (e.g. 7d, 9d, 14d, etc)"))
					}
					warns = []string{}
					return
				},
			},
			"schedule": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The cron schedule of the database backup schedule (0 0 * * *)",
				Default:     "0 0 * * *",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					if !regexp.MustCompile(`^[0-9,\-\*]+ [0-9,\-\*]+ [0-9,\-\*]+ [0-9,\-\*]+ [0-9,\-\*]+$`).MatchString(val.(string)) {
						errs = append(errs, fmt.Errorf("schedule must be in valid cron format (e.g. 0 0 * * *)"))
					}
					warns = []string{}
					return
				},
			},
			"suspended": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the database backup schedule is suspended",
			},
			"backup_target": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "primary",
				Description: "The backup target of the database backup schedule (primary, prefer-standby)",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					if val.(string) != "primary" && val.(string) != "prefer-standby" {
						errs = append(errs, fmt.Errorf("backup_target must be either primary or prefer-standby"))
					}
					warns = []string{}
					return
				},
			},
		},
	}
}

func resourceDbBackupScheduleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	var dbCluster *dbaasalphav1.DbCluster
	dbClusterId := d.Get("db_cluster_id").(string)

	dbCluster, err = client.DbaaSAlphaV1().GetDbCluster(ctx, dbClusterId)
	if err != nil {
		if tcclient.IsNotFound(err) {
			return diag.FromErr(fmt.Errorf("db cluster not found: %w", err))
		}
		return diag.FromErr(fmt.Errorf("error getting db cluster: %w", err))
	}
	switch dbCluster.Status {
	case dbaasalphav1.DbClusterStatusReady, dbaasalphav1.DbClusterStatusUpdating, dbaasalphav1.DbClusterStatusCreating:
		break
	default:
		return diag.FromErr(fmt.Errorf("db cluster is not ready: %s", dbCluster.Status))
	}

	backupTarget := d.Get("backup_target").(string)
	retentionPolicy := d.Get("retention_policy").(string)
	createBackupSchedule := dbaasalphav1.CreatePgBackupScheduleRequest{
		Name:            d.Get("name").(string),
		Schedule:        d.Get("schedule").(string),
		RetentionPolicy: retentionPolicy,
		Target:          dbaasalphav1.DbClusterBackupScheduleTarget(backupTarget),
	}

	if labels, ok := d.GetOk("labels"); ok {
		createBackupSchedule.Labels = dbaasalphav1.Labels(convert.ConvertToMap(labels))
	}
	if annotations, ok := d.GetOk("annotations"); ok {
		createBackupSchedule.Annotations = dbaasalphav1.Annotations(convert.ConvertToMap(annotations))
	}

	createdBackupSchedule, err := client.DbaaSAlphaV1().CreatePgBackupSchedule(ctx, dbCluster.Identity, createBackupSchedule)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating backup schedule: %w", err))
	}

	d.SetId(createdBackupSchedule.Identity)
	d.Set("db_cluster_id", dbClusterId)
	d.Set("name", createdBackupSchedule.Name)
	d.Set("description", createdBackupSchedule.Description)
	d.Set("labels", createdBackupSchedule.Labels)
	d.Set("annotations", createdBackupSchedule.Annotations)
	d.Set("schedule", createdBackupSchedule.Schedule)
	d.Set("retention_policy", createdBackupSchedule.RetentionPolicy)
	d.Set("backup_target", createdBackupSchedule.Target)

	return resourceDbBackupScheduleRead(ctx, d, m)
}

func resourceDbBackupScheduleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	dbClusterId := d.Get("db_cluster_id").(string)
	pgBackupSchedules, err := client.DbaaSAlphaV1().ListPgBackupSchedules(ctx, dbClusterId)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error listing pg backup schedules: %w", err))
	}

	for _, backupSchedule := range pgBackupSchedules {
		if backupSchedule.Identity == d.Id() {
			d.Set("db_cluster_id", dbClusterId)
			d.Set("name", backupSchedule.Name)
			d.Set("schedule", backupSchedule.Schedule)
			d.Set("retention_policy", backupSchedule.RetentionPolicy)
			d.Set("backup_target", backupSchedule.Target)
			d.Set("suspended", backupSchedule.Suspended)
			d.Set("id", backupSchedule.Identity)
			d.Set("description", backupSchedule.Description)
			d.Set("labels", backupSchedule.Labels)
			d.Set("annotations", backupSchedule.Annotations)
			return nil
		}
	}
	return diag.FromErr(fmt.Errorf("backup schedule not found"))
}

func resourceDbBackupScheduleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	dbClusterId := d.Get("db_cluster_id").(string)

	name := d.Get("name").(string)
	schedule := d.Get("schedule").(string)
	retentionPolicy := d.Get("retention_policy").(string)
	backupTarget := d.Get("backup_target").(string)

	updateBackupSchedule := dbaasalphav1.UpdatePgBackupScheduleRequest{
		Name:            name,
		Schedule:        schedule,
		RetentionPolicy: retentionPolicy,
		Target:          dbaasalphav1.DbClusterBackupScheduleTarget(backupTarget),
	}

	var dbCluster *dbaasalphav1.DbCluster
	dbCluster, err = client.DbaaSAlphaV1().GetDbCluster(ctx, dbClusterId)
	if err != nil {
		if tcclient.IsNotFound(err) {
			return diag.FromErr(fmt.Errorf("db cluster not found: %w", err))
		}
		return diag.FromErr(fmt.Errorf("error getting db cluster: %w", err))
	}
	switch dbCluster.Status {
	case dbaasalphav1.DbClusterStatusReady, dbaasalphav1.DbClusterStatusUpdating, dbaasalphav1.DbClusterStatusCreating:
		break
	default:
		return diag.FromErr(fmt.Errorf("db cluster is not ready: %s", dbCluster.Status))
	}

	if labels, ok := d.GetOk("labels"); ok {
		updateBackupSchedule.Labels = dbaasalphav1.Labels(convert.ConvertToMap(labels))
	}
	if annotations, ok := d.GetOk("annotations"); ok {
		updateBackupSchedule.Annotations = dbaasalphav1.Annotations(convert.ConvertToMap(annotations))
	}

	_, err = client.DbaaSAlphaV1().UpdatePgBackupSchedule(ctx, dbCluster.Identity, d.Id(), updateBackupSchedule)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating backup schedule: %w", err))
	}

	return resourceDbBackupScheduleRead(ctx, d, m)
}

func resourceDbBackupScheduleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting db cluster: %w", err))
	}

	dbClusterId := d.Get("db_cluster_id").(string)
	dbCluster, err := client.DbaaSAlphaV1().GetDbCluster(ctx, dbClusterId)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error getting db cluster: %w", err))
	}

	err = client.DbaaSAlphaV1().DeletePgBackupSchedule(ctx, dbCluster.Identity, d.Id())
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error deleting backup schedule: %w", err))
	}
	d.SetId("")
	return nil
}
