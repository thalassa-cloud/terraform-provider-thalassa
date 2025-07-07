# Create a VPC with all optional attributes
resource "thalassa_vpc" "example" {
  name   = "example-vpc"
  region = "nl-01"
  cidrs  = ["10.0.0.0/16", "10.2.0.0/16", "10.3.0.0/16"]

  # Optional attributes
  description = "Example VPC for documentation with all optional attributes"

  # Labels are key-value pairs for organizing resources
  labels = {
    environment = "production"
    project     = "example-project"
    owner       = "team-a"
  }
}

# Create a subnet within the VPC
resource "thalassa_subnet" "example" {
  name        = "example-subnet"
  description = "Example subnet for documentation"
  vpc_id      = thalassa_vpc.example.id
  cidr        = "10.0.1.0/24"
}

# Output the VPC and subnet IDs
output "vpc_id" {
  value = thalassa_vpc.example.id
}

output "subnet_id" {
  value = thalassa_subnet.example.id
}
