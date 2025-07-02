terraform {
  required_providers {
    thalassa = {
      source = "local/thalassa/thalassa"
    }
  }
}

provider "thalassa" {
  # Configuration options
}

# Create a VPC for the virtual machine instance
resource "thalassa_vpc" "example" {
  name            = "example-vpc"
  description     = "Example VPC for virtual machine instance"
  region          = "nl-01"
  cidrs           = ["10.0.0.0/16"]
}

# Create a subnet for the virtual machine instance
resource "thalassa_subnet" "example" {
  name            = "example-subnet"
  description     = "Example subnet for virtual machine instance"
  vpc_id          = thalassa_vpc.example.id
  cidr            = "10.0.1.0/24"
}

# Create a security group for the virtual machine instance
resource "thalassa_security_group" "example" {
  name        = "example-security-group"
  description = "Example security group for virtual machine instance"
  vpc_identity = thalassa_vpc.example.id
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
  name                = "example-instance"
  subnet_id           = thalassa_subnet.example.id
  machine_type        = "pgp-small"  # Available: pgp-small, pgp-medium, pgp-large, pgp-xlarge, pgp-2xlarge, pgp-4xlarge, dgp-small, dgp-medium, dgp-large, dgp-xlarge, dgp-2xlarge, dgp-4xlarge
  machine_image       = data.thalassa_machine_image.ubuntu.name
  availability_zone   = "portable1a"     # Available: nl-01a, nl-01b, nl-01c
  root_volume_size_gb = 20
  root_volume_type    = data.thalassa_volume_type.block.id
}

# Output the virtual machine instance details
output "instance_id" {
  value = thalassa_virtual_machine_instance.example.id
}

output "instance_name" {
  value = thalassa_virtual_machine_instance.example.name
} 


## Expose the virtual machine instance to the internet

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
  load_balancer_id = thalassa_load_balancer.example.id
  protocol        = "tcp"
  port            = 2200
  target_port     = 22 # SSH port on the VM
}

# Create a load balancer target group
resource "thalassa_loadbalancer_target_group" "example" {
  name            = "example-lb-target-group"
  description     = "Example load balancer target group for virtual machine instance"
  load_balancer_id = thalassa_load_balancer.example.id
  protocol        = "tcp"
  port            = 22
  targets         = [thalassa_virtual_machine_instance.example.id]
}

# Output the load balancer details
output "load_balancer_ip" {
  value = thalassa_loadbalancer.example.ip_address
}

output "load_balancer_port" {
  value = thalassa_loadbalancer_listener.example.port
}
