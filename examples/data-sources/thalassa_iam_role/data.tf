terraform {
  required_providers {
    thalassa = {
      source = "thalassa-cloud/thalassa"
    }
  }
}

provider "thalassa" {
  # Configuration options
}

# Get an existing organisation role by name
data "thalassa_iam_role" "example" {
  name = "example-role"
}

# Output the role details
output "role_id" {
  value = data.thalassa_iam_role.example.id
}

output "role_name" {
  value = data.thalassa_iam_role.example.name
}

output "role_slug" {
  value = data.thalassa_iam_role.example.slug
}

output "role_description" {
  value = data.thalassa_iam_role.example.description
}

output "role_is_read_only" {
  value = data.thalassa_iam_role.example.role_is_read_only
}

output "role_system" {
  value = data.thalassa_iam_role.example.system
}

