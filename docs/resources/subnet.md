---
page_title: "thalassa_subnet Resource - terraform-provider-thalassa"
subcategory: "Networking"
description: |-
  Create an subnet in a VPC. Subnets are used to create a network for your resources. A VPC can have multiple subnets, and each subnet must have a different CIDR block. IPv4, IPv6 and Dual-stack subnets are supported. After creationg the CIDR cannot be changed.
---

# thalassa_subnet (Resource)

Create an subnet in a VPC. Subnets are used to create a network for your resources. A VPC can have multiple subnets, and each subnet must have a different CIDR block. IPv4, IPv6 and Dual-stack subnets are supported. After creationg the CIDR cannot be changed.


<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cidr` (String) CIDR of the Subnet
- `name` (String) Name of the Subnet
- `organisation_id` (String) Reference to the Organisation of the Subnet. If not provided, the organisation of the (Terraform) provider will be used.
- `vpc_id` (String) VPC of the Subnet

### Optional

- `annotations` (Map of String) Annotations for the Subnet
- `description` (String) A human readable description about the subnet
- `labels` (Map of String) Labels for the Subnet
- `route_table_id` (String) Route Table of the Subnet

### Read-Only

- `id` (String) The ID of this resource.
- `ipv4_addresses_available` (Number) Number of IPv4 addresses available in the Subnet
- `ipv4_addresses_used` (Number) Number of IPv4 addresses used in the Subnet
- `ipv6_addresses_available` (Number) Number of IPv6 addresses available in the Subnet
- `ipv6_addresses_used` (Number) Number of IPv6 addresses used in the Subnet
- `slug` (String) Slug of the Subnet
- `status` (String) Status of the Subnet
- `type` (String) Type of the Subnet


