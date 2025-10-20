# Create a Kubernetes cluster role for application developers
resource "thalassa_kubernetes_cluster_role" "app_developer" {
  name        = "app-developer"
  description = "Role for application developers with read/write access to applications"

  labels = {
    "team"        = "development"
    "environment" = "production"
  }

  annotations = {
    "managed-by" = "terraform"
  }

  rules {
    resources  = ["pods", "services", "deployments", "configmaps", "secrets"]
    verbs      = ["get", "list", "watch", "create", "update", "delete", "patch"]
    api_groups = ["", "apps"]
    note       = "Full access to application resources"
  }

  rules {
    resources  = ["pods/log"]
    verbs      = ["get", "list"]
    api_groups = [""]
    note       = "Access to pod logs"
  }
}

# Create a read-only cluster role for monitoring
resource "thalassa_kubernetes_cluster_role" "monitoring" {
  name        = "monitoring-readonly"
  description = "Read-only role for monitoring and observability tools"

  labels = {
    "team"    = "platform"
    "purpose" = "monitoring"
  }

  rules {
    resources  = ["*"]
    verbs      = ["get", "list", "watch"]
    api_groups = ["*"]
    note       = "Read-only access to all resources for monitoring"
  }

  rules {
    resources  = ["nodes", "nodes/metrics", "nodes/stats"]
    verbs      = ["get", "list"]
    api_groups = [""]
    note       = "Access to node metrics and stats"
  }
}

# Create a cluster role for database administrators
resource "thalassa_kubernetes_cluster_role" "db_admin" {
  name        = "database-administrator"
  description = "Role for database administrators with access to database-related resources"

  labels = {
    "team" = "database"
    "tier" = "critical"
  }

  rules {
    resources      = ["pods", "services", "configmaps", "secrets", "persistentvolumes", "persistentvolumeclaims"]
    verbs          = ["get", "list", "watch", "create", "update", "delete", "patch"]
    api_groups     = ["", "apps"]
    resource_names = ["postgres-*", "mysql-*", "redis-*"]
    note           = "Access to database-related resources"
  }

  rules {
    resources      = ["pods/exec"]
    verbs          = ["create"]
    api_groups     = [""]
    resource_names = ["postgres-*", "mysql-*", "redis-*"]
    note           = "Execute commands in database pods"
  }
}
