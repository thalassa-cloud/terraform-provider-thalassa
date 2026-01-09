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

# Create a service account first
resource "thalassa_iam_service_account" "example" {
  name        = "example-service-account"
  description = "An example service account for demonstration purposes"

  labels = {
    environment = "development"
    project     = "example"
    type        = "automation"
  }
}

# Create access credentials for the service account
resource "thalassa_iam_service_account_access_credential" "api_credential" {
  service_account_id = thalassa_iam_service_account.example.id
  scopes = [
    "api:read",
    "api:write"
  ]
}

# Create object storage access credentials
resource "thalassa_iam_service_account_access_credential" "storage_credential" {
  service_account_id = thalassa_iam_service_account.example.id
  scopes = [
    "objectStorage"
  ]
}

# Output the credential details
output "api_credential_id" {
  value = thalassa_iam_service_account_access_credential.api_credential.id
}

output "api_credential_access_key" {
  value     = thalassa_iam_service_account_access_credential.api_credential.access_key
  sensitive = true
}

output "api_credential_access_secret" {
  value     = thalassa_iam_service_account_access_credential.api_credential.access_secret
  sensitive = true
}

output "storage_credential_id" {
  value = thalassa_iam_service_account_access_credential.storage_credential.id
}
