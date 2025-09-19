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

# Get a service account by name
data "thalassa_iam_service_account" "by_name" {
  name = "example-service-account"
}

# Get a service account by slug
data "thalassa_iam_service_account" "by_slug" {
  slug = "example-service-account"
}

# Output the service account details
output "service_account_by_name" {
  value = {
    id          = data.thalassa_iam_service_account.by_name.id
    name        = data.thalassa_iam_service_account.by_name.name
    slug        = data.thalassa_iam_service_account.by_name.slug
    description = data.thalassa_iam_service_account.by_name.description
    created_at  = data.thalassa_iam_service_account.by_name.created_at
  }
}

output "service_account_by_slug" {
  value = {
    id          = data.thalassa_iam_service_account.by_slug.id
    name        = data.thalassa_iam_service_account.by_slug.name
    slug        = data.thalassa_iam_service_account.by_slug.slug
    description = data.thalassa_iam_service_account.by_slug.description
    created_at  = data.thalassa_iam_service_account.by_slug.created_at
  }
}

output "service_account_role_bindings" {
  value = data.thalassa_iam_service_account.by_name.role_bindings
}
