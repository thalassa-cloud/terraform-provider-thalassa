resource "thalassa_vpc" "example" {
  name        = "example-vpc"
  description = "Example VPC for security group"
  region      = "nl-01"
  cidrs       = ["10.0.0.0/16"]
}

# Create a security group
resource "thalassa_security_group" "example" {
  name        = "example-security-group"
  description = "Example security group"
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


# Create a security group with rules managed separately
resource "thalassa_security_group" "controlplane" {
  name        = "controlplane-security-group"
  description = "Control plane security group"
  vpc_id      = thalassa_vpc.example.id
  allow_same_group_traffic = false
}

// ingress rules
resource "thalassa_security_group_ingress_rule" "controlplane" {
  security_group_id = thalassa_security_group.controlplane.id
  rule {
    name           = "allow-http"
    ip_version     = "ipv4"
    protocol       = "tcp"
    priority       = 100
    policy         = "allow"
    remote_type    = "securityGroup"
    remote_security_group_identity = thalassa_security_group.cluster.id
    port_range_min = 80
    port_range_max = 80
  }

  rule {
    name           = "allow-ssh"
    ip_version     = "ipv4"
    protocol       = "tcp"
    priority       = 100
    policy         = "allow"
    remote_type    = "address"
    remote_address = "0.0.0.0/0"
    port_range_min = 22
    port_range_max = 22
  }
}

resource "thalassa_security_group_egress_rule" "controlplane" {
  security_group_id = thalassa_security_group.controlplane.id
  rule {
    name           = "allow-all"
    ip_version     = "ipv4"
    protocol       = "all"
    priority       = 100
    policy         = "allow"
    remote_type    = "address"
    remote_address = "0.0.0.0/0"
  }
}

# Create a security group
resource "thalassa_security_group" "cluster" {
  name        = "cluster-security-group"
  description = "Cluster security group"
  vpc_id      = thalassa_vpc.example.id
  allow_same_group_traffic = false
}

resource "thalassa_security_group_egress_rule" "cluster" {
  security_group_id = thalassa_security_group.cluster.id
  rule {
    name           = "allow-controlplane"
    ip_version     = "ipv4"
    protocol       = "tcp"
    priority       = 100
    policy         = "allow"
    remote_type    = "securityGroup"
    remote_security_group_identity = thalassa_security_group.controlplane.id
  }
}
