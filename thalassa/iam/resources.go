package iam

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var (
	ResourcesMap = map[string]*schema.Resource{
		"thalassa_iam_team":                              ResourceTeam(),
		"thalassa_iam_role":                              ResourceRole(),
		"thalassa_iam_role_rule":                         ResourceRoleRule(),
		"thalassa_iam_role_binding":                      ResourceRoleBinding(),
		"thalassa_iam_service_account":                   ResourceServiceAccount(),
		"thalassa_iam_service_account_access_credential": ResourceServiceAccountAccessCredential(),
		"thalassa_iam_federated_identity_provider":       ResourceFederatedIdentityProvider(),
		"thalassa_iam_federated_identity":                ResourceFederatedIdentity(),
	}

	DataSourcesMap = map[string]*schema.Resource{
		"thalassa_iam_team":                        DataSourceTeam(),
		"thalassa_iam_role":                        DataSourceRole(),
		"thalassa_iam_organisation_members":        DataSourceOrganisationMembers(),
		"thalassa_iam_service_account":             DataSourceServiceAccount(),
		"thalassa_iam_federated_identity_provider": DataSourceFederatedIdentityProvider(),
		"thalassa_iam_federated_identity":          DataSourceFederatedIdentity(),
		// "thalassa_iam_user": DataSourceUser(),
	}
)
