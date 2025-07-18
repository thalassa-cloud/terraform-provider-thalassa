---
page_title: "thalassa_virtual_machine_instance Resource - terraform-provider-thalassa"
subcategory: "Compute"
description: |-
  Create an virtual machine instance within a subnet on the Thalassa Cloud platform
---

# thalassa_virtual_machine_instance (Resource)

Create an virtual machine instance within a subnet on the Thalassa Cloud platform

## Example Usage

```terraform
# Create a VPC for the virtual machine instance
resource "thalassa_vpc" "example" {
  name        = "example-vpc"
  description = "Example VPC for virtual machine instance"
  region      = "nl-01"
  cidrs       = ["10.0.0.0/16"]
}

# Create a subnet for the virtual machine instance
resource "thalassa_subnet" "example" {
  name        = "example-subnet"
  description = "Example subnet for virtual machine instance"
  vpc_id      = thalassa_vpc.example.id
  cidr        = "10.0.1.0/24"
}

# Create a security group for the virtual machine instance
resource "thalassa_security_group" "example" {
  name        = "example-security-group"
  description = "Example security group for virtual machine instance"
  vpc_id      = thalassa_vpc.example.id
}

# Create a cloud init template (optional)
resource "thalassa_cloud_init_template" "example" {
  name    = "example-cloud-init-template"
  content = <<-EOT
    #cloud-config
    package_update: true
    package_upgrade: true
    packages:
      - nginx
      - curl
    runcmd:
      - systemctl enable nginx
      - systemctl start nginx
  EOT
}

data "thalassa_volume_type" "block" {
  name = "Block"
}

data "thalassa_machine_image" "ubuntu" {
  name = "ubuntu-22-04-01"
}

# Create a virtual machine instance with Thalassa default values
resource "thalassa_virtual_machine_instance" "example" {
  name                   = "example-instance"
  subnet_id              = thalassa_subnet.example.id
  machine_type           = "pgp-small" # Available: pgp-small, pgp-medium, pgp-large, pgp-xlarge, pgp-2xlarge, pgp-4xlarge, dgp-small, dgp-medium, dgp-large, dgp-xlarge, dgp-2xlarge, dgp-4xlarge
  machine_image          = data.thalassa_machine_image.ubuntu.name
  availability_zone      = "nl-01a" # Available: nl-01a, nl-01b, nl-01c
  root_volume_size_gb    = 20
  root_volume_type       = data.thalassa_volume_type.block.id
  cloud_init_template_id = thalassa_cloud_init_template.example.id
}

# Output the virtual machine instance details
output "instance_id" {
  value = thalassa_virtual_machine_instance.example.id
}

output "instance_name" {
  value = thalassa_virtual_machine_instance.example.name
}

# Create a load balancer for the virtual machine instance
resource "thalassa_loadbalancer" "example" {
  name        = "example-lb"
  region      = "nl-01"
  description = "Example load balancer for virtual machine instance"
  subnet_id   = thalassa_subnet.example.id
}

# Create a load balancer listener
resource "thalassa_loadbalancer_listener" "example" {
  name            = "example-lb-listener"
  description     = "Example load balancer listener for virtual machine instance"
  loadbalancer_id = thalassa_loadbalancer.example.id
  protocol        = "tcp"
  port            = 22
  target_group_id = thalassa_target_group.example.id
}

# Create a load balancer target group
resource "thalassa_target_group" "example" {
  name        = "example-lb-target-group"
  description = "Example load balancer target group for virtual machine instance"
  vpc_id      = thalassa_vpc.example.id
  protocol    = "tcp"
  port        = 22
}

resource "thalassa_target_group_attachment" "example" {
  target_group_id = thalassa_target_group.example.id
  vmi_id          = thalassa_virtual_machine_instance.example.id
}

# Output the load balancer details
output "load_balancer_ip" {
  value = thalassa_loadbalancer.example.ip_address
}

output "load_balancer_port" {
  value = thalassa_loadbalancer_listener.example.port
}
```
<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `machine_image` (String) Machine image of the virtual machine instance
- `machine_type` (String) Machine type of the virtual machine instance
- `name` (String) Name of the Virtual Machine Instance
- `subnet_id` (String) Subnet of the Virtual Machine Instance

### Optional

- `annotations` (Map of String) Annotations for the virtual machine instance
- `availability_zone` (String) Availability zone of the virtual machine instance
- `cloud_init` (String) Cloud init of the virtual machine instance
- `cloud_init_template_id` (String) Cloud init template id of the virtual machine instance. If provided, the cloud init will be set to the content of the template.
- `delete_protection` (Boolean) Delete protection of the virtual machine instance
- `description` (String) A human readable description about the virtual machine instance
- `labels` (Map of String) Labels for the virtual machine instance
- `organisation_id` (String) Reference to the Organisation of the Machine Type. If not provided, the organisation of the (Terraform) provider will be used.
- `root_volume_id` (String) Root volume id of the virtual machine instance. Must be provided if root_volume_type is not set.
- `root_volume_size_gb` (Number) Root volume size of the virtual machine instance. Must be provided if root_volume_id is not set.
- `root_volume_type` (String) Root volume type of the virtual machine instance. Must be provided if root_volume_id is not set.
- `security_group_attachments` (List of String) Security group attached to the virtual machine instance

### Read-Only

- `attached_volume_ids` (List of String) Attached volume ids of the virtual machine instance
- `id` (String) The ID of this resource.
- `ip_addresses` (List of String) IP addresses of the virtual machine instance
- `slug` (String) Slug of the Virtual Machine Instance
- `state` (String) Desired state of the virtual machine instance. Can be 'running', 'stopped', 'deleted'
- `status` (String) Status of the virtual machine instance


