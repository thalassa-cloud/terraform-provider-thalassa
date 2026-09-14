# Observability workspaces are beta and require sign up to the beta program
# before they can be used by your organisation.
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
