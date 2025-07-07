# Create a VPC for the resources
resource "thalassa_vpc" "example" {
  name        = "example-vpc"
  description = "Example VPC for block volume attachment"
  region      = "nl-01"
  cidrs       = ["10.0.0.0/16"]
}

# Create a subnet for the resources
resource "thalassa_subnet" "example" {
  name        = "example-subnet"
  description = "Example subnet for block volume attachment"
  vpc_id      = thalassa_vpc.example.id
  cidr        = "10.0.1.0/24"
}

# Create a block volume
resource "thalassa_block_volume" "example" {
  name        = "example-block-volume"
  description = "Example block volume for attachment"
  region      = "nl-01"
  volume_type = "Block" # Available: Block, Premium Block
  size_gb     = 50
}

# Create a virtual machine instance
resource "thalassa_virtual_machine_instance" "example" {
  name                = "example-instance"
  subnet_id           = thalassa_subnet.example.id
  machine_type        = "pgp-small" # Available: pgp-small, pgp-medium, pgp-large, pgp-xlarge, pgp-2xlarge, pgp-4xlarge, dgp-small, dgp-medium, dgp-large, dgp-xlarge, dgp-2xlarge, dgp-4xlarge
  machine_image       = "ubuntu-22.04"
  availability_zone   = "nl-01a" # Available: nl-01a, nl-01b, nl-01c
  root_volume_size_gb = 20
  root_volume_type    = "Block" # Available: Block, Premium Block
}

# Create a block volume attachment with Thalassa default values
resource "thalassa_block_volume_attachment" "example" {
  volume_id = thalassa_block_volume.example.id
  vmi_id    = thalassa_virtual_machine_instance.example.id
}

# Output the attachment details
output "attachment_id" {
  value = thalassa_block_volume_attachment.example.id
}

output "device_serial" {
  value = thalassa_block_volume_attachment.example.serial
}
