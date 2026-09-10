resource "thalassa_containerregistry_namespace" "example" {
  region      = "nl-01"
  namespace   = "example"
  description = "Example container registry namespace"

  labels = {
    environment = "development"
  }
}

output "namespace_id" {
  value = thalassa_containerregistry_namespace.example.id
}
