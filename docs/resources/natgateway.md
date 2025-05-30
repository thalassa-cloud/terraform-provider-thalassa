---
page_title: "thalassa_natgateway Resource - terraform-provider-thalassa"
subcategory: "Networking"
description: |-
  Create an NAT Gateway within a VPC
---

# thalassa_natgateway (Resource)

Create an NAT Gateway within a VPC

## Example Usage

```terraform
terraform {
  required_providers {
    thalassa = {
      source = "thalassa-cloud/thalassa"
    }
  }
}

provider "thalassa" {
  # Configuration options
}

# Create a NAT gateway
resource "thalassa_natgateway" "example" {
  name        = "example-nat"
  description = "Example NAT gateway for documentation"
  subnet_id   = "subnet-123" # Replace with your subnet ID
}

# Output the NAT gateway ID and IP addresses
output "natgateway_id" {
  value = thalassa_natgateway.example.id
}
```
<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of the NatGateway
- `organisation_id` (String) Reference to the Organisation of the NatGateway. If not provided, the organisation of the (Terraform) provider will be used.
- `subnet_id` (String) Subnet of the NatGateway

### Optional

- `annotations` (Map of String) Annotations for the NatGateway
- `description` (String) A human readable description about the natGateway
- `labels` (Map of String) Labels for the NatGateway

### Read-Only

- `endpoint_ip` (String) Endpoint IP of the NatGateway
- `id` (String) The ID of this resource.
- `slug` (String)
- `status` (String) Status of the NatGateway
- `v4_ip` (String) V4 IP of the NatGateway
- `v6_ip` (String) V6 IP of the NatGateway
- `vpc_id` (String) VPC of the NatGateway


