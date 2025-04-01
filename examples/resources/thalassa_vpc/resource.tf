terraform {
  required_providers {
    thalassa = {
      source = "thalassa-cloud/thalassa"
    }
  }
}

provider "thalassa" {
  # Configuration options
}

# Create a VPC
resource "thalassa_vpc" "example" {
  name        = "example-vpc"
  description = "Example VPC for documentation"
  region      = "eu-west-1"  # Replace with your desired region
  cidrs       = ["10.0.0.0/16", "10.2.0.0/16", "10.3.0.0/16"]
}

# Create a subnet within the VPC
resource "thalassa_subnet" "example" {
  name        = "example-subnet"
  description = "Example subnet for documentation"
  vpc_id      = thalassa_vpc.example.id
  cidr        = ["10.0.1.0/24", "10.2.1.0/24", "10.3.1.0/24"]
  zone        = "nl-01"  # Replace with your desired zone
}

# Output the VPC and subnet IDs
output "vpc_id" {
  value = thalassa_vpc.example.id
}

output "subnet_id" {
  value = thalassa_subnet.example.id
}
