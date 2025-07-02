terraform {
  required_providers {
    thalassa = {
      source = "local/thalassa/thalassa"
    }
  }
}

provider "thalassa" {
  # Configuration options
}

# Create a VPC for the Kubernetes cluster
resource "thalassa_vpc" "example" {
  name            = "example-vpc"
  description     = "Example VPC for Kubernetes cluster"
  region          = "nl-01"
  cidrs           = ["10.0.0.0/16"]
}

# Create a subnet for the Kubernetes cluster
resource "thalassa_subnet" "example" {
  name            = "example-subnet"
  description     = "Example subnet for Kubernetes cluster"
  vpc_id          = thalassa_vpc.example.id
  cidr            = "10.0.1.0/24"
}

# Create a Kubernetes cluster
resource "thalassa_kubernetes_cluster" "example" {
  name            = "example-kubernetes-cluster"
  description     = "Example Kubernetes cluster"
  region          = "nl-01"
  version         = "1.28"
}

# Create a Kubernetes node pool with Thalassa default values
resource "thalassa_kubernetes_node_pool" "example" {
  kubernetes_cluster_id = thalassa_kubernetes_cluster.example.id
  name                  = "example-node-pool"
  machine_type          = "pgp-small"  # Available: pgp-small, pgp-medium, pgp-large, pgp-xlarge, pgp-2xlarge, pgp-4xlarge, dgp-small, dgp-medium, dgp-large, dgp-xlarge, dgp-2xlarge, dgp-4xlarge
  machine_image         = "ubuntu-22.04"
  node_count            = 3
}

# Output the Kubernetes cluster details
output "kubernetes_cluster_id" {
  value = thalassa_kubernetes_cluster.example.id
}

output "kubernetes_cluster_name" {
  value = thalassa_kubernetes_cluster.example.name
}

# Output the node pool details
output "node_pool_id" {
  value = thalassa_kubernetes_node_pool.example.id
}

output "node_pool_name" {
  value = thalassa_kubernetes_node_pool.example.name
} 
