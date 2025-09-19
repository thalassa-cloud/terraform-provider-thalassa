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

# Create a service account
resource "thalassa_iam_service_account" "example" {
  name        = "example-service-account"
  description = "An example service account for demonstration purposes"
  
  labels = {
    environment = "development"
    project     = "example"
    type        = "automation"
  }
  
  annotations = {
    "example.com/created-by" = "terraform"
    "example.com/purpose"    = "ci-cd"
  }
}

# Output the service account details
output "service_account_id" {
  value = thalassa_iam_service_account.example.id
}

output "service_account_name" {
  value = thalassa_iam_service_account.example.name
}

output "service_account_slug" {
  value = thalassa_iam_service_account.example.slug
}

output "service_account_description" {
  value = thalassa_iam_service_account.example.description
}

output "service_account_created_at" {
  value = thalassa_iam_service_account.example.created_at
}
