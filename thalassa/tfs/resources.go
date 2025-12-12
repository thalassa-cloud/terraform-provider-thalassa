package tfs

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var (
	ResourcesMap = map[string]*schema.Resource{
		"thalassa_tfs_instance": resourceTfsInstance(),
	}

	DataSourcesMap = map[string]*schema.Resource{
		"thalassa_tfs_instance": DataSourceTfsInstance(),
	}
)
