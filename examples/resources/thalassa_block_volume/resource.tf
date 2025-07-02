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

# Create a block volume with Thalassa default values
resource "thalassa_block_volume" "example" {
  name        = "example-block-volume"
  description = "Example block volume for documentation"
  region      = "nl-01"
  volume_type = "Block"  # Available: Block, Premium Block
  size_gb     = 20
}

# Output the block volume details
output "block_volume_id" {
  value = thalassa_block_volume.example.id
}

output "block_volume_name" {
  value = thalassa_block_volume.example.name
}
