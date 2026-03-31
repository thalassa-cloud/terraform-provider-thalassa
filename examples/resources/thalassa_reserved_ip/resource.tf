# VPC and subnet in the region where the reserved IP is allocated
resource "thalassa_vpc" "example" {
  name        = "example-vpc-reserved-ip"
  description = "Example VPC for reserved IP"
  region      = "nl-01"
  cidrs       = ["10.0.0.0/16"]
}

resource "thalassa_subnet" "example" {
  name        = "example-subnet-nat"
  description = "Subnet for NAT gateway with a reserved egress IP"
  vpc_id      = thalassa_vpc.example.id
  cidr        = "10.0.1.0/24"
}

resource "thalassa_reserved_ip" "example" {
  name        = "example-reserved-ip"
  description = "Static IP for NAT egress"
  region      = thalassa_vpc.example.region

  labels = {
    environment = "example"
    purpose     = "nat-egress"
  }
}

# Attach the reserved IP when creating the NAT gateway
resource "thalassa_natgateway" "example" {
  name           = "example-nat-gateway"
  subnet_id      = thalassa_subnet.example.id
  description    = "NAT gateway using a reserved IP"
  reserved_ip_id = thalassa_reserved_ip.example.id
}

output "reserved_ip_id" {
  value = thalassa_reserved_ip.example.id
}

output "reserved_ip_ipv4" {
  value = thalassa_reserved_ip.example.ipv4_address
}

output "nat_gateway_id" {
  value = thalassa_natgateway.example.id
}

output "nat_gateway_endpoint_ip" {
  value = thalassa_natgateway.example.endpoint_ip
}
