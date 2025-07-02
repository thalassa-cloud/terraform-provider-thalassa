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

# Create a VPC for the resources
resource "thalassa_vpc" "example" {
  name            = "example-vpc"
  description     = "Example VPC for target group attachment"
  region          = "nl-01"
  cidrs           = ["10.0.0.0/16"]
}

# Create a subnet for the resources
resource "thalassa_subnet" "example" {
  name            = "example-subnet"
  description     = "Example subnet for target group attachment"
  vpc_id          = thalassa_vpc.example.id
  cidr            = "10.0.1.0/24"
}

# Create a target group
resource "thalassa_target_group" "example" {
  name            = "example-target-group"
  description     = "Example target group for attachment"
  vpc_id          = thalassa_vpc.example.id
  protocol        = "http"
  port            = 80
}

# Create a virtual machine instance
resource "thalassa_virtual_machine_instance" "example" {
  name                = "example-instance"
  subnet_id           = thalassa_subnet.example.id
  machine_type        = "pgp-small"  # Available: pgp-small, pgp-medium, pgp-large, pgp-xlarge, pgp-2xlarge, pgp-4xlarge, dgp-small, dgp-medium, dgp-large, dgp-xlarge, dgp-2xlarge, dgp-4xlarge
  machine_image       = "ubuntu-22.04"
  availability_zone   = "nl-01a"     # Available: nl-01a, nl-01b, nl-01c
  root_volume_size_gb = 20
  root_volume_type    = "Block"      # Available: Block, Premium Block
}

# Create a target group attachment with all required attributes
resource "thalassa_target_group_attachment" "example" {
  # Required attributes
  organisation_id = "org-123" # Replace with your organisation ID
  target_group_id = thalassa_target_group.example.id
  vmi_id          = thalassa_virtual_machine_instance.example.id
}

# Output the attachment details
output "attachment_id" {
  value = thalassa_target_group_attachment.example.id
} 