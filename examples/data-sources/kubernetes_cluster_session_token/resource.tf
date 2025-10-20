data "thalassa_kubernetes_cluster" "cluster" {
  name            = "example"
  organisation_id = var.organisation_id
}

data "thalassa_kubernetes_cluster_session_token" "cluster_session_token" {
  cluster_id = data.thalassa_kubernetes_cluster.cluster.id
}

provider "kubernetes" {
  host                   = data.thalassa_kubernetes_cluster_session_token.cluster_session_token.api_server_url
  cluster_ca_certificate = base64decode(data.thalassa_kubernetes_cluster_session_token.cluster_session_token.ca_certificate)
  token                  = data.thalassa_kubernetes_cluster_session_token.cluster_session_token.token
}

resource "kubernetes_namespace" "example" {
  metadata {
    name = "example"
  }
}
