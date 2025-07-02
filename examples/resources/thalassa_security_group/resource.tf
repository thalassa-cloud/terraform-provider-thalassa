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

# Create a VPC for the security group
resource "thalassa_vpc" "example" {
  name            = "example-vpc"
  description     = "Example VPC for security group"
  region          = "nl-01"
  cidrs           = ["10.0.0.0/16"]
}

# Create a security group with Thalassa default values
resource "thalassa_security_group" "example" {
  # Required attributes
  name        = "example-security-group"
  vpc_identity = thalassa_vpc.example.id
  
  # Optional attributes
  description = "Example security group for documentation"
  
  # Allow traffic between instances in the same security group (optional, default: false)
  allow_same_group_traffic = true

  ingress_rules = [
    {
      name = "allow-http"
      ip_version = "ipv4"
      protocol = "tcp"
      priority = 100
      remote_type = "address"
      remote_address = "0.0.0.0/0"
      port_range_min = 80
      port_range_max = 80
      policy = "allow"
    },
    {
      name = "allow-https"
      ip_version = "ipv4"
      protocol = "tcp"
      priority = 101
      remote_type = "address"
      remote_address = "0.0.0.0/0"
      port_range_min = 443
      port_range_max = 443
      policy = "allow"
    }
  ]
}

# Output the security group details
output "security_group_id" {
  value = thalassa_security_group.example.identity
}

output "security_group_name" {
  value = thalassa_security_group.example.name
} 