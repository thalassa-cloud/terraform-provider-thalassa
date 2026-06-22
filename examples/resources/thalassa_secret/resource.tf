resource "thalassa_kms_key" "app" {
  region   = "nl-01"
  name     = "app-secrets"
  key_type = "aes256-gcm96"
}

resource "thalassa_secret" "db_password" {
  region           = "nl-01"
  path             = "/app/prod/db/password"
  kms_key_id = thalassa_kms_key.app.id
}

resource "thalassa_secret_version" "db_password" {
  region        = "nl-01"
  path          = thalassa_secret.db_password.path
  generate_secret {
    byte_length = 32
  }
}
