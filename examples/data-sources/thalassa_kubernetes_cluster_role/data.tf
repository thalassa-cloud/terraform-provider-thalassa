# Get a Kubernetes cluster role by name
data "thalassa_kubernetes_cluster_role" "app_developer" {
  name = "app-developer"
}

# Get a Kubernetes cluster role by slug
data "thalassa_kubernetes_cluster_role" "monitoring" {
  slug = "monitoring-readonly"
}

# Use the retrieved cluster role data
output "app_developer_role" {
  description = "Information about the app developer cluster role"
  value = {
    id          = data.thalassa_kubernetes_cluster_role.app_developer.id
    name        = data.thalassa_kubernetes_cluster_role.app_developer.name
    slug        = data.thalassa_kubernetes_cluster_role.app_developer.slug
    description = data.thalassa_kubernetes_cluster_role.app_developer.description
    system      = data.thalassa_kubernetes_cluster_role.app_developer.system
    rules_count = length(data.thalassa_kubernetes_cluster_role.app_developer.rules)
  }
}

output "monitoring_role" {
  description = "Information about the monitoring cluster role"
  value = {
    id          = data.thalassa_kubernetes_cluster_role.monitoring.id
    name        = data.thalassa_kubernetes_cluster_role.monitoring.name
    slug        = data.thalassa_kubernetes_cluster_role.monitoring.slug
    description = data.thalassa_kubernetes_cluster_role.monitoring.description
    system      = data.thalassa_kubernetes_cluster_role.monitoring.system
    rules_count = length(data.thalassa_kubernetes_cluster_role.monitoring.rules)
  }
}

# Create a role binding using the data source
resource "thalassa_kubernetes_cluster_role_binding" "data_source_binding" {
  name            = "data-source-binding"
  description     = "Binding created using data source"
  cluster_role_id = data.thalassa_kubernetes_cluster_role.app_developer.id
  user_id         = "user-123" # Replace with actual user ID

  note = "Binding created from data source lookup"
}

# Example of using the role in a conditional resource
resource "thalassa_kubernetes_cluster_role_binding" "conditional_binding" {
  count = data.thalassa_kubernetes_cluster_role.monitoring.system ? 0 : 1

  name            = "conditional-binding"
  description     = "Conditional binding based on role type"
  cluster_role_id = data.thalassa_kubernetes_cluster_role.monitoring.id
  team_id         = "team-456" # Replace with actual team ID

  note = "Only created if the role is not a system role"
}
