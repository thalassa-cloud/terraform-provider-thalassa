
variable "region" {
  type        = string
  description = "Region for the VPC"
  default     = "nl-01"
}

# Create a VPC for the NAT gateway
resource "thalassa_vpc" "example" {
  name        = "example-vpc"
  description = "Example VPC for NAT gateway"
  region      = var.region
  cidrs       = ["10.0.0.0/16"]
}

# Create a subnet for the NAT gateway
resource "thalassa_subnet" "example" {
  name        = "example-subnet"
  description = "Example subnet for NAT gateway"
  vpc_id      = thalassa_vpc.example.id
  cidr        = "10.0.1.0/24"
}

resource "thalassa_natgateway" "example" {
  name        = "example-nat-gateway"
  subnet_id   = thalassa_subnet.example.id
  description = "Example NAT gateway"

  labels = {
    environment = "production"
    tier        = "public"
  }
}

# Output the NAT gateway details
output "nat_gateway_id" {
  value = thalassa_natgateway.example.id
}

output "nat_gateway_endpoint_ip" {
  value = thalassa_natgateway.example.endpoint_ip
}
