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

# Create an organisation role
resource "thalassa_iam_role" "example" {
  name        = "example-role"
  description = "An example organisation role for demonstration purposes"
  
  labels = {
    environment = "development"
    project     = "example"
  }
  
  annotations = {
    "example.com/created-by" = "terraform"
  }
}

# Output the role details
output "role_id" {
  value = thalassa_iam_role.example.id
}

output "role_name" {
  value = thalassa_iam_role.example.name
}

output "role_slug" {
  value = thalassa_iam_role.example.slug
}

output "role_description" {
  value = thalassa_iam_role.example.description
}

