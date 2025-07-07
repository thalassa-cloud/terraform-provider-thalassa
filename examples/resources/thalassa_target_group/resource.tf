# Create a VPC for the target group
resource "thalassa_vpc" "example" {
  name        = "example-vpc"
  description = "Example VPC for target group"
  region      = "nl-01"
  cidrs       = ["10.0.0.0/16"]
}

# Create a target group with all optional attributes
resource "thalassa_target_group" "example" {
  name     = "example-target-group"
  vpc_id   = thalassa_vpc.example.id
  protocol = "tcp"
  port     = 80

  # Optional attributes
  description = "Example target group for documentation with all optional attributes"

  labels = {
    environment = "production"
    service     = "web"
    tier        = "backend"
  }

  health_check_protocol = "http"
  health_check_port     = 80
  health_check_path     = "/health"

  health_check_interval = 30 # Check every 30 seconds
  health_check_timeout  = 5  # 5 second timeout
  healthy_threshold     = 3  # 3 successful checks to mark healthy
  unhealthy_threshold   = 3  # 3 failed checks to mark unhealthy
}

# Output the target group ID
output "target_group_id" {
  value = thalassa_target_group.example.id
}
