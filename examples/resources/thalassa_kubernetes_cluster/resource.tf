# Create a VPC for the Kubernetes cluster
variable "region" {
  type        = string
  description = "Region for the Kubernetes cluster"
  default     = "nl-01"
}

resource "thalassa_vpc" "example" {
  name        = "example-vpc"
  description = "Example VPC for Kubernetes cluster"
  region      = var.region
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
  region      = var.region
  subnet_id   = thalassa_subnet.example.id
  # If you wish to restrict the public API endpoint access, you can add API server ACLs. Leaving this empty will allow all traffic.
  # For the VPC internal endpoint, use security groups instead
  # api_server_acls {
  #   allowed_cidrs = ["10.0.0.0/16", "10.0.1.0/24"]
  # }
}

# Create a Kubernetes node pool with Thalassa default values
resource "thalassa_kubernetes_node_pool" "example" {
  name              = "example-node-pool"
  cluster_id        = thalassa_kubernetes_cluster.example.id
  subnet_id         = thalassa_subnet.example.id
  availability_zone = "${var.region}a"
  machine_type      = "pgp-small"
  # replicas              = 2
  enable_autoscaling = true
  min_replicas       = 1
  max_replicas       = 2
  # kubernetes_version    = "1-34-1"
}

# Output the Kubernetes cluster details
output "kubernetes_cluster_id" {
  value = thalassa_kubernetes_cluster.example.id
}

output "kubernetes_cluster_name" {
  value = thalassa_kubernetes_cluster.example.name
}

output "advertise_address" {
  value = thalassa_kubernetes_cluster.example.internal_endpoint
}

# # Output the node pool details
output "node_pool_id" {
  value = thalassa_kubernetes_node_pool.example.id
}

output "node_pool_name" {
  value = thalassa_kubernetes_node_pool.example.name
}
