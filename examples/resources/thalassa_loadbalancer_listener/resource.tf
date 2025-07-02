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

# Create a VPC for the load balancer
resource "thalassa_vpc" "example" {
  name            = "example-vpc"
  description     = "Example VPC for load balancer listener"
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

# Create a load balancer
resource "thalassa_loadbalancer" "example" {
  name            = "example-loadbalancer"
  description     = "Example load balancer for listener"
  vpc_id          = thalassa_vpc.example.id
  region          = "nl-01"
}

# Create a target group for the listener
resource "thalassa_target_group" "example" {
  name            = "example-target-group"
  description     = "Example target group for listener"
  vpc_id          = thalassa_vpc.example.id
  protocol        = "http"
  port            = 80
}

# Create a load balancer listener with all required attributes
resource "thalassa_loadbalancer_listener" "example" {
  loadbalancer_id = thalassa_loadbalancer.example.id
  name            = "example-listener"
  protocol        = "http"
  port            = 80
  default_action  = "forward"
}

# Output the listener details
output "listener_id" {
  value = thalassa_loadbalancer_listener.example.id
}

output "listener_name" {
  value = thalassa_loadbalancer_listener.example.name
} 