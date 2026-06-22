resource "thalassa_secret_access_policy" "db_password" {
  region = "nl-01"
  path   = thalassa_secret.db_password.path

  statement {
    effect     = "Allow"
    actions    = ["read"]
    principals = ["team:platform"]
  }
}
