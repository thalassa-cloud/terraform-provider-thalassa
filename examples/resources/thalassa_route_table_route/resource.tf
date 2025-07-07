
# Create a VPC for the route table
resource "thalassa_vpc" "example" {
  name        = "example-vpc"
  description = "Example VPC for route table route"
  region      = "nl-01"
  cidrs       = ["10.0.0.0/16"]
}

# Create a subnet for the NAT gateway
resource "thalassa_subnet" "example" {
  name        = "example-subnet"
  description = "Example subnet for NAT gateway"
  vpc_id      = thalassa_vpc.example.id
  cidr        = "10.0.1.0/24"
}

# Create a route table
resource "thalassa_route_table" "example" {
  name        = "example-route-table"
  description = "Example route table for route"
  vpc_id      = thalassa_vpc.example.id
}

# Create a NAT gateway for the route
resource "thalassa_natgateway" "example" {
  name      = "example-nat-gateway"
  subnet_id = thalassa_subnet.example.id
}

# Create a route table route
resource "thalassa_route_table_route" "example" {
  route_table_id    = thalassa_route_table.example.id
  destination_cidr  = "0.0.0.0/0"
  target_natgateway = thalassa_natgateway.example.id
}

# Output the route details
output "route_id" {
  value = thalassa_route_table_route.example.id
}

output "route_destination" {
  value = thalassa_route_table_route.example.destination_cidr
}
