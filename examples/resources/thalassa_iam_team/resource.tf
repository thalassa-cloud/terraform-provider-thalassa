# Configure the Thalassa Cloud Provider
terraform {
  required_providers {
    thalassa = {
      source = "thalassa.cloud/thalassa/thalassa"
    }
  }
}

provider "thalassa" {
  # Configure your provider here
  # token = "your-api-token"
  # api = "https://api.thalassa.cloud"
  # organisation_id = "your-organisation-id"
}

# Create a team
resource "thalassa_iam_team" "example" {
  name        = "example-team"
  description = "An example team for demonstration purposes"
  
  labels = {
    environment = "development"
    project     = "example"
  }
  
  annotations = {
    contact = "team@example.com"
    owner   = "devops"
  }
}

# Output the team details
output "team_id" {
  value = thalassa_iam_team.example.id
}

output "team_name" {
  value = thalassa_iam_team.example.name
}

output "team_slug" {
  value = thalassa_iam_team.example.slug
}

output "team_description" {
  value = thalassa_iam_team.example.description
}
