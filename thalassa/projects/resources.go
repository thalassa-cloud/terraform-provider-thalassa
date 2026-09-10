package projects

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var (
	ResourcesMap = map[string]*schema.Resource{
		"thalassa_project": ResourceProject(),
	}

	DataSourcesMap = map[string]*schema.Resource{
		"thalassa_project": DataSourceProject(),
	}
)
