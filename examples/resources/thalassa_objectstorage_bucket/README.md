# Thalassa Object Storage Bucket Example

This example demonstrates how to create and manage object storage buckets using the Thalassa Cloud provider.

## Overview

The `thalassa_objectstorage_bucket` resource allows you to create and manage object storage buckets in the Thalassa Cloud platform. Object storage buckets are used to store and organize objects (files) in the cloud.

## Features

- **Basic Bucket Creation**: Create simple object storage buckets
- **Public Access Control**: Configure buckets for public or private access
- **Custom Policies**: Apply S3-compatible bucket policies for fine-grained access control
- **Region Selection**: Deploy buckets in specific regions
- **Status Monitoring**: Track bucket status and endpoint information

## Example Usage

### Basic Bucket

```hcl
resource "thalassa_objectstorage_bucket" "basic" {
  name   = "my-basic-bucket"
  region = "nl-01"
}
```

### Public Bucket

```hcl
resource "thalassa_objectstorage_bucket" "public" {
  name   = "my-public-bucket"
  region = "nl-01"
  public = true
}
```

### Bucket with Custom Policy

```hcl
resource "thalassa_objectstorage_bucket" "with_policy" {
  name   = "my-policy-bucket"
  region = "nl-01"
  public = false
  
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid       = "AllowReadAccess"
        Effect    = "Allow"
        Principal = {
          AWS = "*"
        }
        Action = [
          "s3:GetObject"
        ]
        Resource = [
          "arn:aws:s3:::my-policy-bucket/*"
        ]
        Condition = {
          StringEquals = {
            "aws:PrincipalOrgID" = "o-exampleorgid"
          }
        }
      }
    ]
  })
}
```

## Resource Arguments

| Argument | Type | Required | Description |
|----------|------|----------|-------------|
| `name` | string | Yes | The name of the bucket (1-63 characters) |
| `region` | string | Yes | The region where the bucket will be created |
| `public` | bool | No | Whether the bucket is publicly accessible (default: false) |
| `policy` | string | No | JSON-encoded bucket policy for access control |
| `organisation_id` | string | No | Organisation ID (uses provider default if not specified) |

## Resource Attributes

| Attribute | Type | Description |
|-----------|------|-------------|
| `id` | string | The unique identifier of the bucket |
| `name` | string | The name of the bucket |
| `region` | string | The region of the bucket |
| `public` | bool | Whether the bucket is publicly accessible |
| `policy` | string | The bucket policy as a JSON string |
| `status` | string | The current status of the bucket |
| `endpoint` | string | The endpoint URL for accessing the bucket |

## Data Source Usage

You can also use the `thalassa_objectstorage_bucket` data source to reference existing buckets:

```hcl
data "thalassa_objectstorage_bucket" "existing" {
  name   = "existing-bucket-name"
  region = "nl-01"
}

# Reference the existing bucket
resource "thalassa_objectstorage_bucket" "new" {
  name   = "new-bucket"
  region = data.thalassa_objectstorage_bucket.existing.region
}
```

## Best Practices

1. **Naming**: Use descriptive, unique bucket names that follow your organization's naming conventions
2. **Security**: Only set `public = true` when absolutely necessary
3. **Policies**: Use bucket policies for fine-grained access control instead of making buckets public
4. **Regions**: Choose regions close to your users for better performance
5. **Lifecycle**: Consider implementing lifecycle policies for object management

## Running the Example

1. Configure your Thalassa Cloud provider credentials
2. Initialize Terraform: `terraform init`
3. Review the plan: `terraform plan`
4. Apply the configuration: `terraform apply`

## Cleanup

To destroy the resources created by this example:

```bash
terraform destroy
```

## Related Resources

- [Thalassa Cloud Documentation](https://docs.thalassa.cloud)
- [S3 Bucket Policy Examples](https://docs.aws.amazon.com/AmazonS3/latest/userguide/example-bucket-policies.html) 