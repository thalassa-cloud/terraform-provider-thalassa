resource "thalassa_project" "example" {
  name        = "example-project"
  description = "An example project"

  labels = {
    environment = "development"
  }
}

output "project_id" {
  value = thalassa_project.example.id
}

output "project_slug" {
  value = thalassa_project.example.slug
}
