---
page_title: "thalassa_loadbalancer Resource - terraform-provider-thalassa"
subcategory: "Networking"
description: |-
  Create an loadbalancer within a VPC
---

# thalassa_loadbalancer (Resource)

Create an loadbalancer within a VPC


<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of the Loadbalancer
- `organisation_id` (String) Reference to the Organisation of the Loadbalancer. If not provided, the organisation of the (Terraform) provider will be used.
- `subnet_id` (String) Subnet of the Loadbalancer

### Optional

- `annotations` (Map of String) Annotations for the Loadbalancer
- `description` (String) A human readable description about the loadbalancer
- `labels` (Map of String) Labels for the Loadbalancer

### Read-Only

- `id` (String) The ID of this resource.
- `slug` (String)
- `vpc_id` (String) VPC of the Loadbalancer


