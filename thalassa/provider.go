package thalassa

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/dbaas"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/iaas"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/iam"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/kubernetes"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/objectstorage"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/organisation"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The API token for authentication. Can be set via the THALASSA_API_TOKEN environment variable.",
				DefaultFunc: schema.EnvDefaultFunc("THALASSA_API_TOKEN", nil),
			},
			"client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The OIDC client ID for authentication. Can be set via the THALASSA_CLIENT_ID environment variable.",
				DefaultFunc: schema.EnvDefaultFunc("THALASSA_CLIENT_ID", nil),
			},
			"client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The OIDC client secret for authentication. Can be set via the THALASSA_CLIENT_SECRET environment variable.",
				DefaultFunc: schema.EnvDefaultFunc("THALASSA_CLIENT_SECRET", nil),
			},
			"allow_insecure_oidc": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Allow insecure OIDC authentication. Can be set via the THALASSA_ALLOW_INSECURE_OIDC environment variable.",
				DefaultFunc: schema.EnvDefaultFunc("THALASSA_ALLOW_INSECURE_OIDC", false),
			},
			"api": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The API endpoint URL. Can be set via the THALASSA_API_ENDPOINT environment variable.",
				DefaultFunc: schema.EnvDefaultFunc("THALASSA_API_ENDPOINT", "https://api.thalassa.cloud"),
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The organisation ID to use. Can be set via the THALASSA_ORGANISATION environment variable.",
				DefaultFunc: schema.EnvDefaultFunc("THALASSA_ORGANISATION", ""),
			},
		},
		ResourcesMap: JoinMaps(
			iaas.ResourcesMap,
			kubernetes.ResourcesMap,
			organisation.ResourcesMap,
			dbaas.ResourcesMap,
			iam.ResourcesMap,
			objectstorage.ResourcesMap,
		),
		DataSourcesMap: JoinMaps(
			iaas.DataSourcesMap,
			kubernetes.DataSourcesMap,
			organisation.DataSourcesMap,
			dbaas.DataSourcesMap,
			iam.DataSourcesMap,
			objectstorage.DataSourcesMap,
		),
		ConfigureContextFunc: provider.ProviderConfigure,
	}
}

func JoinMaps(maps ...map[string]*schema.Resource) map[string]*schema.Resource {
	result := make(map[string]*schema.Resource)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}
