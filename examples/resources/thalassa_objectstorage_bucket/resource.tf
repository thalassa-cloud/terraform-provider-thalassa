# Create a basic object storage bucket
resource "thalassa_objectstorage_bucket" "basic" {
  name   = "my-basic-bucket"
  region = "nl-01"
}

# Create a public object storage bucket
resource "thalassa_objectstorage_bucket" "public" {
  name   = "my-public-bucket"
  region = "nl-01"
  public = true
}

# Create a bucket with a custom policy
resource "thalassa_objectstorage_bucket" "with_policy" {
  name   = "my-policy-bucket"
  region = "nl-01"
  public = false

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "AllowReadAccess"
        Effect = "Allow"
        Principal = {
          Thalassa = "*"
        }
        Action = [
          "s3:GetObject"
        ]
        Resource = [
          "arn:thalassa:s3:::my-policy-bucket/*"
        ]
        Condition = {
          StringEquals = {
            "thalassa:User" = "u-exampleuserid"
          }
        }
      }
    ]
  })
}

# Output the bucket details
output "basic_bucket_id" {
  value = thalassa_objectstorage_bucket.basic.id
}

output "basic_bucket_name" {
  value = thalassa_objectstorage_bucket.basic.name
}

output "basic_bucket_endpoint" {
  value = thalassa_objectstorage_bucket.basic.endpoint
}

output "public_bucket_id" {
  value = thalassa_objectstorage_bucket.public.id
}

output "public_bucket_name" {
  value = thalassa_objectstorage_bucket.public.name
}

output "policy_bucket_id" {
  value = thalassa_objectstorage_bucket.with_policy.id
}

output "policy_bucket_name" {
  value = thalassa_objectstorage_bucket.with_policy.name
} 