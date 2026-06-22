resource "thalassa_kms_key" "app" {
  region   = "nl-01"
  name     = "app-secrets"
  key_type = "aes256-gcm96"

  key_rotation_enabled    = true
  rotation_period_in_days = 90
}
