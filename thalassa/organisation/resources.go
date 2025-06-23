package organisation

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var (
	ResourcesMap = map[string]*schema.Resource{}

	DataSourcesMap = map[string]*schema.Resource{
		"thalassa_organisation": DataSourceOrganisations(),
	}
)
