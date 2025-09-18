data "thalassa_iam_team" "example" {
  # organisation_id is optional - if not provided, the organisation from the provider will be used
  # organisation_id = "org-123456"
  
  # You can search by name or slug
  name = "example-team"
  # slug = "example-team"
}

output "team_details" {
  value = data.thalassa_iam_team.example
}

output "team_members" {
  value = data.thalassa_iam_team.example.members
}

output "team_member_count" {
  value = length(data.thalassa_iam_team.example.members)
}

output "team_admin_members" {
  value = [
    for member in data.thalassa_iam_team.example.members : member
    if member.role == "admin"
  ]
}
