data "thalassa_observability_workspace" "example" {
  name   = "example-observability"
  region = "nl-01"
}

output "workspace_id" {
  value = data.thalassa_observability_workspace.example.id
}

output "prometheus_query_url" {
  value = data.thalassa_observability_workspace.example.prometheus_query_url
}
