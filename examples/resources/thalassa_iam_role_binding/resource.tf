terraform {
  required_providers {
    thalassa = {
      source = "thalassa-cloud/thalassa"
    }
  }
}

provider "thalassa" {
  # Configuration options
}

# Create a role first
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
  rules = [
    {
      resources   = ["cloud_vpc", "cloud_subnet"]
      permissions = ["read", "list"]
      note        = "Allow read access to VPCs and subnets"
    }
  ]
}

# Create a role binding for a user
resource "thalassa_iam_role_binding" "user_binding" {
  role_id     = thalassa_iam_role.example.id
  name        = "admin-user-binding"
  description = "Bind admin user to the example role"
  user_id     = "user-id-123"

  labels = {
    purpose = "admin-access"
  }

  annotations = {
    "example.com/binding-type" = "user"
  }
}

# Create a role binding for a team
resource "thalassa_iam_role_binding" "team_binding" {
  role_id     = thalassa_iam_role.example.id
  name        = "devops-team-binding"
  description = "Bind devops team to the example role"
  team_id     = "team-id-456"

  labels = {
    purpose = "team-access"
  }

  annotations = {
    "example.com/binding-type" = "team"
  }
}

# Output the role binding details
output "user_binding_id" {
  value = thalassa_iam_role_binding.user_binding.id
}

output "user_binding_name" {
  value = thalassa_iam_role_binding.user_binding.name
}

output "team_binding_id" {
  value = thalassa_iam_role_binding.team_binding.id
}

output "team_binding_name" {
  value = thalassa_iam_role_binding.team_binding.name
}
