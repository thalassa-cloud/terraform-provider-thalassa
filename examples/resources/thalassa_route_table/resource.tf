# Create a VPC for the route table
resource "thalassa_vpc" "example" {
  name            = "example-vpc"
  description     = "Example VPC for route table"
  region          = "nl-01"
  cidrs           = ["10.0.0.0/16"]
}

# Create a route table
resource "thalassa_route_table" "example" {
  name            = "example-route-table"
  vpc_id          = thalassa_vpc.example.id
  
  description = "Example route table for documentation"
  
  labels = {
    environment = "production"
    service     = "networking"
    tier        = "private"
  }
  
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
