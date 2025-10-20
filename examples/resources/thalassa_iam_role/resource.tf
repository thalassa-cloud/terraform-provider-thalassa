# Create an organisation role
resource "thalassa_iam_role" "example" {
  name        = "example-role"
  description = "An example organisation role for demonstration purposes"

  labels = {
    environment = "development"
    project     = "example"
  }

  annotations = {
    "example.com/created-by" = "terraform"
  }

  # Define permission rules
  rules {
    resources   = ["cloud_vpc", "cloud_subnet"]
    permissions = ["read", "list"]
    note        = "Allow read access to VPCs and subnets"
  }
  rules {
    resources           = ["cloud_vpc"]
    resource_identities = ["vpc-123", "vpc-456"]
    permissions         = ["update", "delete"]
    note                = "Allow update/delete for specific VPCs"
  }
}

# Output the role details
output "role_id" {
  value = thalassa_iam_role.example.id
}

output "role_name" {
  value = thalassa_iam_role.example.name
}

output "role_slug" {
  value = thalassa_iam_role.example.slug
}

output "role_description" {
  value = thalassa_iam_role.example.description
}

output "role_rules" {
  value = thalassa_iam_role.example.rules
}

