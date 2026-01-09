# Create a block volume to snapshot
resource "thalassa_block_volume" "example" {
  name        = "example-block-volume"
  description = "Example block volume for snapshot"
  region      = "nl-01"
  volume_type = "Block"
  size_gb     = 20
}

# Create a snapshot from the block volume
resource "thalassa_snapshot" "example" {
  name        = "example-snapshot"
  description = "Example snapshot created from block volume"

  # Optional attributes
  delete_protection = false
  wait_until_available = true

  # Labels for organizing the snapshot
  labels = {
    environment = "production"
    backup      = "daily"
    service     = "database"
  }

  # Annotations for additional metadata
  annotations = {
    cost-center = "cc-12345"
    retention   = "30d"
  }
}

# Output the snapshot details
output "snapshot_id" {
  value = thalassa_snapshot.example.id
}

output "snapshot_name" {
  value = thalassa_snapshot.example.name
}

output "snapshot_status" {
  value = thalassa_snapshot.example.status
}

output "snapshot_size_gb" {
  value = thalassa_snapshot.example.size_gb
}

output "source_volume_id" {
  value = thalassa_snapshot.example.source_volume_id
}

