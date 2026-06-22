package secrets

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var ResourcesMap = map[string]*schema.Resource{
	"thalassa_secret":               ResourceSecret(),
	"thalassa_secret_version":       ResourceSecretVersion(),
	"thalassa_secret_access_policy": ResourceSecretAccessPolicy(),
}

var DataSourcesMap = map[string]*schema.Resource{}
