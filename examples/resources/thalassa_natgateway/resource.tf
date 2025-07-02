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

# Create a VPC for the NAT gateway
resource "thalassa_vpc" "example" {
  name            = "example-vpc"
  description     = "Example VPC for NAT gateway"
  region          = "nl-01"
  cidrs           = ["10.0.0.0/16"]
}

# Create a subnet for the NAT gateway
resource "thalassa_subnet" "example" {
  name            = "example-subnet"
  description     = "Example subnet for NAT gateway"
  vpc_id          = thalassa_vpc.example.id
  cidr            = "10.0.1.0/24"
}

# Create a NAT gateway with all optional attributes
resource "thalassa_natgateway" "example" {
  # Required attributes
  organisation_id = "org-123" # Replace with your organisation ID
  name            = "example-nat-gateway"
  subnet_id       = thalassa_subnet.example.id
  
  # Optional attributes
  description = "Example NAT gateway for documentation with all optional attributes"
  
  # Labels are key-value pairs for organizing resources
  labels = {
    environment = "production"
    service     = "networking"
    tier        = "public"
  }
  
  # Annotations are additional metadata for resources
  annotations = {
    cost-center = "cc-12345"
    backup-policy = "none"
    monitoring = "enabled"
  }
}

# Output the NAT gateway details
output "nat_gateway_id" {
  value = thalassa_natgateway.example.id
}

output "nat_gateway_endpoint_ip" {
  value = thalassa_natgateway.example.endpoint_ip
}
