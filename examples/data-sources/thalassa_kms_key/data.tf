data "thalassa_kms_key" "existing" {
  region   = "nl-01"
  identity = "kms-abc123"
}
