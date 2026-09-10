resource "thalassa_iam_federated_identity_provider" "github" {
  name            = "github-actions"
  description     = "GitHub Actions OIDC provider"
  provider_issuer = "https://token.actions.githubusercontent.com"
  status          = "active"

  labels = {
    environment = "ci"
  }
}

output "provider_id" {
  value = thalassa_iam_federated_identity_provider.github.id
}
