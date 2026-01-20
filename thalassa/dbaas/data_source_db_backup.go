package dbaas

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thalassa-cloud/client-go/dbaas"
	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func dataSourceDbBackup() *schema.Resource {
	return &schema.Resource{
		Description: "Get a database backup by db_cluster_id and/or label selector. Always returns the newest matching backup.",
		ReadContext: dataSourceDbBackupRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity of the backup",
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"db_cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter backups by database cluster ID",
			},
			"label_selector": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Match backups that have all provided labels (exact key=value match)",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"identity": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity of the backup",
			},
			"labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Labels of the backup",
			},
			"annotations": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Annotations of the backup",
			},
			"backup_trigger": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Trigger of the backup (manual, schedule, system)",
			},
			"engine_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the database engine",
			},
			"engine_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Version of the database engine",
			},
			"backup_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the backup",
			},
			"online": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the backup is an online backup",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the backup",
			},
			"status_message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status message of the backup",
			},
			"delete_protection": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the backup is protected from deletion",
			},
			"started_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When the backup started (RFC3339 format)",
			},
			"stopped_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When the backup stopped (RFC3339 format)",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When the backup was created (RFC3339 format)",
			},
			"begin_lsn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Starting LSN of the backup (PostgreSQL only)",
			},
			"end_lsn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Ending LSN of the backup (PostgreSQL only)",
			},
			"begin_wal": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Starting WAL of the backup (PostgreSQL only)",
			},
			"end_wal": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Ending WAL of the backup (PostgreSQL only)",
			},
		},
	}
}

func dataSourceDbBackupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	dbClusterId := d.Get("db_cluster_id").(string)
	labelSelector, hasLabelSelector := d.GetOk("label_selector")

	var backups []dbaas.DbClusterBackup

	// Build filters for label selector
	var requestFilters []filters.Filter
	if hasLabelSelector {
		requestFilters = append(requestFilters, &filters.LabelFilter{
			MatchLabels: convert.ConvertToMap(labelSelector),
		})
	}

	// Use ListDbBackupsForDbCluster if db_cluster_id is provided, otherwise use ListDbBackupsForOrganisation
	if dbClusterId != "" {
		// If db_cluster_id is provided, use ListDbBackupsForDbCluster (already filters by cluster)
		// Then filter by label_selector client-side if provided
		backups, err = client.DBaaS().ListDbBackupsForDbCluster(ctx, dbClusterId, &dbaas.ListDbBackupsRequest{
			Filters: requestFilters,
		})
		if err != nil {
			return diag.FromErr(fmt.Errorf("error listing backups for cluster: %w", err))
		}
	} else {
		// If db_cluster_id is not provided, use ListDbBackupsForOrganisation with label filters
		backups, err = client.DBaaS().ListDbBackupsForOrganisation(ctx, &dbaas.ListDbBackupsRequest{
			Filters: requestFilters,
		})
		if err != nil {
			return diag.FromErr(fmt.Errorf("error listing backups: %w", err))
		}
	}

	if len(backups) == 0 {
		return diag.Errorf("no backup found matching the specified criteria")
	}

	// Sort by creation time (newest first) - use StoppedAt if available, otherwise CreatedAt
	sort.Slice(backups, func(i, j int) bool {
		var timeI, timeJ time.Time
		if backups[i].StoppedAt != nil {
			timeI = *backups[i].StoppedAt
		} else if backups[i].StartedAt != nil {
			timeI = *backups[i].StartedAt
		} else {
			timeI = backups[i].CreatedAt
		}

		if backups[j].StoppedAt != nil {
			timeJ = *backups[j].StoppedAt
		} else if backups[j].StartedAt != nil {
			timeJ = *backups[j].StartedAt
		} else {
			timeJ = backups[j].CreatedAt
		}

		return timeI.After(timeJ)
	})

	// Select the newest backup
	backup := backups[0]

	// Set the resource data
	d.SetId(backup.Identity)
	d.Set("identity", backup.Identity)
	d.Set("delete_protection", backup.DeleteProtection)
	d.Set("backup_trigger", backup.BackupTrigger)
	d.Set("engine_type", backup.EngineType)
	d.Set("engine_version", backup.EngineVersion)
	d.Set("backup_type", backup.BackupType)
	d.Set("online", backup.Online)
	d.Set("status", backup.Status)
	d.Set("status_message", backup.StatusMessage)
	d.Set("begin_lsn", backup.BeginLSN)
	d.Set("end_lsn", backup.EndLSN)
	d.Set("begin_wal", backup.BeginWAL)
	d.Set("end_wal", backup.EndWAL)
	d.Set("created_at", backup.CreatedAt.Format(time.RFC3339))

	if backup.Labels != nil {
		d.Set("labels", backup.Labels)
	}
	if backup.Annotations != nil {
		d.Set("annotations", backup.Annotations)
	}
	if backup.StartedAt != nil {
		d.Set("started_at", backup.StartedAt.Format(time.RFC3339))
	}
	if backup.StoppedAt != nil {
		d.Set("stopped_at", backup.StoppedAt.Format(time.RFC3339))
	}

	return nil
}
