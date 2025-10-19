package kubernetes

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var (
	ResourcesMap = map[string]*schema.Resource{
		"thalassa_kubernetes_cluster":   resourceKubernetesCluster(),
		"thalassa_kubernetes_node_pool": resourceKubernetesNodePool(),
	}

	DataSourcesMap = map[string]*schema.Resource{
		"thalassa_kubernetes_version":               DataSourceKubernetesVersion(),
		"thalassa_kubernetes_cluster":               DataSourceKubernetesCluster(),
		"thalassa_kubernetes_cluster_session_token": dataSourceKubernetesClusterSessionToken(),
	}
)
