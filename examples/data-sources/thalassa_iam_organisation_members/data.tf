data "thalassa_iam_organisation_members" "example" {
  # organisation_id is optional - if not provided, the organisation from the provider will be used
  # organisation_id = "org-123456"
}

data "thalassa_iam_organisation_members" "filtered" {
  # Filter members by email address
  email_filter = "admin@example.com"
}

output "organisation_members" {
  value = data.thalassa_iam_organisation_members.example.members
}

output "member_count" {
  value = length(data.thalassa_iam_organisation_members.example.members)
}

output "owners" {
  value = [
    for member in data.thalassa_iam_organisation_members.example.members : member
    if member.member_type == "OWNER"
  ]
}

output "members" {
  value = [
    for member in data.thalassa_iam_organisation_members.example.members : member
    if member.member_type == "MEMBER"
  ]
}

output "user_emails" {
  value = [
    for member in data.thalassa_iam_organisation_members.example.members : member.user[0].email
    if length(member.user) > 0
  ]
}

output "filtered_members" {
  value = data.thalassa_iam_organisation_members.filtered.members
}

output "filtered_member_count" {
  value = length(data.thalassa_iam_organisation_members.filtered.members)
}
