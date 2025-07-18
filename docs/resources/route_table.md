---
page_title: "thalassa_route_table Resource - terraform-provider-thalassa"
subcategory: "Networking"
description: |-
  Create an routeTable
---

# thalassa_route_table (Resource)

Create an routeTable

## Example Usage

```terraform
# Create a VPC for the route table
resource "thalassa_vpc" "example" {
  name        = "example-vpc"
  description = "Example VPC for route table"
  region      = "nl-01"
  cidrs       = ["10.0.0.0/16"]
}

# Create a route table
resource "thalassa_route_table" "example" {
  name   = "example-route-table"
  vpc_id = thalassa_vpc.example.id

  description = "Example route table for documentation"

  labels = {
    environment = "production"
    service     = "networking"
    tier        = "private"
  }

  annotations = {
    cost-center   = "cc-12345"
    backup-policy = "none"
    monitoring    = "enabled"
  }
}

# Output the route table ID
output "route_table_id" {
  value = thalassa_route_table.example.id
}
```
<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of the RouteTable
- `vpc_id` (String) VPC of the RouteTable

### Optional

- `annotations` (Map of String) Annotations for the RouteTable
- `description` (String) A human readable description about the routeTable
- `labels` (Map of String) Labels for the RouteTable
- `organisation_id` (String) Reference to the Organisation of the RouteTable. If not provided, the organisation of the (Terraform) provider will be used.

### Read-Only

- `id` (String) The ID of this resource.
- `slug` (String)


