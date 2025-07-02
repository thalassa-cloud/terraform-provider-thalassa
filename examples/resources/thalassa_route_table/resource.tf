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

# Create a VPC for the route table
resource "thalassa_vpc" "example" {
  name            = "example-vpc"
  description     = "Example VPC for route table"
  region          = "nl-01"
  cidrs           = ["10.0.0.0/16"]
}

# Create a route table with all optional attributes
resource "thalassa_route_table" "example" {
  # Required attributes
  organisation_id = "org-123" # Replace with your organisation ID
  name            = "example-route-table"
  vpc_id          = thalassa_vpc.example.id
  
  # Optional attributes
  description = "Example route table for documentation with all optional attributes"
  
  # Labels are key-value pairs for organizing resources
  labels = {
    environment = "production"
    service     = "networking"
    tier        = "private"
  }
  
  # Annotations are additional metadata for resources
  annotations = {
    cost-center = "cc-12345"
    backup-policy = "none"
    monitoring = "enabled"
  }
}

# Output the route table ID
output "route_table_id" {
  value = thalassa_route_table.example.id
} 