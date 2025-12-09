# Create a VPC for the security group
resource "thalassa_vpc" "example" {
  name        = "example-vpc"
  description = "Example VPC for security group ingress rules"
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

# Create ingress rules for the security group
resource "thalassa_security_group_ingress_rule" "example" {
  security_group_id = thalassa_security_group.example.id

  rule {
    name           = "allow-http"
    ip_version     = "ipv4"
    protocol       = "tcp"
    priority       = 100
    remote_type    = "address"
    remote_address = "0.0.0.0/0"
    port_range_min = 80
    port_range_max = 80
    policy         = "allow"
  }

  rule {
    name           = "allow-https"
    ip_version     = "ipv4"
    protocol       = "tcp"
    priority       = 101
    remote_type    = "address"
    remote_address = "0.0.0.0/0"
    port_range_min = 443
    port_range_max = 443
    policy         = "allow"
  }

  rule {
    name           = "allow-ssh"
    ip_version     = "ipv4"
    protocol       = "tcp"
    priority       = 102
    remote_type    = "address"
    remote_address = "10.0.0.0/8"
    port_range_min = 22
    port_range_max = 22
    policy         = "allow"
  }
}

# Output the security group ID
output "security_group_id" {
  value = thalassa_security_group.example.id
}

