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

  # Add team members
  # members {
  #   email = "example@example.com"
  # }
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

output "team_members" {
  value = thalassa_iam_team.example.members
}

output "team_member_count" {
  value = length(thalassa_iam_team.example.members)
}
