# Get a NAT Gateway by name
data "thalassa_natgateway" "by_name" {
  name = "my-nat-gateway"
}

# Get a NAT Gateway by slug
data "thalassa_natgateway" "by_slug" {
  slug = "my-nat-gateway-slug"
}

# Get a NAT Gateway by VPC ID
data "thalassa_natgateway" "by_vpc" {
  vpc_id = "vpc-123"
}

# Get a NAT Gateway by subnet ID
data "thalassa_natgateway" "by_subnet" {
  subnet_id = "subnet-456"
}

# Get a NAT Gateway by region
data "thalassa_natgateway" "by_region" {
  region = "nl-01"
  name   = "my-nat-gateway"
}

# Get a NAT Gateway by labels
data "thalassa_natgateway" "by_labels" {
  labels = {
    "environment" = "production"
    "team"        = "platform"
  }
}

# Get a NAT Gateway with multiple filters
data "thalassa_natgateway" "filtered" {
  name   = "my-nat-gateway"
  vpc_id = "vpc-production"
  region = "nl-01"
  labels = {
    "environment" = "production"
    "team"        = "platform"
  }
}
