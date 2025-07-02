package dbaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func dataSourceDbBackupSchedule() *schema.Resource {

	return &schema.Resource{
		ReadContext: dataSourceDbBackupScheduleRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the Db Backup Schedule. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"db_cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"schedule": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"suspended": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"backup_target": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"retention_policy": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDbBackupScheduleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	dbClusterId := d.Get("db_cluster_id").(string)
	backupSchedules, err := client.DbaaSAlphaV1().ListPgBackupSchedules(ctx, dbClusterId)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, backupSchedule := range backupSchedules {
		if backupSchedule.Name == d.Get("name").(string) {
			d.SetId(backupSchedule.Identity)
			d.Set("name", backupSchedule.Name)
			d.Set("db_cluster_id", dbClusterId)
			d.Set("schedule", backupSchedule.Schedule)
			d.Set("suspended", backupSchedule.Suspended)
			d.Set("backup_target", backupSchedule.Target)
			d.Set("retention_policy", backupSchedule.RetentionPolicy)
			return nil
		}
	}

	return nil
}
