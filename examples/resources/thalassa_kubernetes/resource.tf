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

# Create a Kubernetes cluster
resource "thalassa_kubernetes_cluster" "example" {
  name        = "example-cluster"
  description = "Example Kubernetes cluster for documentation"
  region      = "nl-01"  # Replace with your desired region
  subnet_id   = "subnet-123" # Replace with your subnet ID
  cluster_version = "1.28.0" # Replace with your desired Kubernetes version
  cluster_type    = "managed"
  networking_cni  = "cilium"
  networking_service_cidr = "10.96.0.0/12"
  networking_pod_cidr     = "10.244.0.0/16"
}

# Create a node pool for the cluster
resource "thalassa_kubernetes_node_pool" "example" {
  name              = "example-node-pool"
  description       = "Example node pool for documentation"
  cluster_id        = thalassa_kubernetes_cluster.example.id
  region            = thalassa_kubernetes_cluster.example.region
}

# Output the cluster and node pool IDs
output "cluster_id" {
  value = thalassa_kubernetes_cluster.example.id
}

output "node_pool_id" {
  value = thalassa_kubernetes_node_pool.example.id
} 
