# Thalassa Object Storage Bucket Data Source Example

This example demonstrates how to use the `thalassa_objectstorage_bucket` data source to reference existing object storage buckets in your Terraform configurations.

## Overview

The `thalassa_objectstorage_bucket` data source allows you to retrieve information about existing object storage buckets without creating new ones. This is useful for referencing existing infrastructure or getting bucket details for use in other resources.

## Use Cases

- **Reference Existing Buckets**: Use existing buckets in your Terraform configuration
- **Cross-Resource References**: Reference bucket properties in other resources
- **Infrastructure Discovery**: Get details about buckets created outside of Terraform
- **Regional Consistency**: Create new resources in the same region as existing buckets

## Example Usage

### Basic Data Source Reference

```hcl
data "thalassa_objectstorage_bucket" "existing" {
  name   = "existing-bucket-name"
  region = "nl-01"
}
```

### Reference by Name Only

```hcl
data "thalassa_objectstorage_bucket" "by_name" {
  name = "another-existing-bucket"
}
```

### Using Data Source in Other Resources

```hcl
# Reference an existing bucket
data "thalassa_objectstorage_bucket" "existing" {
  name   = "existing-bucket-name"
  region = "nl-01"
}

# Create a new bucket in the same region
resource "thalassa_objectstorage_bucket" "new_in_same_region" {
  name   = "new-bucket-in-same-region"
  region = data.thalassa_objectstorage_bucket.existing.region
  public = false
}
```

## Data Source Arguments

| Argument | Type | Required | Description |
|----------|------|----------|-------------|
| `name` | string | Yes | The name of the bucket to reference |
| `region` | string | No | The region of the bucket (optional, will search across regions if not specified) |
| `organisation_id` | string | No | Organisation ID (uses provider default if not specified) |

## Data Source Attributes

| Attribute | Type | Description |
|-----------|------|-------------|
| `id` | string | The unique identifier of the bucket |
| `name` | string | The name of the bucket |
| `region` | string | The region of the bucket |
| `public` | bool | Whether the bucket is publicly accessible |
| `policy` | string | The bucket policy as a JSON string |
| `status` | string | The current status of the bucket |
| `endpoint` | string | The endpoint URL for accessing the bucket |

## Common Patterns

### 1. Regional Consistency

Use existing buckets to ensure new resources are created in the same region:

```hcl
data "thalassa_objectstorage_bucket" "production" {
  name   = "production-bucket"
  region = "nl-01"
}

resource "thalassa_objectstorage_bucket" "staging" {
  name   = "staging-bucket"
  region = data.thalassa_objectstorage_bucket.production.region
}
```

### 2. Policy Reference

Reference existing bucket policies for consistency:

```hcl
data "thalassa_objectstorage_bucket" "template" {
  name = "template-bucket"
}

resource "thalassa_objectstorage_bucket" "new" {
  name   = "new-bucket"
  region = "nl-01"
  policy = data.thalassa_objectstorage_bucket.template.policy
}
```

### 3. Cross-Resource Integration

Use bucket information in other resources:

```hcl
data "thalassa_objectstorage_bucket" "logs" {
  name = "application-logs"
}

resource "thalassa_virtual_machine_instance" "app" {
  name = "application-server"
  # ... other configuration ...
  
  # Use bucket endpoint in user data or configuration
  cloud_init = <<-EOF
    #!/bin/bash
    echo "Log bucket endpoint: ${data.thalassa_objectstorage_bucket.logs.endpoint}"
  EOF
}
```

## Best Practices

1. **Specificity**: Always specify the region when you know it to avoid ambiguity
2. **Error Handling**: Handle cases where the referenced bucket doesn't exist
3. **Documentation**: Document why you're referencing specific buckets
4. **Consistency**: Use consistent naming patterns for referenced buckets

## Running the Example

1. Ensure you have existing buckets in your Thalassa Cloud account
2. Update the bucket names in the example to match your existing buckets
3. Configure your Thalassa Cloud provider credentials
4. Initialize Terraform: `terraform init`
5. Review the plan: `terraform plan`
6. Apply the configuration: `terraform apply`

## Troubleshooting

### Bucket Not Found

If you get a "bucket not found" error:

1. Verify the bucket name is correct
2. Check if the bucket exists in the specified region
3. Ensure you have the necessary permissions to access the bucket
4. Try removing the region argument to search across all regions

### Multiple Buckets Found

If multiple buckets with the same name exist across regions:

1. Specify the region argument to narrow down the search
2. Use a more specific bucket name
3. Check your bucket naming conventions

## Related Resources

- [Thalassa Object Storage Bucket Resource](../thalassa_objectstorage_bucket/)
- [Thalassa Cloud Documentation](https://docs.thalassa.cloud)
- [Terraform Data Sources](https://www.terraform.io/docs/language/data-sources/index.html) 