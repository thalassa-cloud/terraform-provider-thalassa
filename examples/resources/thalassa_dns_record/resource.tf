resource "thalassa_dns_record" "www" {
  zone_id = thalassa_dns_zone.example.id
  name    = "www"
  type    = "A"
  ttl     = 300
  values  = ["192.0.2.1"]
}

resource "thalassa_dns_record" "apex" {
  zone_id = thalassa_dns_zone.example.id
  name    = "@"
  type    = "A"
  values  = ["192.0.2.1"]
}
