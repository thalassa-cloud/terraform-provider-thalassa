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

# Create a VPC for the target group
resource "thalassa_vpc" "example" {
  name            = "example-vpc"
  description     = "Example VPC for target group"
  region          = "nl-01"
  cidrs           = ["10.0.0.0/16"]
}

# Create a target group with all optional attributes
resource "thalassa_target_group" "example" {
  # Required attributes
  organisation_id = "org-123" # Replace with your organisation ID
  name            = "example-target-group"
  vpc_id          = thalassa_vpc.example.id
  protocol        = "http"
  port            = 80
  
  # Optional attributes
  description = "Example target group for documentation with all optional attributes"
  
  # Labels are key-value pairs for organizing resources
  labels = {
    environment = "production"
    service     = "web"
    tier        = "backend"
  }
  
  # Annotations are additional metadata for resources
  annotations = {
    cost-center = "cc-12345"
    backup-policy = "none"
    monitoring = "enabled"
  }
  
  # Health check configuration (optional)
  health_check_protocol = "http"
  health_check_port     = 80
  health_check_path     = "/health"
  
  # Health check timing configuration (optional)
  health_check_interval = 30      # Check every 30 seconds
  health_check_timeout  = 5       # 5 second timeout
  healthy_threshold     = 3       # 3 successful checks to mark healthy
  unhealthy_threshold   = 3       # 3 failed checks to mark unhealthy
}

# Output the target group ID
output "target_group_id" {
  value = thalassa_target_group.example.id
} 