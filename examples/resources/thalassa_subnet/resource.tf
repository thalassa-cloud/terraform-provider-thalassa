# Create a VPC for the subnet
resource "thalassa_vpc" "example" {
  name        = "example-vpc"
  description = "Example VPC for subnet"
  region      = "nl-01"
  cidrs       = ["10.0.0.0/16"]
}

# Create a route table for the subnet (optional)
resource "thalassa_route_table" "example" {
  name        = "example-route-table"
  description = "Example route table for subnet"
  vpc_id      = thalassa_vpc.example.id
}

# Create a subnet with all optional attributes
resource "thalassa_subnet" "example" {
  name   = "example-subnet"
  vpc_id = thalassa_vpc.example.id
  cidr   = "10.0.1.0/24"

  # Optional attributes
  description = "Example subnet for documentation with all optional attributes"

  labels = {
    environment = "production"
    tier        = "web"
    network     = "private"
  }

  # Annotations are additional metadata for resources
  annotations = {
    cost-center   = "cc-12345"
    backup-policy = "none"
    monitoring    = "enabled"
  }

  # Associate with a route table (optional)
  route_table_id = thalassa_route_table.example.id
}

# Output the subnet ID
output "subnet_id" {
  value = thalassa_subnet.example.id
} 