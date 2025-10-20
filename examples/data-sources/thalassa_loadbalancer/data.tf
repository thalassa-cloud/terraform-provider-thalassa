# Get a Load Balancer by name
data "thalassa_loadbalancer" "by_name" {
  name = "my-loadbalancer"
}

# Get a Load Balancer by slug
data "thalassa_loadbalancer" "by_slug" {
  slug = "my-loadbalancer-slug"
}

# Get a Load Balancer by VPC ID
data "thalassa_loadbalancer" "by_vpc" {
  vpc_id = "vpc-123"
}

# Get a Load Balancer by labels
data "thalassa_loadbalancer" "by_labels" {
  labels = {
    "environment" = "production"
    "team"        = "platform"
  }
}
