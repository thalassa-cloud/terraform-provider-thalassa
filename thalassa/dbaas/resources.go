package dbaas

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var (
	ResourcesMap = map[string]*schema.Resource{
		"thalassa_db_cluster":        resourceDbCluster(),
		"thalassa_pg_database":       resourcePgDatabase(),
		"thalassa_pg_roles":          resourcePgRoles(),
		"thalassa_db_backupschedule": resourceDbBackupSchedule(),
	}

	DataSourcesMap = map[string]*schema.Resource{
		"thalassa_db_cluster":        dataSourceDbCluster(),
		"thalassa_pg_database":       dataSourcePgDatabase(),
		"thalassa_pg_roles":          dataSourcePgRoles(),
		"thalassa_db_backupschedule": dataSourceDbBackupSchedule(),
	}
)
