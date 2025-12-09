# Create a VPC for the security group
resource "thalassa_vpc" "example" {
  name        = "example-vpc"
  description = "Example VPC for security group egress rules"
  region      = "nl-01"
  cidrs       = ["10.0.0.0/16"]
}

# Create a security group
resource "thalassa_security_group" "example" {
  name        = "example-security-group"
  description = "Example security group"
  vpc_id      = thalassa_vpc.example.id
  allow_same_group_traffic = false
}

# Create egress rules for the security group
resource "thalassa_security_group_egress_rule" "example" {
  security_group_id = thalassa_security_group.example.id

  rule {
    name           = "allow-all-outbound"
    ip_version     = "ipv4"
    protocol       = "all"
    priority       = 100
    remote_type    = "address"
    remote_address = "0.0.0.0/0"
    policy         = "allow"
  }

  rule {
    name           = "allow-https-outbound"
    ip_version     = "ipv4"
    protocol       = "tcp"
    priority       = 101
    remote_type    = "address"
    remote_address = "0.0.0.0/0"
    port_range_min = 443
    port_range_max = 443
    policy         = "allow"
  }
}

# Output the security group ID
output "security_group_id" {
  value = thalassa_security_group.example.id
}

