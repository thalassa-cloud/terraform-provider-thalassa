resource "thalassa_objectstorage_bucket" "logs" {
  name       = "my-logs-bucket"
  region     = "nl-01"
  versioning = true
}

resource "thalassa_objectstorage_bucket_lifecycle" "logs" {
  bucket_name = thalassa_objectstorage_bucket.logs.name

  rule {
    id     = "expire-logs"
    prefix = "logs/"
    status = "Enabled"
    expiration {
      days = 30
    }
  }
}
