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

# Create a NAT gateway
resource "thalassa_natgateway" "example" {
  name        = "example-nat"
  description = "Example NAT gateway for documentation"
  subnet_id   = "subnet-123" # Replace with your subnet ID
}

# Output the NAT gateway ID and IP addresses
output "natgateway_id" {
  value = thalassa_natgateway.example.id
}
