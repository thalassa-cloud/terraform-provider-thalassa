# Create a block volume that will be selected by the snapshot policy
resource "thalassa_block_volume" "example" {
  name        = "example-block-volume"
  description = "Example block volume for snapshot policy"
  region      = "nl-01"
  volume_type = "Block"
  size_gb     = 50

  # Labels to match the snapshot policy selector
  labels = {
    backup     = "enabled"
    environment = "production"
    service     = "database"
  }
}

# Example 1: Snapshot policy using label selector
resource "thalassa_snapshot_policy" "selector_example" {
  name        = "example-snapshot-policy-selector"
  description = "Example snapshot policy using label selector"
  region      = "nl-01"

  # Schedule: Daily at 2 AM UTC
  schedule = "0 2 * * *"
  timezone = "UTC"

  # TTL: Keep snapshots for 30 days
  ttl = "30d"

  # Optional: Keep maximum of 10 snapshots
  keep_count = 10

  # Enable the policy
  enabled = true

  # Target volumes using label selector
  target {
    type = "selector"
    selector = {
      backup      = "enabled"
      environment = "production"
    }
  }
}

# # Example 2: Snapshot policy targeting specific volumes
resource "thalassa_snapshot_policy" "explicit_example" {
  name        = "example-snapshot-policy-explicit"
  description = "Example snapshot policy targeting specific volumes"
  region      = "nl-01"

  # Schedule: Every 6 hours
  schedule = "0 */6 * * *"
  timezone = "UTC"

  # TTL: Keep snapshots for 7 days
  ttl = "7d"

  # Enable the policy
  enabled = true

  # Target specific volumes by their identities
  target {
    type              = "explicit"
    volume_identities = [thalassa_block_volume.example.id]
  }

  labels = {
    backup-type = "frequent"
  }
}

# Output the snapshot policy details
output "snapshot_policy_selector_id" {
  value = thalassa_snapshot_policy.selector_example.id
}

output "snapshot_policy_selector_name" {
  value = thalassa_snapshot_policy.selector_example.name
}

output "snapshot_policy_selector_next_snapshot_at" {
  value = thalassa_snapshot_policy.selector_example.next_snapshot_at
}

output "snapshot_policy_explicit_id" {
  value = thalassa_snapshot_policy.explicit_example.id
}

