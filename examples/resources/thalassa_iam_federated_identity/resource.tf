resource "thalassa_iam_service_account" "ci" {
  name        = "ci-deployer"
  description = "Service account used by CI federation"
}

resource "thalassa_iam_federated_identity_provider" "github" {
  name            = "github-actions"
  provider_issuer = "https://token.actions.githubusercontent.com"
  status          = "active"
}

resource "thalassa_iam_federated_identity" "ci" {
  name                     = "github-ci"
  description              = "Federated identity for GitHub Actions"
  service_account_id = thalassa_iam_service_account.ci.id
  provider_id        = thalassa_iam_federated_identity_provider.github.id
  provider_subject         = "repo:example-org/example-repo:ref:refs/heads/main"
  trusted_audiences        = ["https://api.thalassa.cloud"]
  audience_match_mode      = "any"
  allowed_scopes = [
    "api:read",
    "api:write",
  ]
  status = "active"
}

output "federated_identity_id" {
  value = thalassa_iam_federated_identity.ci.id
}
