package dbaas

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var (
	ResourcesMap = map[string]*schema.Resource{
		"thalassa_db_cluster": resourceDbCluster(),
	}

	DataSourcesMap = map[string]*schema.Resource{
		"thalassa_db_cluster": DataSourceDbCluster(),
	}
)
