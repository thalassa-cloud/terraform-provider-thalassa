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
  description     = "Example VPC for route table route"
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

# Create a route table
resource "thalassa_route_table" "example" {
  name            = "example-route-table"
  description     = "Example route table for route"
  vpc_id          = thalassa_vpc.example.id
}

# Create a NAT gateway for the route
resource "thalassa_natgateway" "example" {
  name            = "example-nat-gateway"
  subnet_id       = thalassa_subnet.example.id
}

# Create a route table route with all required attributes
resource "thalassa_route_table_route" "example" {
  route_table_id = thalassa_route_table.example.id
  destination    = "0.0.0.0/0"
  target         = "internet-gateway"
  target_id      = "igw-123" # Replace with your internet gateway ID
}

# Output the route details
output "route_id" {
  value = thalassa_route_table_route.example.id
}

output "route_destination" {
  value = thalassa_route_table_route.example.destination
} 