---
page_title: "thalassa_route_table_route Resource - terraform-provider-thalassa"
subcategory: "Networking"
description: |-
  Create an route table route with a destination cidr block, target gateway, target nat gateway and gateway address within a route table.
---

# thalassa_route_table_route (Resource)

Create an route table route with a destination cidr block, target gateway, target nat gateway and gateway address within a route table.

## Example Usage

```terraform
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
```
<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `destination_cidr` (String) Destination CIDR of the Route
- `route_table_id` (String) RouteTable of the Route

### Optional

- `gateway_address` (String) Gateway Address of the Route
- `notes` (String) Notes for the Route
- `organisation_id` (String) Organisation of the RouteTable
- `target_gateway` (String) Target Gateway of the Route
- `target_natgateway` (String) Target NAT Gateway of the Route

### Read-Only

- `id` (String) The ID of this resource.


