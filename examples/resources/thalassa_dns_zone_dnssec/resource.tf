resource "thalassa_dns_zone_dnssec" "example" {
  zone_id = thalassa_dns_zone.example.id
  region  = "nl-01"
  # kms_key_id = thalassa_kms_key.dns.id  # optional; auto-provisioned if omitted
}
