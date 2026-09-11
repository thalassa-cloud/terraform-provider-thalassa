resource "thalassa_containerregistry_namespace" "example" {
  region    = "nl-01"
  namespace = "example"
}

resource "thalassa_containerregistry_namespace_configuration" "example" {
  namespace_id = thalassa_containerregistry_namespace.example.id
  visibility   = "private"

  retention_policy {
    enabled                = true
    delete_untagged_images = true

    rules {
      days  = 30
      count = 10
      scope = "tags"
      tag_patterns = [
        "v*",
      ]
    }
  }
}
