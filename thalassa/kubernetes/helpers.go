package kubernetes

import kubernetes "github.com/thalassa-cloud/client-go/kubernetes"

func clusterVersionReferenceMatches(configured string, version kubernetes.KubernetesVersion) bool {
	return version.Name == configured ||
		version.Slug == configured ||
		version.Identity == configured
}

func resolvedClusterVersionReference(configured string, version kubernetes.KubernetesVersion) string {
	if clusterVersionReferenceMatches(configured, version) {
		return configured
	}
	return version.Slug
}
