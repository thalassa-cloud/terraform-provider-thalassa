package containerregistry

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var (
	ResourcesMap = map[string]*schema.Resource{
		"thalassa_containerregistry_namespace":               ResourceNamespace(),
		"thalassa_containerregistry_namespace_configuration": ResourceNamespaceConfiguration(),
	}

	DataSourcesMap = map[string]*schema.Resource{
		"thalassa_containerregistry_namespace":               DataSourceNamespace(),
		"thalassa_containerregistry_namespace_configuration": DataSourceNamespaceConfiguration(),
	}
)
