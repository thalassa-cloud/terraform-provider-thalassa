terraform {
  required_providers {
    thalassa = {
      source = "local/thalassa/thalassa"
    }
  }
}

provider "thalassa" {
  # Configuration options
}

# Reference an existing object storage bucket
data "thalassa_objectstorage_bucket" "existing" {
  name   = "existing-bucket-name"
  region = "nl-01"
}

# Reference a bucket by name only (will search across regions)
data "thalassa_objectstorage_bucket" "by_name" {
  name = "another-existing-bucket"
}

# Output the bucket details
output "existing_bucket_id" {
  value = data.thalassa_objectstorage_bucket.existing.id
}

output "existing_bucket_name" {
  value = data.thalassa_objectstorage_bucket.existing.name
}

output "existing_bucket_region" {
  value = data.thalassa_objectstorage_bucket.existing.region
}

output "existing_bucket_endpoint" {
  value = data.thalassa_objectstorage_bucket.existing.endpoint
}

output "existing_bucket_public" {
  value = data.thalassa_objectstorage_bucket.existing.public
}

output "existing_bucket_status" {
  value = data.thalassa_objectstorage_bucket.existing.status
}

output "existing_bucket_policy" {
  value = data.thalassa_objectstorage_bucket.existing.policy
}

# Example of using the data source to create a new bucket in the same region
resource "thalassa_objectstorage_bucket" "new_in_same_region" {
  name   = "new-bucket-in-same-region"
  region = data.thalassa_objectstorage_bucket.existing.region
  public = false
} 