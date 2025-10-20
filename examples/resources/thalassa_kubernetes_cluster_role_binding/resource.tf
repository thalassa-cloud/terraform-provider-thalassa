# First, create a cluster role
resource "thalassa_kubernetes_cluster_role" "app_developer" {
  name        = "app-developer"
  description = "Role for application developers"

  rules {
    resources  = ["pods", "services", "deployments", "configmaps"]
    verbs      = ["get", "list", "watch", "create", "update", "delete", "patch"]
    api_groups = ["", "apps"]
  }
}

# Bind a user to the cluster role
resource "thalassa_kubernetes_cluster_role_binding" "user_binding" {
  name            = "app-developer-user-binding"
  description     = "Bind user to app developer role"
  cluster_role_id = thalassa_kubernetes_cluster_role.app_developer.id
  user_id         = "user-123" # Replace with actual user ID

  labels = {
    "team" = "development"
  }

  note = "Binding for application developer access"
}

# Bind a team to the cluster role
resource "thalassa_kubernetes_cluster_role_binding" "team_binding" {
  name            = "app-developer-team-binding"
  description     = "Bind development team to app developer role"
  cluster_role_id = thalassa_kubernetes_cluster_role.app_developer.id
  team_id         = "team-456" # Replace with actual team ID

  labels = {
    "team" = "development"
  }

  note = "Team binding for development team"
}

# Bind a service account to the cluster role
resource "thalassa_kubernetes_cluster_role_binding" "service_account_binding" {
  name               = "app-developer-sa-binding"
  description        = "Bind service account to app developer role"
  cluster_role_id    = thalassa_kubernetes_cluster_role.app_developer.id
  service_account_id = "sa-789" # Replace with actual service account ID

  labels = {
    "purpose" = "automation"
  }

  note = "Service account binding for automated deployments"
}

# Create a read-only role for monitoring
resource "thalassa_kubernetes_cluster_role" "monitoring" {
  name        = "monitoring-readonly"
  description = "Read-only role for monitoring tools"

  rules {
    resources  = ["*"]
    verbs      = ["get", "list", "watch"]
    api_groups = ["*"]
  }
}

# Bind monitoring service account to read-only role
resource "thalassa_kubernetes_cluster_role_binding" "monitoring_binding" {
  name               = "monitoring-sa-binding"
  description        = "Bind monitoring service account to read-only role"
  cluster_role_id    = thalassa_kubernetes_cluster_role.monitoring.id
  service_account_id = "monitoring-sa-001" # Replace with actual service account ID

  labels = {
    "team"    = "platform"
    "purpose" = "monitoring"
  }

  note = "Monitoring service account with read-only access"
}
