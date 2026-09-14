resource "thalassa_observability_workspace" "example" {
  name           = "example-observability"
  region         = "nl-01"
  description    = "Metrics and logs for example workloads"
  retention_days = 30

  # Deletion is asynchronous; set wait_for_deleted = true to block until data is gone.
  # wait_for_deleted         = true
  # wait_for_deleted_timeout = 20

  labels = {
    environment = "development"
  }
}

output "workspace_id" {
  value = thalassa_observability_workspace.example.id
}

output "remote_write_url" {
  value = thalassa_observability_workspace.example.remote_write_url
}

output "loki_query_url" {
  value = thalassa_observability_workspace.example.loki_query_url
}
