package kubernetes

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var (
	ResourcesMap = map[string]*schema.Resource{
		"thalassa_kubernetes_cluster":              resourceKubernetesCluster(),
		"thalassa_kubernetes_node_pool":            resourceKubernetesNodePool(),
		"thalassa_kubernetes_cluster_role":         resourceKubernetesClusterRole(),
		"thalassa_kubernetes_cluster_role_binding": resourceKubernetesClusterRoleBinding(),
	}

	DataSourcesMap = map[string]*schema.Resource{
		"thalassa_kubernetes_version":               DataSourceKubernetesVersion(),
		"thalassa_kubernetes_cluster":               DataSourceKubernetesCluster(),
		"thalassa_kubernetes_cluster_session_token": dataSourceKubernetesClusterSessionToken(),
		"thalassa_kubernetes_cluster_role":          dataSourceKubernetesClusterRole(),
	}
)
