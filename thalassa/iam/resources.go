package iam

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var (
	ResourcesMap = map[string]*schema.Resource{
		"thalassa_iam_team":                              ResourceTeam(),
		"thalassa_iam_role":                              ResourceRole(),
		"thalassa_iam_role_binding":                      ResourceRoleBinding(),
		"thalassa_iam_service_account":                   ResourceServiceAccount(),
		"thalassa_iam_service_account_access_credential": ResourceServiceAccountAccessCredential(),
	}

	DataSourcesMap = map[string]*schema.Resource{
		"thalassa_iam_team":                 DataSourceTeam(),
		"thalassa_iam_role":                 DataSourceRole(),
		"thalassa_iam_organisation_members": DataSourceOrganisationMembers(),
		"thalassa_iam_service_account":      DataSourceServiceAccount(),
		// "thalassa_iam_user": DataSourceUser(),
	}
)
