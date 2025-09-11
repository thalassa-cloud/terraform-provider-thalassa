package iam

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var (
	ResourcesMap = map[string]*schema.Resource{
		"thalassa_iam_team": ResourceTeam(),
		"thalassa_iam_role": ResourceRole(),
	}

	DataSourcesMap = map[string]*schema.Resource{
		"thalassa_iam_team": DataSourceTeam(),
		"thalassa_iam_role": DataSourceRole(),
		// "thalassa_iam_user": DataSourceUser(),
	}
)
