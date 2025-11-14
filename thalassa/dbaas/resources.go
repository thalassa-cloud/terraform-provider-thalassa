package dbaas

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var (
	ResourcesMap = map[string]*schema.Resource{
		"thalassa_dbaas_db_cluster":        resourceDbCluster(),
		"thalassa_dbaas_pg_database":       resourcePgDatabase(),
		"thalassa_dbaas_pg_roles":          resourcePgRoles(),
		"thalassa_dbaas_pg_grant":          resourcePgGrant(),
		"thalassa_dbaas_db_backupschedule": resourceDbBackupSchedule(),
	}

	DataSourcesMap = map[string]*schema.Resource{
		"thalassa_dbaas_db_cluster":        dataSourceDbCluster(),
		"thalassa_dbaas_pg_database":       dataSourcePgDatabase(),
		"thalassa_dbaas_pg_roles":          dataSourcePgRoles(),
		"thalassa_dbaas_db_backupschedule": dataSourceDbBackupSchedule(),
	}
)
