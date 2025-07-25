---
page_title: "thalassa_target_group Resource - terraform-provider-thalassa"
subcategory: "Networking"
description: |-
  Create a target group for a load balancer
---

# thalassa_target_group (Resource)

Create a target group for a load balancer

## Example Usage

```terraform
# Create a VPC for the target group
resource "thalassa_vpc" "example" {
  name        = "example-vpc"
  description = "Example VPC for target group"
  region      = "nl-01"
  cidrs       = ["10.0.0.0/16"]
}

# Create a target group with all optional attributes
resource "thalassa_target_group" "example" {
  name     = "example-target-group"
  vpc_id   = thalassa_vpc.example.id
  protocol = "tcp"
  port     = 80

  # Optional attributes
  description = "Example target group for documentation with all optional attributes"

  labels = {
    environment = "production"
    service     = "web"
    tier        = "backend"
  }

  health_check_protocol = "http"
  health_check_port     = 80
  health_check_path     = "/health"

  health_check_interval = 30 # Check every 30 seconds
  health_check_timeout  = 5  # 5 second timeout
  healthy_threshold     = 3  # 3 successful checks to mark healthy
  unhealthy_threshold   = 3  # 3 failed checks to mark unhealthy
}

# Output the target group ID
output "target_group_id" {
  value = thalassa_target_group.example.id
}
```
<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of the Target Group
- `port` (Number) The port on which the targets receive traffic
- `protocol` (String) The protocol to use for routing traffic to the targets. Must be one of: tcp, udp.
- `vpc_id` (String) The VPC this target group belongs to

### Optional

- `annotations` (Map of String) Annotations for the Target Group
- `attachments` (Block List) The targets to attach to the target group. If provided, the targets will be attached to the target group when the resource is created. Overwrites the target group attachment resource. (see [below for nested schema](#nestedblock--attachments))
- `description` (String) A human readable description about the target group
- `health_check_interval` (Number) The approximate amount of time, in seconds, between health checks of an individual target
- `health_check_path` (String) The path to use for health checks (only for HTTP/HTTPS)
- `health_check_port` (Number) The port to use for health checks
- `health_check_protocol` (String) The protocol to use for health checks. Must be one of: tcp, http.
- `health_check_timeout` (Number) The amount of time, in seconds, during which no response means a failed health check
- `healthy_threshold` (Number) The number of consecutive health checks successes required before considering an unhealthy target healthy
- `labels` (Map of String) Labels for the Target Group
- `organisation_id` (String) Reference to the Organisation of the Target Group. If not provided, the organisation of the (Terraform) provider will be used.
- `unhealthy_threshold` (Number) The number of consecutive health check failures required before considering a target unhealthy

### Read-Only

- `id` (String) The ID of this resource.
- `slug` (String)

<a id="nestedblock--attachments"></a>
### Nested Schema for `attachments`

Required:

- `id` (String) The ID of the target (e.g. instance ID)

 