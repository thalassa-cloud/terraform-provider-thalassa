
data "thalassa_organisation" "example" {
  slug = "your-organisation-slug"
}

resource "thalassa_block_volume" "example" {
  name         = "example"
  organisation_id = data.thalassa_organisation.example.id
  region       = "your-region"
  size_gb      = 15
  volume_type  = "your-volume-type-id"
}
