
# Create a service account first
resource "thalassa_iam_service_account" "example" {
  name        = "cluster-service-account"
  description = "Service account for cluster access"
  labels = {
    environment = "production"
    project     = "cluster"
  }
}

# Create object storage access credentials
resource "thalassa_iam_service_account_access_credential" "storage_credential" {
  service_account_id = thalassa_iam_service_account.example.id
  scopes = [
    "objectStorage"
  ]
}

# # random uuid
resource "random_uuid" "bucket_name" {
}

data "thalassa_organisation" "org" {
  slug = var.organisation_slug
}

# # Create a bucket with a custom policy
resource "thalassa_objectstorage_bucket" "cluster_bucket" {
  name   = "cluster-bucket-${random_uuid.bucket_name.result}"
  region = "nl-01"
  public = false

  policy = jsonencode({
    "Version" : "2012-10-17",
    "Statement" : [
      {
        "Sid" : "Statement1",
        "Action" : [
          "s3:GetObject",
          "s3:GetObjectVersion",
          "s3:PutObject",
          "s3:GetObjectAcl",
          "s3:GetObjectVersionAcl",
          "s3:PutObjectAcl",
          "s3:PutObjectVersionAcl",
          "s3:DeleteObject",
          "s3:DeleteObjectVersion",
          "s3:ListMultipartUploadParts",
          "s3:AbortMultipartUpload",
          "s3:RestoreObject",
          "s3:ListBucket",
          "s3:ListBucketVersions",
          "s3:ListBucketMultipartUploads",
          "s3:GetBucketAcl",
          "s3:PutBucketAcl",
          "s3:GetBucketCORS",
          "s3:PutBucketCORS",
          "s3:GetBucketVersioning",
          "s3:PutBucketVersioning",
          "s3:GetBucketRequestPayment",
          "s3:PutBucketRequestPayment",
          "s3:GetLifecycleConfiguration",
          "s3:PutLifecycleConfiguration",
          "s3:GetObjectTagging",
          "s3:PutObjectTagging",
          "s3:DeleteObjectTagging",
          "s3:GetObjectVersionTagging",
          "s3:PutObjectVersionTagging",
          "s3:DeleteObjectVersionTagging",
          "s3:PutBucketObjectLockConfiguration",
          "s3:GetBucketObjectLockConfiguration",
          "s3:PutObjectRetention",
          "s3:GetObjectRetention",
          "s3:PutObjectLegalHold",
          "s3:GetObjectLegalHold",
          "s3:BypassGovernanceRetention",
          "s3:GetBucketPolicyStatus"
        ],
        "Effect" : "Allow",
        "Resource" : [
          "arn:thalassa:s3:::cluster-bucket-${random_uuid.bucket_name.result}",
          "arn:thalassa:s3:::cluster-bucket-${random_uuid.bucket_name.result}/*"
        ],
        "Principal" : {
          "Thalassa" : [
            "arn:thalassa:iam:::serviceaccount/${data.thalassa_organisation.org.id}:${thalassa_iam_service_account.example.id}"
          ]
        }
      }
    ]
  })
}
