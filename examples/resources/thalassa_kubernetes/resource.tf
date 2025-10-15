# Create a VPC for the Kubernetes cluster
resource "thalassa_vpc" "example" {
  name        = "example-vpc"
  description = "Example VPC for Kubernetes cluster"
  region      = "nl-01"
  cidrs       = ["10.0.0.0/16"]
}

# Create a subnet for the Kubernetes cluster
resource "thalassa_subnet" "example" {
  name        = "example-subnet"
  description = "Example subnet for Kubernetes cluster"
  vpc_id      = thalassa_vpc.example.id
  cidr        = "10.0.1.0/24"
}

# Create a Kubernetes cluster
resource "thalassa_kubernetes_cluster" "example" {
  name        = "example-kubernetes-cluster"
  description = "Example Kubernetes cluster"
  region      = "nl-01"
  subnet_id   = thalassa_subnet.example.id
  api_server_acls {
    allowed_cidrs = ["10.0.0.0/16", "178.85.83.115/32"]
  }
}

# Create a Kubernetes node pool with Thalassa default values
resource "thalassa_kubernetes_node_pool" "example" {
  name                  = "example-node-pool"
  cluster_id = thalassa_kubernetes_cluster.example.id
  subnet_id = thalassa_subnet.example.id
  availability_zone     = "nl-01a"
  machine_type          = "pgp-small"
  # replicas              = 2
  enable_autoscaling    = true
  min_replicas          = 1
  max_replicas          = 2
  # kubernetes_version    = "1-34-1"
}

# Output the Kubernetes cluster details
output "kubernetes_cluster_id" {
  value = thalassa_kubernetes_cluster.example.id
}

output "kubernetes_cluster_name" {
  value = thalassa_kubernetes_cluster.example.name
}

# # Output the node pool details
output "node_pool_id" {
  value = thalassa_kubernetes_node_pool.example.id
}

output "node_pool_name" {
  value = thalassa_kubernetes_node_pool.example.name
}
