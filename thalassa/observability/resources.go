package observability

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var (
	ResourcesMap = map[string]*schema.Resource{
		"thalassa_observability_workspace": ResourceWorkspace(),
	}

	DataSourcesMap = map[string]*schema.Resource{
		"thalassa_observability_workspace": DataSourceWorkspace(),
	}
)
