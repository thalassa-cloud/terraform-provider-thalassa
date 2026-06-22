resource "thalassa_secret_version" "db_password" {
  region        = "nl-01"
  path          = thalassa_secret.db_password.path
  secret_string = "initial-password"
}
