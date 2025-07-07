resource "thalassa_vpc" "example" {
  name        = "example-vpc"
  description = "Example VPC for security group"
  region      = "nl-01"
  cidrs       = ["10.0.0.0/16"]
}

# Create a security group
resource "thalassa_security_group" "example" {
  name        = "example-security-group"
  description = "Example security group for documentation"
  vpc_id      = thalassa_vpc.example.id

  allow_same_group_traffic = false

  ingress_rule {
    name           = "allow-http"
    ip_version     = "ipv4"
    protocol       = "tcp"
    priority       = 100
    remote_type    = "address"
    remote_address = "10.0.0.0/0"
    port_range_min = 80
    port_range_max = 80
    policy         = "allow"
  }
  ingress_rule {
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

  egress_rule {
    name           = "allow-all"
    ip_version     = "ipv4"
    protocol       = "all"
    priority       = 100
    remote_type    = "address"
    remote_address = "0.0.0.0/0"
    policy         = "allow"
  }
}

# Output the security group details
output "security_group_id" {
  value = thalassa_security_group.example.id
}

output "security_group_name" {
  value = thalassa_security_group.example.name
}
