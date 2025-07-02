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

# Create a VPC for the loadbalancer
resource "thalassa_vpc" "example" {
  name            = "example-vpc"
  description     = "Example VPC for loadbalancer"
  region          = "nl-01"
  cidrs           = ["10.0.0.0/16"]
}

# Create a subnet for the loadbalancer
resource "thalassa_subnet" "example" {
  name            = "example-subnet"
  description     = "Example subnet for loadbalancer"
  vpc_id          = thalassa_vpc.example.id
  cidr            = "10.0.1.0/24"
}

# Create a loadbalancer with all optional attributes
resource "thalassa_loadbalancer" "example" {
  # Required attributes
  organisation_id = "org-123" # Replace with your organisation ID
  name            = "example-loadbalancer"
  subnet_id       = thalassa_subnet.example.id
  
  # Optional attributes
  description = "Example loadbalancer for documentation with all optional attributes"
  
  # Labels are key-value pairs for organizing resources
  labels = {
    environment = "production"
    service     = "web"
    tier        = "frontend"
  }
  
  # Annotations are additional metadata for resources
  annotations = {
    cost-center = "cc-12345"
    ssl-cert = "wildcard.example.com"
    health-check = "enabled"
  }
}

# Output the loadbalancer ID
output "loadbalancer_id" {
  value = thalassa_loadbalancer.example.id
} 